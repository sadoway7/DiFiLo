package app

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// reAssetAttr matches src="/images/..." or href="/images/..." (local asset
	// references inside HTML attributes).
	reAssetAttr = regexp.MustCompile(`"(\/(?:images|media|vendor)\/[^"]*)"`)
	// reAssetURL matches CSS url(/images/...).
	reAssetURL = regexp.MustCompile(`url\((\/(?:images|media|vendor)\/[^)]+)\)`)
)

// inlineAssets rewrites local /images, /media and /vendor references (in
// src/href attributes and CSS url()) to base64 data URIs, making a document
// fully self-contained. Missing files are left untouched. Results are cached
// per call so a repeated reference is encoded only once.
func (s *Server) inlineAssets(b []byte) []byte {
	cache := map[string][]byte{}
	resolve := func(path string) ([]byte, bool) {
		if v, ok := cache[path]; ok {
			return v, len(v) > 0
		}
		var abs string
		switch {
		case strings.HasPrefix(path, "/images/"):
			abs = filepath.Join(s.imageDir, filepath.FromSlash(strings.TrimPrefix(path, "/images/")))
		case strings.HasPrefix(path, "/media/"):
			abs = filepath.Join(s.mediaDir, filepath.FromSlash(strings.TrimPrefix(path, "/media/")))
		case strings.HasPrefix(path, "/vendor/"):
			abs = filepath.Join(s.vendorDir, filepath.FromSlash(strings.TrimPrefix(path, "/vendor/")))
		default:
			return nil, false
		}
		data, err := os.ReadFile(abs)
		if err != nil {
			cache[path] = nil
			return nil, false
		}
		uri := "data:" + mimeByExt(filepath.Ext(abs)) + ";base64," + base64.StdEncoding.EncodeToString(data)
		out := []byte(uri)
		cache[path] = out
		return out, true
	}

	b = reAssetAttr.ReplaceAllFunc(b, func(m []byte) []byte {
		inner := m[1 : len(m)-1] // strip surrounding quotes
		if uri, ok := resolve(string(inner)); ok {
			out := make([]byte, 0, len(uri)+2)
			out = append(out, '"')
			out = append(out, uri...)
			out = append(out, '"')
			return out
		}
		return m
	})
	b = reAssetURL.ReplaceAllFunc(b, func(m []byte) []byte {
		sub := reAssetURL.FindSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		if uri, ok := resolve(string(sub[1])); ok {
			out := append([]byte("url("), uri...)
			out = append(out, ')')
			return out
		}
		return m
	})
	return b
}

// mimeByExt returns the MIME type for common static-asset extensions, falling
// back to application/octet-stream.
func mimeByExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	case ".tif", ".tiff":
		return "image/tiff"
	case ".css":
		return "text/css"
	case ".js":
		return "text/javascript"
	default:
		return "application/octet-stream"
	}
}

// sanitizeFilename turns an arbitrary title into a safe, readable filename
// stem (alphanumerics, dash, underscore; spaces collapse to dashes).
func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "page"
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	out := b.String()
	if out == "" {
		return "page"
	}
	return out
}
