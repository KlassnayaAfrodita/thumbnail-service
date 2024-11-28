package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg" // Для поддержки JPEG
	_ "image/png"  // Для поддержки PNG
	"log"
	"os"
	"sync"
	"thumbnail-service/proto"

	"google.golang.org/grpc"
)

func main() {
	async := flag.Bool("async", false, "Enable async download")
	flag.Parse()

	videoURLs := flag.Args()
	if len(videoURLs) == 0 {
		log.Fatalf("No video URLs provided")
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewThumbnailServiceClient(conn)

	if *async {
		var wg sync.WaitGroup
		for _, url := range videoURLs {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				downloadThumbnail(client, url)
			}(url)
		}
		wg.Wait()
	} else {
		for _, url := range videoURLs {
			downloadThumbnail(client, url)
		}
	}
}

func downloadThumbnail(client proto.ThumbnailServiceClient, videoURL string) {
	req := &proto.ThumbnailRequest{VideoUrl: videoURL}
	resp, err := client.GetThumbnail(context.Background(), req)
	if err != nil {
		log.Printf("Failed to download thumbnail for %s: %v", videoURL, err)
		return
	}

	// Логируем содержимое ответа для диагностики
	log.Printf("Received response: %v", resp)

	// Получаем базовое имя файла (с использованием base64 для уникальности)
	fileName := fmt.Sprintf("%s", base64.URLEncoding.EncodeToString([]byte(videoURL)))

	// Сохраняем изображение с правильным расширением
	err = saveImageToFile(fileName, resp.ImageData)
	if err != nil {
		log.Printf("Failed to save thumbnail for %s: %v", videoURL, err)
		return
	}

	log.Printf("Thumbnail for %s saved as %s", videoURL, fileName)
}

// Функция для сохранения изображения с правильным расширением
func saveImageToFile(fileName string, data []byte) error {
	// Проверка на пустые данные
	if len(data) == 0 {
		return fmt.Errorf("empty image data")
	}

	// Логируем длину данных изображения
	log.Printf("Received image data of size: %d bytes", len(data))

	// Печатаем первые несколько байтов для диагностики
	if len(data) < 64 {
		log.Printf("Image data is smaller than 64 bytes: %v", data)
	} else {
		log.Printf("Received image data (first 64 bytes): %v", data[:64])
	}

	// Пробуем декодировать изображение для определения его формата
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Логируем информацию о формате
	log.Printf("Image format detected: %s", format)

	// Устанавливаем расширение файла в зависимости от формата
	var extension string
	switch format {
	case "jpeg":
		extension = ".jpg"
	case "png":
		extension = ".png"
	default:
		extension = ".jpg" // По умолчанию сохраняем как JPG
	}

	// Формируем окончательное имя файла с расширением
	fileNameWithExtension := fileName + extension

	// Сохраняем изображение в файл
	err = os.WriteFile(fileNameWithExtension, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
