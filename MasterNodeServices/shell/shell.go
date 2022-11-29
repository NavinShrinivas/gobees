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

func help() {
	fmt.Println("\nGENERAL COMMANDS")
	fmt.Println("HELP                                   : Show this help")
	fmt.Println("CAT (file)                             : Print file contents")
	fmt.Println("LS (directory)                         : List files in directory")
	fmt.Println("PWD                                    : Print working directory")
	fmt.Println("CD <directory>                         : Change directory")
	fmt.Println("CLEAR                                  : Clear screen")
	fmt.Println("EXIT                                   : Exit the shell")

	fmt.Println("\nWORKER NODE COMMANDS")
	fmt.Println("SHOW <config/nodes/files>              : Show config/nodes/files")
	fmt.Println("PUT <local_file_path>                  : Upload file to shared storage")
	fmt.Println("FETCH <file_name> <local_file_path>    : Fetch file from shared storage to local storage")
	fmt.Println("RENAME <old_file_name> <new_file_name> : Rename file in shared storage")
	fmt.Println("DELETE <file_name>                     : Delete file from shared storage")

	fmt.Println("\nMAPPER REDUCE COMMAND")
	fmt.Println("JOB -mapper=<file_path> -reducer=<file_path> -IN=<file_name> -OUT=<file_name>")
	fmt.Println("(-PARTITION=<file_path>)")
}

func cat(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println("Syntax : CAT <file>")
		return
	}
	file_name := command_parse[1]
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("ERROR: File not found")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func listFiles(command_parse []string) {
	path := "./"
	if len(command_parse) > 1 {
		path = command_parse[1]
	}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Could not read directory!"))
		return
	}
	fmt.Println("\nFiles in", path, ":")
	for _, f := range files {
		if f.IsDir() {
			fmt.Println(color.Colorize(color.Blue, f.Name()))
		} else {
			fmt.Println(f.Name())
		}
	}
}

func printWorkingDirectory() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(color.Colorize(color.Red, "Could not get working directory!"))
		return
	}
	fmt.Println(path)
}

func changeDirectory(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "Syntax : CD <directory>"))
		return
	}
	path := command_parse[1]
	err := os.Chdir(path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Could not change directory!"))
		return
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func show_commands(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println(color.Colorize(color.Yellow, "Syntax : SHOW <config/nodes/files>"))
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
	// if command_parse[1] == "JOBS" {
	// 	fmt.Println("Past job history : ")
	// 	for w, v := range globals.JobHistory {
	// 		fmt.Println(w+1, ":", v)
	// 	}
	// 	return
	// }
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
				log.Println(color.Colorize(color.Red, "File with same name exists! Please rename the file and try again."))
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
	var command_map map[string]string = make(map[string]string)
	var re = regexp.MustCompile(`[\-](\w*)=([^\s]+[^\-]*\b)`)
	for _, match := range re.FindAllString(command, -1) {
		sub_res := re.FindStringSubmatch(match)
		command_map[sub_res[1]] = sub_res[2] //"mapper" will actually contain the whole mapper command and not just mapper file name
	}
	return command_map, nil
}

func newHandler(command_parse []string, command string) {
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
	out_file_name, ok := command_vars["OUT"]
	if !ok {
		out_file_name = input_file_ss_name + "_job_out"
	}
	err = jobs.StartReduce(command_vars["reducer"], out_file_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Reduce job failed :("))
		return
	}

	//Below part is already done in reducer function
	// //We need to update metadata about the new reducer output
	// new_output_file_meta := globals.File{
	// 	File_name: command_vars["IN"],
	// }
	// for _,v := range globals.FileMetadata{
	// 	if v.File_name == command_vars["IN"]{
	// 		new_output_file_meta.Nodes = v.Nodes
	// 		new_output_file_meta.Splits = int32(len(v.Nodes))
	// 	}
	// }
}

func commandProcessor(command string) {
	command_parse := strings.Split(strings.Trim(strings.Trim(command, "\n"), "\r"), " ")
	command_parse[0] = strings.ToUpper(command_parse[0])
	if command_parse[0] == "HELP" {
		help()
		return
	}
	if command_parse[0] == "CAT" {
		cat(command_parse)
		return
	}
	if command_parse[0] == "LS" || command_parse[0] == "LIST" {
		listFiles(command_parse)
		return
	}
	if command_parse[0] == "PWD" {
		printWorkingDirectory()
		return
	}
	if command_parse[0] == "CD" {
		changeDirectory(command_parse)
		return
	}
	if command_parse[0] == "CLEAR" || command_parse[0] == "CLS" || command_parse[0] == "\x0C" {
		clearScreen()
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
	if command_parse[0] == "FETCH" {
		fetchHandler(command_parse)
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
	if command_parse[0] == "JOB" {
		newHandler(command_parse, command)
		return
	}
	if command_parse[0] == "EXIT" {
		utils.ExitSequence()
		//There is no coming back from exit sequence
	}
	fmt.Println("Invalid command. Type HELP for list of commands")
}

func fetchHandler(command_parse []string) {
	if len(command_parse) < 3 {
		log.Println(color.Colorize(color.Red, "SYNTAX ERROR: FETCH <FILE_NAME_IN_SS> <FILE_NAME_IN_LOCAL>"))
		return
	}
	ss_file_name := command_parse[1]
	local_to_file_name := command_parse[2]
	var all_nodes_with_split []globals.WorkerNode
	for _, v := range globals.FileMetadata {
		if v.File_name == ss_file_name {
			all_nodes_with_split = v.Nodes
		}
	}
	if len(all_nodes_with_split) == 0 {
		log.Println(color.Colorize(color.Red, "File with name does not exists in Share Storage"))
		return
	}
	//Need to make request to each node
	//[WARNING] We SHOULD NOT PARELLELISE THIS PART, as we cant have multiple thread write to a single file
	for _, v := range all_nodes_with_split {
		err := data.FetchAndMergeFile(v, local_to_file_name, ss_file_name)
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error fetching files :("))
			os.Remove(local_to_file_name)
			return
		}
	}
	log.Println(color.Colorize(color.Green, "Finished fetching file"))
}

func PrintToShell(str string) {
	fmt.Println()
	fmt.Println(str)
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wd = "X"
	}
	fmt.Print("[" + wd + "] master> ")
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
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			wd = "X"
		}
		fmt.Print("\n[" + wd + "] master> ")
	}
}

func Initialize() {
	reader = bufio.NewReader(os.Stdin)
	log.Println("Starting shell...")
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wd = "X"
	}
	fmt.Print("[" + wd + "] master> ")
	go CommandListner()
}
