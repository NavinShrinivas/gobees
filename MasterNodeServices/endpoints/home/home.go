package home

import (
	"MasterGobees/shell"
	"MasterGobees/utils"
	"github.com/TwiN/go-color"
	"net/http"
)

func MainHome(w http.ResponseWriter, r *http.Request) {
	shell.PrintToShell(color.Colorize(color.Yellow, "[ENDPOINT] Access on home endpoint."))
	utils.SimpleSuccesssStatus("This is the MasterNode home end point. Master node is alive.", w)
	return
}

func Health(w http.ResponseWriter, r *http.Request) {
	utils.SimpleSuccesssStatus("Master node is alive.", w)
	return
}
