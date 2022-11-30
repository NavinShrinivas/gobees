package jobs

import (
	"MasterGobees/globals"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	log.Println(color.Colorize(color.Yellow, "Starting Map Job"))
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
	var nodes_with_split []globals.WorkerNode
	for _, v := range globals.FileMetadata {
		if v.File_name == input_file {
			nodes_with_split = v.Nodes
		}
	}
	for _, v := range nodes_with_split {
		temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/mapjob"
		local_wg.Add(1)
		go NodeMapJob(body, temp_node, local_wg, err_chan, formtype)
	}
	local_wg.Wait() //After all the threads are done, we can see if any errors in channel
	for err := range err_chan {
		if err != nil {
			log.Println(color.Colorize(color.Red, "Map job failed :("))
			return err
		}
		if err == nil {
			//Meaning out map job passed!
			log.Println(color.Colorize(color.Green, "Map Job completed Successfully!"))
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

	res_body, err := io.ReadAll(res.Body)
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
		log.Println(color.Colorize(color.Red, "One of the Worker Nodes ran into an error while running map, node :"+node))
		log.Println("Error from node : ")
		fmt.Println(res_body_obj["message"])
		err_chan <- errors.New("Error from WorkerNode")
		wg.Done()
		return
	}
	log.Println(color.Colorize(color.Green, "Done with Map job on : "+node))
	err_chan <- nil
	wg.Done()
	return
}

func StartShuffle(custom_function bool, shuffle_file_path string) error {
	log.Println(color.Colorize(color.Yellow, "Starting Shuffle Job"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if custom_function {
		writer.WriteField("custom", "true")
		file, _ := os.Open(shuffle_file_path)
		part, _ := writer.CreateFormFile("ShuffleFile", filepath.Base(shuffle_file_path))
		io.Copy(part, file)
		defer file.Close()
	} else {
		writer.WriteField("custom", "false")
	}
	metadata_bytes, err := json.Marshal(globals.WorkerNodesMetadata)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error packaging Node metadata for shuffle job"))
	}
	writer.WriteField("NodeInfo", string(metadata_bytes))

	formtype := writer.FormDataContentType()
	writer.Close()

	local_wg := new(sync.WaitGroup)
	err_chan := make(chan error, 1000)
	for _, v := range globals.WorkerNodesMetadata {
		temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/startshuffle"
		local_wg.Add(1)
		go NodeShuffleJob(body, temp_node, local_wg, err_chan, formtype)
	}
	local_wg.Wait() //After all the threads are done, we can see if any errors in channel
	for err := range err_chan {
		if err != nil {
			log.Println(color.Colorize(color.Red, "Shuffle init job failed :("))
			return err
		}
		if err == nil {
			//Meaning out map job passed!
			log.Println(color.Colorize(color.Green, "Shuffle Job completed Successfully!"))
			return nil
		}
	}
	for _, v := range globals.WorkerNodesMetadata {
		temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/sortshuffle"
		local_wg.Add(1)
		go func(){
			http.Get(temp_node)
		}()
	}
	return nil
}

func NodeShuffleJob(body *bytes.Buffer, node string, wg *sync.WaitGroup, err_chan chan error, formtype string) {
	r, err := http.NewRequest("POST", node, body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong uploading shuffle file to worker : "+node))
		err_chan <- err
	}
	r.Header.Add("Content-Type", formtype)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker : "+node))
		err_chan <- err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker during shuffle job :"+node))
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
		log.Println(color.Colorize(color.Red, "One of the Worker ran into an error while running shuffle, node :"+node))
		log.Println("Error from node : ")
		fmt.Println(res_body_obj["message"])
		err_chan <- errors.New("Error from WorkerNode")
		wg.Done()
		return
	}
	log.Println(color.Colorize(color.Green, "Done with Shuffle job on : "+node))
	err_chan <- nil
	wg.Done()
	//Shuffler function is also made to sort, after these three steps the file remaining is
	//INTERPART00003
	return
}

func StartReduce(reducer_command string, output_file string) error {
	//Need to upload reducer file along with output file name
	log.Println(color.Colorize(color.Yellow, "Starting Reducer Job"))

	command_split := strings.Split(reducer_command, " ")
	reducer_file_path := command_split[0]
	reducer_args := ""
	for i := 1; i < len(command_split); i++ {
		reducer_args += command_split[i] + " "
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := os.Open(reducer_file_path)
	defer file.Close()
	writer.WriteField("ReducerArgs", reducer_args)
	writer.WriteField("OutputFileSS", output_file)
	part, _ := writer.CreateFormFile("ReducerFile", filepath.Base(reducer_file_path))
	io.Copy(part, file)
	formtype := writer.FormDataContentType()
	writer.Close()

	local_wg := new(sync.WaitGroup)
	err_chan := make(chan error, 1000)
	for _, v := range globals.WorkerNodesMetadata {
		temp_node := "http://" + v.Ip_addr + ":" + v.Port + "/reducejob"
		local_wg.Add(1)
		go NodeReduceJob(body, temp_node, local_wg, err_chan, formtype)
	}
	local_wg.Wait() //After all the threads are done, we can see if any errors in channel
	for err := range err_chan {
		if err != nil {
			log.Println(color.Colorize(color.Red, "Reduce job failed :("))
			return err
		}
		if err == nil {
			//Meaning out map job passed!

			//This way of adding meta from master node needs to be changed in the future
			log.Println(color.Colorize(color.Green, "Reduce Job completed Successfully!"))
			mr_job_out_file_meta := globals.File{
				File_name: output_file,
				Splits:    int32(len(globals.WorkerNodesMetadata)),
				Nodes:     globals.WorkerNodesMetadata,
			}
			globals.FileMetadata = append(globals.FileMetadata, mr_job_out_file_meta)
			return nil
		}
	}
	return nil
}

func NodeReduceJob(body *bytes.Buffer, node string, wg *sync.WaitGroup, err_chan chan error, formtype string) {
	r, err := http.NewRequest("POST", node, body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong uploading reduce file to worker : "+node))
		err_chan <- err
	}
	r.Header.Add("Content-Type", formtype)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker : "+node))
		err_chan <- err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red, "Something went wrong reading response from worker during reduce job :"+node))
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
		log.Println(color.Colorize(color.Red, "One of the Worker ran into an error while running reduce, node :"+node))
		log.Println("Error from node : ")
		fmt.Println(res_body_obj["message"])
		err_chan <- errors.New("Error from WorkerNode")
		wg.Done()
		return
	}
	log.Println(color.Colorize(color.Green, "Done with Reduce job on : "+node))
	err_chan <- nil
	wg.Done()
	//Outputs in SS after successful reduce stage should be the "OUT" file name
	return
}
