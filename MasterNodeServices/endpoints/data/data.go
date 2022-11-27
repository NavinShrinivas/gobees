package data

import (
	// "MasterGobees/shell"
	"MasterGobees/globals"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"mime"

	"path/filepath"

	"github.com/TwiN/go-color"
)

func SplitAndUploadFile(file_path string, delimiter string) error {
	file_record := globals.File{
		File_name: filepath.Base(file_path),
	}
	fd, err := os.Open(file_path)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading file."))
		return err
	}
	delimiter_count := 0
	//Stepping through file inorder to not peak the ram usage!
	r := bufio.NewReader(fd)
	b := make([]byte, 1)
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf(color.Colorize(color.Red, "Error reading file"))
			return err
		}
		if string(b) == delimiter {
			delimiter_count += 1
		}
	}
	fd.Close()
	number_of_split := len(globals.WorkerNodesMetadata)
	file_record.Splits = int32(number_of_split)
	if number_of_split == 0 {
		log.Println(color.Colorize(color.Red, "There are not worker node to store data in!"))
		return errors.New("No Worker nodes!")
	}
	delimiter_per_split := int32(delimiter_count / number_of_split)
	if delimiter_per_split == 0 {
		delimiter_per_split = 1 //Minimum amount, but now we need to handle empty splits
	}
	//[MUST]After we have the above three metrics in theory we can parellelise the code

	fd, err = os.Open(file_path)
	r = bufio.NewReader(fd)
	b = make([]byte, 1)
	var i int
	for i = 1; i <= number_of_split; i++ {
		number_of_delims_in_split := 0
		var temp_split_content []byte
		for {
			if number_of_delims_in_split == int(delimiter_per_split) && i != number_of_split {
				//Second notation in if condition to cover entire file in the last split
				break
			}
			_, err := r.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf(color.Colorize(color.Red, "Error reading file mid way"))
				return err
			}
			temp_split_content = append(temp_split_content, b...)
			if string(b) == delimiter {
				number_of_delims_in_split += 1
			}
		}
		os.Mkdir("./temp_splits", 0777)
		f_split, err := os.Create("./temp_splits/" + filepath.Base(file_path) + "_PART00000")
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error creatring split"))
			return err
		}
		f_split.Write(temp_split_content)
		f_split.Close()
		temp_node := "http://" + globals.WorkerNodesMetadata[i-1].Ip_addr + ":" + globals.WorkerNodesMetadata[i-1].Port + "/storefile"
		err = UploadData("./temp_splits/"+filepath.Base(file_path)+"_PART00000", temp_node, file_path)
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error uploading file"))
			return err
		}
		file_record.Nodes = append(file_record.Nodes, globals.WorkerNodesMetadata[i-1])
	}

	//Below for loop only to handle empty files, if we have not send a file to every workernode but have reached end of file
	for i = i; i <= number_of_split; i++ {
		f_split, err := os.Create("./temp_splits/" + filepath.Base(file_path) + "_PART00000")
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error creatring split"))
			return err
		}
		f_split.Close()
		temp_node := "http://" + globals.WorkerNodesMetadata[i-1].Ip_addr + ":" + globals.WorkerNodesMetadata[i-1].Port + "/storefile"
		err = UploadData("./temp_splits/"+filepath.Base(file_path)+"_PART00000", temp_node, file_path)
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error uploading file"))
			return err
		}
		file_record.Nodes = append(file_record.Nodes, globals.WorkerNodesMetadata[i-1])
	}
	globals.FileMetadata = append(globals.FileMetadata, file_record) //Storing to file metadata only if everyhting went fine!
	os.RemoveAll("./temp_splits")
	return nil
}

func UploadData(TempFilePath string, node string, fileName string) error {
	file, _ := os.Open(TempFilePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("File", filepath.Base(fileName))
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("POST", node, body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong uploading file to worker : "+node))
		return err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker : "+node))
		return err
	}

	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker : "+node))
		return err
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "[ENDPOINT ERROR] Error parsing request from Worker node addition request."))
		return err
	}
	if res_body_obj["status"] == false {
		log.Println(color.Colorize(color.Red, "One of the Worker ran into an error while uploading file, node :"+node))
		return errors.New("Error from WorkerNode")
	}
	log.Println(color.Colorize(color.Green, "Done writing part file to : "+node))
	return nil
}

type FileFetchRequest struct{
	SS_file string `json:"SS_file"`
}

func FetchAndMergeFile(v globals.WorkerNode, to_file_path string, ss_file string) error{
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(to_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
    log.Println(color.Colorize(color.Red, "Error creating/opening files to write in local."))
    return err
	}
	temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/fetchfile"
	request := FileFetchRequest{
		SS_file: ss_file,
	}
	request_bytes, err := json.Marshal(request)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error Marshalling response on a fail status"))
		return err
	}
	r, err := http.Post(temp_node, "application/json", bytes.NewBuffer(request_bytes))
	if err!=nil{
		log.Println(color.Colorize(color.Red, "Error sending fetch request to worker"))
		return err
	}
	_, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	mr := multipart.NewReader(r.Body,params["boundary"])
	//We let multipart form take care of the parts it wants to split into
	//meaning we need to collect all part here and combine them
	var part_data []byte
	for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
		value, _ := ioutil.ReadAll(part)
		part_data = append(part_data, value...)
	}
	if err!=nil{
		log.Println(color.Colorize(color.Red, "Error getting response from workernode when fetching files : "+temp_node))
		return err
	}
	f.Write(part_data)
	f.Close()
	log.Println(color.Colorize(color.Green, "Succefully recived split from : "+temp_node))
	return nil
}

func RenameFile(oldpath string, newpath string) error {
	for i := 0; i < len(globals.WorkerNodesMetadata); i++ {
		worker := globals.WorkerNodesMetadata[i]
		url := "http://" + worker.Ip_addr + ":" + worker.Port + "/renamefile"
		data := map[string]string{"oldpath": oldpath, "newpath": newpath}
		jsonValue, _ := json.Marshal(data)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error while renaming file in worker node"))
			return err
		}
		defer response.Body.Close()
	}
	return nil
}

func DeleteFile(file_path string) error {
	for i := 0; i < len(globals.WorkerNodesMetadata); i++ {
		worker := globals.WorkerNodesMetadata[i]
		url := "http://" + worker.Ip_addr + ":" + worker.Port + "/deletefile"
		data := map[string]string{"filepath": file_path}
		jsonValue, _ := json.Marshal(data)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Println(color.Colorize(color.Red, "Error while deleting file in worker node"))
			return err
		}
		defer response.Body.Close()
	}
	return nil
}
