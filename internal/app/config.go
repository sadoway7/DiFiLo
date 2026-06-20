package app

// Config holds tunable application parameters. Defaults come from
// DefaultConfig; the values can be overridden at construction time.
type Config struct {
	SearchLimit   int   // max results returned by a single search query
	CommentMaxLen int   // maximum characters allowed in a comment body
	GalleryMax    int   // maximum sidebar gallery images per wiki page
	LinkMax       int   // maximum sidebar related-link cards per wiki page
	HeroMinSize   int64 // minimum byte size for a hero-background candidate image
	SessionMaxAge int   // session cookie lifetime, in seconds
}

// DefaultConfig returns the standard configuration with sane defaults that
// match the behaviour of the original flat codebase.
func DefaultConfig() *Config {
	return &Config{
		SearchLimit:   100,
		CommentMaxLen: 5000,
		GalleryMax:    8,
		LinkMax:       25,
		HeroMinSize:   80 * 1024,
		SessionMaxAge: 30 * 24 * 3600,
	}
}
