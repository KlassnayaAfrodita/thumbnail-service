package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var youtubeThumbnailURLRegex = regexp.MustCompile(`https?://(?:www\.)?youtube\.com/watch\?v=([^&]+)`)

func DownloadThumbnail(videoURL string) ([]byte, error) {
	matches := youtubeThumbnailURLRegex.FindStringSubmatch(videoURL)
	if len(matches) < 2 {
		return nil, errors.New("invalid YouTube URL")
	}

	videoID := matches[1]
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download thumbnail")
	}

	return io.ReadAll(resp.Body)
}
