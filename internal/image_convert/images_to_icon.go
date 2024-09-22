package image_convert

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"rory-pearson/pkg/util"

	ico "github.com/Kodeworks/golang-image-ico"
	"golang.org/x/image/draw"
)

const storageDirectory = "image_convert_storage/"

// SupportedImageType represents the type of image supported
type SupportedImageType int

const (
	PNG SupportedImageType = iota
	JPEG
	WEBP
)

var Sizes = []image.Point{
	{16, 16},
	{24, 24},
	{32, 32},
	{48, 48},
	{64, 64},
	{128, 128},
	{256, 256},
}

func Convert(imagePath string) (string, error) {
	// Generate a unique UUID for the directory
	outputDir := util.GenerateUUIDv4()
	dirPath := filepath.Join(storageDirectory, outputDir)

	// Create the output directory
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// Open and decode the image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("could not open image file: %v", err)
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return "", fmt.Errorf("could not decode image: %v", err)
	}

	// Process and save each resized image
	for _, size := range Sizes {
		resizedImg := resizeImage(img, size)
		rgbaImg := convertToRGBA(resizedImg)

		icoFilePath := filepath.Join(dirPath, fmt.Sprintf("%dx%d.ico", size.X, size.Y))
		icoFile, err := os.Create(icoFilePath)
		if err != nil {
			return "", fmt.Errorf("could not create output file: %v", err)
		}

		if err := ico.Encode(icoFile, rgbaImg); err != nil {
			icoFile.Close()
			return "", fmt.Errorf("could not encode image: %v", err)
		}
		icoFile.Close()
	}

	// Compress the directory and clean up
	zipName := outputDir + ".zip"
	if err := util.CompressDirectoryAndDelete(dirPath, storageDirectory, zipName); err != nil {
		return "", fmt.Errorf("could not compress directory: %v", err)
	}

	return zipName, nil
}

func GetConvertedZipPath(zipName string, delete bool) (string, error) {
	zipPath := filepath.Join(storageDirectory, zipName)
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %v", err)
	}

	if delete {
		if err := os.Remove(zipPath); err != nil {
			return "", fmt.Errorf("could not delete file: %v", err)
		}
	}

	return zipPath, nil
}

// convertToRGBA converts the image to RGBA format
func convertToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

// resizeImage resizes the image to the specified size
func resizeImage(img image.Image, size image.Point) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Src, nil)
	return resizedImg
}
