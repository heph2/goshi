package cmd

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"

	"git.mrkeebs.eu/goshi/models"
)

func ZipFiles(zipName string, in <-chan models.Page, done chan<- struct{}) {
	var pages []models.Page

	// Add every page to pages
	for page := range in {
		pages = append(pages, page)
	}

	// Create the zip files that will stores all the chapter
	zipFile, _ := os.Create(zipName)

	// Close the zip after we finish
	defer zipFile.Close()

	// Create the writer that actually write files to the archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to ZipFiles
	for _, file := range pages {
		AddFile(zipWriter, file.Img, file.Name)
		log.Println("Added")
	}

	done <- struct{}{}
}

func AddFile(zipWriter *zip.Writer, file bytes.Buffer, nameFile string) {

	header := zip.FileHeader{
		Name:   nameFile + ".jpg",
		Method: zip.Store,
	}

	writer, _ := zipWriter.CreateHeader(&header)

	_, _ = io.Copy(writer, &file)
}
