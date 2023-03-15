package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"
)

type Calibre struct {
	backend *Backend
}

func NewCalibre(backend *Backend) *Calibre {
	return &Calibre{
		backend: backend,
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
	Exists  bool     `json:"exists"`
}

func (c *Calibre) GetCalibreBooks() ([]CalibreBook, error) {
	books := []CalibreBook{}
	// TODO have user specify calibre dictionary (optional)
	calibre := exec.Command("calibredb", "list",
		"--for-machine",
		"--fields", "cover,authors,title,formats")

	var stdout, stderr bytes.Buffer
	calibre.Stdout = &stdout
	calibre.Stderr = &stderr
	err := calibre.Run()
	if err != nil {
		return books, errors.New(stderr.String())
	}
	err = json.Unmarshal(stdout.Bytes(), &books)
	if err != nil {
		return books, err
	}

	for i := range books {
		exists, err := c.backend.BookLibrary.BookExists(
			books[i].Author, books[i].Title)
		if err != nil {
			return nil, err
		}
		books[i].Exists = exists
	}

	return books, nil
}

func (c *Calibre) ImportCalibreBooks(books []CalibreBook) error {
	c.backend.setupProgress("Processing calibre books", len(books))
	for _, book := range books {
		log.Println("Trying", book.Author, book.Title)
		exists, err := c.backend.BookLibrary.BookExists(book.Author, book.Title)
		if err != nil {
			log.Println("error ", err)
			return err
		}
		if !exists {
			log.Println("Potential new book", book.Author, book.Title)
			for _, format := range book.Formats {
				if strings.HasSuffix(format, ".txt") {
					bookId, err := c.backend.BookLibrary.AddBook(book.Author, book.Title, book.Cover, format)
					if err != nil {
						log.Println("error ", err)
						return err
					}
					book, err := c.backend.BookLibrary.GetBook(bookId)
					if err != nil {
						log.Println("error ", err)
						return err
					}
					c.backend.Generator.GenerateSentenceTableForBook(book)

					break
				}
			}
		} else {
			log.Println("It exists")
		}
		c.backend.KnownWords.SyncFrequency()
		c.backend.progress()
	}
	return nil
}
