package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.mrkeebs.eu/goshi/models"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// This function download the jpg given an URL, it uses
// a custom http-header for avoiding cloudflare 404 Error
func DownloadFile(in <-chan models.Page, out chan<- models.Page) {

	for page := range in {

		client := &http.Client{}
		req, _ := http.NewRequest("GET", page.URL, nil)
		req.Header.Set("Referer", "https://www2.mangaeden.com")
		res, _ := client.Do(req)

		var data bytes.Buffer

		_, _ = io.Copy(&data, res.Body)

		page.Img = data

		time.Sleep(1 * time.Second)
		go spinner(100 * time.Millisecond)
		out <- page

		defer res.Body.Close()
	}
}
