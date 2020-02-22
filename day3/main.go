package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "net/http"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error{
			{
				address := "http://google.com"
				if c.NArg()>0{
					address = c.Args().Get(0)
				}
				resp, err:= http.Get(address)
				if err!=nil{
					log.Fatal(err)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
				}
		return nil
		},

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}