package data

import (
	"WorkerGobees/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/TwiN/go-color"
)

func StoreFile(w http.ResponseWriter, r *http.Request) {
    log.Println(color.Colorize(color.Yellow,"[ENDPOINT] Reciving a split file..."))

    file, handler, err := r.FormFile("File")
    if err != nil {
        log.Println(color.Colorize(color.Red, "Error recieving file, please check."))
        return
    }
    defer file.Close()
    log.Printf("Uploaded File: %+v\n", handler.Filename)
    log.Printf("File Size: %+v\n", handler.Size)

    fileBytes, err := ioutil.ReadAll(file)

    //[TODO] Later on we should be able to dynamically change path on worker nodes as well
    new_fd, err := os.Create("./SS/"+handler.Filename)
    defer new_fd.Close()
    if err != nil {
    	log.Println(err)
    }

    _, err = new_fd.Write(fileBytes)
    if err != nil {
        log.Println(color.Colorize(color.Red,"Error storing file!"))
    }
    log.Println(color.Colorize(color.Green, "Storing split file success."))
    utils.SimpleSuccesssStatus("",w)
}
