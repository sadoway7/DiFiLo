// Package app is the application wiring layer: it holds the Server
// struct, HTTP routing, all request handlers, middleware, logging,
// configuration, and the CurrentUser helper. It imports and composes
// all other internal packages (db, auth, content, ui, textutil).
package app
