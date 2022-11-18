package jobs

import (
	"MasterGobees/globals"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/TwiN/go-color"
)

func SendMapJobs(mapper_command string, input_file string) error {
	command_split := strings.Split(mapper_command, " ")
	mapper_file_path := command_split[0]
	mapper_args := ""
	for i := 1; i < len(command_split); i++ {
		mapper_args += command_split[i] + " "
	}
	file, _ := os.Open(mapper_file_path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("MapperArgs", mapper_args)
	writer.WriteField("InputFileSS", input_file)
	part, _ := writer.CreateFormFile("MapperFile", filepath.Base(mapper_file_path))
	io.Copy(part, file)
	formtype := writer.FormDataContentType()
	writer.Close()

	local_wg := new(sync.WaitGroup)
	err_chan := make(chan error, 1000)
	for _, v := range globals.WorkerNodesMetadata {
		temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/mapjob"
		local_wg.Add(1)
		go NodeMapJob(body, temp_node, local_wg, err_chan, formtype)
	}
	local_wg.Wait() //After all the threads are done, we can see if any errors in channel
	for err := range err_chan {
		if err != nil {
			log.Println(color.Colorize(color.Red, "Map job failed :<("))
			return err
		}
		if err == nil {
			//Meaning out map job passed!
			log.Println(color.Colorize(color.Green, "Map Job completed Succesfully!"))
			return nil
		}
	}
	return nil

}

func NodeMapJob(body *bytes.Buffer, node string, wg *sync.WaitGroup, err_chan chan error, formtype string) {
	r, err := http.NewRequest("POST", node, body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong uploading map file to worker : "+node))
		err_chan <- err
	}
	r.Header.Add("Content-Type", formtype)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker : "+node))
		err_chan <- err
	}

	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker during map job :"+node))
		err_chan <- err
		wg.Done()
		return
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Fatal(color.Colorize(color.Red, "Error parsing response from worker during Map reduce job"))
		err_chan <- err
		wg.Done()
		return
	}
	if res_body_obj["status"] == false {
		log.Println(color.Colorize(color.Red, "One of the Worker ran into an error while uploading file, node :"+node))
		err_chan <- errors.New("Error from WorkerNode")
		wg.Done()
		return
	}
	log.Println(color.Colorize(color.Green, "Done with Map job on : "+node))
	err_chan <- nil
	wg.Done()
	return
}
