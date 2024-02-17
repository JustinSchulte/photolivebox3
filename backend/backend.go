package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type ImageResponse struct {
	Files       []string `json:"files"`
	NewLastDate string   `json:"newLastDate"`
}

func ListImages(c echo.Context) error {
	images, err := os.ReadDir("./backend/images")
	if err != nil {
		fmt.Printf("cant read images: %v\n", err)
	}
	imageList := []string{}
	for _, img := range images {
		imageList = append(imageList, img.Name())
	}
	return c.JSON(http.StatusOK, imageList)
}

func GetImage(c echo.Context) error {
	filename := c.Param("filename")
	filepath := fmt.Sprintf("./backend/images/%v", filename)
	return c.File(filepath)
}

func UploadImage(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("Error reading form file: %v\n", err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		fmt.Printf("Error open file: %v\n", err)
		return err
	}
	defer src.Close()

	// Destination
	timestamp := time.Now().Format(time.RFC3339)
	filename := fmt.Sprintf("./backend/images/%v_%v", timestamp, file.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error create new file: %v\n", err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Printf("Error copy file: %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, true)
}
