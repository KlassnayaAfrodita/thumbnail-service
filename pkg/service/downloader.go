package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadThumbnail(videoID string) (string, error) {
	url := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch thumbnail: %v", err)
	}
	defer resp.Body.Close()

	filePath := fmt.Sprintf("%s.jpg", videoID)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("Failed to save thumbnail: %v", err)
	}

	return filePath, nil
}
