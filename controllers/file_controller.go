package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

func UploadFile(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	files := form.File["file"]

	urls := make([]string, 0, len(files))

	for _, file := range files {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		fileData, err := io.ReadAll(fileBytes)
		if err != nil {
			fileBytes.Close()
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		ext := filepath.Ext(file.Filename)
		filePath := "./uploads/" + uuid.New().String() + ext

		err = os.WriteFile(filePath, fileData, 0666)
		if err != nil {
			fileBytes.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fileBytes.Close()
		urls = append(urls, "http://localhost:8888/"+filePath)
	}
	return c.JSON(fiber.Map{
		"urls": urls,
	})
}
