package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"git.mrkeebs.eu/goshi/goshi"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// DownloadFile download the jpg given an URL, it uses
// a custom http-header for avoiding cloudflare 404 Error
func DownloadFile(in <-chan goshi.Page, out chan<- goshi.Page) {

	for page := range in {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", page.URL, nil)
		req.Header.Set("Referer", page.Referer)
		res, err := client.Do(req)
		if err != nil {
			log.Println("Client error:", err)
		}
		defer res.Body.Close()

		var data bytes.Buffer

		_, err = io.Copy(&data, res.Body)

		if err != nil {
			fmt.Println(err)
		}

		page.Img = data

		time.Sleep(1 * time.Second)
		go spinner(100 * time.Millisecond)
		out <- page

	}
}
