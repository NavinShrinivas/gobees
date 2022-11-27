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

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func help() {
	fmt.Println("Commands : ")
	fmt.Println("SHOW <config/nodes/files/jobs>         : Show config/nodes/files/jobs")
	fmt.Println("PUT <local_file_path> [delimiter]      : Upload file to shared storage")
	fmt.Println("RENAME <old_file_name> <new_file_name> : Rename file in shared storage")
	fmt.Println("DELETE <file_name>                     : Delete file from shared storage")
	fmt.Println("EXIT                                   : Exit the shell")
	fmt.Println("HELP                                   : Show this help")
}

func show_commands(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println(color.Colorize(color.Yellow, "Syntax : SHOW <config/nodes/files/jobs>"))
		return
	}
	command_parse[1] = strings.ToUpper(command_parse[1])
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
			space := 20 - len(v.File_name)
			fmt.Println(w+1, ":", v.File_name, strings.Repeat(" ", space), "| splits : ", v.Splits)
		}
		return
	}
	if command_parse[1] == "JOBS" {
		fmt.Println("Past job history : ")
		return
	}
	fmt.Println("Invalid SHOW command")
}

func handleFiles(command_parse []string) {
	if command_parse[0] == "PUT" {
		log.Println(color.Colorize(color.Yellow, "Uploading file..."))
		//Upload files
		local_file_path := command_parse[1]

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
			log.Println(color.Colorize(color.Red, "Couldn't upload file successfully :("))
			return
			//As for deliberable we do not have to go back and handle clearing part files, but that a good [TODO]
		}
		log.Println(color.Colorize(color.Green, "File upload complete :)"))
		return
	}
}

func renameFile(command_parse []string) {
	if len(command_parse) < 3 {
		log.Println(color.Colorize(color.Yellow, "Syntax : RENAME <old_file_name> <new_file_name>"))
		return
	}
	old_name := command_parse[1]
	new_name := command_parse[2]

	for i, v := range globals.FileMetadata {
		if v.File_name == old_name {
			err := data.RenameFile(old_name, new_name)
			if err != nil {
				log.Println(color.Colorize(color.Red, "Could not rename file!"))
				return
			}
			log.Println(color.Colorize(color.Green, "File renamed successfully on nodes!"))

			globals.FileMetadata = append(globals.FileMetadata[:i], globals.FileMetadata[i+1:]...)
			v.File_name = new_name
			globals.FileMetadata = append(globals.FileMetadata, v)

			log.Println(color.Colorize(color.Green, "File metadata updated successfully!"))
			return
		}
	}
	log.Println(color.Colorize(color.Red, "File does not exist!"))
}

func deleteFile(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "Syntax : DELETE <file_name>"))
		return
	}
	file_name := command_parse[1]

	for i, v := range globals.FileMetadata {
		if v.File_name == file_name {
			err := data.DeleteFile(file_name)
			if err != nil {
				log.Println(color.Colorize(color.Red, "Could not delete file!"))
				return
			}
			log.Println(color.Colorize(color.Green, "File deleted successfully on nodes!"))

			globals.FileMetadata = append(globals.FileMetadata[:i], globals.FileMetadata[i+1:]...)

			log.Println(color.Colorize(color.Green, "File metadata updated successfully!"))
			return
		}
	}
	log.Println(color.Colorize(color.Red, "File does not exist!"))
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
		err = jobs.StartShuffle(custom_partition, partition_file_path)
		if err != nil {
			log.Println(color.Colorize(color.Red, "Partition job failed :("))
			return
		}
		err = jobs.StartReduce(command_vars["reducer"], command_vars["OUT"])
		if err != nil {
			log.Println(color.Colorize(color.Red, "Reduce job failed :("))
			return
		}
	}
}

func commandProcessor(command string) {
	command_parse := strings.Split(strings.Trim(strings.Trim(command, "\n"), "\r"), " ")
	command_parse[0] = strings.ToUpper(command_parse[0])
	if command_parse[0] == "CLEAR" {
		clearScreen()
		return
	}
	if command_parse[0] == "HELP" {
		help()
		return
	}
	if command_parse[0] == "SHOW" {
		show_commands(command_parse)
		return
	}
	if command_parse[0] == "PUT" {
		handleFiles(command_parse)
		return
	}
	if command_parse[0] == "RENAME" {
		renameFile(command_parse)
		return
	}
	if command_parse[0] == "DELETE" {
		deleteFile(command_parse)
		return
	}
	if command_parse[0] == "EXIT" {
		utils.ExitSequence()
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
