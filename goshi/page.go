package goshi

import "bytes"

// Page struct list referer http referer headers for download
// name of the page, URL that should me passed into the channel to
// download. Img is the downloaded img in bytes, that should be
// passed to the archiver function.
type Page struct {
	Referer string
	Name    string
	URL     string
	Img     bytes.Buffer
	Num     int
}
