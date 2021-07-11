package network

import (
	"fmt"
	"github.com/niceb5y/heart-extractor/pkg/db"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Download(mediaURL string) {
	filename := getFilenameFromURL(mediaURL)

	if db.IsExistInDownloadHistory(filename) {
		fmt.Println("Skipping", filename)
	} else {
		resp, err := http.Get(mediaURL)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Download", filename)
		path := filepath.Join(homeDir, "Downloads", filename)

		file, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Println(err)
		}

		db.AddDownloadHistory(filename)
	}
}

func getFilenameFromURL(mediaURL string) string {
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		log.Println(err)
	}

	path := parsedURL.Path
	segments := strings.Split(path, "/")

	return segments[len(segments)-1]
}
