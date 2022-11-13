package home

import (
	"WorkerGobees/utils"
	"log"
	"net/http"

	"github.com/TwiN/go-color"
)


func MainHome(w http.ResponseWriter,r *http.Request){
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Access on home endpoint."))
  utils.SimpleSuccesssStatus("This is one of the Worker node home end point. Worker node is alive.", w)
  return
}
