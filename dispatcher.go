package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func dispatch(serviceUrl string, args string, responseUrl string) {
	payload := []byte(args)

	req, err := http.NewRequest("POST", serviceUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating a request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 2000,
	}

	var resp *http.Response
	failCounter := 1

	for failCounter < 60 {
		resp, err = client.Do(req)
		if err != nil {
			failCounter += 1
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	if failCounter == 60 {
		fmt.Println("Error making request", err)
		return
	}

	fmt.Println("Request successful")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
