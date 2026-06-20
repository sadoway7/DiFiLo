package content

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// Shared markdown reference patterns used across the import pipeline.
var (
	// mdImgRe captures the image src of the first ![alt](src) reference.
	mdImgRe = regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	// mdLinkRe matches a markdown link, capturing the link text.
	mdLinkRe = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
)

// GunzipIfCompressed transparently decompresses gzip data (detected by the
// magic header 1f 8b). Non-gzip data is returned unchanged. Mirror files may
// be stored gzipped to save space.
func GunzipIfCompressed(b []byte) []byte {
	if len(b) < 2 || b[0] != 0x1f || b[1] != 0x8b {
		return b
	}
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return b
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		return b
	}
	return out
}

// ExtractThumb finds the first image reference in markdown text and returns
// it as a /images/... path (or "" if no image is found).
func ExtractThumb(md string) string {
	m := mdImgRe.FindStringSubmatch(md)
	if m == nil {
		return ""
	}
	return ResolveImagePath(m[1])
}

// RouteOf converts a full URL (e.g. https://digitalfire.com/material/925) to
// a site-relative route (e.g. /material/925).
func RouteOf(rawURL string) string {
	u := rawURL
	// Strip protocol.
	if i := strings.Index(u, "://"); i >= 0 {
		u = u[i+3:]
	}
	// Strip domain.
	if i := strings.IndexByte(u, '/'); i >= 0 {
		u = u[i:]
	} else {
		u = "/"
	}
	// Strip query string and fragment.
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		u = u[:i]
	}
	if u == "" {
		u = "/"
	}
	return u
}

// ResolveImagePath converts a markdown image reference to a /images/... path.
// Handles relative paths (../../images/..., images/...) and bare filenames.
func ResolveImagePath(ref string) string {
	ref = strings.TrimSpace(ref)
	// Strip alt text or title: ![alt](src "title") → src
	if i := strings.IndexByte(ref, ' '); i >= 0 {
		ref = ref[:i]
	}
	ref = strings.Trim(ref, `"`)
	// Already absolute path starting with /images/
	if strings.HasPrefix(ref, "/images/") {
		return ref
	}
	// Relative paths: ../../images/foo → /images/foo
	ref = strings.TrimPrefix(ref, "./")
	for strings.HasPrefix(ref, "../") {
		ref = strings.TrimPrefix(ref, "../")
	}
	if strings.HasPrefix(ref, "images/") {
		return "/" + ref
	}
	// Bare filename — assume it's under images/
	if ref != "" && !strings.HasPrefix(ref, "http") {
		return "/images/" + ref
	}
	return ""
}

// stripInlineFormatting removes markdown emphasis and link noise from a blurb.
// Used when extracting clean link blurbs from markdown table cells.
func stripInlineFormatting(s string) string {
	s = mdLinkRe.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "`", "")
	return strings.TrimSpace(s)
}

// --- Path derivation helpers (ported from the scraper / rewrite layer) ---
//
// These mirror the on-disk layout the scraper produced so the import pipeline
// can locate the md/html/image files for a given source URL.

var (
	// reMultiSlash collapses accidental double slashes in a URL path.
	reMultiSlash = regexp.MustCompile(`/+`)
	// reUnsafe mirrors the scraper's safe_path_part sanitiser: anything that is
	// not alphanumeric, dot, underscore, hyphen or slash becomes an underscore.
	reUnsafe = regexp.MustCompile(`[^A-Za-z0-9._\-/]`)
	// reImgExt detects an image extension embedded in a URL.
	reImgExt = regexp.MustCompile(`(?i)\.(jpe?g|png|gif|webp|svg|bmp|tiff?)`)
)

// splitExt returns the path split into base and extension (including the dot).
func splitExt(p string) (base, ext string) {
	ext = path.Ext(p)
	return strings.TrimSuffix(p, ext), ext
}

// stripExt returns p with its extension removed.
func stripExt(p string) string {
	return strings.TrimSuffix(p, path.Ext(p))
}

// SafePathPart ports the scraper's transform of a digitalfire URL into a
// filesystem-safe relative path (no host).
func SafePathPart(rawurl string) string {
	p, err := url.Parse(rawurl)
	if err != nil {
		return "_index"
	}
	pa := p.Path
	if pa == "" {
		pa = "/"
	}
	pa = reMultiSlash.ReplaceAllString(pa, "/")
	pa = strings.TrimPrefix(pa, "/")
	if pa == "" {
		pa = "_index"
	}
	pa = reUnsafe.ReplaceAllString(pa, "_")
	return pa
}

// DeriveHTMLPath mirrors the scraper's derive_html_path: the on-disk file under
// html/ for a given page URL.
func DeriveHTMLPath(rawurl string) string {
	rel := SafePathPart(rawurl)
	base, ext := splitExt(rel)
	el := strings.ToLower(ext)
	if el != ".htm" && el != ".html" && el != ".php" {
		if rel == "" {
			rel = "_index"
		}
		rel = rel + ".html"
	} else {
		if !strings.HasSuffix(base, ".html") {
			rel = base + ".html"
		}
		if !strings.HasSuffix(rel, ".html") {
			rel += ".html"
		}
	}
	return rel
}

// DeriveMDPath mirrors the scraper's derive_md_path: the on-disk .md file for
// a given page URL.
func DeriveMDPath(rawurl string) string {
	rel := SafePathPart(rawurl)
	base := stripExt(rel)
	if base == "" {
		base = "_index"
	}
	return base + ".md"
}

// DeriveImagePath mirrors the scraper's derive_image_path: host/path under
// images/ for a given image URL.
func DeriveImagePath(imgurl string) string {
	p, err := url.Parse(imgurl)
	if err != nil {
		return ""
	}
	host := p.Hostname()
	if host == "" {
		host = "local"
	}
	pa := p.Path
	if pa == "" {
		pa = "/"
	}
	pa = reMultiSlash.ReplaceAllString(pa, "/")
	pa = strings.TrimPrefix(pa, "/")
	if pa == "" {
		pa = "image"
	}
	base, ext := splitExt(pa)
	if ext == "" {
		m := reImgExt.FindStringSubmatch(imgurl)
		if len(m) > 1 {
			ext = "." + strings.ToLower(m[1])
		} else {
			ext = ".bin"
		}
		pa = base + ext
	}
	pa = reUnsafe.ReplaceAllString(pa, "_")
	return path.Join(host, pa)
}
