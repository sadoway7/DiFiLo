// Package auth handles authentication and authorization: password
// hashing (bcrypt), session token creation and verification, cookie
// management, role constants, and permission helpers.
//
// This is a leaf package with zero dependencies on any other DiFiLo
// package. The session secret is generated at init time; for
// persistent sessions across restarts, call InitSecret with a
// persisted value.
package auth
