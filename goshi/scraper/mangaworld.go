package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"git.mrkeebs.eu/goshi/goshi"
	"github.com/PuerkitoBio/goquery"
)

type MangaWorldScraper struct {
}

const MangaWorldURL = "https://www.mangaworld.io"

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

func (m *MangaWorldScraper) FetchChapter(chapterURL string, out chan<- goshi.Page) {

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
