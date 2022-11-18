package globals

import (
	"sync"
)

var Config_obj map[string]interface{}
var Config_file_path string
var Debug_flag bool
var MainWg *sync.WaitGroup
var ServerPort string
var WorkerNodesMetadata []WorkerNode
var FileMetadata []File
var NewCluster bool

type WorkerNode struct {
	Ip_addr string `json:"ip_addr"`
	Port    string `json:"port"`
}

type File struct {
	File_name string
	Splits    int32
	Nodes     []WorkerNode
}
