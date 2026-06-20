package content

import (
	"regexp"
	"strings"

	"difilo/internal/db"
)

// reLinkTable matches a "Related Links" table row:
//   | Type | [Title](url) blurb |
// Group 1 = type, group 2 = title, group 3 = url, group 4 = trailing blurb.
var reLinkTable = regexp.MustCompile(`(?m)^\|\s*([^|]+)\s*\|\s*\[([^\]]*)\]\(([^)]+)\)(.*?)\s*\|\s*$`)

// ExtractPageLinks parses the "Related Links" tables at the bottom of a page
// body into link rows. Each link's target route is derived from its URL.
func ExtractPageLinks(body string) []db.LinkRow {
	var out []db.LinkRow
	for _, m := range reLinkTable.FindAllStringSubmatch(body, -1) {
		targetURL := strings.TrimSpace(m[3])
		out = append(out, db.LinkRow{
			TargetType:  strings.TrimSpace(m[1]),
			TargetTitle: strings.TrimSpace(m[2]),
			TargetURL:   targetURL,
			Blurb:       strings.TrimSpace(stripInlineFormatting(m[4])),
			TargetRoute: RouteOf(targetURL),
		})
	}
	return out
}
