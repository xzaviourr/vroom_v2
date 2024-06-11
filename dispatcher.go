package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// Launches an API call on the running workload to execute the request
func dispatch(funcReq *FuncReq, serviceUrl string, variantId string, variantAccuracy float32, logger *Logger) {
	payload := []byte(funcReq.Args) // Send the args as the payload

	errorMessage := "None"
	var resp *http.Response
	failCounter := 1

	req, err := http.NewRequest("POST", serviceUrl, bytes.NewBuffer(payload))
	if err != nil {
		errorMessage = "Error creating a request from payload : " + err.Error()
		funcReq.State = "error"
	} else {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: time.Second * 1200, // Request timeout is 20 minutes
		}

		for failCounter < 60 { // Attempt 60 times before failing the request
			funcReq.SentForExecutionTs = time.Now()
			funcReq.State = "running"
			resp, err = client.Do(req)
			if err != nil {
				failCounter += 1
				time.Sleep(2 * time.Second) // Time between attempts
			} else {
				break
			}
		}
	}

	funcReq.ResponseTs = time.Now()

	if failCounter == 60 {
		errorMessage = "Error in executing the request on the variant : " + err.Error()
	} else if funcReq.State != "error" {
		defer resp.Body.Close()
		funcReq.State = "completed"

		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorMessage = "Error reading response of the: " + err.Error()
			funcReq.State = "error"
		}
	}

	logger.newLog(funcReq, variantId, variantAccuracy, errorMessage, serviceUrl)
}
