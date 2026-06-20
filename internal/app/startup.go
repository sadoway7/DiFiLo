package app

import (
	"os"
	"path/filepath"
	"strings"
)

// BuildHeroImages gathers a pool of larger, photo-like images under images/
// to use as random homepage hero backgrounds. Logos and small thumbnails are
// filtered out using s.Config.HeroMinSize.
func (s *Server) BuildHeroImages() {
	minSize := s.Config.HeroMinSize
	if minSize <= 0 {
		minSize = 80 * 1024
	}
	var imgs []string
	_ = filepath.WalkDir(s.imageDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())
		if strings.HasPrefix(name, "logo") {
			return nil
		}
		if !strings.HasSuffix(name, ".webp") && !strings.HasSuffix(name, ".jpg") &&
			!strings.HasSuffix(name, ".jpeg") && !strings.HasSuffix(name, ".png") {
			return nil
		}
		fi, err := d.Info()
		if err != nil || fi.Size() < minSize {
			return nil
		}
		rel, err := filepath.Rel(s.imageDir, p)
		if err != nil {
			return nil
		}
		imgs = append(imgs, "/images/"+filepath.ToSlash(rel))
		return nil
	})
	s.heroImages = imgs
	logger.Info("startup", "hero_images", len(imgs))
}

// BuildAliases loads the slug→canonical-route alias map from the database.
// It is used by handlePage and handleDownload to resolve legacy slug URLs
// (e.g. /material/cobalt+oxide) to their canonical numeric routes.
func (s *Server) BuildAliases() {
	s.aliases = s.DB.LoadAliases()
	logger.Info("startup", "aliases", len(s.aliases))
}
