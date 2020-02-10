package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fullacc/darintern/day6hw/book_store"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	datapath = "/home/fullacc/go/src/github.com/fullacc/day6hw/books.json"
	port="8080"
	config="/home/fullacc/go/src/github.com/fullacc/day6hw/config.json"
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "config /filepath",
			Required:    true,
			Destination: &config,
		},
	}
)

func main() {
	app := &cli.App{
		/*		Action: func(c *cli.Context) error{
				{
					if c.NArg()>0 {
						address = c.Args().Get(1)
					}
					DirPrinter(address,0)
				}
				return nil
			},*/
		Commands: []*cli.Command{
			{
				Name:    "launch",
				Aliases: []string{"l"},
				Usage:   "launch",
				Action:  run,
			},
		},
	}
	app.Flags=flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err:=LaunchServer(config);err!=nil {
		return err
	}

	return nil
}

func LaunchServer(configpath string) error{
	file, err := os.Open(configpath)
	if err != nil {
		return err
	}
	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}
	var configfile *book_store.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		return err
	}
	file.Close()
	bookStore, err := book_store.NewBookStore(configfile.JsonFilePath)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(GetBook(bookStore, "id"))
	router.Methods("POST").Path("/").HandlerFunc(CreateBook(bookStore))
//	router.Methods("GET").Path("/exit").HandlerFunc(ExitWithSave(configfile.JsonFilePath,&bookStore))
	router.Methods("GET").Path("/").HandlerFunc(ListBooks(bookStore))
	router.Methods("PUT").Path("/{id}").HandlerFunc(UpdateBook(bookStore, "id"))
	router.Methods("DELETE").Path("/{id}").HandlerFunc(DeleteBook(bookStore, "id"))
	fmt.Println("Server started")
	http.ListenAndServe("0.0.0.0:" + configfile.Port, router)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	go func(filepath string, bookStore *book_store.BookStore){
		for sig := range c {
			log.Printf("captured %v",sig)
			ExitWithSave(filepath,bookStore)
			os.Exit(1)
		}
	}(configfile.JsonFilePath,&bookStore)

	return nil
}

func GetBook(store book_store.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		book, err := store.GetBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("I'm sorry: " + err.Error()))
			return
		}
		data, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func CreateBook(store book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		book := &book_store.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := store.Create(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func ExitWithSave(filepath string, store *book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := store.SaveBooks(filepath)
		if err != nil {
			return err
		}
	}
}

func ListBooks(store book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		books, err := store.ListBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("I'm sorry: " + err.Error()))
			return
		}
		data, err := json.Marshal(books)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		/*for _,v := range books {
			data, err := json.Marshal(v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error: " + err.Error()))
				return
			}
			w.Write(data)
		}*/
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateBook(store book_store.BookStore,idParam string) func (w http.ResponseWriter,r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		book := &book_store.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := store.UpdateBook(id,book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteBook(store book_store.BookStore, idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		vars := mux.Vars(r)
		id,ok := vars[idParam]
		if !ok {
			w.Write([]byte("Error: Not Found"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := store.DeleteBook(id)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}