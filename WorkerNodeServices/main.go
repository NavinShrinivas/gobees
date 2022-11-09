package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TwiN/go-color"
)

func main(){
  log.Println(color.Colorize(color.Yellow,"Starting worker node..."))

  //Simple testing network request
  res, err := http.Get("http://0.0.0.0:3001/")
  if err!=nil{
    log.Fatal(err)
  }
  res_body, err:= ioutil.ReadAll(res.Body)
  if err!=nil{
    log.Fatal(err) 
  }
  fmt.Println(string(res_body))
}
