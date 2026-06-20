// Package ui generates HTML, CSS, and JavaScript fragments for the
// DiFiLo web interface: the navigation panel, page shell, comment
// section, bookmark/pin toolbar, and all CSS styles.
//
// This is a leaf package with zero dependencies on any other DiFiLo
// package. Functions that need user context accept a *Viewer DTO
// rather than a *db.User, keeping the package fully decoupled.
package ui
