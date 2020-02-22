package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/urfave/cli"
	"strconv"
)
const (
	DirColor = "\033[1;33m%s\033[0m\n"
	FileColor   = "\033[1;31m%s\033[0m\n"
	DirColor2 = "\033[1;35m%s\033[0m\n"
	FileColor2   = "\033[1;36m%s\033[0m\n"
	)
var (
	address = "/home/fullacc/go/src/fullacc"
	colorscheme ="1"
	deepness = "10000"
	files="1"

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "depth",
			Aliases:     []string{"d"},
			Destination: &deepness,
			Usage: "depth of path print",
		},
		&cli.StringFlag{
			Name:        "color",
			Aliases:     []string{"c"},
			Destination: &colorscheme,
			Usage: "color scheme of path print",
		},
		&cli.StringFlag{
			Name:        "files",
			Aliases:     []string{"f"},
			Destination: &files,
			Usage: "print files or not",
		},
	}
)
func main(){
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
				Name:    "print",
				Aliases: []string{"p"},
				Usage:   "printpath",
				Action:  func(c *cli.Context) error {
					if c.NArg()>0 {
						address = c.Args().Get(0)
					}
					DirPrinter(address,0,0)
					fmt.Printf("\n")
					return nil
				},
			},
		},
	}
	app.Flags=flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func DirPrinter(a string, depth int,last int){
	deep, err := strconv.Atoi(deepness)
	if depth>deep {
		return
	}
	f,err := ioutil.ReadDir(a)
	if err != nil {
		panic(err)
	}
	for j,r := range f {
		if  j==(len(f)-1) {
			last=1
		} else {
			last=0
		}
		if depth>0 {
			if last<1 {
				for i := 0; i < depth; i++ {
					fmt.Printf("⏐    ")
				}
			}	else {
				for i := 0; i < depth; i++ {
					fmt.Printf("     ")
				}
			}

		}
		fmt.Printf("⎿")
		if r.IsDir()==true {
			if(colorscheme=="1") {
				fmt.Printf(DirColor, r.Name())
			} else {
					fmt.Printf(DirColor2, r.Name())
			}
			DirPrinter(a+string("/")+r.Name(),depth+1,last)
		} else {
			if files=="1" {
				if (colorscheme == "1") {
					fmt.Printf(FileColor, r.Name())
				} else {
					fmt.Printf(FileColor2, r.Name())
				}
			}
		}
	}
}