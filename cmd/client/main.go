package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
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

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
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

	fileName := fmt.Sprintf("%s.jpg", base64.URLEncoding.EncodeToString([]byte(videoURL)))
	err = os.WriteFile(fileName, resp.ImageData, 0644)
	if err != nil {
		log.Printf("Failed to save thumbnail for %s: %v", videoURL, err)
		return
	}

	log.Printf("Thumbnail for %s saved as %s", videoURL, fileName)
}
