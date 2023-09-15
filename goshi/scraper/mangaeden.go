package scraper

import (
	"fmt"
	"regexp"

	"git.mrkeebs.eu/goshi/goshi"
	"github.com/PuerkitoBio/goquery"
)

// Use an empty struct for the interface and implement the methods
type MangaEdenScraper struct {
}

// BaseURL
const MangaEdenURL = "https://www.mangaeden.com"

func (m *MangaEdenScraper) SearchManga(input string) ([]goshi.Manga, error) {
	urlIta := "/it/it-directory/"
	search := MangaEdenURL + urlIta + "/?title=" + input

	var mangas []goshi.Manga
	doc, err := goquery.NewDocument(search)
	if err != nil {
		return mangas, err
	}

	doc.Find("#mangaList tr").Each(func(i int, s *goquery.Selection) {
		var manga goshi.Manga

		name, _ := s.Find("td a").Attr("href")
		re := regexp.MustCompile(`manga\/(.*?)\/`)
		match := re.FindStringSubmatch(name)

		if len(match) > 1 {
			manga.Name = match[1]
		}

		mangas = append(mangas, manga)
	})

	return mangas, nil
}

// ScrapeChapters take the user search input, scrape with goquery all the
// the availables chapters and return a slice with all the chapters
func (m *MangaEdenScraper) ScrapeChapters(input string) ([]goshi.Chapter, error) {
	// search for Italian manga
	urlIta := "/it/it-manga/"
	search := MangaEdenURL + urlIta + input

	var chs []goshi.Chapter
	doc, err := goquery.NewDocument(search)
	if err != nil {
		return chs, err
	}

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		var ch goshi.Chapter

		ch.Name = s.Find("b").Eq(0).First().Text()
		link, _ := s.Find("a").Eq(0).First().Attr("href")
		ch.Link = MangaEdenURL + string(link)

		// Append each chapter to the slice
		chs = append(chs, ch)

	})

	return chs, nil
}

// FetchChapter get the chapterURL as input, fetch with goquery all the urls for the
// jpg images, and send them as a page struct into the channel
// Then another function will elaborate those structs
func (m *MangaEdenScraper) FetchChapter(chapterURL string, out chan<- goshi.Page) error {

	// Get the Url for all the pages
	doc, err := goquery.NewDocument(chapterURL)
	if err != nil {
		return err
	}

	doc.Find("#pageSelect option").Each(func(i int, s *goquery.Selection) {

		page, _ := s.Attr("value")

		var p goshi.Page
		pageURL := MangaEdenURL + page
		d, _ := goquery.NewDocument(pageURL)
		img, _ := d.Find("a.next img").Attr("src")
		p.Referer = MangaEdenURL
		p.URL = "https:" + img
		p.Name = fmt.Sprintf("%03d", i)
		p.Num = i

		out <- p
	})

	return nil
}
