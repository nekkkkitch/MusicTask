package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{}
	body, _ := json.Marshal(map[string]string{"group": "Polniy", "song": "Kaif"})
	r := bytes.NewReader(body)
	request, err := http.NewRequest("POST", "http://localhost:8082/song", r)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
}
