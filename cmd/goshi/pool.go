package main

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"git.mrkeebs.eu/goshi/goshi"
)

const workers = 5

// Pool spawn n goroutines that will download concurrently
// all the pages of a selected chapter and add them to a zip archive.
func Pool(scraper goshi.Scraper, chapterURL string) {
	in := make(chan goshi.Page)
	out := make(chan goshi.Page)

	var wg sync.WaitGroup

	// Start spawning the goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			DownloadFile(in, out)
			wg.Done()
		}()
	}

	done := make(chan struct{})

	// Create a name for the zip files, based on the chapterURL
	zipName := func(chapterURL string) string {
		u, _ := url.Parse(chapterURL)
		zipName := strings.Replace(u.Path, "/", "-", -1) + ".cbz"
		return zipName[1:]
	}(chapterURL)

	go ZipFiles(zipName, out, done)

	// Start fetching
	scraper.FetchChapter(chapterURL, in)
	close(in)

	log.Println("before wait")
	// Wait until all the goroutines ends
	wg.Wait()

	// Signal that all the downloads are finished
	close(out)

	<-done
	log.Println("DONE")

}
