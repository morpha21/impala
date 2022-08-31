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

// does the kernel convolution over a given pixel (x, y) of an image
func kernelConvolutedPixel(img *image.Image, x, y int) color.Color {
	var blur_kernel [3][3]int
	blur_kernel[0][0], blur_kernel[0][2], blur_kernel[2][0], blur_kernel[2][2] = 1, 1, 1, 1
	blur_kernel[0][1], blur_kernel[1][0], blur_kernel[1][2], blur_kernel[2][1] = 2, 2, 2, 2
	blur_kernel[1][1] = 4

	kernelSum := 0

	for i := 0; i < 9; i++ {
		kernelSum += blur_kernel[i/3][i%3]
	}
	SumR, SumG, SumB := 0, 0, 0

	for j := 0; j <= 2; j++ {
		for i := 0; i <= 2; i++ {
			col := color.RGBAModel.Convert((*img).At(x-1+i, y-1+j)).(color.RGBA)

			SumR += blur_kernel[i][j] * int(col.R)
			SumG += blur_kernel[i][j] * int(col.G)
			SumB += blur_kernel[i][j] * int(col.B)
		}
	}

	R := uint8(float32(SumR) / float32(kernelSum))
	G := uint8(float32(SumG) / float32(kernelSum))
	B := uint8(float32(SumB) / float32(kernelSum))

	new_pixel := color.RGBA{R, G, B, 255}
	return new_pixel
}
