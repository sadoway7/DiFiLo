package ui

// Viewer is a lightweight DTO carrying the user context that UI
// functions need, without coupling ui to the db package.
type Viewer struct {
	LoggedIn bool
	ID       int64
	Username string
	Email    string
	Role     string
}
