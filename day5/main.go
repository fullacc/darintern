package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var (
	configFile = ""
	port = "8080"

	flags = []cli.Flag{
		&cli.StringFlag{
			Name: "config",
			Aliases: []string{"c"},
			Usage: "config file path",
			Destination: &configFile,
		},
	}
)


func main(){
	app := &cli.App{
		Action: func(c *cli.Context) error{
			{
				goserv()
			}
			return nil
		},
	}
	app.Flags = flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func goserv(){
	r:=mux.NewRouter()
	r.HandleFunc("/createfile", CreateFile)
	r.HandleFunc("/deletefile", DeleteFile)
	r.HandleFunc("/updatefile", Updatefile)
	r.HandleFunc("/",listHandler)
	http.Handle("/",r)
	http.ListenAndServe(":" + port, r)
}

type Message struct {
	Body string
}

type AutoGenerated struct {
	Books []struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
		Genre  string `json:"genre"`
	} `json:"books"`
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
 	var books AutoGenerated
	f, err := os.Open("package.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	file, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	// Unmarshaling given json
	err = json.Unmarshal(file, &books)
	if err != nil {
		panic(err)
	}

	// From JSON to struct
	fmt.Println(books)


}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/list", 301)
	}
	var bk Book
	strconv.
	book.ElectricAmount, _ = strconv.ParseInt(r.FormValue("ElectricAmount"), 10, 64)
	book.ElectricPrice, _ = strconv.ParseFloat(r.FormValue("ElectricPrice"), 64)
	book.WaterAmount, _ = strconv.ParseInt(r.FormValue("WaterAmount"), 10, 64)
	book.WaterPrice, _ = strconv.ParseFloat(r.FormValue("WaterPrice"), 64)
	.CheckedDate = r.FormValue("CheckedDate")
	fmt.Println(cost)

	// Save to database
	stmt, err := db.Prepare(`
		INSERT INTO cost(electric_amount, electric_price, water_amount, water_price, checked_date)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Prepare query error")
		panic(err)
	}
	_, err = stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
		cost.WaterAmount, cost.WaterPrice, cost.CheckedDate)
	if err != nil {
		fmt.Println("Execute query error")
		panic(err)
	}
	http.Redirect(w, r, "/list", 301)
}

func CreateFile(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nm, ok := vars[name]
		if !ok {
			resp, _ := json.Marshal(&Message{Body: "Error"})
			w.Write(resp)
			return
		}
		a := "Answer"
		w.Write([]byte(a))
		fmt.Println(nm)
		b := "result"
		w.Write([]byte(b))
		//out, err := exec.Command("ls", "-a").Output()
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(string(out))
		//jsonAnswer, err := json.Marshal(out)
		//if err != nil {
		//	panic(err)
		//}
		//w.Write(jsonAnswer)
	}
}

func CreateHandler(w http.ResponseWriter,r *http.Request) {
	out,err := exec.Command("mkdir",name).Output()
	if err!=nil{
		panic (err)
	}
	fmt.Println(string(out))
	jsonAnswer, err := json.Marshal(out)
	if err!=nil{
		panic(err)
	}
	w.Write(jsonAnswer)
}

func DeleteHandler(w http.ResponseWriter,r *http.Request) {
	out,err := exec.Command("rmdir",name).Output()
	if err!=nil{
		panic (err)
	}
	fmt.Println(string(out))
	jsonAnswer, err := json.Marshal(out)
	if err!=nil{
		panic(err)
	}
	w.Write(jsonAnswer)
}