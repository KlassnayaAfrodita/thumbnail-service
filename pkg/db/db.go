package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS thumbnails (
		id TEXT PRIMARY KEY,
		video_id TEXT UNIQUE NOT NULL,
		file_path TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("Failed to create table: %v", err)
	}

	return db, nil
}
