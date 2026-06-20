package textutil

import (
	"html"
	"regexp"
	"strings"
)

// reChrome strips leftover site navigation chrome from text snippets.
var (
	reCodeBlock = regexp.MustCompile("(?s)```.*?```")
	reHTMLTag   = regexp.MustCompile(`<[^>]+>`)
	reHeading   = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reHorizRule = regexp.MustCompile(`(?m)^---+\s*$`)
)

// CleanProse strips markdown formatting and HTML tags from text, producing
// clean prose suitable for snippets and previews.
func CleanProse(s string) string {
	// Remove code blocks
	s = reCodeBlock.ReplaceAllString(s, " ")
	// Remove HTML tags
	s = reHTMLTag.ReplaceAllString(s, "")
	// Remove heading markers
	s = reHeading.ReplaceAllString(s, "")
	// Remove horizontal rules
	s = reHorizRule.ReplaceAllString(s, "")
	// Remove markdown table pipes
	s = StripTablePipes(s)
	// Decode HTML entities
	s = DecodeEntities(s)
	return s
}

// StripTablePipes removes markdown table formatting from text: it drops table
// separator rows (| --- | --- |) and trims leading/trailing pipes from each
// remaining line.
func StripTablePipes(s string) string {
	lines := strings.Split(s, "\n")
	var out []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip table separator rows (| --- | --- |)
		if strings.HasPrefix(line, "|") && regexp.MustCompile(`^\|[\s\-:|]+\|?$`).MatchString(line) {
			continue
		}
		// Strip leading/trailing pipes
		line = strings.TrimPrefix(line, "|")
		line = strings.TrimSuffix(line, "|")
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// DecodeEntities replaces common HTML entities with their literal characters.
func DecodeEntities(s string) string {
	// Common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", `"`)
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s
}

// Excerpt returns a clean, single-line preview of a page's body text for
// homepage cards and search snippets. CleanProse strips markdown/HTML noise.
// If title appears in the cleaned body, the excerpt starts after the title.
// The result is truncated to n characters at a word boundary with a trailing …
// when truncated.
func Excerpt(body, title string, n int) string {
	s := CleanProse(body)
	if title != "" {
		lt := strings.ToLower(title)
		if i := strings.Index(strings.ToLower(s), lt); i >= 0 {
			if rest := strings.TrimSpace(s[i+len(title):]); rest != "" {
				s = rest
			}
		}
	}
	// Collapse whitespace
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > n {
		// Try to cut at a word boundary
		if i := strings.LastIndex(s[:n], " "); i > n-30 {
			s = s[:i]
		} else {
			s = s[:n]
		}
		s += "…"
	}
	return s
}

// MakeSnippet extracts a window of words around the first matched term and
// wraps matched terms in <mark>. Starts from just after the page title.
// window controls the number of words included; ellipses are added when the
// snippet does not start at the first word or end at the last.
func MakeSnippet(body, title string, terms []string, window int) string {
	if body == "" {
		return ""
	}
	body = CleanProse(body)
	if title != "" {
		lt := strings.ToLower(title)
		if i := strings.Index(strings.ToLower(body), lt); i >= 0 {
			if rest := strings.TrimSpace(body[i+len(title):]); rest != "" {
				body = rest
			}
		}
	}
	tset := map[string]bool{}
	for _, t := range terms {
		tset[strings.ToLower(t)] = true
	}
	words := strings.Fields(body)
	start := 0
	for i, w := range words {
		clean := strings.ToLower(strings.Trim(w, ".,;:!?()[]\"'"))
		if tset[clean] {
			start = i - window/2
			break
		}
	}
	if start < 0 {
		start = 0
	}
	end := start + window
	if end > len(words) {
		end = len(words)
	}
	var b strings.Builder
	if start > 0 {
		b.WriteString("… ")
	}
	for i := start; i < end; i++ {
		w := words[i]
		clean := strings.ToLower(strings.Trim(w, ".,;:!?()[]\"'"))
		if tset[clean] {
			b.WriteString("<mark>")
			b.WriteString(html.EscapeString(w))
			b.WriteString("</mark> ")
		} else {
			b.WriteString(html.EscapeString(w))
			b.WriteByte(' ')
		}
	}
	if end < len(words) {
		b.WriteString("…")
	}
	return strings.TrimSpace(b.String())
}
