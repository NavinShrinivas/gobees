package utils

import (
	"encoding/json"
	"log"
	"net/http"

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
