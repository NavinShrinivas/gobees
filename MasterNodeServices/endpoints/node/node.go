package node

import (
	"MasterGobees/utils"
	"MasterGobees/shell"
	"MasterGobees/globals"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/TwiN/go-color"
)

func MainNodeBirth(w http.ResponseWriter, r *http.Request){
  //Need to store statefully node info
  if r.Method != "POST"{
    utils.SimpleInvalidPath("Ivalid path", w)
    return
  }
  //Let's not store worker node info staefully, entirely in memory
  shell.PrintToShell(color.Colorize(color.Yellow,"[ENDPOINT] Request recived to add new worker node."))
	res_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err!=nil{
		log.Fatal(color.Colorize(color.Red,"[ENDPOINT ERROR] Error parsing request from Worker node addition request."))
	}
	new_worker_node := globals.WorkerNode{
		Ip_addr : res_body_obj["ip_addr"].(string),
		Port : res_body_obj["port"].(string),
		Files : make([]string,10000), //Limit of at max 10000 files
	}
	globals.WorkerNodesMetadata = append(globals.WorkerNodesMetadata,new_worker_node)
	utils.SimpleSuccesssStatus("Successfully Added node to cluster!", w)
	shell.PrintToShell(color.Colorize(color.Green,"Added one Node to cluster : "+new_worker_node.Ip_addr+":"+new_worker_node.Port))
	return
}


// [TODO] As for the deliverable, Node failures are not to be handled

// func MainNodeDeath(w http.ResponseWriter, r *http.Request){
//   //Need to remove Node info
// }
//
// func MainNodeAlive(w http.ResponseWriter, r *http.Request){
//   //Need to somehow spawn a function that Kills node if no reponse for 3 rounds
// }

