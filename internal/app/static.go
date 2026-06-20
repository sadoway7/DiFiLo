package app

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// serveStatic serves a file from root, with a .webp fallback for old .jpg/.png
// URLs. Path traversal is prevented by checking that the resolved path stays
// within root. Static assets served here (/images/, /vendor/, /media/) get a
// long-lived Cache-Control header (Gap 16) so browsers and proxies can cache
// them aggressively.
func (s *Server) serveStatic(w http.ResponseWriter, r *http.Request, root, rel string) {
	// Long-cache static assets — they are content-addressed by URL and never
	// change without a path change.
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	rel = strings.TrimPrefix(rel, "/")
	clean := filepath.Clean(filepath.FromSlash(rel))
	abs := filepath.Join(root, clean)
	if !strings.HasPrefix(abs, filepath.Clean(root)+string(os.PathSeparator)) && abs != filepath.Clean(root) {
		http.NotFound(w, r)
		return
	}
	// If the exact file is missing, try a converted .webp sibling so that old
	// .jpg/.png URLs keep resolving after the batch WebP conversion.
	if _, err := os.Stat(abs); err != nil {
		if base := strings.TrimSuffix(clean, filepath.Ext(clean)); base != clean {
			if webpAbs := filepath.Join(root, base+".webp"); nonEmptyFile(webpAbs) {
				w.Header().Set("Content-Type", "image/webp")
				http.ServeFile(w, r, webpAbs)
				return
			}
		}
	}
	http.ServeFile(w, r, abs)
}

// nonEmptyFile reports whether p exists, is a regular file, and is non-empty.
func nonEmptyFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir() && fi.Size() > 0
}
