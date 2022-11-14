package backend

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type Calibre struct {
	bookLibrary BookLibrary
}

func NewCalibre(b BookLibrary) *Calibre {
	return &Calibre{
		bookLibrary: b,
	}
}

func calibreExists() bool {
	_, err := exec.LookPath("calibredb")
	return err != nil
}

type CalibreBook struct {
	Author  string   `json:"authors"`
	Cover   string   `json:"cover"`
	Formats []string `json:"formats"`
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
}

func getCalibreBooks() ([]CalibreBook, error) {
	books := []CalibreBook{}
	// TODO have user specify calibre dictionary (optional)
	calibre := exec.Command("calibredb", "list",
		"--for-machine",
		"--fields", "cover,authors,title,formats")
	output, err := calibre.Output()
	if err != nil {
		return books, err
	}
	err = json.Unmarshal(output, &books)
	if err != nil {
		return books, err
	}

	return books, nil
}

func (c *Calibre) ImportCalibreBooks() error {
	log.Println("Loading calibre")
	books, err := getCalibreBooks()
	if err != nil {
		log.Println("Failed", err)
		return err
	}
	for _, book := range books {
		log.Println("Trying", book.Author, book.Title)
		exists, err := c.bookLibrary.BookExists(book.Author, book.Title)
		if err != nil {
			log.Println("error ", err)
			return err
		}
		if !exists {
			log.Println("Potential new book", book.Author, book.Title)
			for _, format := range book.Formats {
				if strings.HasSuffix(format, ".txt") {
					err := c.bookLibrary.AddBook(book.Author, book.Title, book.Cover, format)
					if err != nil {
						log.Println("error ", err)
						return err
					}
					break
				}
			}
		} else {
			log.Println("It exists")

		}
	}
	return nil
}
