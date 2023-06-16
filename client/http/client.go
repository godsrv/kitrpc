package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func main() {
	data := map[string]string{"key": "lipper", "val": "30"}
	jsonData, err := json.Marshal(data)
	prep, err := http.NewRequest("POST", "http://127.0.0.1:7001/set", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	prep.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	presp, err := client.Do(prep)
	var res response

	if err = json.NewDecoder(presp.Body).Decode(&res); err != nil {
		log.Fatalf("json.NewDecoder: %v", err)
	}
	log.Println(res)

}
