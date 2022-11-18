package jobs

import (
	"WorkerGobees/utils"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/TwiN/go-color"
)

func MapJob(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Reciving a Map job..."))

	file, handler, err := r.FormFile("MapperFile")
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error recieving file, please check."))
		utils.SimpleFailStatus("Failed storing file in worker for mapper proc", w)
		return
	}
	map_args := r.FormValue("MapperArgs")
	input_file_name := r.FormValue("InputFileSS")

	fileBytes, err := ioutil.ReadAll(file)

	//[TODO] Later on we should be able to dynamically change path on worker nodes as well
	new_fd, err := os.Create("./" + handler.Filename)
	os.Chmod("./"+handler.Filename, 0777)
	if err != nil {
		log.Println(err)
		utils.SimpleFailStatus("Failed storing file in wokrker", w)
		return
	}

	_, err = new_fd.Write(fileBytes)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error storing file!"))
		utils.SimpleFailStatus("Failed storing file in wokrker", w)
		return
	}
	log.Println(color.Colorize(color.Green, "Recieved Map file!"))
	new_fd.Close()
	file.Close()
	log.Println(color.Colorize(color.Yellow, "Running Map job"))
	cmd := exec.Command("python3", handler.Filename, map_args)
	stdinPipe, err := cmd.StdinPipe()

	go func() {
		defer stdinPipe.Close()
		map_input_file, err := os.ReadFile("./SS/" + input_file_name)
		if err != nil {
			log.Println(err)
			utils.SimpleFailStatus("Error reading split File", w)
			return
		}
		stdinPipe.Write(map_input_file)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		utils.SimpleFailStatus("Error storing output after map job", w)
		return
	}
	err = os.WriteFile("./MAPPART00000", out, 0777)
	if err != nil {
		log.Println(err)
		log.Println(color.Colorize(color.Red, "Error storing map output"))
		return
	}
	log.Println(color.Colorize(color.Green, "Succesfully completed assigned map task"))
	os.Remove("./MAPPART00000")
	os.Remove("./" + handler.Filename)
	utils.SimpleSuccesssStatus("Finished map taks!", w)
}
func copyStdin(stdin io.Writer, input_file_name string, error_chan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	split_file, err := os.ReadFile("./SS/" + input_file_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "split not found!"))
		error_chan <- err
		return
	}
	io.WriteString(stdin, string(split_file))
	error_chan <- nil
}

func copyStdout(r io.Reader, file_name string, error_chan chan error, wg *sync.WaitGroup) {
	wg.Done()
	os.Remove("./" + file_name)
	file, err := os.Create("./" + file_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "error storing intermediate output"))
		error_chan <- err
		return
	}
	bufio_reader := bufio.NewReader(r)
	for {
		lineout, _, _ := bufio_reader.ReadLine()
		fmt.Println(string(lineout))
		file.Write(lineout)
		if err == io.EOF {
			break
		}
	}
	fmt.Println("stdout done")
	error_chan <- nil
}
