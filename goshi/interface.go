package goshi

// This interface implement the following methods for fetching and
// downloading each page of a Manga Chapter
type Scraper interface {
	// This method fetch all the pages of a chapter and send the
	// struct Page populated with name and URL of the jpg
	FetchChapter(chapterURL string, out chan<- Page) error
	// This method scrape a manga and return the list of the
	// available chapters
	ScrapeChapters(url string) ([]Chapter, error)
	// This method search for a manga and return the list of manga
	// availables
	SearchManga(search string) ([]Manga, error)
}
