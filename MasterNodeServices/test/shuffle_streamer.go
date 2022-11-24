// V2 Shuffling function streamers, the function that needs to be change are :
// -

package main

import (
	"bufio"
	"fmt"
	"os"
	//Only std libs can be used
)

func streamer(key string) string {
	//This function gets each key as a string, returns back some string
	//Anything below this can be modified-----------------------------
	ans := key
	//Anything above this can be modified-----------------------------
	return ans
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	main_out := ""
	for scanner.Scan() {
		stdin_temp_key := scanner.Text()
		main_out += streamer(stdin_temp_key) + ","
	}
	fmt.Print(main_out)
}
