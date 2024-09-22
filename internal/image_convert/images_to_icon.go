package image_convert

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"rory-pearson/environment"
	"rory-pearson/pkg/util"

	ico "github.com/Kodeworks/golang-image-ico"
	"golang.org/x/image/draw"
)

var storageDirectory = environment.CreateStorageDirectory("image_convert_storage")

var Sizes = []image.Point{
	{16, 16},
	{24, 24},
	{32, 32},
	{48, 48},
	{64, 64},
	{128, 128},
	{256, 256},
}

// Convert converts the image into multiple icon sizes and compresses them into a zip file.
func Convert(imagePath string) (string, error) {
	outputDir := util.GenerateUUIDv4()
	dirPath := filepath.Join(storageDirectory, outputDir)

	// Create the directory for storing converted icons
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("could not create output directory: %v", err)
	}

	// Open and decode the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("could not open image file: %v", err)
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return "", fmt.Errorf("could not decode image: %v", err)
	}

	// Resize and save the icons
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

	// Compress the directory into a zip file and remove the icons folder
	zipName := outputDir + ".zip"
	if err := util.CompressDirectoryAndDelete(dirPath, storageDirectory, zipName); err != nil {
		return "", fmt.Errorf("could not compress directory: %v", err)
	}

	return zipName, nil
}

// GetConvertedZipPath returns the path to the converted zip file and optionally deletes it.
func GetConvertedZipPath(zipName string) (string, error) {
	// Ensure the zipName ends with .zip
	if filepath.Ext(zipName) != ".zip" {
		zipName += ".zip"
	}

	zipPath := filepath.Join(storageDirectory, zipName)
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %v", err)
	}

	return zipPath, nil
}

func DeleteConvertedZip(zipName string) error {
	if err := os.Remove(filepath.Join(storageDirectory, zipName)); err != nil {
		return fmt.Errorf("could not delete zip file: %v", err)
	}

	return nil
}

// Utility functions for image conversion and resizing
func convertToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

func resizeImage(img image.Image, size image.Point) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Src, nil)
	return resizedImg
}
