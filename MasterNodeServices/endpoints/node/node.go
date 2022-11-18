package node

import (
	"MasterGobees/globals"
	"MasterGobees/shell"
	"MasterGobees/utils"
	"encoding/json"
	"github.com/TwiN/go-color"
	"io/ioutil"
	"log"
	"net/http"
)

func MainNodeBirth(w http.ResponseWriter, r *http.Request) {
	//Need to store statefully node info
	if r.Method != "POST" {
		utils.SimpleInvalidPath("Ivalid path", w)
		return
	}
	//Let's not store worker node info staefully, entirely in memory
	shell.PrintToShell(color.Colorize(color.Yellow, "[ENDPOINT] Request recived to add new worker node."))
	res_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "[ENDPOINT ERROR] Error parsing request from Worker node addition request."))
	}
	new_worker_node := globals.WorkerNode{
		Ip_addr: res_body_obj["ip_addr"].(string),
		Port:    res_body_obj["port"].(string),
	}
	for _, v := range globals.WorkerNodesMetadata {
		if new_worker_node.Ip_addr == v.Ip_addr && v.Port == new_worker_node.Port {
			shell.PrintToShell(color.Colorize(color.Green, "Node already registered!"))
			utils.SimpleSuccesssStatus("Node already register from previous run", w)
			return
		}
	}
	globals.WorkerNodesMetadata = append(globals.WorkerNodesMetadata, new_worker_node)
	utils.SimpleSuccesssStatus("Successfully Added node to cluster!", w)
	shell.PrintToShell(color.Colorize(color.Green, "Added one Node to cluster : "+new_worker_node.Ip_addr+":"+new_worker_node.Port))
	return
}

func ResetNode(w http.ResponseWriter, r *http.Request) {
	//If already in meta we need to make sure the node "re-joining does not reset itself!"
	var res_body_obj map[string]interface{}
	res_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "[ENDPOINT ERROR] Error parsing request from Worker node reset question."))
	}
	new_worker_node := globals.WorkerNode{
		Ip_addr: res_body_obj["ip_addr"].(string),
		Port:    res_body_obj["port"].(string),
	}

	w.WriteHeader(http.StatusOK)
	for _, v := range globals.WorkerNodesMetadata {
		if new_worker_node.Ip_addr == v.Ip_addr && v.Port == new_worker_node.Port {
			shell.PrintToShell(color.Colorize(color.Green,"node rejoining, hence no node reset."))
			body, _:= json.Marshal("false")
			w.Write(body)
			return
		}
	}
	body, err := json.Marshal(globals.NewCluster)
	if err != nil {
		shell.PrintToShell(color.Colorize(color.Red, "Not able to comm with worker :("))
	}
	w.Write(body)
	return
}

// func KillNode(){
// 	for _,v := range globals.WorkerNodesMetadata{
// 		request_body, err := json.Marshal(true)
// 		if err!=nil{
// 			log.Println(color.Colorize(color.Red,"Error turning off on of the WorkerNode"))
// 		}
// 	}
// }

// [TODO] As for the deliverable, Node failures are not to be handled

// func MainNodeDeath(w http.ResponseWriter, r *http.Request){
//   //Need to remove Node info
// }
//
// func MainNodeAlive(w http.ResponseWriter, r *http.Request){
//   //Need to somehow spawn a function that Kills node if no reponse for 3 rounds
// }
