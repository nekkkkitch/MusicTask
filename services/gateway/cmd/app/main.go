package main

import (
	"MusicTask/services/gateway/internal/db"
	"MusicTask/services/gateway/internal/router"
	mls "MusicTask/services/gateway/internal/services/ml"
	ses "MusicTask/services/gateway/internal/services/se"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB  *db.Config     `yaml:"db"`
	RTR *router.Config `yaml:"rtr"`
	SES *ses.Config    `yaml:"ses"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := readConfig("./cfg.yml")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Config file read successfully")
	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB connected successfully")
	ses, err := ses.New(cfg.SES)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Song Enricher Service connected successfully")
	mls := mls.New(db, ses)
	log.Println("Music Library Service connected successfully")
	rtr, err := router.New(cfg.RTR, mls)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Router created")
	err = rtr.Listen()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Gateway started successfully")
}
