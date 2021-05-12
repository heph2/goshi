package models

import "bytes"

type Chapter struct {
	Name string
	Link string
}

type Page struct {
	Name string
	URL  string
	Img  bytes.Buffer
}
