package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	group = "Polniy"
	song  = "Kaif"
	page  = 1
	limit = 2
)

func main() {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/song?group=%v&song=%v&page=%v&limit=%v",
		group, song, page, limit), nil)
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
