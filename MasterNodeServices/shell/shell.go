package shell

import (
	"MasterGobees/utils"
	"MasterGobees/endpoints/data"
	"MasterGobees/globals"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
)

var reader *bufio.Reader 

func show_commands(command_parse []string){
	if command_parse[1] == "CONFIG"{
		fmt.Println(globals.Config_obj)
		return
	}
	if command_parse[1] == "NODES"{
		for w,v := range globals.WorkerNodesMetadata{
			if v.Ip_addr != ""{
				fmt.Println("Node",w+1,":",v.Ip_addr, v.Port)
			}
		}
		return
	}
	if command_parse[1] == "FILES"{
		for w,v := range globals.FileMetadata{
			fmt.Println(w,":",v.File_name,"| splits : ",v.Splits)
		}
		return
	}
	if command_parse[1] == "JOBS"{
		fmt.Println("Past job history : ")
		return
	}
	fmt.Println("Invalid command")
	return
}

func handleFiles(command_parse []string){
	if command_parse[0] == "PUT"{
		log.Println(color.Colorize(color.Yellow, "Uploading file..."))
		//Upload files 
		local_file_path := strings.Trim(command_parse[1], "\n")
		//Check if file exists
		_, err := os.Stat(local_file_path)
		
		if err!=nil{
			log.Println(color.Colorize(color.Red, "File does not exist!"))
			return
		}

		//Check if file already in SS [Shared Storage]
		for _,v :=  range globals.FileMetadata{
			if v.File_name == filepath.Base(local_file_path){
				log.Println(color.Colorize(color.Red, "File with same name exists!"))
				return
			}
		}
		//Split files by delimiter, default is "\n"
		delimiter := "\n"
		if len(command_parse) > 2{
			delimiter = command_parse[2]
		}
		err = data.SplitAndUploadFile(local_file_path, delimiter)
		if err!=nil{
	    log.Println(color.Colorize(color.Red, "Couldln't upload file successfully :,("))
	    return
	    //As for deliberable we do not have to go back and handle clearing part files, but that a good [TODO]
		}
		log.Println(color.Colorize(color.Green, "File upload complete :)"))
		return
	}
}


func commandParsing(command_parse []string)( map[string]string,error ){
	var command_map map[string]string
	command_map = make(map[string]string)
	for _,v := range command_parse{
		if len(v) == 0{
			continue
		}
		if string(v[0]) == "-" && string(v[1]) != "-"{
			key_val_string := v[1:]
			array := strings.Split(key_val_string,"=")
			command_map[array[0]] = array[1]
		} else if string(v[0])+string(v[1]) == "--"{
			key_val_string := v[2:]
			array := strings.Split(key_val_string,"=")
			command_map[array[0]] = array[1]
		}
	}
	return command_map,nil
}

func newHandler(command_parse []string){
	if command_parse[1] == "JOB"{
		//Map reduce job incoming
		//Time to write command parsing tool
		command_vars, _ := commandParsing(command_parse)
		//Need to check if everything we need is here and if it is valid
		//Mapper //Reducer //IN file //OUT file
		input_file_ss_name, ok:= command_vars["IN"]
		if !ok{
			log.Println(color.Colorize(color.Red,"Not all needed parameters for Map Reduce jobs given"))
			return
		}
		file_in_SS := false;
		for _,v := range globals.FileMetadata{
			if v.File_name == input_file_ss_name{
				file_in_SS = true
				break
			}
		}
		if !file_in_SS{
			//If we have reached last but not found file in SS record
			log.Println(color.Colorize(color.Red,"Given file not found in Share Storage (SS)"))
			return
		}
	}
}

func commandParser(command string){
	command_parse := strings.Split(strings.Trim(command,"\n"), " ")
	if command_parse[0] == "SHOW" && len(command_parse) >= 2{
		show_commands(command_parse)
		return
	}
	if command_parse[0] == "PUT"{
		handleFiles(command_parse);
		return
	}
	if command_parse[0] == "EXIT"{
		utils.ExistSequence()
		//There is no coming back from exit sequence
	}
	if command_parse[0] == "NEW"{
		newHandler(command_parse)
		return
	}
	fmt.Println("Invalid command")
}

func PrintToShell(str string){
  fmt.Println()
  fmt.Println(str)
  fmt.Print("master>")
}

func CommandListner(){
	for{
		command := ""
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		commandParser(command)
		fmt.Print("master> ")
	}
}

func Initialize(){
	reader = bufio.NewReader(os.Stdin)
  log.Println("Starting shell...")
	fmt.Print("master> ")
  go CommandListner()
}
