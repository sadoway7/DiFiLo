package content

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmHtml "github.com/yuin/goldmark/renderer/html"
)

// markdownEngine is the shared GFM markdown renderer used to turn body_md into
// HTML. Tables, strikethrough and autolinks are enabled; raw HTML in the
// markdown is passed through (unsafe) since the source is trusted content.
var markdownEngine = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(gmHtml.WithUnsafe()),
)

// pre-processing patterns for URL + image rewriting
var (
	// strip digitalfire.com domain from links: https://digitalfire.com/material/925 -> /material/925
	reDigitalfireHTTPS = regexp.MustCompile(`\((?:https?://)?(?:www\.)?digitalfire\.com`)
	// bare domain without protocol: (digitalfire.com/glossary/x -> (/glossary/x
	reDigitalfireBare = regexp.MustCompile(`\(digitalfire\.com`)
	// relative image paths: ../../images/ or ../images/ -> /images/
	reRelImages = regexp.MustCompile(`\((?:\.\.\/)+(images\/)`)
	// bare image ref without leading slash: (images/ -> (/images/
	reBareImages = regexp.MustCompile(`\((images\/)`)
	// leading H1 heading: # Title
	reLeadingH1 = regexp.MustCompile(`(?m)^#\s+.+\s*\n`)
)

// StripLeadingH1 removes the first "# Title" line from markdown, since the
// wiki template header already displays the page title as an <h1>.
func StripLeadingH1(mdText string) string {
	// Only strip if it's at the very start (after optional whitespace).
	trimmed := strings.TrimLeft(mdText, " \t\r\n")
	if !strings.HasPrefix(trimmed, "# ") {
		return mdText
	}
	loc := reLeadingH1.FindStringIndex(trimmed)
	if loc == nil {
		return mdText
	}
	return strings.TrimSpace(trimmed[loc[1]:])
}

// RenderMarkdown converts body_md (GFM markdown) to HTML, rewriting all
// internal links and image paths so they resolve inside the offline app.
//
// Link rewriting rules:
//   - https://digitalfire.com/material/925  ->  /material/925
//   - digitalfire.com/glossary/plasticity   ->  /glossary/plasticity
//   - ../../images/foo.jpg                  ->  /images/foo.jpg
//
// External links (anything not digitalfire.com) get target="_blank" rel="noopener".
func RenderMarkdown(mdText string) string {
	processed := rewriteMarkdownURLs(mdText)

	var buf bytes.Buffer
	if err := markdownEngine.Convert([]byte(processed), &buf); err != nil {
		return strings.ReplaceAll(processed, "\n", "<br>\n")
	}
	htmlOut := buf.String()

	// Post-process: add target="_blank" to external links.
	htmlOut = makeExternalLinksBlank(htmlOut)

	// Post-process: remove empty section headings (h2-h6 with no content
	// before the next heading or end of document).
	htmlOut = stripEmptyHTMLSections(htmlOut)

	return htmlOut
}

// rewriteMarkdownURLs fixes all internal links and image paths in the raw
// markdown text before it is parsed by goldmark.
func rewriteMarkdownURLs(mdText string) string {
	// Strip digitalfire.com domains so links become relative.
	mdText = reDigitalfireHTTPS.ReplaceAllString(mdText, "(")
	mdText = reDigitalfireBare.ReplaceAllString(mdText, "(")
	// Fix relative image paths.
	mdText = reRelImages.ReplaceAllString(mdText, "(/$1")
	mdText = reBareImages.ReplaceAllString(mdText, "(/$1")
	return mdText
}

// reExtLink matches <a href="http..."> (external links only).
var reExtLink = regexp.MustCompile(`<a href="(https?://[^"]+)"`)

// makeExternalLinksBlank adds target="_blank" rel="noopener" to all external links.
func makeExternalLinksBlank(htmlStr string) string {
	return reExtLink.ReplaceAllString(htmlStr, `<a href="$1" target="_blank" rel="noopener"`)
}

// reHeadingTag matches any heading tag opening (h1-h6, with optional id attr).
var reHeadingTag = regexp.MustCompile(`(?i)<h[2-6][^>]*>.*?</h[2-6]>\s*`)

// stripEmptyHTMLSections removes heading elements (h2-h6) that have no content
// between them and the next heading or end of string. Runs in a loop because
// removing one empty heading may expose the previous one as also empty.
func stripEmptyHTMLSections(htmlStr string) string {
	for {
		locs := reHeadingTag.FindAllStringIndex(htmlStr, -1)
		if len(locs) == 0 {
			break
		}
		changed := false
		for i := 0; i < len(locs); i++ {
			start := locs[i][0]
			end := locs[i][1]
			// Check what's between this heading and the next heading (or end)
			var nextStart int
			if i+1 < len(locs) {
				nextStart = locs[i+1][0]
			} else {
				nextStart = len(htmlStr)
			}
			between := strings.TrimSpace(htmlStr[end:nextStart])
			if between == "" {
				// Empty section — remove this heading
				htmlStr = htmlStr[:start] + htmlStr[end:]
				changed = true
				break // restart scan since indices shifted
			}
		}
		if !changed {
			break
		}
	}
	return htmlStr
}
