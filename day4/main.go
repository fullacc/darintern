package fourth

import (
	"fmt"
	"io/ioutil"
)

func main(){
	f,err := ioutil.ReadDir("/")
	if err != nil {
		panic(err)
	}
	fmt.Println(f)
}
