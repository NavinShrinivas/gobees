package home

import (
	"MasterGobees/utils"
	"log"
	"net/http"

	"github.com/TwiN/go-color"
)


func MainHome(w http.ResponseWriter,r *http.Request){
  log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Access on home endpoint."))
  utils.SimpleSuccesssStatus("This is the MasterNode home end point. Master node is alive.", w)
  return
}
