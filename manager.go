package gorms

import (
	"path/filepath"
)

type Manager map[string]*DB

// Get returns the DB represented by the given file path.
func (m Manager) Get(file string) *DB {
	if db, ok := m[file]; ok {
		return db
	}
	if db, ok := m[filepath.Clean(file)]; ok {
		return db
	}
	if db, ok := m[filepath.Dir(file)]; ok {
		return db
	}
	panic("current file does not have a database")
}

// Set sets the DB in different paths.
func (m Manager) Set(db *DB, pkg ...string) {
	// set private DB
	m[db.file] = db
	m[filepath.Clean(db.file)] = db
	// set public DB
	for _, file := range pkg {
		if _, ok := GlobalDB[file]; !ok {
			GlobalDB[file] = db
		}
	}
}
