package shell

import (
	"MasterGobees/endpoints/data"
	"MasterGobees/endpoints/jobs"
	"MasterGobees/globals"
	"MasterGobees/utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TwiN/go-color"
)

var reader *bufio.Reader

func help() {
	fmt.Println()
	fmt.Println(color.Colorize(color.Underline, "GENERAL COMMANDS"))
	fmt.Println("HELP                                   : Show this help")
	fmt.Println("CAT (file)                             : Print file contents")
	fmt.Println("RM <file>                              : Delete file")
	fmt.Println("PWD                                    : Print working directory")
	fmt.Println("LS (directory)                         : List files in directory")
	fmt.Println("CD <directory>                         : Change directory")
	fmt.Println("MKDIR <directory>                      : Create directory")
	fmt.Println("RMDIR <directory>                      : Delete directory")
	fmt.Println("CLEAR                                  : Clear screen")
	fmt.Println("EXIT                                   : Exit the shell")

	fmt.Println()
	fmt.Println(color.Colorize(color.Underline, "WORKER NODE COMMANDS"))
	fmt.Println("SHOW <config/nodes/files>              : Show config/nodes/files")
	fmt.Println("PUT <local_file_path>                  : Upload file to shared storage")
	fmt.Println("FETCH <file_name> <local_file_path>    : Fetch file from shared storage to local storage")
	fmt.Println("RENAME <old_file_name> <new_file_name> : Rename file in shared storage")
	fmt.Println("DELETE <file_name>                     : Delete file from shared storage")

	fmt.Println()
	fmt.Println(color.Colorize(color.Underline, "MAPPER REDUCE COMMAND"))
	fmt.Println("JOB -mapper=<file_path> -reducer=<file_path> -IN=<file_name> -OUT=<file_name> (-PARTITION=<file_path>)")
	fmt.Println("MAPPER    : Local Path to mapper file")
	fmt.Println("REDUCER   : Local Path to reducer file")
	fmt.Println("IN        : Name of input file in shared storage")
	fmt.Println("OUT       : Name of output file in shared storage")
	fmt.Println("PARTITION : Local Path to partition file")
}

func cat(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println(color.Colorize(color.Yellow, "SYNTAX: cat <file>"))
		return
	}
	file_name := command_parse[1]
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println(color.Colorize(color.Red, "ERROR: File not found"))
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func rm(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println(color.Colorize(color.Yellow, "SYNTAX: rm <file>"))
		return
	}
	file_name := command_parse[1]
	fmt.Println(color.Colorize(color.Yellow, "Are you sure you want to delete file "+file_name+"? (y/n)"))
	var response string = ""
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		err := os.Remove(file_name)
		if err != nil {
			fmt.Println(color.Colorize(color.Red, "ERROR: Could not delete file"))
			return
		}
		fmt.Println(color.Colorize(color.Green, "File deleted successfully"))
	} else {
		fmt.Println(color.Colorize(color.Green, "File deletion aborted. File not deleted"))
	}
}

func printWorkingDirectory() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Could not get working directory"))
		return
	}
	fmt.Println(path)
}

func listFiles(command_parse []string) {
	path := "./"
	if len(command_parse) > 1 {
		path = command_parse[1]
	}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Could not list files in directory"))
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

func changeDirectory(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "SYNTAX: cd <directory>"))
		return
	}
	path := command_parse[1]
	err := os.Chdir(path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Could not change directory"))
		return
	}
}

func makeDirectory(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "SYNTAX: mkdir <directory>"))
		return
	}
	path := command_parse[1]
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Could not create directory"))
		return
	}
	fmt.Println(color.Colorize(color.Green, "Directory created successfully"))
}

