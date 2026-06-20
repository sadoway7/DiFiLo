package content

import (
	"regexp"
	"strings"
)

// Frontmatter holds the YAML frontmatter fields parsed from a .md file.
type Frontmatter struct {
	URL       string
	Title     string
	WaybackTS string
	HTMLSHA1  string
}

var (
	reFMBlock = regexp.MustCompile(`(?s)^\s*---\n(.*?)\n---\n`)
	reFMField = regexp.MustCompile(`(?m)^(\w+):\s*(.*)$`)
	// reMetaDesc extracts the <meta name="description" content="..."> value.
	reMetaDesc = regexp.MustCompile(`(?i)<meta\s+name\s*=\s*["']description["']\s+content\s*=\s*["']([^"']*)`)

	// Maintainer extraction patterns.
	reMaintainerItalic = regexp.MustCompile(`By\s+\*([^*]+)\*\s*Follow`)
	reMaintainerPlain  = regexp.MustCompile(`By\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+){1,2})\s+Follow`)
	reMaintainerByLine = regexp.MustCompile(`(?:^|\n)\s*By\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+){1,2})\s*(?:\n|$)`)
)

// ParseFrontmatter extracts YAML frontmatter fields from a markdown file.
func ParseFrontmatter(mdText string) Frontmatter {
	var fm Frontmatter
	m := reFMBlock.FindStringSubmatch(mdText)
	if m == nil {
		return fm
	}
	for _, line := range strings.Split(m[1], "\n") {
		fmMatch := reFMField.FindStringSubmatch(line)
		if fmMatch == nil {
			continue
		}
		key := fmMatch[1]
		val := strings.TrimSpace(fmMatch[2])
		val = strings.Trim(val, `"'`)
		switch key {
		case "url":
			fm.URL = val
		case "title":
			fm.Title = val
		case "wayback_ts":
			fm.WaybackTS = val
		case "html_sha1":
			fm.HTMLSHA1 = val
		}
	}
	return fm
}

// extractMetaDescription pulls the <meta name="description" content="..."> value
// from a raw HTML byte slice. Returns "" when absent.
func extractMetaDescription(htmlData []byte) string {
	if m := reMetaDesc.FindSubmatch(htmlData); m != nil {
		return strings.TrimSpace(string(m[1]))
	}
	return ""
}

// StripFooterChrome removes the recurring footer that appears at the bottom of
// captured pages: ko-fi link, "Got a Question?", author follow line, copyright,
// privacy policy, reference library image, etc.
func StripFooterChrome(s string) string {
	cutPatterns := []string{
		"| By *",
		"By *Tony Hansen* Follow me on",
		"By Tony Hansen Follow me on",
		"### Got a Question?",
		"Got a Question?",
		"Buy me a coffee and we can talk",
		"<https://digitalfire.com>, All Rights Reserved",
		"https://digitalfire.com, All Rights Reserved",
		"All Rights Reserved",
		"ReferenceLibrary.svg",
		"[Privacy Policy]",
		"Privacy Policy](/",
	}
	lower := strings.ToLower(s)
	cutAt := len(s)
	for _, p := range cutPatterns {
		if idx := strings.Index(lower, strings.ToLower(p)); idx >= 0 && idx < cutAt {
			// Back up to the start of the table row or blank line before the marker
			cutAt = idx
		}
	}
	if cutAt < len(s) {
		return strings.TrimSpace(s[:cutAt])
	}
	return s
}

// ExtractMaintainer finds the author/maintainer name from the footer chrome
// before it's stripped. Looks for "By *Name*" and "By Name Follow" patterns.
// Defaults to "Tony Hansen" if no other author is found.
func ExtractMaintainer(body string) string {
	// Pattern 1: | By *Name* Follow me on | (italicized markdown)
	if m := reMaintainerItalic.FindStringSubmatch(body); m != nil {
		name := strings.TrimSpace(m[1])
		if isValidAuthorName(name) {
			return name
		}
	}
	// Pattern 2: By Name Follow (plain text, stop before "Follow")
	if m := reMaintainerPlain.FindStringSubmatch(body); m != nil {
		name := strings.TrimSpace(m[1])
		if isValidAuthorName(name) {
			return name
		}
	}
	// Pattern 3: By Name on its own line near the bottom
	if m := reMaintainerByLine.FindStringSubmatch(body); m != nil {
		name := strings.TrimSpace(m[1])
		if isValidAuthorName(name) {
			return name
		}
	}
	// Default: Tony Hansen created essentially everything on Digitalfire
	return "Tony Hansen"
}

// isValidAuthorName filters out false positives from the author regex.
func isValidAuthorName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 3 || len(name) > 40 {
		return false
	}
	// Reject if it contains non-name words
	lower := strings.ToLower(name)
	bad := []string{"the bay", "follow", "related", "question", "rights", "privacy"}
	for _, b := range bad {
		if strings.Contains(lower, b) {
			return false
		}
	}
	return true
}
