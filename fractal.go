package main

import (
	"errors"
	"image"
	"image/gif"
	"image/png"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	width  int //  = 1020
	height int //= 660
	window bool

	outLoc        string
	scale         float64 //      = 67.0
	frameCount    int
	depth         float64
	maxIterations int
	bottomLeft    = point{-0.535, 0.586}
)

var frames []*image.Paletted

func Image() error {
	c := make(chan bool)
	if window {
		go showWindow(c)
	}

	frames = make([]*image.Paletted, 1)
	im := fractal(scale)
	frames[0] = im

	c <- true

	outputPath := checkOutputPath(outLoc, "png")
	imageFile, err := os.Create(outputPath)
	if err != nil {
		return errors.New("Failed")
	}

	defer imageFile.Close()
	png.Encode(imageFile, frames[0])

	_, err = getConfirm("quit?")
	return err
}

func Anim() error {
	c := make(chan bool)
	if window {
		go showWindow(c)
	}

	frames = make([]*image.Paletted, frameCount)

	g := new(errgroup.Group)
	g.SetLimit(30)
	// check that the errors returned by the goroutines are all nil
	for i := 0; i < frameCount; i++ {
		i := i
		g.Go(func() error {
			im := fractal(scale + (depth * float64(i)))

			frames[i] = im
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	c <- true

	delay := make([]int, frameCount)
	for i := 0; i < frameCount; i++ {
		delay[i] = 0
	}

	anim := gif.GIF{Delay: delay, Image: frames}

	outputPath := checkOutputPath(outLoc, "gif")
	imageFile, err := os.Create(outputPath)
	if err != nil {
		return errors.New("Failed")
	}

	defer imageFile.Close()

	err = gif.EncodeAll(imageFile, &anim)
	if err != nil {
		return errors.New("Failed")
	}
	_, err = getConfirm("quit?")
	return err
}
