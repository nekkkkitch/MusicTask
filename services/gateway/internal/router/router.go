package router

import (
	"MusicTask/pkg/models"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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

func New(cfg *Config, ml IMLService) (*Router, error) {
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
