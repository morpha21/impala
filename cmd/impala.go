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
	//////////////////////////////////////// defines the path of the image and the effect to be applied
	command := os.Args[1]
	imagePath := os.Args[2]
	//////////////////////////////////////// loads the image
	img := imp.LoadImage(imagePath)

	var new_img *image.RGBA
	//////////////////////////////////////// decides what to do
	if command == "gray" {
		fmt.Println("making image grayscale...")

		//////////////////////////////////////// defines image dimensions
		new_img = image.NewRGBA(img.Bounds())
		x_max := new_img.Bounds().Max.X
		y_max := new_img.Bounds().Max.Y
		//////////////////////////////////////// creates the new image
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
		new_img = imp.GaussianBlur(&img)
	} else if command == "sharpen" {
		fmt.Println("sharpening image...")
		new_img = imp.Sharpen(&img)
	} else {
		fmt.Println("failed to do anything.")
		return
	}

	//////////////////////////////////////// saves the processed image
	fmt.Println("saving image...")
	imp.SaveImage(&new_img, command+"_image.png")
	fmt.Println("done.")
	fmt.Println()
}
