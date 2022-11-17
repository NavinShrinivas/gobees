package home

import (
	"MasterGobees/shell"
	"MasterGobees/utils"
	"net/http"
	"github.com/TwiN/go-color"
)


func MainHome(w http.ResponseWriter,r *http.Request){
	shell.PrintToShell(color.Colorize(color.Yellow, "[ENDPOINT] Access on home endpoint."))
  utils.SimpleSuccesssStatus("This is the MasterNode home end point. Master node is alive.", w)
  shell.PrintToShell(r.RemoteAddr)
  return
}
