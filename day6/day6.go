package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var books = []string{"first","second"}
panic(err)
func main() {
	router := mux.NewRouter()
	router.Methods("GET").HandlerFunc(ReadBooks)
	router.Methods("POST").HandlerFunc(CreateBooks)
	router.Methods("POST").HandlerFunc(DeleteBooks)
	fmt.Println("Server started")

	http.ListenAndServe("0.0.0.0:8080",router)
}

func ReadBooks(w http.ResponseWriter, r *http.Request){

	data, err := json.Marshal(books)
	if err!=nil{
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	w.Write(data)
}

func CreateBooks(w http.ResponseWriter,r *http.Request){
	var data string
	decoder:=json.NewDecoder(r.Body)

	if err:=decoder.Decode(&data); err!=nil {
		w.Write([]byte("Error: " + err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	books = append(books, data)
	w.WriteHeader(http.StatusCreated)
}

func DeleteBooks(){

}
