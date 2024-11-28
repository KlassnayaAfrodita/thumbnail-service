package db

import (
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	// Создаем временный файл для тестовой базы
	tempDBFile := "test.db"
	defer os.Remove(tempDBFile)

	// Инициализируем базу
	database, err := NewDB(tempDBFile)
	if err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}

	// Проверяем сохранение данных
	videoURL := "https://example.com/video"
	imageData := []byte("test_image_data")
	err = database.SaveThumbnail(videoURL, imageData)
	if err != nil {
		t.Fatalf("Failed to save thumbnail: %v", err)
	}

	// Проверяем получение данных
	retrievedData, err := database.GetThumbnail(videoURL)
	if err != nil {
		t.Fatalf("Failed to retrieve thumbnail: %v", err)
	}
	if string(retrievedData) != string(imageData) {
		t.Errorf("Expected %v, got %v", string(imageData), string(retrievedData))
	}

	// Проверяем случай, когда данных нет
	_, err = database.GetThumbnail("https://example.com/unknown")
	if err != nil {
		t.Fatalf("Error while retrieving non-existent thumbnail: %v", err)
	}
}
