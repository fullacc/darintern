package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type BookStore interface {

	CreateBook (book *Book) (*Book,error)

	GetBook (id string) (*Book, error)

	ListBooks () ([]*Book,error)

	UpdateBook (id string, book *Book) (*Book,error)

	DeleteBook(id string) error

}

func NewBookStore(filename string) (BookStore, error){
	f, err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	file, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// Unmarshaling given json
	err = json.Unmarshal(file, &books)
	if err != nil {
		return nil, err
	}
	return &BookStore{},nil
}

type BookStore struct{
	books []*Book
}

func(bsc *BookStore)GetBook (id string)(*Book,error){
	for _,v := range bsc.books{
		if v.ID==id {
			return v,nil
		}
	}
	return nil,errors.New("Not found")
}

func (bsc *BookStore)

type Book struct{
	Title string 'json:"title,omitempty"'
	Author string 'json:"author,omitempty"'

}