func removeDirectory(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "SYNTAX: rmdir <directory>"))
		return
	}
	path := command_parse[1]
	fmt.Println(color.Colorize(color.Yellow, "Are you sure you want to delete directory "+path+"? (y/n)"))
	var response string = ""
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println(color.Colorize(color.Red, "ERROR: Could not delete directory"))
			return
		}
		fmt.Println(color.Colorize(color.Green, "Directory deleted successfully"))
	} else {
		fmt.Println(color.Colorize(color.Green, "Directory deletion aborted. Directory not deleted"))
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func showCommands(command_parse []string) {
	if len(command_parse) < 2 {
		fmt.Println(color.Colorize(color.Yellow, "SYNTAX : SHOW <config/nodes/files>"))
		return
	}
	command_parse[1] = strings.ToUpper(command_parse[1])
	if command_parse[1] == "CONFIG" {
		fmt.Println(color.Colorize(color.Underline, "CONFIGURATIONS"))
		fmt.Println(globals.Config_obj)
		return
	}
	if command_parse[1] == "NODES" {
		if len(globals.WorkerNodesMetadata) == 0 {
			fmt.Println(color.Colorize(color.Yellow, "No nodes found registered with master"))
			return
		}
		fmt.Println("Nodes registered with master node :")
		for w, v := range globals.WorkerNodesMetadata {
			if v.Ip_addr != "" {
				fmt.Println("Node", w+1, ":", v.Ip_addr, v.Port)
			}
		}
		return
	}
	if command_parse[1] == "FILES" {
		if len(globals.FileMetadata) == 0 {
			fmt.Println(color.Colorize(color.Yellow, "No files uploaded to distributed file system"))
			return
		}
		fmt.Println("Files in distributed file system :")
		for w, v := range globals.FileMetadata {
			space := 20 - len(v.File_name)
			if space < 0 {
				space = 0
			}
			fmt.Println(w+1, ":", v.File_name, strings.Repeat(" ", space), "| splits : ", v.Splits)
		}
		return
	}
	fmt.Println(color.Colorize(color.Yellow, "SYNTAX : SHOW <config/nodes/files>"))
}

func putFileSS(command_parse []string) {
	log.Println(color.Colorize(color.Yellow, "Starting file upload to shared storage ..."))
	//Upload files
	local_file_path := command_parse[1]

	//Check if file exists
	_, err := os.Stat(local_file_path)

	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: File does not exist"))
		return
	}

	//Check if file already in SS [Shared Storage]
	for _, v := range globals.FileMetadata {
		if v.File_name == filepath.Base(local_file_path) {
			log.Println(color.Colorize(color.Red, "ERROR: File with same name exists in shared storage! Please rename the file and try again."))
			return
		}
	}

	// check if file name is valid with regex
	regex := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !regex.MatchString(filepath.Base(local_file_path)) {
		log.Println(color.Colorize(color.Red, "ERROR: File name is not valid. Please rename the file and try again."))
		return
	}

	//Split files by delimiter, default is "\n"
	delimiter := "\n"
	if len(command_parse) > 2 {
		delimiter = command_parse[2]
	}
	err = data.SplitAndUploadFile(local_file_path, delimiter)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Could not upload file to shared storage"))
		return
		//As for deliberable we do not have to go back and handle clearing part files, but that a good [TODO]
	}
	log.Println(color.Colorize(color.Green, "File uploaded successfully to shared storage"))
}

func fetchFileSS(command_parse []string) {
	if len(command_parse) < 3 {
		log.Println(color.Colorize(color.Red, "SYNTAX : FETCH <FILE_NAME_IN_SS> <FILE_NAME_IN_LOCAL>"))
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
		log.Println(color.Colorize(color.Red, "ERROR: File does not exist in shared storage"))
		return
	}
	//Need to make request to each node
	//[WARNING] We SHOULD NOT PARELLELISE THIS PART, as we cant have multiple thread write to a single file
	for _, v := range all_nodes_with_split {
		err := data.FetchAndMergeFile(v, local_to_file_name, ss_file_name)
		if err != nil {
			log.Println(color.Colorize(color.Red, "ERROR: Could not fetch file from shared storage"))
			os.Remove(local_to_file_name)
			return
		}
	}
	log.Println(color.Colorize(color.Green, "File fetched successfully from shared storage"))
}

