package textutil

import (
	"strings"
	"testing"
)

func TestCleanProseStripsHTMLAndMarkdown(t *testing.T) {
	in := "# Heading\n\n" +
		"<p>Some <b>bold</b> text &amp; more</p>\n\n" +
		"---\n\n" +
		"```code block```" +
		"Plain paragraph.\n"

	got := CleanProse(in)

	// No leftover HTML tags or markdown noise.
	if strings.ContainsAny(got, "<>") {
		t.Errorf("expected no HTML tags, got %q", got)
	}
	if strings.Contains(got, "```") {
		t.Errorf("expected code fences removed, got %q", got)
	}
	if strings.Contains(got, "```") || strings.Contains(got, "code block") {
		t.Errorf("expected code-block content removed, got %q", got)
	}
	// HTML entity should be decoded.
	if !strings.Contains(got, "bold text & more") {
		t.Errorf("expected decoded entities and inner text, got %q", got)
	}
	// Heading marker removed but word retained.
	if !strings.Contains(got, "Heading") {
		t.Errorf("expected heading text retained, got %q", got)
	}
	if strings.Contains(got, "# Heading") {
		t.Errorf("expected heading marker stripped, got %q", got)
	}
	// Final paragraph survives.
	if !strings.Contains(got, "Plain paragraph.") {
		t.Errorf("expected plain paragraph retained, got %q", got)
	}
}

func TestPrettySection(t *testing.T) {
	if got := PrettySection("material"); got != "Material" {
		t.Errorf("PrettySection(material) = %q, want %q", got, "Material")
	}
	if got := PrettySection(""); got != "pages" {
		t.Errorf("PrettySection(\"\") = %q, want %q", got, "pages")
	}
}

func TestOrDefault(t *testing.T) {
	if got := OrDefault("", "fallback"); got != "fallback" {
		t.Errorf("OrDefault(\"\", fallback) = %q, want %q", got, "fallback")
	}
	if got := OrDefault("kept", "fallback"); got != "kept" {
		t.Errorf("OrDefault(kept, fallback) = %q, want %q", got, "kept")
	}
}

func TestAZKey(t *testing.T) {
	if got := AZKey("Ball Clay"); got != 'b' {
		t.Errorf("AZKey(\"Ball Clay\") = %q, want %q", got, 'b')
	}
	if got := AZKey("123numeric"); got != '#' {
		t.Errorf("AZKey(\"123numeric\") = %q, want %q", got, '#')
	}
	if got := AZKey("   "); got != '#' {
		t.Errorf("AZKey(\"   \") = %q, want %q", got, '#')
	}
	if got := AZKey("Zinc"); got != 'z' {
		t.Errorf("AZKey(\"Zinc\") = %q, want %q", got, 'z')
	}
}

func TestSectionFromRoute(t *testing.T) {
	if got := SectionFromRoute("/material/925"); got != "material" {
		t.Errorf("SectionFromRoute(\"/material/925\") = %q, want %q", got, "material")
	}
	if got := SectionFromRoute("/material"); got != "material" {
		t.Errorf("SectionFromRoute(\"/material\") = %q, want %q", got, "material")
	}
	if got := SectionFromRoute(""); got != "" {
		t.Errorf("SectionFromRoute(\"\") = %q, want %q", got, "")
	}
}
