package jobs

import (
	"WorkerGobees/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/TwiN/go-color"
)

func MapJob(w http.ResponseWriter, r *http.Request) {
	os.Remove("./MAPPART00000")
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
		log.Println(color.Colorize(color.Red,"error running mapper file"))
	  os.WriteFile("./MAPERROR", out, 0777)
	  log.Println(string(out))
		utils.SimpleFailStatus(string(out), w)
		return
	}
	err = os.WriteFile("./MAPPART00000", out, 0777)
	if err != nil {
		log.Println(err)
		log.Println(color.Colorize(color.Red, "Error storing map output"))
		return
	}
	log.Println(color.Colorize(color.Green, "Succesfully completed assigned map task"))
	os.Remove("./" + handler.Filename)
	utils.SimpleSuccesssStatus("Finished map taks!", w)
}
