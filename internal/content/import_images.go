package content

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"difilo/internal/db"
)

// reImgRef matches a markdown image with alt text and source:
// ![alt](src) → group 1 = alt (caption), group 2 = src.
var reImgRef = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// ExtractPageImages scans the page body markdown for image references, resolves
// each to a /images/... web path, de-duplicates them, and marks the page's
// thumbnail (thumb) as primary. Returns image rows ordered by first appearance.
func ExtractPageImages(body, thumb string, imgIndex map[string]string) []db.ImageRow {
	var out []db.ImageRow
	order := 0
	seen := map[string]bool{}
	for _, m := range reImgRef.FindAllStringSubmatch(body, -1) {
		caption := strings.TrimSpace(m[1])
		ref := strings.TrimSpace(m[2])
		resolved := resolveImageRef(ref, imgIndex)
		if resolved == "" || seen[resolved] {
			continue
		}
		seen[resolved] = true
		out = append(out, db.ImageRow{
			ImagePath:   resolved,
			OriginalRef: ref,
			Caption:     caption,
			IsPrimary:   resolved == thumb,
			SortOrder:   order,
		})
		order++
	}
	return out
}

// resolveImageRef maps a raw image reference from markdown to a /images/...
// web path. Uses imgIndex (basename → relative webp path) to remap the
// scraper's ../media/images/ references to their webp siblings.
func resolveImageRef(ref string, imgIndex map[string]string) string {
	// ../../images/... → /images/...
	if strings.HasPrefix(ref, "../../images/") {
		return "/images/" + strings.TrimPrefix(ref, "../../images/")
	}
	// ../media/images/<basename>.<ext> → look up basename in s3 webp
	if strings.HasPrefix(ref, "../media/images/") {
		basename := filepath.Base(ref)
		stem := strings.TrimSuffix(basename, filepath.Ext(basename))
		if resolved, ok := imgIndex[stem]; ok {
			return "/images/" + resolved
		}
		return "" // genuinely missing
	}
	// http(s)://s3-us-west-2.amazonaws.com/... → /images/s3-us-west-2.amazonaws.com/...
	if strings.HasPrefix(ref, "http") {
		rel := DeriveImagePath(ref)
		if rel != "" {
			return "/images/" + rel
		}
		return ""
	}
	// ../media/videos/... → /media/...
	if strings.HasPrefix(ref, "../media/") {
		return "/" + strings.TrimPrefix(ref, "../")
	}
	// Fallback: already a /images/ path
	if strings.HasPrefix(ref, "/images/") {
		return ref
	}
	return ""
}

// buildImageBasenameIndex walks images/ and maps every file's basename
// (without extension) to its relative path. Used to resolve ../media/images/
// references to their webp siblings stored under s3.
func buildImageBasenameIndex(imageDir string) map[string]string {
	index := map[string]string{}
	_ = filepath.WalkDir(imageDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(imageDir, p)
		if err != nil {
			return nil
		}
		stem := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		// Prefer the first occurrence (s3 content pictures)
		if _, exists := index[stem]; !exists {
			index[stem] = filepath.ToSlash(rel)
		}
		return nil
	})
	return index
}
