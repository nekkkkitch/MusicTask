package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	group = "Polniy"
	song  = "Kaif"
)

func main() {
	client := &http.Client{}
	body, _ := json.Marshal(map[string]string{"releaseDate": "16.09.2015", "text": "123123\n123321\n2222222\n545454\n\n123123123\n111"})
	r := bytes.NewReader(body)
	request, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8082/song?group=%v&song=%v", group, song), r)
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
