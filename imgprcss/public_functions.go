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

////////////////////////////////////////////////////////////// calculates the mean value of R, G and B of a given color
func GrayscalePixel(c color.Color) uint8 {
	pixel_color := color.RGBAModel.Convert(c).(color.RGBA)
	r, g, b := pixel_color.R, pixel_color.G, pixel_color.B

	return uint8((float32(r) + float32(g) + float32(b)) / 3)
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
