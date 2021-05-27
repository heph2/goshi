package scraper

import (
	"fmt"

	"git.mrkeebs.eu/goshi/goshi"
	"github.com/PuerkitoBio/goquery"
)

// Use an empty struct for the interface and implement the methods
type MangaNeloScraper struct {
}

// BaseURL
const MangaNeloURL = "https://manganelo.tv"

// ScrapeChapters take the user search input, scrape with goquery all the
// the availables chapters and return a slice with all the chapters
func (m *MangaNeloScraper) ScrapeChapters(input string) []goshi.Chapter {
	search := MangaNeloURL + "/search/" + input

	doc, _ := goquery.NewDocument(search)
	manga, _ := doc.Find("a.a-h.text-nowrap.item-title").Eq(0).First().Attr("href")

	// Search all the chapter of `manga`
	chapterSearch := MangaNeloURL + manga

	var chs []goshi.Chapter
	d, _ := goquery.NewDocument(chapterSearch)
	d.Find("a.chapter-name.text-nowrap").Each(func(i int, s *goquery.Selection) {
		var chapter goshi.Chapter

		chapter.Name = s.Text()
		link, _ := s.Attr("href")
		chapter.Link = MangaNeloURL + string(link)

		// Append each chapter
		chs = append(chs, chapter)

	})
	return chs
}

// FetchChapter get the chapterURL as input, fetch with goquery all the urls for the
// jpg images, and send them as a page struct into the channel
// Then another function will elaborate those structs
func (m *MangaNeloScraper) FetchChapter(chapterURL string, out chan<- goshi.Page) {

	doc, _ := goquery.NewDocument(chapterURL)

	//	index := 1
	doc.Find("img.img-loading").Each(func(i int, s *goquery.Selection) {
		var p goshi.Page
		img, _ := s.Attr("data-src")
		p.URL = img
		p.Name = fmt.Sprintf("%03d", i)
		p.Num = i
		out <- p

	})
}
