package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"sort"

	"git.mrkeebs.eu/goshi/goshi"
)

func ZipFiles(zipName string, in <-chan goshi.Page, done chan<- struct{}) {
	var pages []goshi.Page

	// Add every page to pages and sort them
	for page := range in {
		pages = append(pages, page)
	}
	// sorting
	sort.SliceStable(pages, func(i, j int) bool {
		return pages[i].Num < pages[j].Num
	})

	// Create the zip files that will stores all the chapter
	zipFile, err := os.Create(zipName)
	if err != nil {
		log.Println(err)
	}

	// Close the zip after we finish
	defer zipFile.Close()

	// Create the writer that actually write files to the archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to ZipFiles
	for _, file := range pages {
		AddFile(zipWriter, file.Img, file.Name)
	}

	done <- struct{}{}
}

func AddFile(zipWriter *zip.Writer, file bytes.Buffer, nameFile string) {

	header := zip.FileHeader{
		Name:   nameFile + ".jpg",
		Method: zip.Store,
	}

	writer, _ := zipWriter.CreateHeader(&header)

	_, err := io.Copy(writer, &file)
	if err != nil {
		log.Println("Error durant zipping", err)
	}
}
