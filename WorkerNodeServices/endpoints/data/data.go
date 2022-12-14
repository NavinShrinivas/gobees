package data

import (
	"WorkerGobees/utils"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TwiN/go-color"
)

func StoreFile(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Receiving split file..."))
	os.Mkdir("./SS", 0777)
	file, handler, err := r.FormFile("File")
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error receiving file, please check."))
		utils.SimpleFailStatus("Failed storing file in worker", w)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error receiving file, please check."))
		utils.SimpleFailStatus("Failed storing file in worker", w)
		return
	}

	//[TODO] Later on we should be able to dynamically change path on worker nodes as well
	new_fd, err := os.Create("./SS/" + handler.Filename)
	if err != nil {
		log.Println(err)
		utils.SimpleFailStatus("Failed storing file in wokrker", w)
		return
	}

	defer new_fd.Close()

	_, err = new_fd.Write(fileBytes)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error storing file!"))
		utils.SimpleFailStatus("Failed storing file in wokrker", w)
		return
	}
	log.Println(color.Colorize(color.Green, "Storing split file success."))
	utils.SimpleSuccesssStatus("", w)
}

func FetchFile(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Reciving a fetch file request..."))
	if r.Method != "POST" {
		utils.SimpleInvalidPath("Invalid path", w)
		return
	}
	//Let's not store worker node info staefully, entirely in memory
	res_body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res_body_obj map[string]string
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "[ENDPOINT ERROR] Error parsing request from Worker node addition request."))
	}
	file_in_ss_name := string(res_body_obj["SS_file"])
	_, err = os.Open("./SS/" + file_in_ss_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error opening file from SS"))
		utils.SimpleFailStatus("Error reading file in SS", w)
		return
	}
	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())
	fw, err := mw.CreateFormFile("Splitfile", filepath.Base(file_in_ss_name))
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error creating form file"))
		utils.SimpleFailStatus("Error creating form file", w)
		return
	}
	bytes, err := os.ReadFile("./SS/" + file_in_ss_name)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error reading file from SS"))
		utils.SimpleFailStatus("Error reading file in SS", w)
		return
	}
	_, err = fw.Write(bytes)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Network transmission error"))
		utils.SimpleFailStatus("Error writing back multipart form data!", w)
		return
	}
}

func RenameFile(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Renaming a split file..."))
	os.Mkdir("./SS", 0777)
	// read json body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error renaming file, please check."))
		utils.SimpleFailStatus("Failed renaming file in worker", w)
		return
	}
	// unmarshal json body
	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error renaming file, please check."))
		utils.SimpleFailStatus("Failed renaming file in worker", w)
		return
	}

	log.Println("Renaming file: " + data["oldpath"] + " to " + data["newpath"])

	_, err = os.Stat("./SS/" + data["oldpath"])
	if os.IsNotExist(err) {
		log.Println(color.Colorize(color.Red, "Error renaming file, file does not exist in SS."))
		utils.SimpleFailStatus("Failed renaming file in worker, file does not exist in SS", w)
		return
	}

	err = os.Rename("./SS/"+data["oldpath"], "./SS/"+data["newpath"])
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error renaming file, please check."))
		utils.SimpleFailStatus("Failed renaming file in worker", w)
		return
	}
	log.Println(color.Colorize(color.Green, "Renamed split file success."))
	utils.SimpleSuccesssStatus("Renaming split file success", w)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Deleting a split file..."))
	os.Mkdir("./SS", 0777)
	// read json body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error deleting file, please check."))
		utils.SimpleFailStatus("Failed deleting file in worker", w)
		return
	}
	// unmarshal json body
	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error deleting file, please check."))
		utils.SimpleFailStatus("Failed deleting file in worker", w)
		return
	}
	// delete file
	err = os.Remove("./SS/" + data["filepath"])
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error deleting file, please check."))
		utils.SimpleFailStatus("Failed deleting file in worker", w)
		return
	}
	log.Println(color.Colorize(color.Green, "Deleted split file success."))
	utils.SimpleSuccesssStatus("Deleted split file success", w)
}
