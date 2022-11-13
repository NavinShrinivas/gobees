package lifestatus

import (
	"WorkerGobees/globals"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/TwiN/go-color"
)

type NodeInfo struct{
  Ip_addr string `json:"ip_addr"`
  Port string `json:"port"`
}

func NodeBirthRegister(){
  new_node_info :=NodeInfo{
    Ip_addr: globals.Ip,
    Port: globals.Port,
  }
  request_bytes, err :=  json.Marshal(new_node_info)
  if err!=nil{
  	log.Fatal(color.Colorize(color.Red,"Error registering node with master node, please check IP address of master"))
  }
  request_stream := bytes.NewBuffer(request_bytes)
	response, err := http.Post(globals.MasterUrl+"nodebirth", "application/json", request_stream)

res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err!=nil{
		log.Fatal(color.Colorize(color.Red,"Error parsing request from Master node addition request."))
	}
	if res_body_obj["status"] == true{
		log.Println(color.Colorize(color.Green,"Succesfully register with Master node in cluster!"))
	}else{
  	log.Fatal(color.Colorize(color.Red,"Error registering node with master node, please check status of Master node."))
	}
}
