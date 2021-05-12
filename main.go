/*
   The main package get the manga name from arguments, then print
   results, ask again to the user which chapter choose do download.
   Scrape how many pages has the selected chapter, and then spawn
   n goroutines for each pages to download. After that, send through
   a channel the jpg to another goroutines that sequentually add
   those jpg to a zip archive.
*/
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"git.mrkeebs.eu/goshi/cmd"
	"git.mrkeebs.eu/goshi/models"
	"git.mrkeebs.eu/goshi/scraper"
)

// Those are the numbers of goroutines that will work in parallel
const workers = 5

// This function take the chapter selected by the user and build
// the name for the zip Archive
func parseInput(userInput string) string {
	u, _ := url.Parse(userInput)
	zipName := strings.Replace(u.Path, "/", "-", -1) + ".cbz"

	return zipName
}

// This function spawn n goroutines that will download concurrently
// all the pages of a selected chapter and add them to a zip archive.
func pool(selChapter, name string) {
	in := make(chan models.Page)
	out := make(chan models.Page)

	var wg sync.WaitGroup

	// Start spawning the goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			cmd.DownloadFile(in, out)
			wg.Done()
		}()
	}

	// Use a channel for check if the zipping
	// process is done
	done := make(chan struct{})
	go cmd.ZipFiles(name, out, done)

	// Now start fetching the urls from mangaeden
	scraper.MangaEdenChapter(selChapter, in)

	// Signal the end of the fetching
	close(in)

	log.Println("before wait")
	// Wait until all the goroutines ends
	wg.Wait()

	// Signal that all the downloads are finished
	close(out)

	<-done
	log.Println("DONE")

}

func main() {
	// Get the manga's name to search from arguments.
	// Each words must be separated by a '-' for building
	// the correct URL
	search := strings.Join(os.Args[1:], "-")

	// Now scrape the results and print it on stdout
	// Then associate each chapter to its index, so
	// the user can choose which capter download
	choose := make(map[int]string)

	table := scraper.MangaEden(search)
	for i, chapter := range table {
		fmt.Printf("%d\t%s\n", i, chapter.Name)
		choose[i] = chapter.Link
	}

	// Now ask the user to choose which chapter download
	var input string
	fmt.Print("Choose which chapter download: ")
	fmt.Scan(&input)

	// Chapters to download
	ch := strings.Split(input, ",")
	var num []int
	for i := 0; i < len(ch); i++ {
		v, _ := strconv.Atoi(ch[i])
		num = append(num, v)
	}

	for _, n := range num {
		// Get the name for the zip Archives
		name := parseInput(choose[n])

		// Start the pools of goroutines
		pool(choose[n], name)
	}
}
