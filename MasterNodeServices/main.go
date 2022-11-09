package main

import (
	"MasterGobees/configuration"
	"MasterGobees/endpoints/home"
	"MasterGobees/endpoints/node"
	"MasterGobees/globals"
	"flag"
	"log"
	"net/http"
	"sync"
	"github.com/TwiN/go-color"
)

func NetworkEndpoints() {
	http.HandleFunc("/", home.MainHome)
	http.HandleFunc("/nodebirth", node.MainNodeBirth)
	http.HandleFunc("/nodedeath", node.MainNodeBirth)
	master_node_url := "0.0.0.0:"+globals.ServerPort
	log.Fatal(http.ListenAndServe(master_node_url, nil))
}

func main() {
	log.Println(color.Colorize(color.Yellow, "Starting Master node..."))

  //Some gloabls inits : 
  globals.MainWg = new(sync.WaitGroup)
  //--------------------

  //Command line arguments flag configs : 
  flag.StringVar(&globals.Config_file_path, "config_path", "./config.json", "Path to configuration file")
  debug_flag_local := flag.Bool("debug", false, "Whether to print debug outputs or not")
 	flag.StringVar(&globals.ServerPort,"port","3001","Port in which Master Node listens to.") 
  flag.Parse()
  globals.Debug_flag = *debug_flag_local //Pushing value to global variable
  //-------------------------------------

	//Spawn configuration management routines
	err := configuration.ConfigurationMain()
	if err != nil {
		log.Fatal(color.Colorize(color.Red, err.Error()))
	}

  //Network endpoint routines
  globals.MainWg.Add(1)
  go NetworkEndpoints()
  log.Println(color.Colorize(color.Green,"Listen on port "+globals.ServerPort))
  globals.MainWg.Wait()

	//Shell routines

}
