package service

import (
	"context"
	"fmt"
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
	// Получаем миниатюру из базы данных
	imageData, err := s.db.GetThumbnail(req.VideoUrl)
	if err != nil {
		log.Printf("Failed to get thumbnail for %s: %v", req.VideoUrl, err)
		return nil, fmt.Errorf("failed to get thumbnail: %w", err)
	}

	if imageData == nil {
		// Если нет в кэше, скачиваем и сохраняем
		// Здесь можно добавить логику скачивания миниатюры
		log.Printf("Thumbnail not found for %s, downloading...\n", req.VideoUrl)
		imageData = []byte("dummy_image_data") // Это просто пример, замените на реальную логику

		err = s.db.SaveThumbnail(req.VideoUrl, imageData)
		if err != nil {
			log.Printf("Failed to save thumbnail for %s: %v", req.VideoUrl, err)
			return nil, fmt.Errorf("failed to save thumbnail: %w", err)
		}
	}

	// Возвращаем успешный ответ
	return &proto.ThumbnailResponse{ImageData: imageData}, nil
}
