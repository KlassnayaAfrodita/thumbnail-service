package db

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Интерфейс для работы с хранилищем миниатюр
type ThumbnailStorage interface {
	GetThumbnail(videoURL string) ([]byte, error)
	SaveThumbnail(videoURL string, imageData []byte) error
}

type DB struct {
	conn *sql.DB
}

func NewDB(filePath string) (*DB, error) {
	conn, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	// Создание таблицы, если ее нет
	query := `
	CREATE TABLE IF NOT EXISTS thumbnails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		video_url TEXT UNIQUE,
		image_data BLOB,
		created_at TIMESTAMP
	);
	`
	_, err = conn.Exec(query)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) GetThumbnail(videoURL string) ([]byte, error) {
	query := `SELECT image_data FROM thumbnails WHERE video_url = ?`
	row := db.conn.QueryRow(query, videoURL)

	var imageData []byte
	err := row.Scan(&imageData)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Логирование размера извлеченных данных
	fmt.Printf("Retrieved thumbnail size: %d bytes\n", len(imageData))
	return imageData, nil
}

func (db *DB) SaveThumbnail(videoURL string, imageData []byte) error {
	query := `
	INSERT INTO thumbnails (video_url, image_data, created_at)
	VALUES (?, ?, ?)
	ON CONFLICT(video_url) DO UPDATE SET image_data = excluded.image_data, created_at = excluded.created_at;
	`
	_, err := db.conn.Exec(query, videoURL, imageData, time.Now())
	if err != nil {
		return err
	}

	// Логирование хэша для сохраненных данных
	hash := calculateHash(imageData)
	fmt.Printf("Saved thumbnail hash: %s\n", hash)
	return nil
}

// calculateHash вычисляет MD5-хэш данных
func calculateHash(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
