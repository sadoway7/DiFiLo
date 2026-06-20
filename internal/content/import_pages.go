package content

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"difilo/internal/db"
	"difilo/internal/textutil"
)

// reH1 matches the start of a level-1 heading line. Used to locate the first
// real content heading so nav chrome above it can be dropped.
var reH1 = regexp.MustCompile(`(?m)^#\s+`)

// ProcessedPage bundles a fully-built page row with the images and links
// extracted from its body, ready for batch insertion by ImportContent.
type ProcessedPage struct {
	Page   db.ContentPage
	Images []db.ImageRow
	Links  []db.LinkRow
}

// ProcessPage reads a single page's md (+ optional html) files from disk,
// strips the captured-site chrome, and assembles a ProcessedPage (page row +
// extracted images and links). It does not touch the database; the caller
// (ImportContent) performs the batched transactional inserts. A non-nil error
// means the page should be skipped.
func ProcessPage(pg PageManifest, mdDir, htmlDir string, imgIndex map[string]string) (*ProcessedPage, error) {
	// --- Read the md file ---
	mdRel := DeriveMDPath(pg.URL)
	mdPath := filepath.Join(mdDir, filepath.FromSlash(mdRel))
	mdRaw, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, fmt.Errorf("read md %s: %w", mdRel, err)
	}
	mdRaw = GunzipIfCompressed(mdRaw)
	mdFull := string(mdRaw)

	// Parse frontmatter (manifest fields take a back-seat to frontmatter).
	fm := ParseFrontmatter(mdFull)
	fm.URL = pg.URL
	if fm.Title == "" {
		fm.Title = pg.Title
	}

	// Strip frontmatter.
	body := reFMBlock.ReplaceAllString(mdFull, "")

	// Strip nav chrome: everything before the first "# " heading.
	contentStart := reH1.FindStringIndex(body)
	var contentBody string
	if contentStart != nil {
		contentBody = body[contentStart[0]:]
	} else {
		contentBody = body
	}

	// Extract meta description from html (needed before footer strip).
	metaDesc := ""
	htmlRel := DeriveHTMLPath(pg.URL)
	htmlPath := filepath.Join(htmlDir, filepath.FromSlash(htmlRel))
	if htmlRaw, err := os.ReadFile(htmlPath); err == nil {
		htmlRaw = GunzipIfCompressed(htmlRaw)
		metaDesc = extractMetaDescription(htmlRaw)
	}

	// Extract maintainer from footer BEFORE stripping it.
	maintainer := ExtractMaintainer(contentBody)

	// Strip trailing footer chrome and the empty structures it leaves behind.
	contentBody = StripFooterChrome(contentBody)
	contentBody = stripTrailingEmpty(contentBody)
	contentBody = stripEmptySections(contentBody)

	// Strip leading H1 (template renders the title separately).
	contentBody = StripLeadingH1(contentBody)

	// Strip leading paragraph if it duplicates the meta description.
	if metaDesc != "" {
		contentBody = stripLeadingDupe(contentBody, metaDesc)
	}

	// Build cleaned prose text for search and the raw markdown body.
	bodyText := textutil.CleanProse(contentBody)
	bodyMd := contentBody

	// Extract thumbnail from cleaned body (after footer strip).
	thumb := ExtractThumb(contentBody)

	// Build route + slug.
	route := RouteOf(pg.URL)
	if route == "" {
		route = "/"
	}
	slug := strings.TrimPrefix(route, "/")
	if idx := strings.Index(slug, "/"); idx >= 0 {
		slug = slug[idx+1:]
	}

	sortTitle := strings.ToLower(strings.TrimSpace(fm.Title))
	wordCount := len(strings.Fields(bodyText))

	page := db.ContentPage{
		Section:         pg.Section,
		Title:           fm.Title,
		Slug:            slug,
		SourceURL:       pg.URL,
		Route:           route,
		WaybackTS:       fm.WaybackTS,
		HTMLSHA1:        fm.HTMLSHA1,
		BodyText:        bodyText,
		BodyMD:          bodyMd,
		Thumb:           thumb,
		MetaDescription: metaDesc,
		AuthorByline:    maintainer,
		Status:          "published",
		SortTitle:       sortTitle,
		WordCount:       wordCount,
	}

	return &ProcessedPage{
		Page:   page,
		Images: ExtractPageImages(bodyMd, thumb, imgIndex),
		Links:  ExtractPageLinks(bodyMd),
	}, nil
}

