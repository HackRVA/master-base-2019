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
     	// fifteenMin := time.Now().Local().Add(time.Minute * time.Duration(15))
	fifteenMin := time.Now().Local().Add(time.Duration(1))
	url := "http://localhost:8000/api/newgame"

	var jsonStr = []byte(fmt.Sprintf(`{
		 	"body":123,
		 	"AbsStart": %d,
		 	"Duration": 120,
		 	"Variant": 0
			 }`, fifteenMin.Unix()))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Info().Msgf("response Status: %s", resp.Status)
	logger.Info().Msgf("response Headers: %s", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Info().Msgf("response Body: %s", string(body))
}
