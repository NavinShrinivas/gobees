package shell

import (
	"MasterGobees/globals"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var reader *bufio.Reader 

func show_commands(command_parse []string){
	if command_parse[1] == "CONFIG"{
		fmt.Println(globals.Config_obj)
		return
	}
	if command_parse[1] == "JOBS"{
		fmt.Println("Past job history : ")
		return
	}
	fmt.Println("Invalid command")
	return
}

func commandParser(command string){
	command_parse := strings.Split(strings.Trim(command,"\n"), " ")
	if command_parse[0] == "SHOW" && len(command_parse) >= 2{
		show_commands(command_parse)
	}else{
		fmt.Println("Invalid command")
	}
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
