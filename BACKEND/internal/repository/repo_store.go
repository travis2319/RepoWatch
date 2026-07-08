package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB opens a SQLite connection and runs schema.sql
func InitDB(path string) *sql.DB {
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS repos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		full_name TEXT,
		url TEXT,
		last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS collaborators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_id INTEGER,
		username TEXT NOT NULL,
		has_access BOOLEAN,
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(repo_id) REFERENCES repos(id)
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	migrateRepoTable(db)

	return db
}

// migrateRepoTable adds columns to an existing repos table created by an
// older version of this schema (e.g. your existing db/collabs.db that only
// has url, not full_name).
func migrateRepoTable(db *sql.DB) {
	rows, err := db.Query("PRAGMA table_info(repos)")
	if err != nil {
		log.Printf("failed to inspect repos table: %v", err)
		return
	}
	defer rows.Close()

	hasFullName := false
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			continue
		}
		if name == "full_name" {
			hasFullName = true
		}
	}

	if !hasFullName {
		log.Println("migrating repos table: adding full_name column")
		if _, err := db.Exec("ALTER TABLE repos ADD COLUMN full_name TEXT"); err != nil {
			log.Printf("migration failed: %v", err)
		}
	}
}
