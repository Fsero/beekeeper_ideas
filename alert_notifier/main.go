package main

import (
	//"bufio"
	//"os"
	//"fmt"
	"encoding/json"
	"log"
	"time"
	"fmt"
)

//{"output":"17:20:45.212076717: Alert Shell spawned in a container other than entrypoint)",
// "priority":"Alert","rule":"Run shell in container","time":"2017-02-26T17:20:45.212076717Z"}

type FalcoNotification struct {
	RawOutput         string    `json:"output"`
	Priority          string    `json:"priority"`
	RuleNameTriggered string    `json:"rule"`
	Time              time.Time `json:"time"`
}

func main() {
	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}
	b := []byte(`{"output":"17:20:45.212076717: Alert Shell spawned in a container other than entrypoint)","priority":"Alert","rule":"Run shell in container","time":"2017-02-26T17:20:45.212076717Z"}`)
	var f FalcoNotification
	if err := json.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", f)
	fmt.Printf("\n")
	fmt.Printf(f.Time.String())



}