func renameFileSS(command_parse []string) {
	if len(command_parse) < 3 {
		log.Println(color.Colorize(color.Yellow, "SYNTAX : RENAME <old_file_name> <new_file_name>"))
		return
	}
	old_name := command_parse[1]
	new_name := command_parse[2]

	regex := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !regex.MatchString(new_name) {
		log.Println(color.Colorize(color.Red, "ERROR: File name is not valid. Please rename the file and try again."))
		return
	}

	if strings.Contains(new_name, "/") {
		log.Println(color.Colorize(color.Red, "ERROR: File name is not valid. Please rename the file and try again."))
		return
	}

	for i, v := range globals.FileMetadata {
		if v.File_name == old_name {
			err := data.RenameFile(old_name, new_name)
			if err != nil {
				log.Println(color.Colorize(color.Red, "	ERROR: Could not rename file"))
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
	log.Println(color.Colorize(color.Red, "ERROR: File does not exist in shared storage"))
}

func deleteFileSS(command_parse []string) {
	if len(command_parse) < 2 {
		log.Println(color.Colorize(color.Yellow, "SYNTAX : DELETE <file_name>"))
		return
	}
	file_name := command_parse[1]

	for i, v := range globals.FileMetadata {
		if v.File_name == file_name {
			err := data.DeleteFile(file_name)
			if err != nil {
				log.Println(color.Colorize(color.Red, "ERROR: Could not delete file"))
				return
			}
			log.Println(color.Colorize(color.Green, "File deleted successfully on nodes!"))

			globals.FileMetadata = append(globals.FileMetadata[:i], globals.FileMetadata[i+1:]...)

			log.Println(color.Colorize(color.Green, "File metadata updated successfully!"))
			return
		}
	}
	log.Println(color.Colorize(color.Red, "ERROR: File does not exist in shared storage"))
}

func jobcommandParsing(command string) (map[string]string, error) {
	command = strings.Replace(command, "\r", "", -1)
	command = strings.Replace(command, "\n", "", -1)
	var command_map map[string]string = make(map[string]string)
	var re = regexp.MustCompile(`[\-](\w*)=([^\s]+[^\-]*\b)`)
	for _, match := range re.FindAllString(command, -1) {
		sub_res := re.FindStringSubmatch(match)
		command_map[sub_res[1]] = sub_res[2] //"mapper" will actually contain the whole mapper command and not just mapper file name
	}
	return command_map, nil
}

func jobHandlerSS(command_parse []string, command string) {
	//Map reduce job incoming
	//Time to write command parsing tool
	command_vars, _ := jobcommandParsing(command)
	//Need to check if everything we need is here and if it is valid
	//Mapper //Reducer //IN file //OUT file
	input_file_ss_name, ok := command_vars["IN"]
	if !ok {
		log.Println(color.Colorize(color.Yellow, "SYNTAX ERROR : No input file specified!"))
		return
	}

	output_file_ss_name, ok := command_vars["OUT"]
	if !ok {
		log.Println(color.Colorize(color.Yellow, "SYNTAX ERROR : No output file specified!"))
		return
	}

	mapper_file_name, ok := command_vars["mapper"]
	if !ok {
		log.Println(command_vars)
		log.Println(color.Colorize(color.Yellow, "SYNTAX ERROR : No mapper file specified!"))
		return
	}

	reducer_file_name, ok := command_vars["reducer"]
	if !ok {
		log.Println(color.Colorize(color.Yellow, "SYNTAX ERROR : No reducer file specified!"))
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
		log.Println(color.Colorize(color.Red, "ERROR: Input file does not exist in shared storage"))
		return
	}

	//Check if mapper and reducer file exists
	_, err := os.Stat(mapper_file_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Mapper file does not exist"))
		return
	}
	_, err = os.Stat(reducer_file_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Reducer file does not exist"))
		return
	}
	err = jobs.SendMapJobs(mapper_file_name, input_file_ss_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Map jobs failed"))
		return
	}
	partition_file_path, ok := command_vars["partition"]
	custom_partition := true
	if !ok {
		custom_partition = false
	}
	err = jobs.StartShuffle(custom_partition, partition_file_path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Partitioning/Shuffling failed"))
		return
	}

	err = jobs.StartReduce(reducer_file_name, output_file_ss_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "ERROR: Reduce jobs failed"))
		return
	}
	log.Println(color.Colorize(color.Green, "Map Reduce Job completed! :)"))
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
	command_parse := strings.Split(command, " ")
	command_parse[0] = strings.ToUpper(command_parse[0])
	switch command_parse[0] {
	case "HELP":
		help()
	case "CAT":
		cat(command_parse)
	case "RM":
		rm(command_parse)
	case "PWD":
		printWorkingDirectory()
	case "LS", "LIST":
		listFiles(command_parse)
	case "CD":
		changeDirectory(command_parse)
	case "MKDIR":
		makeDirectory(command_parse)
	case "RMDIR":
		removeDirectory(command_parse)
	case "CLEAR", "CLS", "CLR":
		clearScreen()
	case "SHOW":
		showCommands(command_parse)
	case "PUT":
		putFileSS(command_parse)
	case "FETCH":
		fetchFileSS(command_parse)
	case "RENAME":
		renameFileSS(command_parse)
	case "DELETE":
		deleteFileSS(command_parse)
	case "JOB":
		jobHandlerSS(command_parse, command)
	case "EXIT":
		utils.ExitSequence()
	case "VI", "VIM":
		fmt.Println(color.Colorize(color.Yellow, "VI/VIM is not supported (yet *wink wink*)"))
	default:
		fmt.Println(color.Colorize(color.Red, "Invalid command. Type HELP for list of commands"))
	}
}

func PrintToShell(str string) {
	fmt.Println()
	fmt.Println(str)
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wd = "X"
	}
	fmt.Println()
	fmt.Print(color.Colorize(color.Blue, "["+wd+"] master> "))
}

func CommandListner() {
	for {
		command := ""
		command, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				utils.ExitSequence()
			}
			log.Println(color.Colorize(color.Red, "ERROR: Command not read properly"))
			return
		}
		command = strings.Join(strings.Fields(command), " ")
		command = strings.Replace(command, "\r", "", -1)
		command = strings.Replace(command, "\t", "", -1)
		command = strings.Replace(command, "\v", "", -1)
		command = strings.Replace(command, "\n", "", -1)
		commandProcessor(command)
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			wd = "X"
		}
		fmt.Println()
		fmt.Print(color.Colorize(color.Blue, "["+wd+"] master> "))
	}
}

func Initialize() {
	reader = bufio.NewReader(os.Stdin)
	log.Println(color.Colorize(color.Green, "Starting shell ..."))
	fmt.Println()
	fmt.Println(color.Colorize(color.Bold, "Welcome to the Master Node"))
	fmt.Println(color.Colorize(color.Bold, "Type HELP for list of commands"))
	fmt.Println(color.Colorize(color.Bold, "Type EXIT to exit"))
	fmt.Print()
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wd = "X"
	}
	fmt.Println()
	fmt.Print(color.Colorize(color.Blue, "["+wd+"] master> "))
	go CommandListner()
}
