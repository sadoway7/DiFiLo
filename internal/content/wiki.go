package content

import (
	"fmt"
	"html"
	"strings"

	"difilo/internal/db"
	"difilo/internal/textutil"
)

// RenderWikiPage produces the full HTML body for a wiki-style page rendered
// from the database content: a header (badge, title, description), the
// rendered markdown body, a sidebar (page info, image gallery, related links),
// and the lightbox modal. It is a standalone function — it holds no receiver
// state and depends only on its arguments.
func RenderWikiPage(p *db.ContentPage, images []db.ImageRow, links []db.LinkRow) string {
	var b strings.Builder

	b.WriteString(`<div class="df-wiki">`)

	// ---- main column ----
	b.WriteString(`<div class="df-wiki-main">`)

	// header
	b.WriteString(`<div class="df-wiki-header">`)
	b.WriteString(fmt.Sprintf(`<span class="df-wiki-badge">%s</span>`,
		html.EscapeString(textutil.PrettySection(p.Section))))
	b.WriteString(fmt.Sprintf(`<h1>%s</h1>`,
		html.EscapeString(p.Title)))
	if p.MetaDescription != "" {
		b.WriteString(fmt.Sprintf(`<p class="df-wiki-desc">%s</p>`,
			html.EscapeString(p.MetaDescription)))
	}
	b.WriteString(`</div>`)

	// content
	contentHTML := RenderMarkdown(StripLeadingH1(p.BodyMD))
	b.WriteString(fmt.Sprintf(`<div class="df-wiki-content">%s</div>`, contentHTML))

	b.WriteString(`</div>`)

	// ---- sidebar ----
	b.WriteString(`<aside class="df-wiki-sidebar">`)

	// metadata
	b.WriteString(`<div class="df-wiki-sidebox">`)
	b.WriteString(`<h3>Page Info</h3>`)
	b.WriteString(`<dl>`)
	writeMetaRow(&b, "Section", textutil.PrettySection(p.Section))
	writeMetaRow(&b, "Maintainer", textutil.OrDefault(p.AuthorByline, "Tony Hansen"))
	writeMetaRow(&b, "Words", fmt.Sprintf("%d", p.WordCount))
	if p.SourceURL != "" {
		writeMetaRow(&b, "Source", `<a href="`+html.EscapeString(p.SourceURL)+`" target="_blank" rel="noopener">Original ↗</a>`)
	}
	b.WriteString(`</dl></div>`)

	// image gallery
	galleryImgs := filterGalleryImages(images)
	if len(galleryImgs) > 0 {
		b.WriteString(`<div class="df-wiki-sidebox">`)
		b.WriteString(fmt.Sprintf(`<h3>Images (%d)</h3>`, len(galleryImgs)))
		b.WriteString(`<div class="df-wiki-thumbs">`)
		for _, img := range galleryImgs {
			if img.ImagePath == "" {
				continue
			}
			b.WriteString(fmt.Sprintf(`<img class="df-wiki-gal" loading="lazy" src="%s" alt="%s" data-caption="%s" onclick="dfLightbox(this)">`,
				html.EscapeString(img.ImagePath),
				html.EscapeString(img.Caption),
				html.EscapeString(img.Caption),
			))
		}
		b.WriteString(`</div></div>`)
	}

	// related links
	if len(links) > 0 {
		shown := 0
		var linkHTML strings.Builder
		for _, lnk := range links {
			if shown >= 25 {
				break
			}
			route := lnk.TargetRoute
			if route == "" {
				continue
			}
			title := lnk.TargetTitle
			if title == "" {
				title = route
			}
			sec := textutil.SectionFromRoute(route)
			linkHTML.WriteString(fmt.Sprintf(`<a href="%s"><span>%s</span><span class="df-wiki-link-sec">%s</span></a>`,
				html.EscapeString(route),
				html.EscapeString(title),
				html.EscapeString(sec),
			))
			shown++
		}
		if shown > 0 {
			b.WriteString(`<div class="df-wiki-sidebox df-wiki-links">`)
			b.WriteString(fmt.Sprintf(`<h3>Related Pages (%d)</h3>`, shown))
			b.WriteString(linkHTML.String())
			b.WriteString(`</div>`)
		}
	}

	b.WriteString(`</aside>`)
	b.WriteString(`</div>`)

	// lightbox modal + JS
	b.WriteString(LightboxHTML)

	return b.String()
}

// writeMetaRow emits one sidebar metadata row (label/value).
func writeMetaRow(b *strings.Builder, label, value string) {
	b.WriteString(fmt.Sprintf(`<div class="df-wiki-meta-row"><dt>%s</dt><dd>%s</dd></div>`,
		label, value))
}

// filterGalleryImages returns images suitable for the sidebar gallery: those
// with a non-empty path, capped at eight.
func filterGalleryImages(images []db.ImageRow) []db.ImageRow {
	var out []db.ImageRow
	for _, img := range images {
		if img.ImagePath == "" {
			continue
		}
		out = append(out, img)
	}
	if len(out) > 8 {
		out = out[:8]
	}
	return out
}
