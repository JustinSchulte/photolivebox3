package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ImageResponse struct {
	Files       []string `json:"files"`
	NewLastDate string   `json:"newLastDate"`
}

var imagePath = filepath.Join("backend", "images")

func ListImages(c echo.Context) error {
	images, err := os.ReadDir(imagePath)
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
	filepath := filepath.Join(imagePath, filename)
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
	timestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	filename := fmt.Sprintf("%v_%v", timestamp, file.Filename)
	fullFilepath := filepath.Join(imagePath, filename)
	dst, err := os.Create(fullFilepath)
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
