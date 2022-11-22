package globals

// var Config_obj map[string]interface{}
var MasterUrl string
var Port string
var Ip string

type WorkerNode struct {
	Ip_addr string `json:"ip_addr"`
	Port    string `json:"port"`
}
var ShuffleNodeMetadata []WorkerNode
