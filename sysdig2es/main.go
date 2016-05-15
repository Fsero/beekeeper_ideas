package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rakyll/magicmime"
	"gopkg.in/olivere/elastic.v3"
	"strconv"
	"time"
)

type Trace struct {
	ContainerId         string `json:"container.id"`
	ContainerName       string `json:"container.name"`
	EventCpu            int    `json:"evt.cpu"`
	EventDir            string `json:"evt.dir"`
	EventInfo           string `json:"evt.info"`
	EventNumber         int    `json:"evt.number"`
	EventOutputUnixTime int64  `json:"evt.outputtime"`
	EventType           string `json:"evt.type"`
	ProcName            string `json:"proc.name"`
	ThreadTid           int    `json:"thread.tid"`
	ThreadVTid          int    `json:"thread.vtid"`
}

func (t *Trace) MarshalJSON() ([]byte, error) {
	type Alias Trace
	return json.Marshal(&struct {
		ContainerId         string `json:"ContainerId"`
		ContainerName       string `json:"ContainerName"`
		EventCpu            int    `json:"evtCpu"`
		EventDir            string `json:"evtDir"`
		EventInfo           string `json:"evtInfo"`
		EventNumber         int    `json:"evtNumber"`
		EventOutputUnixTime string `json:"evtOutputtime"` //this would be a date in ES.
		EventType           string `json:"evtType"`
		ProcName            string `json:"procName"`
		ThreadTid           int    `json:"threadTid"`
		ThreadVTid          int    `json:"threadVtid"`
		*Alias
	}{
		ContainerId:   t.ContainerId,
		ContainerName: t.ContainerName,
		EventCpu:      t.EventCpu,
		EventDir:      t.EventDir,
		EventInfo:     t.EventInfo,
		EventNumber:   t.EventNumber,
		//sysdig gives us an unix timestamp of 16 digits that we cannot represent.
		EventOutputUnixTime: (strconv.FormatInt(t.EventOutputUnixTime, 10)[0:13]),
		EventType:           t.EventType,
		ProcName:            t.ProcName,
		ThreadTid:           t.ThreadTid,
		ThreadVTid:          t.ThreadVTid,
	})
}

//noinspection ALL
func extractJson(f os.FileInfo) []Trace {
	fmt.Println(f.Name())
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	//TODO: parametrize this.
	var container_name string = "ssh_ssh_1"
	defer magicmime.Close()
	path, err := filepath.Abs("./data/srv/capture/" + f.Name())
	fmt.Println("%s", path)
	mimetype, err := magicmime.TypeByFile(path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}

	log.Printf("mime-type: %v", mimetype)
	if (mimetype != "application/gzip") && (mimetype != "application/octet-stream") {
		log.Printf(" Invalid format in file: %s \n", path)
	}

	containerName := fmt.Sprintf("container.name = %s", container_name)

	tmp_path := fmt.Sprintf("/tmp/%s.json", f.Name())

	cmd := fmt.Sprintf("/usr/bin/sysdig -N -A -r %s -pc %s -j > %s", path, containerName,tmp_path)
	fmt.Printf(cmd)
	out, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		fmt.Sprintf("Failed to execute command: %s\n Output%s", cmd,out)
	}
	jsonFile,err := os.Open(tmp_path)
	if err != nil {
		log.Fatal(err)
	}
	var traces []Trace
	if err := json.NewDecoder(jsonFile).Decode(&traces); err != nil {
		log.Fatal(err)
	}


	return traces

}
func createIndexIfNotExists(client *elastic.Client){
	// Create an index
	exists, err := client.IndexExists("sinker01").Do()
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create a new index.

		mapping := `{
		    "settings":{
			"number_of_shards":3,
			"number_of_replicas":0
		    },
		    "mappings": {
			 "trace": {
			        "dynamic": "strict",
				"properties": {
					  "ContainerId": {
					    "type": "string"
					  },
					  "ContainerName": {
					    "type": "string"
					  },
					  "evtCpu": {
					    "type": "long"
					  },
					  "evtDir": {
					    "type": "string"
					  },
					  "evtInfo": {
					    "type": "string"
					  },
					  "evtNumber": {
					    "type": "long"
					  },
					  "evtOutputtime": {
					    "type": "date",
					    "format": "strict_date_optional_time||epoch_millis||epoch_second"
					  },
					  "evtType": {
					    "type": "string"
					  },
					  "procName": {
					    "type": "string"
					  },
					  "threadTid": {
					    "type": "long"
					  },
					  "threadVtid": {
					    "type": "long"
					  }
				}
			 }
		    }
		}`
		createIndex, err := client.CreateIndex("sinker01").BodyString(mapping).Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic(err)
		}

	}
}
func putJsonInElastic(trace *Trace, service *elastic.BulkProcessor) {



	// Add a document to the index
	r := elastic.NewBulkIndexRequest().
		Index("sinker01").
		Type("trace").
		Id(fmt.Sprintf("%d_%d_%d_%d", trace.EventNumber, trace.ThreadTid, trace.ThreadVTid, trace.EventOutputUnixTime)).
		Doc(trace)

	service.Add(r)


}

func dealWithFiles(files chan os.FileInfo, service *elastic.BulkProcessor) {

	for f := range files {
		fmt.Printf("reading file %s\n", f.Name())
		traces := extractJson(f)
		for _, value := range traces {
			putJsonInElastic(&value, service)
		}

	}
}
func main() {
	path, err := exec.LookPath("sysdig")
	if err != nil {
		log.Fatal("cannot found sysdig, unable to fetch data")
	}
	fmt.Printf("sysdig is available at %s\n", path)

	files, _ := ioutil.ReadDir("./data/srv/capture/")
	chanFiles := make(chan os.FileInfo)
	elastic.SetURL("http://172.17.0.10:9200")
	client, err := elastic.NewClient()
	createIndexIfNotExists(client)
	if err != nil {
		panic(err)
	}
	// Setup a bulk processor
	service, err := client.BulkProcessor().
	Name("TracesUploader").
	Workers(10).
	BulkActions(1000).              // commit if # requests >= 1000
	BulkSize(2 << 20).              // commit if size of requests >= 2 MB
	FlushInterval(30*time.Second).  // commit every 30s
	Do()
	if err != nil { panic(err) }
	var MAX_WORKERS = 1
	for w := 0; w < MAX_WORKERS; w++ {
		go dealWithFiles(chanFiles, service)

	}
	for _, f := range files {
		fmt.Println("sending file ")
		chanFiles <- f

	}
}
