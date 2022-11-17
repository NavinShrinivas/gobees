package data

import (
	// "MasterGobees/shell"
	"MasterGobees/globals"
	"bufio"
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"errors"

	"github.com/TwiN/go-color"
)

func SplitAndUploadFile(file_path string, delimiter string) error{
  file_record := globals.File{
    File_name: filepath.Base(file_path),
  }
  fd, err := os.Open(file_path)
  if err!=nil{
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
        log.Printf(color.Colorize(color.Red,"Error reading file"))
        return err
    }
    if string(b) == delimiter{
      delimiter_count+=1
    }
  }
  fd.Close()
  number_of_split := len(globals.WorkerNodesMetadata)
  file_record.Splits = int32(number_of_split)
  delimiter_per_split := int32(delimiter_count/number_of_split)
  //[MUST]After we have the above three metrics in theory we can parellelise the code

  fd, err = os.Open(file_path)
  r = bufio.NewReader(fd)
  b = make([]byte, 1)
  for i:=1;i<=number_of_split;i++{
    number_of_delims_in_split := 0
    var temp_split_content []byte
     for {
        if number_of_delims_in_split == int(delimiter_per_split) && i!=number_of_split{
          //Second notation in if condition to cover entire file in the last split
          break
        }
        _, err := r.Read(b)
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf(color.Colorize(color.Red,"Error reading file mid way"))
            return err
        }
        temp_split_content = append(temp_split_content,b...)
        if string(b) == delimiter{
          number_of_delims_in_split+=1
        }
      }
      f_split, err:= os.Create("./temp_splits/"+filepath.Base(file_path)+"_PART00000")
      if err!=nil{
        log.Println(color.Colorize(color.Red, "Error creatring split"))
        return err
      }
      f_split.Write(temp_split_content)
      f_split.Close()
	    temp_node := "http://"+globals.WorkerNodesMetadata[i-1].Ip_addr+":"+globals.WorkerNodesMetadata[i-1].Port+"/storefile"
		  err = UploadData("./temp_splits/"+filepath.Base(file_path)+"_PART00000",temp_node,file_path)
      if err!=nil{
        log.Println(color.Colorize(color.Red, "Error uploading file"))
        return err
      }
      os.Remove("./temp_splits/"+filepath.Base(file_path)+"_PART00000") //Removing split for next split
      file_record.Nodes = append(file_record.Nodes, globals.WorkerNodesMetadata[i-1])
  }
  globals.FileMetadata = append(globals.FileMetadata, file_record) //Storing to file metadata only if everyhting went fine!
  return nil
}


func UploadData(TempFilePath string, node string, fileName string) error{
  file, _ := os.Open(TempFilePath)
  defer file.Close()

  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, _ := writer.CreateFormFile("File", filepath.Base(fileName))
  io.Copy(part, file)
  writer.Close()

  r, err := http.NewRequest("POST",node, body)
  if err!=nil{
    log.Println(color.Colorize(color.Red,"Something went wrong uploading file to worker : "+node))
    return err
  }
  r.Header.Add("Content-Type", writer.FormDataContentType())
  client := &http.Client{}
  res,err := client.Do(r)
  if err!=nil{
		log.Println(color.Colorize(color.Red,"Something went wrong reading response from worker : "+node))
		return err
  }

	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(color.Colorize(color.Red,"Something went wrong reading response from worker : "+node))
		return err
	}
	var res_body_obj map[string]interface{}
	err = json.Unmarshal(res_body, &res_body_obj)
	if err!=nil{
		log.Fatal(color.Colorize(color.Red,"[ENDPOINT ERROR] Error parsing request from Worker node addition request."))
		return err
	}
	if res_body_obj["status"] == false{
	  log.Println(color.Colorize(color.Red, "One of the Worker ran into an error while uploading file, node :"+node))
	  return errors.New("Error from WorkerNode")
	}
	log.Println(color.Colorize(color.Green,"Done writing part file to : "+node))
  return nil
}
