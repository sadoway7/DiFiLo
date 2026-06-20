// Package content handles content processing: the import pipeline
// (mirror files → SQLite), markdown rendering, wiki page layout
// assembly, and the image lightbox.
//
// This package depends on db (for writes/reads) and textutil (for
// prose helpers used during rendering).
package content
