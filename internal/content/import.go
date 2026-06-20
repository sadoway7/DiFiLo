package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"difilo/internal/db"
)

// PageManifest mirrors the fields we need from mirror/pages.json.
type PageManifest struct {
	URL     string `json:"url"`
	Section string `json:"section"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	View    string `json:"view"`
}

// ContentImportStats reports the results of an import run.
type ContentImportStats struct {
	Pages    int
	Images   int
	Links    int
	Skipped  int
	Duration time.Duration
}

const (
	pageInsertSQL = `INSERT INTO pages
		(section, title, slug, source_url, route, wayback_ts, html_sha1,
		 body_text, body_md, thumb, meta_description, author_byline, status,
		 sort_title, word_count)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	imgInsertSQL = `INSERT INTO page_images
		(page_id, image_path, original_ref, caption, is_primary, sort_order)
		VALUES (?,?,?,?,?,?)`
	linkInsertSQL = `INSERT INTO page_links
		(page_id, target_type, target_url, target_route, target_title, blurb)
		VALUES (?,?,?,?,?,?)`
)

// ImportContent walks mirror/pages.json + mirror/md/ + mirror/html/ and loads
// every page into the SQLite content tables.
//
// Unlike the legacy importer, this does NOT drop the content tables. It deletes
// only previously-imported pages (those with a non-empty source_url), preserving
// any user-created pages and all user data (comments, bookmarks, settings, …).
// The pages_fts triggers clean up the search index; a full FTS rebuild is run at
// the end as a safety net. Inserts are batched inside transactions.
func ImportContent(database *db.DB, mirrorDir string) (*ContentImportStats, error) {
	start := time.Now()
	stats := &ContentImportStats{}

	// --- Clear previously-imported content (preserves user-created pages). ---
	// Only pages with a non-empty source_url are imported content; user-created
	// pages keep source_url = ''. The pages_fts triggers clean up the FTS index.
	if _, err := database.Exec("DELETE FROM page_links WHERE page_id IN (SELECT id FROM pages WHERE source_url != '')"); err != nil {
		return nil, fmt.Errorf("clear page_links: %w", err)
	}
	if _, err := database.Exec("DELETE FROM page_images WHERE page_id IN (SELECT id FROM pages WHERE source_url != '')"); err != nil {
		return nil, fmt.Errorf("clear page_images: %w", err)
	}
	if _, err := database.Exec("DELETE FROM pages WHERE source_url != ''"); err != nil {
		return nil, fmt.Errorf("clear pages: %w", err)
	}

	// --- Load pages.json manifest ---
	pagesPath := filepath.Join(mirrorDir, "pages.json")
	data, err := os.ReadFile(pagesPath)
	if err != nil {
		return nil, fmt.Errorf("reading pages.json: %w", err)
	}
	var pages []PageManifest
	if err := json.Unmarshal(data, &pages); err != nil {
		return nil, fmt.Errorf("parsing pages.json: %w", err)
	}

	// --- Build image basename index for ../media/images/ remapping ---
	imgIndex := buildImageBasenameIndex(filepath.Join(mirrorDir, "images"))

	mdDir := filepath.Join(mirrorDir, "md")
	htmlDir := filepath.Join(mirrorDir, "html")

	// --- Batch inserts inside transactions (re-prepare every N pages) ---
	sqlDB := database.SQL()
	tx, err := sqlDB.Begin()
	if err != nil {
		return nil, err
	}
	pageStmt, err := tx.Prepare(pageInsertSQL)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	imgStmt, err := tx.Prepare(imgInsertSQL)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	linkStmt, err := tx.Prepare(linkInsertSQL)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	batchSize := 0
	const batchCommitEvery = 500

	for _, pg := range pages {
		if pg.Status != "ok" || pg.URL == "" {
			stats.Skipped++
			continue
		}

		processed, perr := ProcessPage(pg, mdDir, htmlDir, imgIndex)
		if perr != nil {
			stats.Skipped++
			continue
		}

		// --- Insert page ---
		res, err := pageStmt.Exec(
			processed.Page.Section, processed.Page.Title, processed.Page.Slug,
			processed.Page.SourceURL, processed.Page.Route, processed.Page.WaybackTS,
			processed.Page.HTMLSHA1, processed.Page.BodyText, processed.Page.BodyMD,
			processed.Page.Thumb, processed.Page.MetaDescription, processed.Page.AuthorByline,
			processed.Page.Status, processed.Page.SortTitle, processed.Page.WordCount,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "import: error inserting %s: %v\n", processed.Page.Route, err)
			continue
		}
		pageID, _ := res.LastInsertId()
		stats.Pages++

		// --- Insert images ---
		for _, img := range processed.Images {
			isPrimary := 0
			if img.IsPrimary {
				isPrimary = 1
			}
			imgStmt.Exec(pageID, img.ImagePath, img.OriginalRef, img.Caption, isPrimary, img.SortOrder)
			stats.Images++
		}

		// --- Insert links ---
		for _, lnk := range processed.Links {
			linkStmt.Exec(pageID, lnk.TargetType, lnk.TargetURL, lnk.TargetRoute, lnk.TargetTitle, lnk.Blurb)
			stats.Links++
		}

		batchSize++
		if batchSize >= batchCommitEvery {
			pageStmt.Close()
			imgStmt.Close()
			linkStmt.Close()
			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("batch commit: %w", err)
			}
			tx, err = sqlDB.Begin()
			if err != nil {
				return nil, err
			}
			pageStmt, err = tx.Prepare(pageInsertSQL)
			if err != nil {
				return nil, err
			}
			imgStmt, err = tx.Prepare(imgInsertSQL)
			if err != nil {
				return nil, err
			}
			linkStmt, err = tx.Prepare(linkInsertSQL)
			if err != nil {
				return nil, err
			}
			batchSize = 0
			fmt.Fprintf(os.Stderr, "\rimported %d pages...", stats.Pages)
		}
	}

	pageStmt.Close()
	imgStmt.Close()
	linkStmt.Close()
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("final commit: %w", err)
	}

	// --- Rebuild FTS index (safety net in case triggers didn't fire during bulk import) ---
	database.Exec("INSERT INTO pages_fts(pages_fts) VALUES('rebuild')")

	stats.Duration = time.Since(start)
	fmt.Fprintf(os.Stderr, "\nimport complete: %d pages, %d images, %d links in %s (%d skipped)\n",
		stats.Pages, stats.Images, stats.Links, stats.Duration, stats.Skipped)
	return stats, nil
}