// stripTrailingEmpty removes empty markdown structures left at the end of the
// body after footer chrome stripping: empty table rows, empty table headers,
// dangling section headings with no content after them, and excess blank lines.
func stripTrailingEmpty(s string) string {
	lines := strings.Split(s, "\n")
	for len(lines) > 0 {
		last := strings.TrimSpace(lines[len(lines)-1])
		if last == "" {
			lines = lines[:len(lines)-1]
			continue
		}
		// Empty table rows: |  |  |  or | | |
		if isTableRow(last) {
			lines = lines[:len(lines)-1]
			continue
		}
		// Table separator: | --- | --- |
		if isTableSeparator(last) {
			lines = lines[:len(lines)-1]
			continue
		}
		// Dangling empty headings: ### Links, ## Related Information
		if isHeading(last) {
			lines = lines[:len(lines)-1]
			continue
		}
		// Horizontal rules
		if isHorizRule(last) {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func isTableRow(s string) bool       { b, _ := regexp.MatchString(`^\|[\s|]*\|$`, s); return b }
func isTableSeparator(s string) bool { b, _ := regexp.MatchString(`^\|[\s\-:|]+\|$`, s); return b }
func isHeading(s string) bool        { b, _ := regexp.MatchString(`^#{1,6}\s+`, s); return b }
func isHorizRule(s string) bool      { b, _ := regexp.MatchString(`^-{3,}$`, s); return b }

// stripEmptySections removes any markdown heading (## or deeper) that has no
// content between it and the next heading or end of document. Runs repeatedly
// until no more empty sections are found.
func stripEmptySections(s string) string {
	reHeadingLine := regexp.MustCompile(`^(#{1,6})\s+`)
	for {
		lines := strings.Split(s, "\n")
		var out []string
		changed := false
		for i := 0; i < len(lines); i++ {
			if reHeadingLine.MatchString(strings.TrimSpace(lines[i])) {
				// Look ahead — is there content before the next heading?
				hasContent := false
				for j := i + 1; j < len(lines); j++ {
					trimmed := strings.TrimSpace(lines[j])
					if trimmed == "" {
						continue
					}
					if reHeadingLine.MatchString(trimmed) {
						// Hit next heading with no content — this heading is empty
						break
					}
					// Found actual content
					hasContent = true
					break
				}
				if !hasContent {
					changed = true
					continue // skip this heading
				}
			}
			out = append(out, lines[i])
		}
		s = strings.Join(out, "\n")
		if !changed {
			break
		}
	}
	return strings.TrimSpace(s)
}

// stripLeadingDupe removes the first paragraph of the body if it closely
// matches the meta description (which the template already shows in the header).
func stripLeadingDupe(body, metaDesc string) string {
	body = strings.TrimSpace(body)
	if body == "" || metaDesc == "" {
		return body
	}
	// Strip leading "#### Description" heading if present.
	body = strings.TrimPrefix(body, "#### Description")
	body = strings.TrimPrefix(body, "### Description")
	body = strings.TrimSpace(body)

	// Get first paragraph.
	firstPara := body
	if i := strings.Index(body, "\n\n"); i >= 0 {
		firstPara = body[:i]
	}
	// Strip markdown formatting for comparison.
	cleanFirst := strings.TrimSpace(stripInlineFormatting(firstPara))
	cleanMeta := strings.TrimSpace(metaDesc)

	if len(cleanFirst) > 20 && len(cleanMeta) > 20 {
		if cleanFirst == cleanMeta ||
			strings.Contains(cleanFirst, cleanMeta) ||
			strings.Contains(cleanMeta, cleanFirst) {
			rest := body[len(firstPara):]
			return strings.TrimSpace(rest)
		}
	}
	return body
}
