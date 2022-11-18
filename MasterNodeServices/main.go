package main

import (
	"MasterGobees/configuration"
	"MasterGobees/endpoints/home"
	"MasterGobees/endpoints/node"
	"MasterGobees/globals"
	"MasterGobees/shell"
	"MasterGobees/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/TwiN/go-color"
)

func NetworkEndpoints() {
	http.HandleFunc("/", home.MainHome)
	http.HandleFunc("/nodebirth", node.MainNodeBirth)
	http.HandleFunc("/resetnode", node.ResetNode)
	// http.HandleFunc("/nodedeath", node.MainNodeBirth)
	master_node_url := "0.0.0.0:" + globals.ServerPort
	log.Fatal(http.ListenAndServe(master_node_url, nil))
}

func MasterStartupSequence() {
	//Some gloabls inits :
	globals.MainWg = new(sync.WaitGroup)
	globals.NewCluster = false

	//--------------------
	//Ressting cluster settings :
	fmt.Print(color.Colorize(color.Red, "Resetting previous cluster data [Y]/n : "))
	var option string
	fmt.Scanln(&option)
	if option != "n" {
		globals.NewCluster = true
		os.RemoveAll("./temp_splits")
		os.Remove("./NodeMeta.json")
		os.Remove("./FileMeta.json")
	} else {
		os.RemoveAll("./temp_splits")
		_, err := os.Stat("./NodeMeta.json")
		if err != nil {
			log.Println(color.Colorize(color.Red, "No previous meta data found"))
			// globals.NewCluster = true
			os.Remove("./FileMeta.json")
			return
		}
		_, err = os.Stat("./FileMeta.json")
		if err != nil {
			log.Println(color.Colorize(color.Red, "No previous meta data found"))
			// globals.NewCluster = true
			return
		}
		//Need to initlise metadata from stored files
		json1, _ := os.ReadFile("./NodeMeta.json")
		json.Unmarshal(json1, &globals.WorkerNodesMetadata)
		json2, _ := os.ReadFile("./FileMeta.json")
		json.Unmarshal(json2, &globals.FileMetadata)
		os.Remove("./NodeMeta.json")
		os.Remove("./FileMeta.json")
	}
	//--------------------------
}

func main() {
	log.Println(color.Colorize(color.Yellow, "Starting Master node..."))
	MasterStartupSequence()

	//Command line arguments flag configs :
	flag.StringVar(&globals.Config_file_path, "config_path", "./config.json", "Path to configuration file")
	debug_flag_local := flag.Bool("debug", false, "Whether to print debug outputs or not")
	flag.StringVar(&globals.ServerPort, "port", "3001", "Port in which Master Node listens to.")
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
	go cntrlc()
	log.Println(color.Colorize(color.Green, "Listen on port "+globals.ServerPort))

	//Shell routines, all prints after shell Initialize must be printed only using print function in shell
	shell.Initialize()
	globals.MainWg.Wait()
}

func cntrlc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			utils.ExistSequence()
		}
	}()
}
