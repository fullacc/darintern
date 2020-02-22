package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fullacc/darintern/day7/book_store"
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
	datapath = "/home/fullacc/go/src/github.com/fullacc/darintern/day7/books.json"
	port=""
	config="/home/fullacc/go/src/github.com/fullacc/darintern/day7/config.json"
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
	postgreBookStore, err := book_store.NewPostgreBookStore(configpath)
	if err != nil {
		panic(err)
	}

	endpoints := book_store.NewEndpointsFactory(bookStore)
	postgreendpoints := book_store.NewEndpointsFactory(postgreBookStore)
	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(endpoints.GetBook("id"))
	router.Methods("POST").Path("/").HandlerFunc(endpoints.CreateBook())
	router.Methods("GET").Path("/").HandlerFunc(endpoints.ListBooks())
	router.Methods("PUT").Path("/{id}").HandlerFunc(endpoints.UpdateBook("id"))
	router.Methods("DELETE").Path("/{id}").HandlerFunc(endpoints.DeleteBook("id"))

	router.Methods("GET").Path("/db/{id}").HandlerFunc(postgreendpoints.GetBook("id"))
	router.Methods("POST").Path("/db/").HandlerFunc(postgreendpoints.CreateBook())
	router.Methods("GET").Path("/db/").HandlerFunc(postgreendpoints.ListBooks())
	router.Methods("PUT").Path("/db/{id}").HandlerFunc(postgreendpoints.UpdateBook("id"))
	router.Methods("DELETE").Path("/db/{id}").HandlerFunc(postgreendpoints.DeleteBook("id"))
	fmt.Println("Server started")
	go func(port string, rtr *mux.Router) {
		http.ListenAndServe("0.0.0.0:" + port, rtr)
	}(configfile.Port, router)

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	go func() {
		<-c
		done <- true
	}()

	<- done
	log.Printf("server shutdown")
	ExitWithSave(configfile.JsonFilePath, bookStore)
	os.Exit(1)

	return nil
}


func ExitWithSave(filepath string, store book_store.BookStore) {
	err := store.SaveBooks(filepath)
	if err != nil {
		fmt.Println(err)
	}
}