package book_store

import "testing"

var bookStore BookStore

func TestPostgreStore_Create(t *testing.T) {
	var err error
	bookStore, err = NewPostgreBookStore("../config.json")
	if err != nil {
		t.Error(err)
	}

	newbook := Book{Title:"GOO",Author:"GOAUTHOR",Description:"GOBOOK",NumberOfPages:123,ID:123}

	book, err := bookStore.Create(&newbook)
	if err != nil {
		t.Error(err)
	}
	book, err = bookStore.GetBook(123)
	if err != nil {
		t.Error(err)
	}
	if book.Title != "GOO" {
		t.Error("Title is not equal")
	}
	if book.Author != "GOAUTHOR" {
		t.Error("FAKE AUTHOR")
	}
	if book.Description != "GOBOOK" {
		t.Error("Wrong description")
	}
	if book.NumberOfPages != 1222 {
		t.Error("Fake pages")
	}
}

func TestPostgreStore_GetBook(t *testing.T) {
	var err error
	bookStore, err = NewPostgreBookStore("../config.json")
	if err != nil {
		t.Error(err)
	}
	book, err := bookStore.GetBook(10)
	if err != nil {
		t.Error(err)
	}
	if book.Title != "JUMNJI" {
		t.Error("Title is not equal")
	}
	if book.Author != "POP" {
		t.Error("NOT THE SET AUTHOR")
	}
	if book.Description != "book" {
		t.Error("Wrong description")
	}
	if book.NumberOfPages != 1222 {
		t.Error("Fake pages")
	}
}


func TestPostgreStore_DeleteBook(t *testing.T) {
	var err error
	bookStore, err = NewPostgreBookStore("../config.json")
	if err != nil {
		t.Error(err)
	}
	err = bookStore.DeleteBook(123)
	if err != nil {
		t.Error(err)
	}
}