package service

import (
	"context"
	"log"
	"thumbnail-service/pkg/db"
	"thumbnail-service/proto"
)

// ThumbnailService реализует интерфейс ThumbnailServiceServer
type ThumbnailService struct {
	proto.UnimplementedThumbnailServiceServer // Встроенный сервис по умолчанию
	db                                        db.ThumbnailStorage
}

// Новый конструктор для создания сервиса
func NewThumbnailService(db db.ThumbnailStorage) *ThumbnailService {
	return &ThumbnailService{db: db}
}

// Реализация метода GetThumbnail
func (s *ThumbnailService) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	videoURL := req.GetVideoUrl()
	log.Printf("Received request for video URL: %s\n", videoURL)

	// Проверяем кэш
	imageData, err := s.db.GetThumbnail(videoURL)
	if err != nil {
		log.Printf("Error retrieving thumbnail from cache: %v", err)
		return nil, err
	}

	// Если данные найдены в кэше
	if imageData != nil {
		log.Printf("Thumbnail retrieved from cache, size: %d bytes\n", len(imageData))
		return &proto.ThumbnailResponse{ImageData: imageData}, nil
	}

	// Загружаем миниатюру с YouTube
	imageData, err = DownloadThumbnail(videoURL)
	if err != nil {
		log.Printf("Error downloading thumbnail: %v", err)
		return nil, err
	}

	// Сохраняем миниатюру в кэш
	err = s.db.SaveThumbnail(videoURL, imageData)
	if err != nil {
		log.Printf("Error saving thumbnail to cache: %v", err)
	}

	log.Printf("Thumbnail downloaded and cached, size: %d bytes\n", len(imageData))
	return &proto.ThumbnailResponse{ImageData: imageData}, nil
}
