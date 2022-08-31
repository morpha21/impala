package imgprcss

import (
	"image"
	"image/color"
	_ "image/jpeg"
)

func ErrorDiffusionDithering(img *image.Image) *image.RGBA {
	dither, x_max, y_max := copyImage(*img)

	/////////////////////////////////////////////  calculates the mean light color of the image
	/////////////////////////////////////////////  to be used to do the dither
	cnt := 0
	sumR, sumG, sumB := 0, 0, 0
	for y := 0; y < y_max; y++ {
		for x := 0; x < x_max; x++ {
			col := color.RGBAModel.Convert(dither.At(x, y)).(color.RGBA)

			if GrayscalePixel(col) >= 128 { // decides which pixels will be considered in the mean color
				sumR += int(col.R)
				sumG += int(col.G)
				sumB += int(col.B)
				cnt++
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////// defines the color to be used
	dither_color := color.RGBA{uint8(sumR / cnt), uint8(sumG / cnt), uint8(sumB / cnt), 255}

	/////////////////////////////////////////////  does the error diffusion dithering
	for y := 0; y < y_max; y++ {
		for x := 0; x < x_max; x++ {

			old_pixel := color.RGBAModel.Convert(dither.At(x, y)).(color.RGBA)
			new_pixel := pixelDecide(old_pixel, &dither_color)
			dither.Set(x, y, new_pixel)

			////////////////////////////////////////////////// allows the error to be calculated considering the white color,
			////////////////////////////////////////////////// so the dithering corresponds to the expected behavior, even using a color different from white
			if new_pixel == dither_color {
				new_pixel = color.RGBA{255, 255, 255, 255}
			}

			pixel_error := int16(GrayscalePixel(old_pixel)) - int16(GrayscalePixel(new_pixel))

			var DL, D, DR, R uint8
			////////////////////////////////////////////////// Calculates the new values of the pixels around the pixel at (x, y),
			DL = GrayscalePixel(dither.At(x-1, y+1)) + uint8(float32(pixel_error)*3.0/16)
			D = GrayscalePixel(dither.At(x, y+1)) + uint8(float32(pixel_error)*5.0/16)
			DR = GrayscalePixel(dither.At(x+1, y+1)) + uint8(float32(pixel_error)/16)
			R = GrayscalePixel(dither.At(x+1, y)) + uint8(float32(pixel_error)*7.0/16)

			////////////////////////////////////////////////// updates the values of the pixels around (x,y)
			dither.Set(x-1, y+1, color.RGBA{DL, DL, DL, 255})
			dither.Set(x, y+1, color.RGBA{D, D, D, 255})
			dither.Set(x+1, y+1, color.RGBA{DR, DR, DR, 255})
			dither.Set(x+1, y, color.RGBA{R, R, R, 255})
		}
	}
	return dither
}

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

////////////////////////////////////////////////////////////// calculates the mean value of R, G and B of a given color
func GrayscalePixel(c color.Color) uint8 {
	pixel_color := color.RGBAModel.Convert(c).(color.RGBA)
	r, g, b := pixel_color.R, pixel_color.G, pixel_color.B

	return uint8((float32(r) + float32(g) + float32(b)) / 3)
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

////////////////////////////////////////////////////////////// does the
func GaussianBlur(img *image.Image) *image.RGBA {
	blurred := image.NewRGBA((*img).Bounds())
	x_max := blurred.Bounds().Max.X
	y_max := blurred.Bounds().Max.Y

	for y := 0; y < y_max; y++ {
		for x := 0; x < x_max; x++ {
			for i := 0; i < 9; i++ {
				blurred.Set(x, y, kernelConvolutedPixel(img, x, y))
			}

		}
	}
	return blurred
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
