package imgprcss

import (
	"image"
	"image/color"
)

////////////////////////////////////////////////////////////// creates a copy of the image, and returns the copy and its dimensions
func copyImage(img image.Image) (*image.RGBA, int, int) {
	copy := image.NewRGBA(img.Bounds())
	x_max := copy.Bounds().Max.X
	y_max := copy.Bounds().Max.Y
	for y := 0; y < y_max; y++ {
		for x := 0; x < x_max; x++ {
			copy.Set(x, y, img.At(x, y))
		}
	}
	return copy, x_max, y_max
}

////////////////////////////////////////////////////////////// decides if the new pixel is black or colored.
////////////////////////////////////////////////////////////// can be altered to use more than 2 colors
func pixelDecide(pixel color.RGBA, dither_color *(color.RGBA)) color.RGBA {
	if GrayscalePixel(pixel) >= 128 {
		return *dither_color
	} else {
		return color.RGBA{0, 0, 0, 255} // change this color to see different results
	}
}

////////////////////////////////////////////////////////////// does the kernel convolution over a given pixel (x, y) of an image
func kernelConvolution(img *image.Image, x, y int, kernel *[3][3]float32) (uint8, uint8, uint8) {

	var SumR, SumG, SumB float32 = 0, 0, 0

	for j := 0; j <= 2; j++ {
		for i := 0; i <= 2; i++ {
			col := color.RGBAModel.Convert((*img).At(x-1+i, y-1+j)).(color.RGBA)

			SumR += (*kernel)[i][j] * float32(col.R)
			SumG += (*kernel)[i][j] * float32(col.G)
			SumB += (*kernel)[i][j] * float32(col.B)
		}
	}

	if SumR < 0 {
		SumR = 0
	} else if SumR > 255 {
		SumR = 255
	}

	if SumG < 0 {
		SumG = 0
	} else if SumG > 255 {
		SumG = 255
	}

	if SumB < 0 {
		SumB = 0
	} else if SumB > 255 {
		SumB = 255
	}
	return uint8(SumR), uint8(SumG), uint8(SumB)
}
