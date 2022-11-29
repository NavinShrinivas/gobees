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

func ConfigurationMain() error {
	log.Println(color.Colorize(color.Yellow, "Reading Configuration..."))
	config_raw_string, err := os.ReadFile(globals.Config_file_path)
	if err != nil {
		return errors.New("invalid configuration file path")
	}
	is_valid_json := json.Valid(config_raw_string)
	if !is_valid_json {
		return errors.New("malformed json content in config file")
	}
	err = json.Unmarshal(config_raw_string, &globals.Config_obj)
	if err != nil {
		return errors.New("error parsing configuration file :(")
	}
	if globals.Debug_flag {
		fmt.Println("Debug print : ", globals.Config_obj["replication-rate"])
		fmt.Println("Debug print : ", globals.Config_obj)
	}
	return nil
}
