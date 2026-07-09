package repository

import (
	"database/sql"
	"log"
	"fmt"
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
			id                   INTEGER PRIMARY KEY AUTOINCREMENT,
			name                 TEXT NOT NULL,
			owner                TEXT NOT NULL,
			full_name            TEXT NOT NULL,
			url                  TEXT,
			visibility           TEXT,
			private              BOOLEAN DEFAULT 0,
			is_fork              BOOLEAN DEFAULT 0,
			forked_from          TEXT,
			forked_from_owner    TEXT,
			forks_count          INTEGER DEFAULT 0,
			forked_to_count      INTEGER DEFAULT 0,
			stargazers_count     INTEGER DEFAULT 0,
			collaborators_count  INTEGER DEFAULT 0,
			who_has_access       TEXT,
			collaborators_list   TEXT,
			language             TEXT,
			size_kb              INTEGER DEFAULT 0,
			created_at           TEXT,
			updated_at           TEXT,
			pushed_at            TEXT,
			default_branch       TEXT,
			archived             BOOLEAN DEFAULT 0,
			disabled             BOOLEAN DEFAULT 0,
			license              TEXT,
			description          TEXT,
			last_checked         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(name, owner)
		);

		CREATE TABLE IF NOT EXISTS collaborators (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id    INTEGER NOT NULL,
			username   TEXT NOT NULL,
			has_access BOOLEAN NOT NULL DEFAULT 0,
			checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(repo_id) REFERENCES repos(id)
		);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	migrateRepoTable(db)
	ensureUniqueIndex(db)

	return db
}

// migrateRepoTable adds any columns present in the current schema but
// missing from an older on-disk repos table (e.g. db/collabs.db created
// by a previous version of this app).
func migrateRepoTable(db *sql.DB) {
	rows, err := db.Query("PRAGMA table_info(repos)")
	if err != nil {
		log.Printf("failed to inspect repos table: %v", err)
		return
	}

	existing := map[string]bool{}
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			continue
		}
		existing[name] = true
	}
	rows.Close()

	wanted := map[string]string{
		"url":                 "TEXT",
		"visibility":          "TEXT",
		"private":             "BOOLEAN DEFAULT 0",
		"is_fork":             "BOOLEAN DEFAULT 0",
		"forked_from":         "TEXT",
		"forked_from_owner":   "TEXT",
		"forks_count":         "INTEGER DEFAULT 0",
		"forked_to_count":     "INTEGER DEFAULT 0",
		"stargazers_count":    "INTEGER DEFAULT 0",
		"collaborators_count": "INTEGER DEFAULT 0",
		"language":            "TEXT",
		"size_kb":             "INTEGER DEFAULT 0",
		"created_at":          "TEXT",
		"updated_at":          "TEXT",
		"pushed_at":           "TEXT",
		"default_branch":      "TEXT",
		"archived":            "BOOLEAN DEFAULT 0",
		"disabled":            "BOOLEAN DEFAULT 0",
		"license":             "TEXT",
		"description":         "TEXT",
		"full_name":           "TEXT",
	}

	for col, ddl := range wanted {
		if existing[col] {
			continue
		}
		log.Printf("migrating repos table: adding %s column", col)
		if _, err := db.Exec(fmt.Sprintf("ALTER TABLE repos ADD COLUMN %s %s", col, ddl)); err != nil {
			log.Printf("migration failed for column %s: %v", col, err)
		}
	}
}

// ensureUniqueIndex creates a UNIQUE index on (name, owner) if it doesn't
// already exist. This is what ON CONFLICT(name, owner) in Upsert relies on.
// A plain ALTER TABLE cannot add a table-level UNIQUE constraint in SQLite,
// but a unique index works identically for ON CONFLICT resolution.
func ensureUniqueIndex(db *sql.DB) {
	// Clean up any duplicate (name, owner) rows first — a unique index
	// creation will fail if duplicates already exist in the table.
	_, err := db.Exec(`
		DELETE FROM repos
		WHERE id NOT IN (
			SELECT MIN(id) FROM repos GROUP BY name, owner
		)
	`)
	if err != nil {
		log.Printf("failed to dedupe repos before indexing: %v", err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_repos_name_owner ON repos(name, owner)`)
	if err != nil {
		log.Printf("failed to create unique index on repos(name, owner): %v", err)
	} else {
		log.Println("ensured unique index on repos(name, owner)")
	}
}