package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Session cookie configuration.
const (
	// SessionCookieName is the name of the browser cookie that carries
	// the signed session token.
	SessionCookieName = "difilo_session"
	// SessionMaxAge is the lifetime of a session, in seconds (30 days).
	SessionMaxAge = 30 * 24 * 3600
)

// sessionSecret signs session tokens. It defaults to a random value
// generated at init time, so sessions do not survive restarts. Call
// InitSecret with a persisted value to keep sessions across restarts.
var sessionSecret []byte

// InitSecret sets the secret used to sign session tokens. If secret is
// nil or empty, a fresh random 32-byte secret is generated instead.
func InitSecret(secret []byte) {
	if len(secret) == 0 {
		// A local offline app: a random per-restart secret is fine.
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		secret = b
	}
	sessionSecret = secret
}

// init guarantees a secret exists even if InitSecret is never called,
// preserving the original "random per restart" behavior.
func init() {
	InitSecret(nil)
}

// CreateSessionToken builds a signed token encoding the userID and an
// expiry timestamp, of the form: base64(userID.expires).hmac
func CreateSessionToken(userID int64) string {
	expires := time.Now().Add(time.Duration(SessionMaxAge) * time.Second).Unix()
	payload := strconv.FormatInt(userID, 10) + "." + strconv.FormatInt(expires, 10)
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

// ParseSessionToken validates the signature and expiry of the token and
// returns the encoded userID. It returns 0 for any malformed, tampered,
// or expired token.
func ParseSessionToken(token string) int64 {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return 0
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0
	}
	payload := string(payloadBytes)
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSig)) {
		return 0
	}
	dot := strings.LastIndex(payload, ".")
	if dot < 0 {
		return 0
	}
	userID, _ := strconv.ParseInt(payload[:dot], 10, 64)
	expires, _ := strconv.ParseInt(payload[dot+1:], 10, 64)
	if time.Now().Unix() > expires {
		return 0
	}
	return userID
}

// SetSessionCookie creates a signed session cookie for the given userID
// and attaches it to the response.
func SetSessionCookie(w http.ResponseWriter, userID int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    CreateSessionToken(userID),
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie expires the session cookie on the response, ending
// the session on the client.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
