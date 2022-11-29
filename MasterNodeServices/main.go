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
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/TwiN/go-color"
)

func NetworkEndpoints() {
	http.HandleFunc("/", home.MainHome)
	http.HandleFunc("/health", home.Health)
	http.HandleFunc("/nodebirth", node.MainNodeBirth)
	http.HandleFunc("/resetnode", node.ResetNode)
	// http.HandleFunc("/nodedeath", node.MainNodeBirth)
	master_node_url := "0.0.0.0:" + globals.ServerPort
	log.Fatal(http.ListenAndServe(master_node_url, nil))
}

func MasterStartupSequence() {
	//Some gloabls inits :
	globals.MainWg = new(sync.WaitGroup)
	//--------------------
	// Reseting cluster settings :
	if globals.NewCluster {
		log.Println(color.Colorize(color.Green, "Starting new cluster..."))
		os.RemoveAll("./temp_splits")
		os.Remove("./NodeMeta.json")
		os.Remove("./FileMeta.json")
	} else {
		log.Println(color.Colorize(color.Green, "Joining existing cluster..."))
		os.RemoveAll("./temp_splits")
		_, err := os.Stat("./NodeMeta.json")
		if err != nil {
			log.Println(color.Colorize(color.Red, "ERROR : No NodeMeta.json file was found."))
			// globals.NewCluster = true
			os.Remove("./FileMeta.json")
			return
		}
		_, err = os.Stat("./FileMeta.json")
		if err != nil {
			log.Println(color.Colorize(color.Red, "ERROR : No FileMeta.json file was found."))
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		dir = "."
	}
	globals.InitialDirectory = dir
	//--------------------------
}

func main() {
	log.Println(color.Colorize(color.Green, "Starting Master Node..."))

	//Command line arguments flag configs :
	flag.StringVar(&globals.Config_file_path, "config_path", "./config.json", "Path to configuration file")
	debug_flag_local := flag.Bool("debug", false, "Whether to print debug outputs or not")
	new_cluster_flag := flag.Bool("new_cluster", false, "Whether to start a new cluster or not")
	flag.StringVar(&globals.ServerPort, "port", "3000", "Port in which Master Node listens to.")

	flag.Parse()
	//Pushing value to global variable
	globals.Debug_flag = *debug_flag_local
	globals.NewCluster = *new_cluster_flag
	//-------------------------------------

	MasterStartupSequence()

	//Spawn configuration management routines
	err := configuration.ConfigurationMain()
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "ERROR: "+err.Error()))
	}

	//Network endpoint routines
	globals.MainWg.Add(1)
	go NetworkEndpoints()
	go cntrlc()
	log.Println(color.Colorize(color.Green, "Listening on port "+globals.ServerPort))

	//Shell routines, all prints after shell Initialize must be printed only using print function in shell
	shell.Initialize()
	globals.MainWg.Wait()
}

func cntrlc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			utils.ExitSequence()
		}
	}()
}
