package main

import (
	"WorkerGobees/endpoints/home"
	"WorkerGobees/endpoints/lifestatus"
	"WorkerGobees/endpoints/data"
	"WorkerGobees/globals"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/TwiN/go-color"
)

func mainHttpHandler(){
	http.HandleFunc("/", home.MainHome)
	http.HandleFunc("/storefile", data.StoreFile)
	// http.HandleFunc("/nodedeath", node.MainNodeBirth)
	master_node_url := globals.Ip+":"+globals.Port
	log.Fatal(http.ListenAndServe(master_node_url, nil))
}

var Mainwg *sync.WaitGroup
func main() {
	log.Println(color.Colorize(color.Yellow, "Starting worker node..."))

	flag.StringVar(&globals.MasterUrl, "master", "http://0.0.0.0:3001/", "Path to Masternode")
	flag.StringVar(&globals.Port, "port", "3002", "Port worker node listens to")
	flag.StringVar(&globals.Ip, "ip", "0.0.0.0", "IP addr of worker node in same network as Master")
	flag.Parse()

	//Checking if given port is free
	ln, err := net.Listen("tcp", ":" + globals.Port)
  if err != nil {
    log.Fatal(color.Colorize(color.Red, "Given port is not free, please give some other port."))
  }
  ln.Close()

	//Checking is Master URL is correct
	log.Println(color.Colorize(color.Yellow, "Checking status of master node..."))
	res, err := http.Get(globals.MasterUrl)
	if err != nil {
		log.Fatal(err)
	}
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err!=nil{
		log.Fatal(color.Colorize(color.Red,"Something went wrong reading reponse from master node, please check master node url."))
	}
	if res_body_obj["status"] == false{
		log.Fatal("Master Node doesn't seem to be ready, please check master nodes status.")
	}
	log.Println(color.Colorize(color.Green, "Worker node is ready, registering worker node."))
	lifestatus.NodeBirthRegister() //Register Worker Node

	//Starting main service
	Mainwg = new(sync.WaitGroup)
	Mainwg.Add(1)
	go mainHttpHandler()
	Mainwg.Wait()
	//Need to start listening for stuff from master node from here on
}
