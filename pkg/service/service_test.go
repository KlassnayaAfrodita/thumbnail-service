package service

import (
	"context"
	"testing"
	"thumbnail-service/proto"
)

// mockDB реализует интерфейс db.ThumbnailStorage
type mockDB struct {
	data map[string][]byte
}

func (m *mockDB) GetThumbnail(videoURL string) ([]byte, error) {
	return m.data[videoURL], nil
}

func (m *mockDB) SaveThumbnail(videoURL string, imageData []byte) error {
	m.data[videoURL] = imageData
	return nil
}

func TestThumbnailService(t *testing.T) {
	// Создаем мокированную базу данных
	mockDatabase := &mockDB{data: make(map[string][]byte)}
	// Создаем сервис с мокированным хранилищем
	service := NewThumbnailService(mockDatabase)

	// Ожидаемые данные
	expectedData := []byte("test_image_data")

	// Сохраняем данные в мокированное хранилище
	err := mockDatabase.SaveThumbnail("https://example.com/video", expectedData)
	if err != nil {
		t.Fatalf("Failed to save thumbnail: %v", err)
	}

	// Тестируем запрос для видео
	req := &proto.ThumbnailRequest{VideoUrl: "https://example.com/video"}
	resp, err := service.GetThumbnail(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to process request: %v", err)
	}

	// Проверяем, что данные не пустые
	if resp.ImageData == nil {
		t.Errorf("Expected non-nil image data")
	}

	// Проверяем, что данные совпадают с ожидаемыми
	if string(resp.ImageData) != string(expectedData) {
		t.Errorf("Expected %v, got %v", string(expectedData), string(resp.ImageData))
	}
}
