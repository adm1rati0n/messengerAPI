package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

type fileResponse struct {
	FileName string `json:"file_name"`
}

func UploadFile(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	files := form.File["file"]

	var filesResponse []fileResponse

	for _, file := range files {
		var fileResponse fileResponse
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
		fileName := uuid.New().String() + ext
		filePath := "./uploads/" + fileName

		err = os.WriteFile(filePath, fileData, 0666)
		if err != nil {
			fileBytes.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fileResponse.FileName = fileName
		fileBytes.Close()
		filesResponse = append(filesResponse, fileResponse)
	}
	return c.Status(fiber.StatusOK).JSON(filesResponse)
}
