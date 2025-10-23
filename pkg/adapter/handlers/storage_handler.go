package handlers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetStorageFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Query("path")
		log.Infof("GetStorageFile: %s", path)
		return c.SendFile(url.PathEscape(path))
	}
}
