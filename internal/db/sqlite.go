package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Conn *sql.DB
}

func NewSQLiteDB(filepath string) (*Database, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{Conn: db}
	if err := database.initSchema(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		hostname TEXT,
		platform TEXT,
		architecture TEXT,
		integrity_hash TEXT,
		last_seen TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS jobs (
		job_id TEXT PRIMARY KEY,
		agent_id TEXT,
		type TEXT,
		status TEXT,
		created_at TIMESTAMP
	);
	`
	_, err := d.Conn.Exec(query)
	return err
}

func (d *Database) RegisterAgent(id, hostname, platform, arch, hash string) error {
	query := `
	INSERT INTO agents (id, hostname, platform, architecture, integrity_hash, last_seen)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		last_seen = excluded.last_seen,
		hostname = excluded.hostname;
	`
	_, err := d.Conn.Exec(query, id, hostname, platform, arch, hash, time.Now())
	return err
}

func (d *Database) Close() error {
	return d.Conn.Close()
}
