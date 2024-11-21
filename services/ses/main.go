package main

import (
	"github.com/gofiber/fiber/v2"
)

// Это заглушка, не воспринимайте её всерьез пожалуйста😌
func main() {
	app := fiber.New()
	app.Get("/info", Enrich)
	app.Listen(":8083")
}

func Enrich(c *fiber.Ctx) error {
	based := map[string]string{"releaseDate": "16.07.2006", "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		"link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"}
	return c.JSON(based)
}
