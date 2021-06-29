package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"git.mrkeebs.eu/goshi/goshi"
	"github.com/PuerkitoBio/goquery"
)

// Use an empty struct for the interface and implement the methods
type MangaWorldScraper struct {
}

// BaseURL
const MangaWorldURL = "https://www.mangaworld.io"

func (m *MangaWorldScraper) SearchManga(input string) []goshi.Manga {
	urlIta := "/it/it-directory/"
	search := MangaEdenURL + urlIta + "/?title=" + input

	var mangas []goshi.Manga
	doc, _ := goquery.NewDocument(search)
	doc.Find(".openManga").Each(func(i int, s *goquery.Selection) {
		var manga goshi.Manga

		manga.Name = s.Text()

		mangas = append(mangas, manga)
	})
	
	return mangas
}

// ScrapeChapters take the user search input, scrape with goquery all the
// the availables chapters and return a slice with all the chapters
func (m *MangaWorldScraper) ScrapeChapters(input string) []goshi.Chapter {
	search := MangaWorldURL + "/manga/1708/" + input

	var chs []goshi.Chapter
	doc, _ := goquery.NewDocument(search)
	doc.Find("a.chap").Each(func(i int, s *goquery.Selection) {
		var chapter goshi.Chapter
		chapter.Name, _ = s.Attr("title")
		chapter.Link, _ = s.Attr("href")

		chs = append(chs, chapter)
	})
	return chs
}

// FetchChapter get the chapterURL as input, fetch with goquery all the urls for the
// jpg images, and send them as a page struct into the channel
// Then another function will elaborate those structs
func (m *MangaWorldScraper) FetchChapter(chapterURL string, out chan<- goshi.Page) {

	// This orrible hacky workaround, found the last page of the chapter to fetch
	doc, _ := goquery.NewDocument(chapterURL)
	pag := doc.Find("select.page.custom-select option").Eq(0).First().Text()
	slashIndex := strings.Index(pag, "/")
	maxPage, _ := strconv.Atoi(pag[slashIndex+1:])

	for i := 1; i <= maxPage; i++ {
		pageUrl := chapterURL + fmt.Sprintf("/%d", i)
		doc, _ := goquery.NewDocument(pageUrl)

		var p goshi.Page
		img, _ := doc.Find("img.img-fluid").Eq(1).Attr("src")
		p.Referer = MangaWorldURL
		p.URL = img
		fmt.Println(p.URL)
		p.Name = fmt.Sprintf("%03d", i)
		p.Num = i
		out <- p
	}

}
