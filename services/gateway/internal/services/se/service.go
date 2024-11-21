package ses

import (
	"MusicTask/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Service struct {
	client *http.Client
	addr   string
}

func New(cfg *Config) (*Service, error) {
	client := &http.Client{}
	return &Service{client: client, addr: cfg.Host + cfg.Port}, nil
}

func (s *Service) EnrichSong(group, title string) (models.Song, error) {
	log.Println("Trying to enrich song")
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%v/info?group=%v&song=%v", s.addr, group, title), nil)
	if err != nil {
		log.Println("Failed to create request:", err)
		return models.Song{}, err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		log.Println("Failed to do request:", err)
		return models.Song{}, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body bytes:", err)
		return models.Song{}, err
	}
	params := map[string]string{}
	err = json.Unmarshal(bodyBytes, &params)
	if err != nil {
		log.Println("Failed to unmarshal body:", err)
		return models.Song{}, err
	}
	date, err := time.Parse("02.01.2006", params["releaseDate"])
	if err != nil {
		log.Println("Failed to parse date:", err)
		return models.Song{}, err
	}
	song := models.Song{Group: group, Song: title, ReleaseDate: date, Text: params["text"], Link: params["link"]}
	return song, nil
}
