package jobs

import (
	"WorkerGobees/globals"
	"WorkerGobees/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	hash "github.com/theTardigrade/golang-hash"
	"github.com/twotwotwo/sorts" //Gives fast as hek in memory Quick Sort impl
)

func MapJob(w http.ResponseWriter, r *http.Request) {
	os.Remove("./INTERPART00001") //Output from mapper stage
	os.Remove("./INTERPART00002") //Output from suffle stage
	os.Remove("./INTERPART00003") //Output from sorting
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
	err = os.WriteFile("./INTERPART00001", out, 0777)
	if err != nil {
		log.Println(err)
		log.Println(color.Colorize(color.Red, "Error storing map output"))
		return
	}
	log.Println(color.Colorize(color.Green, "Succesfully completed assigned map task"))
	os.Remove("./" + handler.Filename)
	utils.SimpleSuccesssStatus("Finished map taks!", w)
}


func StartShuffle(w http.ResponseWriter, r *http.Request){

	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Reciving a Shuffle job..."))
	os.Create("./INTERPART00002")
	custom_function := r.FormValue("custom")
	if custom_function == "true"{
		//[TODO]
		customShuffleFunction(w,r)
	  return
	}
	NodeMetaData := r.FormValue("NodeInfo")
	err := json.Unmarshal([]byte(NodeMetaData),&globals.ShuffleNodeMetadata)
	if err!=nil{
		log.Println(color.Colorize(color.Red, "Error getting Node info from Master node for shuffle taks "))
		utils.SimpleFailStatus("Node Info Error",w)
		return
	}
	mod_value := len(globals.ShuffleNodeMetadata)

	fd, err := os.Open("./INTERPART00001") //Read from mapper file
	if err!=nil{
		log.Println(color.Colorize(color.Red,"Error opening Map output file!!!"))
		utils.SimpleFailStatus("Map outputs file open error!!",w)
		return
	}
  scanner := bufio.NewScanner(fd)
	for scanner.Scan(){
		line := string(scanner.Bytes())
		line_split := strings.Split(line,",")
		key := line_split[0]
		key_hash := hash.UintString(key)
		to_node := uint(key_hash)%uint(mod_value)
		line += "\n"
		var base_url_of_to_node string
		if globals.ShuffleNodeMetadata[to_node].Ip_addr == "0.0.0.0"{
			//If they were part of local cluster then this is activated
			master_url_split := strings.Split(globals.MasterUrl, ":") 
			base_url_of_to_node = "http:"+master_url_split[1]+":"+globals.ShuffleNodeMetadata[to_node].Port
		}else if globals.ShuffleNodeMetadata[to_node].Ip_addr == globals.Ip{
			//If they were part of global cluster but for current node its local
			//Need Inspection, some nodes we failing on access their servers on 0.0.0.0 :skull:
			base_url_of_to_node = "http://" + globals.ShuffleNodeMetadata[to_node].Ip_addr + ":" + globals.ShuffleNodeMetadata[to_node].Port 
		}else{
			//If completetly global
			base_url_of_to_node = "http://" + globals.ShuffleNodeMetadata[to_node].Ip_addr + ":" + globals.ShuffleNodeMetadata[to_node].Port 
		}
		to_node_url := base_url_of_to_node + "/shuffleshare"
		log.Println(to_node_url)
		body_streamer := bytes.NewReader([]byte(line))
		res, err := http.Post(to_node_url,"text/plain",body_streamer)
		if err!=nil{
			log.Println(color.Colorize(color.Red,"Error on network interactions during shuffle"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
		res_body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var res_body_obj map[string]interface{}
		err = json.Unmarshal(res_body, &res_body_obj)
		if err != nil {
			log.Fatal(color.Colorize(color.Red, "Something went wrong reading reponse from other worker node"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
		if res_body_obj["status"] == false {
			log.Fatal(color.Colorize(color.Red, "Something went wrong on to node during shuffle"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
	}
  err = InMemSorter("./INTERPART00002")
  if err!=nil{
  	log.Println(color.Colorize(color.Red, "Failed sorting file"))
  	utils.SimpleFailStatus("Failed sorting file from shuffling stage", w)
  	return 
  }
	utils.SimpleSuccesssStatus("",w)
	return
}

func ShuffleShare(w http.ResponseWriter, r *http.Request){
  fd, err := os.OpenFile("./INTERPART00002", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
  defer fd.Close()
  if err!=nil{
  	log.Println(color.Colorize(color.Red,"Error storing shared suffle kv pair"))
  	utils.SimpleFailStatus("Error storing kv pair",w)
  	return
  }
	res_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
  	log.Println(color.Colorize(color.Red,"Error reading shared suffle kv pair"))
  	utils.SimpleFailStatus("Error reading shared kv pair",w)
  	return
	}
  fd.WriteString(string(res_body))
  utils.SimpleSuccesssStatus("",w)
  return
}


func customShuffleFunction(w http.ResponseWriter,r *http.Request){
	file, handler, err := r.FormFile("ShuffleFile")
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error recieving custom shuffle file, please check."))
		utils.SimpleFailStatus("Failed storing file in worker for shuffle proc", w)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
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
	log.Println(color.Colorize(color.Green, "Recieved Shuffle file!"))
	new_fd.Close()
	file.Close()
	NodeMetaData := r.FormValue("NodeInfo")
	err = json.Unmarshal([]byte(NodeMetaData),&globals.ShuffleNodeMetadata)
	if err!=nil{
		log.Println(color.Colorize(color.Red, "Error getting Node info from Master node for shuffle taks "))
		utils.SimpleFailStatus("Node Info Error",w)
		return
	}
	mod_value := len(globals.ShuffleNodeMetadata)
	fd, err := os.Open("./INTERPART00001") //Read from map output
	if err!=nil{
		log.Println(color.Colorize(color.Red,"Error opening Map output file!!!"))
		utils.SimpleFailStatus("Map outputs file open error!!",w)
		return
	}
  scanner := bufio.NewScanner(fd)
	cmd := exec.Command("go","run","./"+handler.Filename)
	stdinPipe, _ := cmd.StdinPipe()
	for scanner.Scan(){
		line := string(scanner.Bytes())
		line_split := strings.Split(line,",")
		key := line_split[0]+"\n"
		stdinPipe.Write([]byte(key))
	}
	stdinPipe.Close()
	fd.Close()
	com_put, _ := cmd.CombinedOutput()
	hash_values := strings.Split(string(com_put),",")
	fd, err = os.Open("./INTERPART00001")
	if err!=nil{
		log.Println(color.Colorize(color.Red,"Error opening Map output file!!!"))
		utils.SimpleFailStatus("Map outputs file open error!!",w)
		return
	}
  scanner = bufio.NewScanner(fd)
  kv_pair := 0
  for scanner.Scan(){
  	line := scanner.Text()
		key_hash := hash_values[kv_pair]
		kv_pair+=1
		key_hash_int, err := strconv.Atoi(strings.Trim(key_hash,"\n"))
		if err!=nil{
			log.Println(color.Colorize(color.Red,"Invalid Streaming shuffle function"))
			utils.SimpleFailStatus("INVALID SHUFFLE FUNCTION",w)
			return
		}
		to_node := uint(key_hash_int)%uint(mod_value)
		line += "\n"
		to_node_url := "http://" + globals.ShuffleNodeMetadata[to_node].Ip_addr + ":" + globals.ShuffleNodeMetadata[to_node].Port + "/shuffleshare"
		body_streamer := bytes.NewReader([]byte(line))
		res, err := http.Post(to_node_url,"text/plain",body_streamer)
		if err!=nil{
			log.Println(color.Colorize(color.Red,"Error on network interactions during shuffle"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
		res_body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var res_body_obj map[string]interface{}
		err = json.Unmarshal(res_body, &res_body_obj)
		if err != nil {
			log.Fatal(color.Colorize(color.Red, "Something went wrong reading reponse from other worker node"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
		if res_body_obj["status"] == false {
			log.Fatal(color.Colorize(color.Red, "Something went wrong on to node during shuffle"))
			utils.SimpleFailStatus("Error network transffer during shuffle",w)
			return
		}
  }
  os.Remove("./"+handler.Filename)
  err = InMemSorter("./INTERPART00002")
  if err!=nil{
  	log.Println(color.Colorize(color.Red, "Failed sorting file"))
  	utils.SimpleFailStatus("Failed sorting file from shuffling stage", w)
  	return 
  }
	utils.SimpleSuccesssStatus("",w)
	return
}

type record struct{
  Key string
	Value string
}
type InMemFile []record

func (a InMemFile) Len() int{
	return (len(a))
}

func (a InMemFile) Swap(i, j int) { 
	temp := a[j]
	a[j] = a[i]
	a[i] = temp
}

func (a InMemFile) Key(i int) string {
	return a[i].Key
}

func (a InMemFile) Less(i, j int) bool {
	//Understanding of this interface function comes from : 
	//https://github.com/twotwotwo/sorts/blob/master/radixsort.go
	//Line 106, tells us where this function is used and what's expected of it
	//POWER OF OPEN SOURCE!
	if a[i].Key < a[j].Key{
		return true
	}
	return false
}

func InMemSorter(file_name string) error{
	raw_file, err := os.ReadFile(file_name)
	if err!=nil{
		log.Println(color.Colorize(color.Red, "[SORT ERROR] Error reading file from shuffle task, maybe shuffle failed?"))
		return err
	}
	file_array := strings.Split(string(raw_file),"\n")
	var file InMemFile
	for _,v := range file_array{
		kv_pair := strings.Split(v, ",")
		if len(kv_pair) < 2{
			continue
		}
		new_record :=  record{
			Key : kv_pair[0],
			Value : kv_pair[1],
		}
		file = append(file, new_record)
	}
	sorts.ByString(file)
	output_string := ""
	for _,v := range file{
		temp := v.Key+","+v.Value+"\n"
		output_string += temp
	}
	os.WriteFile("./INTERPART00003", []byte(output_string), 0777)
	return nil

}

func ReduceJob(w http.ResponseWriter, r *http.Request) {
	log.Println(color.Colorize(color.Yellow, "[ENDPOINT] Reciving a Reduce job..."))

	file, handler, err := r.FormFile("ReducerFile")
	if err != nil {
		log.Println(color.Colorize(color.Red, "Error recieving file, please check."))
		utils.SimpleFailStatus("Failed storing file in worker for reduce proc", w)
		return
	}
	reduce_args := r.FormValue("ReducerArgs")
	output_file_name := r.FormValue("OutputFileSS")

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
	log.Println(color.Colorize(color.Green, "Recieved reduce file!"))
	new_fd.Close()
	file.Close()
	log.Println(color.Colorize(color.Yellow, "Running Reduce job"))
	cmd := exec.Command("python3", handler.Filename, reduce_args)
	stdinPipe, err := cmd.StdinPipe()

	go func() {
		defer stdinPipe.Close()
		map_input_file, err := os.ReadFile("./INTERPART00003")
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
	  os.WriteFile("./REDUCEERROR", out, 0777)
	  log.Println(string(out))
		utils.SimpleFailStatus(string(out), w)
		return
	}
	err = os.WriteFile("./SS/"+output_file_name, out, 0777)
	if err != nil {
		log.Println(err)
		log.Println(color.Colorize(color.Red, "Error storing reduce output"))
		return
	}
	log.Println(color.Colorize(color.Green, "Succesfully completed assigned reduce task"))
	os.Remove("./" + handler.Filename)
	utils.SimpleSuccesssStatus("Finished reduce taks!", w)
}
