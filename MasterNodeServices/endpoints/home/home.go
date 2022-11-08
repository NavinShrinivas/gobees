package home

import ( 
  "net/http" 
  "MasterGobees/utils"
)


func MainHome(w http.ResponseWriter,r *http.Request){
  utils.SimpleSuccesssStatus("This is the MasterNode home end point. Master node is alive.", w)
  return
}
