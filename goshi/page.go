package goshi

import "bytes"

type Page struct {
	Referer string
	Name    string
	URL     string
	Img     bytes.Buffer
	Num     int
}
