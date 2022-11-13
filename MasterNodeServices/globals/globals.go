package globals

import (
  "sync"
)

var Config_obj map[string]interface{}
var Config_file_path string
var Debug_flag bool
var MainWg *sync.WaitGroup
var ServerPort string
var Shell_busy_counter uint64
var WorkerNodesMetadata[]WorkerNode


type WorkerNode struct{
  Ip_addr string `json:"ip_addr"`
  Port string `json:"port"`
  Files []string `json:"files"`//We store each file split it has in this array of the form  : filename_PARTn
}
