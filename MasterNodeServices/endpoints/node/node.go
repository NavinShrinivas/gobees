package node

import (
	"MasterGobees/utils"
	"net/http"
)

func MainNodeBirth(w http.ResponseWriter, r *http.Request){
  //Need to store statefully node info
  if r.Method != "POST"{
    utils.("Ivalid path", w)
    return
  }
}

func MainNodeDeath(w http.ResponseWriter, r *http.Request){
  //Need to remove Node info
}

func MainNodeAlive(w http.ResponseWriter, r *http.Request){
  //Need to somehow spawn a function that Kills node if no reponse for 3 rounds
}

