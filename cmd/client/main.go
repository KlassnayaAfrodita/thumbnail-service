package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "thumbnail-service/proto"

	"google.golang.org/grpc"
)

func fetchThumbnail(client pb.ThumbnailServiceClient, videoID string, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.GetThumbnail(ctx, &pb.GetThumbnailRequest{VideoId: videoID})
	if err != nil {
		log.Printf("Failed to fetch thumbnail for video %s: %v", videoID, err)
		return
	}

	fmt.Printf("Thumbnail for video %s saved at: %s\n", videoID, resp.FilePath)
}

func main() {
	asyncFlag := flag.Bool("async", false, "Enable asynchronous fetching")
	flag.Parse()

	videoIDs := flag.Args()
	if len(videoIDs) == 0 {
		fmt.Println("Usage: client --async <video_id1> <video_id2> ...")
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewThumbnailServiceClient(conn)

	if *asyncFlag {
		var wg sync.WaitGroup
		for _, videoID := range videoIDs {
			wg.Add(1)
			go fetchThumbnail(client, videoID, &wg)
		}
		wg.Wait()
	} else {
		for _, videoID := range videoIDs {
			fetchThumbnail(client, videoID, nil)
		}
	}
}
