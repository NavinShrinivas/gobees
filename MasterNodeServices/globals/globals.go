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
