package service

import (
	"context"
	"thumbnail-service/pkg/db"
	pb "thumbnail-service/proto"
)

type ThumbnailService struct {
	pb.UnimplementedThumbnailServiceServer
	db *db.DB
}

func NewThumbnailService(db *db.DB) *ThumbnailService {
	return &ThumbnailService{db: db}
}

func (s *ThumbnailService) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	// Проверка кеша
	imageData, err := s.db.GetThumbnail(req.VideoUrl)
	if err != nil {
		return nil, err
	}

	// Если в кэше нет, скачиваем
	if imageData == nil {
		imageData, err = DownloadThumbnail(req.VideoUrl)
		if err != nil {
			return nil, err
		}

		// Сохраняем в кэш
		if err := s.db.SaveThumbnail(req.VideoUrl, imageData); err != nil {
			return nil, err
		}
	}

	return &pb.ThumbnailResponse{ImageData: imageData}, nil
}
