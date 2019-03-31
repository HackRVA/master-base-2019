package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/HackRVA/master-base-2019/filelogging"
)

var logger = log.Ger

func main() {
	fifteenMin := time.Now().Local().Add(time.Minute * time.Duration(15))
	url := "http://localhost:3000/api/newgame"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(fmt.Sprintf(`{
		 	"body":123,
		 	"AbsStart": %d,
		 	"Duration": 13,
		 	"Variant": 1
			 }`, fifteenMin.Unix()))

	fmt.Println(fifteenMin.Unix())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
