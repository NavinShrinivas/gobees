package main

import (
	"WorkerGobees/globals"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TwiN/go-color"
)

func main() {
	log.Println(color.Colorize(color.Yellow, "Starting worker node..."))

	var_temp := ""
	flag.StringVar(&var_temp, "master", "http://0.0.0.0:3001/", "Path to Masternode")
	flag.Parse()
	globals.Master_Url = var_temp

	res, err := http.Get(globals.Master_Url)

	if err != nil {
		log.Fatal(err)
	}
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res_body))
}
