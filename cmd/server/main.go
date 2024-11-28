package main

import (
	"log"
	"net"
	"thumbnail-service/pkg/db"
	"thumbnail-service/pkg/service"
	pb "thumbnail-service/proto"

	"google.golang.org/grpc"
)

func main() {
	// Инициализация базы данных
	db, err := db.NewDB("thumbnails.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Создаем сервис
	thumbnailService := service.NewThumbnailService(db)

	// Настраиваем gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterThumbnailServiceServer(grpcServer, thumbnailService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
