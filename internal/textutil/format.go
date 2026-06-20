package textutil

import "strings"

// PrettySection capitalizes the first letter of a section name. An empty
// section name is rendered as "pages".
func PrettySection(s string) string {
	if s == "" {
		return "pages"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// OrDefault returns v when it is non-empty, otherwise def.
func OrDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// SectionFromRoute extracts the section name from a route like "/material/925".
// For "/material/925" it returns "material"; for a bare "/material" it returns
// "material"; for "" it returns "".
func SectionFromRoute(route string) string {
	route = strings.TrimPrefix(route, "/")
	if i := strings.IndexByte(route, '/'); i >= 0 {
		return route[:i]
	}
	return route
}

// AZKey returns the alphabet bucket for a title: the lower-cased first letter
// for a-z, or '#' for anything else (digits, symbols, empty).
func AZKey(title string) byte {
	t := strings.ToLower(strings.TrimSpace(title))
	if len(t) == 0 {
		return '#'
	}
	c := t[0]
	if c >= 'a' && c <= 'z' {
		return c
	}
	return '#'
}
