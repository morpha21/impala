impala is a Go image processing CLI tool I'm working on.

It currently has an implementation of the Floyd-Steinberg Dithering algorithm, which checks the light pixels of the image, calculates a mean color between those and does the dithering with two colors: black and that mean color. It also has a simple Gaussian Blur, and a function to turn images into grayscale. 

# Usage: 

Currently, it nees two arguments: the command, which can be either "dither", "blur" or "gray"; and the image path. 

For example, if you have an image /home/me/Pictures/MyPicture.png and wants to do the dithering, you need to: 

```go run cmd/impala.go dither /home/me/Pictures/MyPicture.png```



Please, feel free to reach me out to contribute, give some ideas, ask for help, offer me a job or teach me something. 