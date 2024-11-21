package main

import (
	"MusicTask/pkg/models"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//Пришлось сделать вот такую копию рутера из-за того, что библиотека go-swagger умеет читать только main.go файлы :(
// P.S. Данный файл никоим образом не участвует в работе программы 

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Router struct {
	Config    *Config
	App       *fiber.App
	MLService IMLService
}

type IMLService interface {
	GetLib(page, limit int) ([]models.Song, error)
	GetSong(group, title string, page, limit int) (string, error)
	DeleteSong(group, title string) error
	ChangeSong(group, title string, params map[string]string) error
	AddSong(group, title string) error
}

// @title Music Library API
// @version 1.0
// @host localhost:8082
// @schemes http
// @BasePath /
func main(cfg *Config, ml IMLService) (*Router, error) {
	app := fiber.New()
	router := Router{App: app, MLService: ml, Config: cfg}
	router.App.Get("/lib", router.GetLib())
	router.App.Get("/song", router.GetSong())
	router.App.Delete("/song", router.DeleteSong())
	router.App.Put("/song", router.ChangeSong())
	router.App.Post("/song", router.AddSong())
	return &router, nil
}

func (r *Router) Listen() error {
	err := r.App.Listen(r.Config.Host + r.Config.Port)
	return err
}

// @Description Возвращает список песен с пагинацией
// @Param limit query int true "Songs per page"
// @Param page query int true "Page"
// @produce json
// @Success 200
// @Router /lib [get]
func (r *Router) GetLib() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Got page: %v, limit: %v\n", c.Queries()["page"], c.Queries()["limit"])
		page, err := strconv.Atoi(c.Queries()["page"])
		if err != nil || page < 1 {
			log.Println("Wrong page type")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("Bad page request")
			return nil
		}
		limit, err := strconv.Atoi(c.Queries()["limit"])
		if err != nil || page < 1 {
			log.Println("Wrong limit type")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("Bad limit request")
			return nil
		}
		toRet, err := r.MLService.GetLib(page, limit)
		if err != nil {
			log.Println("Cant get library")
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		return c.JSON(toRet)
	}
}

// @Description Возвращает текст песни по куплетам
// @Param song query string true "Song title"
// @Param group query string true "Group"
// @Param limit query int true "Couplets per page"
// @Param page query int true "Page"
// @Success 200
// @Router /song [get]
func (r *Router) GetSong() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Got group: %v, song: %v, page: %v, limit: %v\n", c.Queries()["group"], c.Queries()["song"], c.Queries()["page"], c.Queries()["limit"])
		song := c.Queries()["song"]
		group := c.Queries()["group"]
		if song == "" {
			log.Println("Empty title")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No song given")
			return nil
		}
		if group == "" {
			log.Println("Empty group")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No group given")
			return nil
		}
		page, err := strconv.Atoi(c.Queries()["page"])
		if err != nil || page < 1 {
			log.Println("Wrong page type")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("Bad page request")
			return nil
		}
		limit, err := strconv.Atoi(c.Queries()["limit"])
		if err != nil || page < 1 {
			log.Println("Wrong limit type")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("Bad limit request")
			return nil
		}
		toRet, err := r.MLService.GetSong(group, song, page, limit)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		_, err = c.WriteString(toRet)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}

// @Description Удалить существующую песню
// @Param song query string true "Song title"
// @Param group query string true "Group"
// @Success 200
// @Router /song [delete]
func (r *Router) DeleteSong() fiber.Handler {
	return func(c *fiber.Ctx) error {
		song := c.Queries()["song"]
		group := c.Queries()["group"]
		if song == "" {
			log.Println("Empty title")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No song given")
			return nil
		}
		if group == "" {
			log.Println("Empty group")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No group given")
			return nil
		}
		err := r.MLService.DeleteSong(group, song)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.WriteString("Song deleted successfully")
		c.Status(fiber.StatusOK)
		return nil
	}
}

// @Description Изменить существующую песню
// @Param song query string true "Song title"
// @Param group query string true "Group"
// @Param changes body string true "you can change releaseDate, text or link"
// @Success 200
// @accepts json
// @Router /song [put]
func (r *Router) ChangeSong() fiber.Handler {
	return func(c *fiber.Ctx) error {
		song := c.Queries()["song"]
		group := c.Queries()["group"]
		if song == "" {
			log.Println("Empty title")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No song given")
			return nil
		}
		if group == "" {
			log.Println("Empty group")
			c.Status(fiber.StatusBadRequest)
			c.WriteString("No group given")
			return nil
		}
		params := map[string]string{}
		err := json.Unmarshal(c.Body(), &params)
		if err != nil {
			log.Println("Cant unmarshal body:", err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		log.Println("Got params:", params)
		err = r.MLService.ChangeSong(group, song, params)
		if err != nil {
			log.Println("Cant change song:", err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.WriteString("Song changed successfully")
		c.Status(fiber.StatusOK)
		return nil
	}
}

// @Description Добавление песни
// @Success 200
// @accepts json
// @Param song body string true "Need group and song"
// @Router /song [post]
func (r *Router) AddSong() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := map[string]string{}
		err := json.Unmarshal(c.Body(), &params)
		if err != nil {
			log.Println("Cant unmarshal body:", err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		log.Println("Got params:", params)
		err = r.MLService.AddSong(params["group"], params["song"])
		if err != nil {
			log.Println("Cant add song:", err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.WriteString("Song added successfully")
		c.Status(fiber.StatusOK)
		return nil
	}
}
