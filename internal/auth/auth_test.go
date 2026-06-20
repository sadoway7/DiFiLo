package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	pw := "correct horse battery staple"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == pw {
		t.Fatal("HashPassword returned the plaintext password")
	}
	// Correct password must verify.
	if !CheckPassword(pw, hash) {
		t.Error("CheckPassword failed for the correct password")
	}
	// Wrong password must not verify.
	if CheckPassword("definitely-not-the-password", hash) {
		t.Error("CheckPassword succeeded for a wrong password")
	}
}

func TestCreateAndParseSessionToken(t *testing.T) {
	const userID int64 = 42
	token := CreateSessionToken(userID)
	if token == "" {
		t.Fatal("CreateSessionToken returned empty token")
	}
	// Tokens carry a signature separator.
	if !strings.Contains(token, ".") {
		t.Fatalf("token %q has no '.' separator", token)
	}
	got := ParseSessionToken(token)
	if got != userID {
		t.Errorf("ParseSessionToken = %d, want %d", got, userID)
	}
}

func TestParseSessionTokenTamperDetection(t *testing.T) {
	token := CreateSessionToken(42)

	dot := strings.IndexByte(token, '.')
	if dot < 0 {
		t.Fatal("token has no '.' separator")
	}
	sig := token[dot+1:]
	if len(sig) == 0 {
		t.Fatal("token has empty signature")
	}
	// Flip the first signature character so the HMAC no longer matches.
	c := sig[0]
	flipped := "0"
	if c == '0' {
		flipped = "1"
	}
	tampered := token[:dot+1] + flipped + sig[1:]

	if got := ParseSessionToken(tampered); got != 0 {
		t.Errorf("ParseSessionToken on tampered token = %d, want 0", got)
	}
}

func TestParseSessionTokenMalformed(t *testing.T) {
	cases := []string{"", "no-separator", "abc.def.ghi", "!!!.xyz"}
	for _, tc := range cases {
		if got := ParseSessionToken(tc); got != 0 {
			t.Errorf("ParseSessionToken(%q) = %d, want 0", tc, got)
		}
	}
}

func TestSetSessionCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	SetSessionCookie(rec, 7)
	resp := rec.Result()

	var cookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == SessionCookieName {
			cookie = c
		}
	}
	if cookie == nil {
		t.Fatalf("session cookie %q was not set", SessionCookieName)
	}
	if cookie.Value == "" {
		t.Error("cookie value is empty")
	}
	if cookie.MaxAge != SessionMaxAge {
		t.Errorf("cookie MaxAge = %d, want %d", cookie.MaxAge, SessionMaxAge)
	}
	if !cookie.HttpOnly {
		t.Error("cookie is not HttpOnly")
	}
	if cookie.Path != "/" {
		t.Errorf("cookie Path = %q, want %q", cookie.Path, "/")
	}
	// The cookie value must round-trip back to the userID.
	if got := ParseSessionToken(cookie.Value); got != 7 {
		t.Errorf("cookie value does not parse back: got %d, want 7", got)
	}
}

func TestClearSessionCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	ClearSessionCookie(rec)
	resp := rec.Result()

	var cookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == SessionCookieName {
			cookie = c
		}
	}
	if cookie == nil {
		t.Fatalf("session cookie %q was not set", SessionCookieName)
	}
	if cookie.Value != "" {
		t.Errorf("cookie Value = %q, want empty", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("cookie MaxAge = %d, want -1", cookie.MaxAge)
	}
}

// fakeUser and fakeComment stand in for db.User and db.Comment so the
// permission logic can be exercised without importing the db package.
type fakeUser struct {
	role string
	id   int64
}

func (u *fakeUser) GetRole() string { return u.role }
func (u *fakeUser) GetID() int64    { return u.id }

type fakeComment struct {
	userID int64
}

func (c *fakeComment) GetUserID() int64 { return c.userID }

func TestCanDeleteComment(t *testing.T) {
	cases := []struct {
		name    string
		user    UserRole
		comment CommentOwner
		want    bool
	}{
		{"admin deletes any", &fakeUser{RoleAdmin, 1}, &fakeComment{999}, true},
		{"manager deletes any", &fakeUser{RoleManager, 2}, &fakeComment{999}, true},
		{"general deletes own", &fakeUser{RoleGeneral, 5}, &fakeComment{5}, true},
		{"general deletes other", &fakeUser{RoleGeneral, 5}, &fakeComment{6}, false},
		{"unknown role", &fakeUser{"superuser", 1}, &fakeComment{1}, false},
		{"nil user", nil, &fakeComment{5}, false},
		{"nil comment", &fakeUser{RoleAdmin, 1}, nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CanDeleteComment(tc.user, tc.comment); got != tc.want {
				t.Errorf("CanDeleteComment = %v, want %v", got, tc.want)
			}
		})
	}
}
