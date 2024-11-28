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
	dbConn, err := db.InitDB("thumbnails.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	// Настройка gRPC-сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterThumbnailServiceServer(grpcServer, service.NewServer(dbConn))

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
