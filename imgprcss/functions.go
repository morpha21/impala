package imgprcss

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

// loads an image
func LoadImage(imagePath string) image.Image {
	// opens the image
	f, err := os.Open(imagePath)
	if err != nil {
		err = fmt.Errorf("failed to open file: %w", err)
		log.Fatal(err)
	}
	defer f.Close()

	// decodes the image
	image, _, err := image.Decode(f)
	if err != nil {
		err = fmt.Errorf("failed to decode image: %w", err)
		log.Fatal(err)
	}
	return image
}

// saves an image
func SaveImage(img **image.RGBA, name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, *img)
	if err != nil {
		log.Fatal(err)
	}
}
