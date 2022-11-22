package shell

import (
	"MasterGobees/endpoints/data"
	"MasterGobees/endpoints/jobs"
	"MasterGobees/globals"
	"MasterGobees/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TwiN/go-color"
)

var reader *bufio.Reader

func show_commands(command_parse []string) {
	if command_parse[1] == "CONFIG" {
		fmt.Println(globals.Config_obj)
		return
	}
	if command_parse[1] == "NODES" {
		for w, v := range globals.WorkerNodesMetadata {
			if v.Ip_addr != "" {
				fmt.Println("Node", w+1, ":", v.Ip_addr, v.Port)
			}
		}
		return
	}
	if command_parse[1] == "FILES" {
		for w, v := range globals.FileMetadata {
			fmt.Println(w, ":", v.File_name, "| splits : ", v.Splits)
		}
		return
	}
	if command_parse[1] == "JOBS" {
		fmt.Println("Past job history : ")
		return
	}
	fmt.Println("Invalid command")
	return
}

func handleFiles(command_parse []string) {
	if command_parse[0] == "PUT" {
		log.Println(color.Colorize(color.Yellow, "Uploading file..."))
		//Upload files
		local_file_path := strings.Trim(command_parse[1], "\n")
		//Check if file exists
		_, err := os.Stat(local_file_path)

		if err != nil {
			log.Println(color.Colorize(color.Red, "File does not exist!"))
			return
		}

		//Check if file already in SS [Shared Storage]
		for _, v := range globals.FileMetadata {
			if v.File_name == filepath.Base(local_file_path) {
				log.Println(color.Colorize(color.Red, "File with same name exists!"))
				return
			}
		}
		//Split files by delimiter, default is "\n"
		delimiter := "\n"
		if len(command_parse) > 2 {
			delimiter = command_parse[2]
		}
		err = data.SplitAndUploadFile(local_file_path, delimiter)
		if err != nil {
			log.Println(color.Colorize(color.Red, "Couldln't upload file successfully :,("))
			return
			//As for deliberable we do not have to go back and handle clearing part files, but that a good [TODO]
		}
		log.Println(color.Colorize(color.Green, "File upload complete :)"))
		return
	}
}

func NewcommandParsing(command string) (map[string]string, error) {
	var command_map map[string]string
	command_map = make(map[string]string)
	var re = regexp.MustCompile(`[\-](\w*)=([^\s]+[^\-]*\b)`)
	for _, match := range re.FindAllString(command, -1) {
		sub_res := re.FindStringSubmatch(match)
		command_map[sub_res[1]] = sub_res[2] //"mapper" will actually contain the whole mapper command and not just mapper file name
	}
	return command_map, nil
}

func newHandler(command_parse []string, command string) {
	if command_parse[1] == "JOB" {
		//Map reduce job incoming
		//Time to write command parsing tool
		command_vars, _ := NewcommandParsing(command)
		//Need to check if everything we need is here and if it is valid
		//Mapper //Reducer //IN file //OUT file
		input_file_ss_name, ok := command_vars["IN"]
		if !ok {
			log.Println(color.Colorize(color.Red, "Not all needed parameters for Map Reduce jobs given"))
			return
		}
		file_in_SS := false
		for _, v := range globals.FileMetadata {
			if v.File_name == input_file_ss_name {
				file_in_SS = true
				break
			}
		}
		fmt.Println(input_file_ss_name)
		if !file_in_SS {
			//If we have reached last but not found file in SS record
			log.Println(color.Colorize(color.Red, "Given file not found in Share Storage (SS)"))
			return
		}
		//Check if mapper and reducer file exists
		_, err := os.Stat(strings.Split(command_vars["mapper"], " ")[0])
		if err != nil {
			log.Println(color.Colorize(color.Red, "Mapper file provided does not exists."))
			return
		}
		_, err = os.Stat(strings.Split(command_vars["reducer"], " ")[0])
		if err != nil {
			log.Println(color.Colorize(color.Red, "Reducer file provided does not exists."))
			return
		}
		err = jobs.SendMapJobs(command_vars["mapper"], command_vars["IN"])
		if err != nil {
			log.Println(color.Colorize(color.Red, "Map job failed"))
			return
		}
		partition_file_path, ok := command_vars["PARTITION"]
		custom_partition := true
		if !ok {
			custom_partition = false
		}
		err = jobs.StartShuffle(custom_partition,partition_file_path)
		if err!=nil{
			log.Println(color.Colorize(color.Red,"Partition job failed :("))
			return
		}
		// err = jobs.StartReduce()
	}
}

func commandProcessor(command string) {
	command_parse := strings.Split(strings.Trim(command, "\n"), " ")
	if command_parse[0] == "SHOW" && len(command_parse) >= 2 {
		show_commands(command_parse)
		return
	}
	if command_parse[0] == "PUT" {
		handleFiles(command_parse)
		return
	}
	if command_parse[0] == "EXIT" {
		utils.ExistSequence()
		//There is no coming back from exit sequence
	}
	if command_parse[0] == "NEW" {
		newHandler(command_parse, command)
		return
	}
	fmt.Println("Invalid command")
}

func PrintToShell(str string) {
	fmt.Println()
	fmt.Println(str)
	fmt.Print("master>")
}

func CommandListner() {
	for {
		command := ""
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		commandProcessor(command)
		fmt.Print("master> ")
	}
}

func Initialize() {
	reader = bufio.NewReader(os.Stdin)
	log.Println("Starting shell...")
	fmt.Print("master> ")
	go CommandListner()
}
