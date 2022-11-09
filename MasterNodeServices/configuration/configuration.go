package configuration

import (
	"MasterGobees/globals"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/TwiN/go-color"
)

func ConfigurationMain() error{
  //Read file from path : 
  log.Println(color.Colorize(color.Yellow,"Reading Configuration..."))
  config_raw_string, err := os.ReadFile(globals.Config_file_path) //return back a byte array of file content
  if err!=nil{
    return errors.New("Invalid configuration file path.")
  }
  is_valid_json :=  json.Valid(config_raw_string)
  if !is_valid_json{
    return errors.New("Malformed json content in Config file!")
  }
  err = json.Unmarshal(config_raw_string, &globals.Config_obj)
  if err!=nil{
    return errors.New("Error parsing configuration file :(.")
  }
  if globals.Debug_flag{
    fmt.Println("debug print : ",globals.Config_obj["replication-rate"])
    fmt.Println("debug print : ",globals.Config_obj)
  }
  
  return nil
}
