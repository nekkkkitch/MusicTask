package mls

import (
	"MusicTask/pkg/models"
	"fmt"
	"log"
	"strings"
)

type Service struct {
	db        IDB
	seService ISEService
}

type IDB interface {
	GetLib(page, limit int) ([]models.Song, error)
	GetSong(group, title string) (models.Song, error)
	DeleteSong(group, title string) error
	ChangeSong(group, title string, params map[string]string) error
	AddSong(song models.Song) error
}

type ISEService interface {
	EnrichSong(group, song string) (models.Song, error)
}

func New(db IDB, se ISEService) *Service {
	s := Service{db: db, seService: se}
	return &s
}

func (s *Service) GetLib(page, limit int) ([]models.Song, error) {
	songs, err := s.db.GetLib(page, limit)
	if err != nil {
		log.Println("Failed to get songs:", err)
		return nil, err
	}
	return songs, nil
}

func (s *Service) GetSong(group, title string, page, limit int) (string, error) {
	song, err := s.db.GetSong(group, title)
	if err != nil {
		log.Println("Failed to get song:", err)
		return "", err
	}
	partedSong := strings.Split(song.Text, "\n\n")
	text := ""
	page -= 1
	for i := range limit {
		if page*limit >= len(partedSong) {
			log.Println("Failed to return text")
			return "", fmt.Errorf("page limit for this song exceed")
		}
		if page*limit+i >= len(partedSong) {
			break
		}
		text += partedSong[page*limit+i]
		if i != limit-1 {
			text += "\n\n"
		}
	}
	return text, nil
}

func (s *Service) DeleteSong(group, title string) error {
	err := s.db.DeleteSong(group, title)
	if err != nil {
		log.Println("Failed to delete song:", err)
		return err
	}
	return nil
}

func (s *Service) ChangeSong(group, title string, params map[string]string) error {
	err := s.db.ChangeSong(group, title, params)
	if err != nil {
		log.Println("Failed to change song:", err)
		return err
	}
	return nil
}

func (s *Service) AddSong(group, title string) error {
	song, err := s.seService.EnrichSong(group, title)
	if err != nil {
		log.Println("Failed to enrich song:", err)
		song = models.Song{Group: group, Song: title}
	}
	err = s.db.AddSong(song)
	if err != nil {
		log.Println("Failed to add song:", err)
		return err
	}
	return nil
}
