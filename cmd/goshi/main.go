/*
 * The main package get the manga name from arguments, then print
 * results, ask again to the user which chapter choose to download.
 * Scrape how many pages has the selected chapter, and then spawn
 * n goroutines for each pages to download. After that, send through
 * a channel the jpg to another goroutines that sequentually add
 * those jpg to a zip archive.
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"git.mrkeebs.eu/goshi/goshi"
	"git.mrkeebs.eu/goshi/goshi/scraper"
)

var (
	fetchPtr   = flag.String("fetch", "", "Fetch the available chapters for the manga selected")
	idPtr      = flag.String("down", "", "Download a chapter given an ID")
	scraperPtr = flag.String("scraper", "", "Select scraper")
	searchPtr  = flag.String("search", "", "Search a Manga")
)

// GetScaper use scraperPtr and choose with which scraper fetch
func GetScraper(scr string) (goshi.Scraper, error) {
	switch scr {
	case "mangaeden":
		return &scraper.MangaEdenScraper{}, nil
	case "manganelo":
		return &scraper.MangaNeloScraper{}, nil
	case "mangaworld":
		return &scraper.MangaWorldScraper{}, nil
	}

	return nil, errors.New("Source not supported")
}

// GetChapter use idPtr for search which chapter(s) should be downloaded
func GetChapter(str string) []int {
	var chapters []int
	var s []string
	var t []string

	// check if we need to
	// download from an episode to another one; or
	// we need to download different episodes
	checkComma := strings.Index(str, ",")
	checkSep := strings.Index(str, "-")

	// Here for different episodes!
	if checkComma != -1 {
		s = strings.Split(str, ",")
		for _, value := range s {
			num, _ := strconv.Atoi(value)
			chapters = append(chapters, num)
		}
		return chapters
	}

	// Here for "from - to"
	if checkSep != -1 {
		t = strings.Split(str, "-")
		from, _ := strconv.Atoi(t[0])
		to, _ := strconv.Atoi(t[len(t)-1])
		for i := from; i <= to; i++ {
			chapters = append(chapters, i)
		}
		return chapters
	}

	// Now if the input is a single episode
	num, _ := strconv.Atoi(str)
	chapters = append(chapters, num)

	return chapters

}

func main() {
	flag.Parse()

	// This flag must be passed
	if *scraperPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check which scraper is passed
	scraper, err := GetScraper(*scraperPtr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if *searchPtr != "" {
		available, err := scraper.SearchManga(*searchPtr)
		if err != nil {
			log.Printf("failed to search for %q: %s",
				*searchPtr, err)
			os.Exit(1)
		}

		for _, a := range available {
			fmt.Println(a.Name)
		}
	}

	// Search for the chapters available and print the name of the
	// chapter associated to an ID
	ids := make(map[int]string)
	if *fetchPtr != "" {
		chapters, err := scraper.ScrapeChapters(*fetchPtr)
		if err != nil {
			log.Printf("failed to fetch chapters: %s", err)
			os.Exit(1)
		}

		for i, chapter := range chapters {
			fmt.Printf("ID:%d \t %s\n", i, chapter.Link)
			ids[i] = chapter.Link
		}
	}

	if *idPtr != "" {
		// Get the Chapter to download from the user input
		var chapters []int
		chapters = GetChapter(*idPtr)
		// Start the pool of goroutines with the scraper selected
		for _, chapter := range chapters {
			Pool(scraper, ids[chapter])
		}
	}
}
