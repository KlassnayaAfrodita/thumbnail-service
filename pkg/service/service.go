package service

import (
	"context"
	"database/sql"
	"fmt"

	"thumbnail-service/pkg/downloader"
	pb "thumbnail-service/proto"
)

type server struct {
	db *sql.DB
	pb.UnimplementedThumbnailServiceServer
}

func NewServer(db *sql.DB) *server {
	return &server{db: db}
}

func (s *server) GetThumbnail(ctx context.Context, req *pb.GetThumbnailRequest) (*pb.GetThumbnailResponse, error) {
	videoID := req.GetVideoId()

	// Проверка кэша
	row := s.db.QueryRow("SELECT file_path FROM thumbnails WHERE video_id = ?", videoID)
	var filePath string
	if err := row.Scan(&filePath); err == nil {
		return &pb.GetThumbnailResponse{FilePath: filePath}, nil
	}

	// Загрузка превью
	filePath, err := downloader.DownloadThumbnail(videoID)
	if err != nil {
		return nil, fmt.Errorf("Failed to download thumbnail: %v", err)
	}

	// Кэширование в базе данных
	_, err = s.db.Exec("INSERT INTO thumbnails (id, video_id, file_path) VALUES (?, ?, ?)", videoID, videoID, filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to cache thumbnail: %v", err)
	}

	return &pb.GetThumbnailResponse{FilePath: filePath}, nil
}
