/*
   This package scrape results from mangaeden.com
*/

package scraper

import (
	"fmt"
	"log"

	"git.mrkeebs.eu/goshi/goshi"
	"github.com/PuerkitoBio/goquery"
)

type MangaEdenScraper struct {
}

const MangaEdenURL = "https://www.mangaeden.com"

func (m *MangaEdenScraper) ScrapeChapters(input string) []goshi.Chapter {
	// search for Italian manga
	urlIta := "/it/it-manga/"
	search := MangaEdenURL + urlIta + input

	var chs []goshi.Chapter
	doc, _ := goquery.NewDocument(search)
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		var ch goshi.Chapter

		ch.Name = s.Find("b").Eq(0).First().Text()
		link, _ := s.Find("a").Eq(0).First().Attr("href")
		ch.Link = MangaEdenURL + string(link)

		// Append each chapter to the slice
		chs = append(chs, ch)

	})

	return chs
}

// This Function get the chapterUrl in input and return a slice
// of string with all the pages URL e.g:
// Chapter 825, we got : all the pages in the chapter 825 as URL
func (m *MangaEdenScraper) FetchChapter(chapterURL string, out chan<- goshi.Page) {

	// Get the Url for all the pages
	doc, err := goquery.NewDocument(chapterURL)
	if err != nil {
		log.Println("Error with the URL:", err)
	}

	// Concurrently get the images that has to be downloaded
	doc.Find("#pageSelect option").Each(func(i int, s *goquery.Selection) {

		page, _ := s.Attr("value")

		var p goshi.Page
		pageURL := MangaEdenURL + page
		func(pageURL string) {
			d, _ := goquery.NewDocument(pageURL)
			img, _ := d.Find("a.next img").Attr("src")
			p.Referer = MangaEdenURL
			p.URL = "https:" + img
			p.Name = fmt.Sprintf("%03d", i)
			p.Num = i

			out <- p
		}(pageURL)

	})
}
