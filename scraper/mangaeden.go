/*
   This package scrape results from mangaeden.com
*/

package scraper

import (
	"git.mrkeebs.eu/goshi/models"
	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "https://www2.mangaeden.com"
const baseURLIta string = baseURL + "/it/it-manga/"

// This function use goquery for searching through mangaeden
// and scrape all the chapter of the searched manga.
// Then add it to a slice and returns all the chapters available.
func MangaEden(input string) (results []models.Chapter) {

	search := baseURLIta + input

	doc, _ := goquery.NewDocument(search)
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		var ch models.Chapter

		ch.Name = s.Find("b").Eq(0).First().Text()
		link, _ := s.Find("a").Eq(0).First().Attr("href")
		ch.Link = baseURL + string(link)

		// Append each chapter to the slice
		results = append(results, ch)

	})

	return results
}

// This Function get the chapterUrl in input and return a slice
// of string with all the pages URL e.g:
// Chapter 825, we got : all the pages in the chapter 825 as URL
func MangaEdenChapter(chapter string, out chan<- models.Page) {

	// Get the Url for all the pages
	doc, _ := goquery.NewDocument(chapter)

	// Concurrently get the images that has to be downloaded
	doc.Find("#pageSelect option").Each(func(i int, s *goquery.Selection) {

		page, _ := s.Attr("value")

		var p models.Page
		pageURL := baseURL + page
		func(pageURL string) {
			d, _ := goquery.NewDocument(pageURL)
			img, _ := d.Find("a.next img").Attr("src")
			p.URL = "https:" + img
			p.Name = page

			out <- p
		}(pageURL)

	})
}
