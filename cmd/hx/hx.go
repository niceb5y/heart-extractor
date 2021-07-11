package main

import (
	"flag"
	"fmt"
	"github.com/niceb5y/heart-extractor/pkg/db"
	"github.com/niceb5y/heart-extractor/pkg/network"
	"github.com/niceb5y/heart-extractor/pkg/twitterclient"
	"sync"
)

var (
	version string
	build   string
)

func main() {
	fmt.Printf("HeartExtractor %s (%s)\n", version, build)

	db.InitDatabase()

	resetToken := flag.Bool("reset-token", false, "Reset access token")
	resetDownloadHistory := flag.Bool("reset-history", false, "Reset download history")
	flag.Parse()

	if *resetToken {
		db.ResetAccessToken()
		fmt.Println("Access token reset complete")
	}

	if *resetDownloadHistory {
		fmt.Println("History reset complete")
		db.ResetDownloadHistory()
	}

	client := twitterclient.GetClient()

	fmt.Println("Fetch favorite items...")
	var mediaURLs []string
	twitterclient.Fetch(client, nil, &mediaURLs)

	waitGroup := sync.WaitGroup{}
	for _, mediaURL := range mediaURLs {
		waitGroup.Add(1)
		mediaURL := mediaURL
		go func() {
			defer waitGroup.Done()
			network.Download(mediaURL)
		}()
	}
	waitGroup.Wait()

	fmt.Println("Done")
}
