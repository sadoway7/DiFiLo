package app

import (
	"net/http"
	"regexp"
	"strings"

	"difilo/internal/auth"
)

// reEmail validates a basic email address shape for registration.
var reEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// handleRegister: GET = form, POST = create account. The first account ever
// created becomes the admin; all subsequent accounts are general users.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.renderShell(w, r, "Register", "", registerHTML())
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if !reEmail.MatchString(email) {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Please enter a valid email address.</p>`)
		return
	}
	if len(username) < 2 {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Username must be at least 2 characters.</p>`)
		return
	}
	if len(password) < 6 {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Password must be at least 6 characters.</p>`)
		return
	}
	if password != password2 {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Passwords do not match.</p>`)
		return
	}
	existing, _ := s.DB.GetUserByEmail(email)
	if existing != nil {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>An account with that email already exists.</p>`)
		return
	}
	existingU, _ := s.DB.GetUserByUsername(username)
	if existingU != nil {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>That username is already taken.</p>`)
		return
	}

	// First user becomes admin; others are general.
	count, _ := s.DB.UserCount()
	role := auth.RoleGeneral
	if count == 0 {
		role = auth.RoleAdmin
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Server error creating account.</p>`)
		return
	}
	u, err := s.DB.CreateUser(email, username, hash, role)
	if err != nil {
		s.renderShell(w, r, "Register", "", registerHTML()+`<p class='df-err'>Server error creating account.</p>`)
		return
	}
	auth.SetSessionCookie(w, u.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleLogin: GET = form, POST = authenticate.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.renderShell(w, r, "Login", "", loginHTML())
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	u, err := s.DB.GetUserByEmail(email)
	if err != nil || u == nil || !auth.CheckPassword(password, u.PasswordHash) {
		s.renderShell(w, r, "Login", "", loginHTML()+`<p class='df-err'>Invalid email or password.</p>`)
		return
	}
	auth.SetSessionCookie(w, u.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleLogout clears the session cookie and returns to the homepage.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- HTML forms ---

func loginHTML() string {
	return `<div class="df-auth">
<h1>Log in</h1>
<form method="post" action="/login" class="df-auth-form">
<label>Email<input type="email" name="email" required autofocus></label>
<label>Password<input type="password" name="password" required></label>
<button type="submit">Log in</button>
</form>
<p class="df-muted">No account? <a href="/register">Register</a></p>
</div>`
}

func registerHTML() string {
	return `<div class="df-auth">
<h1>Register</h1>
<form method="post" action="/register" class="df-auth-form">
<label>Email<input type="email" name="email" required autofocus></label>
<label>Username (public display name)<input type="text" name="username" required minlength="2" maxlength="30"></label>
<label>Password<input type="password" name="password" required></label>
<label>Confirm password<input type="password" name="password2" required></label>
<button type="submit">Create account</button>
</form>
<p class="df-muted">Already have an account? <a href="/login">Log in</a></p>
</div>`
}
