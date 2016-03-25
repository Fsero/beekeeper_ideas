package main

import (
	"os/exec"
	"fmt"
	"log"
	"encoding/json"
	 "io/ioutil"
	"os"
	"github.com/rakyll/magicmime"
	"path/filepath"
)

type Trace struct {
	ContainerId string `json:"container.id"`
	ContainerName string `json:"container.name"`
	EventCpu int `json:"evt.cpu"`
	EventDir string `json:"evt.dir"`
	EventInfo string `json:"evt.info"`
	EventNumber int `json:"evt.number"`
	EventOutputUnixTime int64 `json:"evt.outputtime"`
	EventType string `json:"evt.type"`
	ProcName string `json:"proc.name"`
	ThreadTid int `json:"thread.tid"`
	ThreadVTid int `json:"thread.vtid"`
}

func extractJson(f os.FileInfo) []Trace {
	fmt.Println(f.Name())
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	var container_name string = "ssh_ssh_1"
	defer magicmime.Close()
	path, err := filepath.Abs("./test/" + f.Name())
	fmt.Println("%s",path)
	mimetype, err := magicmime.TypeByFile(path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}

	log.Printf("mime-type: %v", mimetype)
	if mimetype != "application/gzip" {
		log.Fatalf(" Invalid format in file: %s \n", path)
	}

	containerName := fmt.Sprintf("container.name = %s",container_name)

	cmd := exec.Command("/usr/bin/sysdig","-N", "-b", "-r", path , "-pc" , containerName, "-j")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var traces []Trace
	if err := json.NewDecoder(stdout).Decode(&traces); err != nil {
		log.Fatal(err)
	}

	return traces


}

func dealwithfiles(files chan os.FileInfo) {
	for f := range files {
		fmt.Println("reading file ")
		traces := extractJson(f)
		for _, value := range traces {
			fmt.Printf("%+v", value)
			fmt.Println("\n")
		}
	}
}
func main() {
	path, err := exec.LookPath("sysdig")
	if err != nil {
		log.Fatal("cannot found sysdig, unable to fetch data")
	}
	fmt.Printf("sysdig is available at %s\n", path)
	var MAX_WORKERS = 2
	files, _ := ioutil.ReadDir("./test")
	chanFiles := make(chan os.FileInfo)
	for w:=1; w <= MAX_WORKERS; w++ {
		go dealwithfiles(chanFiles)
	}
	for _, f := range files {
		fmt.Println("sending file ")
		chanFiles <- f

	}




	/*var jsonBlob = []byte(`[
	 {"container.id":"154a173602f6",
	  "container.name":"ssh_ssh_1",
	  "evt.cpu":0,"evt.dir":"<","evt.info":"fd=5(<4t>81.218.170.178:44549->172.17.0.2:ssh) tuple=81.218.170.178:44549->172.17.0.2:ssh queuepct=0 queuelen=0 queuemax=128 ",
	  "evt.num":298012,"evt.outputtime":1457286491203423434,"evt.type":"accept","proc.name":"sshd","thread.tid":30662,"thread.vtid":9}
	]`)
	var traces []Trace
	err = json.Unmarshal(jsonBlob, &traces)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", traces)
	fmt.Printf("\n")
	fmt.Printf("\n")

	data, err := json.Marshal(traces)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
*/
}
