package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	imp "impala/imgprcss"
	"os"
)

func main() {
	//////////////////////////////////////// loads the image to be processed
	command := os.Args[1]
	imagePath := os.Args[2]

	img := imp.LoadImage(imagePath)

	//////////////////////////////////////// process the image
	var new_img *image.RGBA

	if command == "gray" {
		fmt.Println("making image grayscale...")
		new_img = image.NewRGBA(img.Bounds())
		x_max := new_img.Bounds().Max.X
		y_max := new_img.Bounds().Max.Y

		for y := 0; y < y_max; y++ {
			for x := 0; x < x_max; x++ {
				num := imp.GrayscalePixel(img.At(x, y))
				new_img.Set(x, y, color.RGBA{num, num, num, 255})
			}
		}
	} else if command == "dither" {
		fmt.Println("making an error diffusion dithering out of your image...")
		new_img = imp.ErrorDiffusionDithering(&img)
	} else if command == "blur" {
		fmt.Println("blurring image...")
		new_img = imp.KCblur(&img)
	} else {
		fmt.Println("failed to do anything.")
		return
	}

	//////////////////////////////////////// saves the processed image
	fmt.Println("saving image...")
	imp.SaveImage(&new_img, "neW_image.png")
}
