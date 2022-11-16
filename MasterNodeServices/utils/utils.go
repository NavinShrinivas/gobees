package utils

import (
	"MasterGobees/globals"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/TwiN/go-color"
)

type SimpleResponse struct{
  Status bool `json:"status"`
  Message string `json:"message"`
}

func SimpleFailStatus(res string, w http.ResponseWriter){
  response := SimpleResponse{
    Status : false,
    Message : res,
  }
  response_bytes, err := json.Marshal(response)
  if err!=nil{
    log.Println(color.Colorize(color.Red,"Error Marshalling response on a fail status"))
    return
  }
  w.WriteHeader(http.StatusForbidden)
  w.Write(response_bytes)
  return
}

func SimpleSuccesssStatus(res string, w http.ResponseWriter){
  response := SimpleResponse{
    Status : true,
    Message : res,
  }
  response_bytes, err := json.Marshal(response)
  if err!=nil{
    log.Println(color.Colorize(color.Red,"Error Marshalling response on a fail status"))
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(response_bytes)
  return
}

func SimpleInvalidPath(res string, w http.ResponseWriter){
	response := SimpleResponse{
		Status: false,
		Message: res,
	}
  response_bytes, err := json.Marshal(response)
  if err!=nil{
    log.Println(color.Colorize(color.Red,"Error Marshalling response on a fail status"))
    return
  }
  w.WriteHeader(http.StatusNotFound)
  w.Write(response_bytes)
  return
}

func ExistSequence(){
	fd1, _ := os.Create("./NodeMeta.json")
	byte_buffer1, _:= json.Marshal(globals.WorkerNodesMetadata)
	fd1.Write(byte_buffer1)
	fd1.Close()

	fd2, _ := os.Create("./FileMeta.json")
	byte_buffer2, _:= json.Marshal(globals.WorkerNodesMetadata)
	fd2.Write(byte_buffer2)
	fd2.Close()

	log.Println(color.Colorize(color.Green,"Storing meta for future use. Adios."))
	os.Exit(0)
}
