package db

import (
	"MusicTask/pkg/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type DB struct {
	Config *Config
	db     *pgx.Conn
}

func New(cfg *Config) (*DB, error) {
	d := &DB{Config: cfg}
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	log.Println("Connecting to DB on", connection)
	db, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Println("Failed to connect DB:", err)
		return nil, err
	}
	d.db = db
	return d, nil
}

func (d *DB) GetLib(page, limit int) ([]models.Song, error) {
	log.Printf("Trying to get songs from db, page %v, limit %v\n", page, limit)
	row, err := d.db.Query(context.Background(), `select "group", song, releasedate, "text", link from public.songs order by id limit $1 offset $2`,
		limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	songs := make([]models.Song, 0, limit)
	for row.Next() {
		var group pgtype.Text
		var song pgtype.Text
		var release_date pgtype.Date
		var song_text pgtype.Text
		var link pgtype.Text
		err := row.Scan(&group, &song, &release_date, &song_text, &link)
		if err != nil {
			log.Println("Failed to scan song:", err)
			return nil, err
		}
		songs = append(songs, models.Song{Group: group.String, Song: song.String, ReleaseDate: release_date.Time, Text: song_text.String, Link: link.String})
	}
	log.Println("Successfully got songs from db")
	return songs, nil
}

func (d *DB) GetSong(group, title string) (models.Song, error) {
	log.Println("Trying to get song titled", title)
	var song_text pgtype.Text
	err := d.db.QueryRow(context.Background(), `select "text" from public.songs where song = $1 and "group" = $2`,
		title, group).Scan(&song_text)
	if err != nil {
		log.Println("Failed to scan song:", err)
		return models.Song{}, err
	}
	log.Println("Successfully got song text from db")
	return models.Song{Text: song_text.String}, nil
}

func (d *DB) DeleteSong(group, title string) error {
	log.Println("Trying to delete song titled", title)
	cmd, err := d.db.Exec(context.Background(), `delete from public.songs where song = $1 and "group" = $2`, title, group)
	if err != nil {
		log.Println("Failed to delete song:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("No songs were deleted")
		return fmt.Errorf("no song with such group and title")
	}
	log.Println("Successfully deleted song")
	return nil
}

func (d *DB) ChangeSong(group, title string, params map[string]string) error {
	log.Printf("Trying to change song titled %v with parameters: %v", title, params)
	for key, val := range params {
		if val == "" {
			continue
		}
		if key == "releaseDate" {
			date, err := time.Parse("02.01.2006", val)
			if err != nil {
				log.Println("Failed to parse time:", err)
				return err
			}
			rowsAffected, err := d.db.Exec(context.Background(), fmt.Sprintf(`update public.songs set "%v" = $1 where song = $2 and "group" = $3`,
				strings.ToLower(key)), date, title, group)
			if err != nil {
				log.Printf("Failed to update song %v:%v\n", key, err)
				return err
			}
			if rowsAffected.RowsAffected() == 0 {
				log.Println("No rows affected therefore there is no song", group, title)
				return fmt.Errorf("no songs with such name and group")
			}
			continue
		}
		rowsAffected, err := d.db.Exec(context.Background(), fmt.Sprintf(`update public.songs set "%v" = $1 where song = $2 and "group" = $3`,
			strings.ToLower(key)), val, title, group)
		if err != nil {
			log.Printf("Failed to update song %v:%v\n", key, err)
			return err
		}
		if rowsAffected.RowsAffected() == 0 {
			log.Println("No rows affected therefore there is no song", group, title)
			return fmt.Errorf("no songs with such name and group")
		}
	}
	log.Println("Song successfully changed")
	return nil
}

func (d *DB) AddSong(song models.Song) error {
	log.Println("Trying to add new song:", song)
	if d.SongInDB(song.Group, song.Song) {
		log.Println("Song already in DB")
		return fmt.Errorf("song with such group and title already in db")
	}
	d.TryAddGroup(song.Group)
	_, err := d.db.Exec(context.Background(), `insert into public.songs("group", song, releasedate, "text", link) values($1, $2, $3, $4, $5)`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		log.Println("Failed to add new song:", err)
		return err
	}
	log.Println("Song succesfully added")
	return nil
}

func (d *DB) TryAddGroup(group string) error {
	log.Println("Trying to add new group:", group)
	_, err := d.db.Exec(context.Background(), `insert into public.groups(title) values($1) returning id`, group)
	if err != nil {
		log.Println("Failed to add new group")
		return err
	}
	log.Println("Group successfully added")
	return nil
}

func (d *DB) SongInDB(group, song string) bool {
	log.Println("Checking if song is in db")
	var id pgtype.Int4
	_ = d.db.QueryRow(context.Background(), `select id from public.songs where "group"=$1 and song=$2`, group, song).Scan(&id)
	log.Println(id)
	return id.Int32 > 0
}
