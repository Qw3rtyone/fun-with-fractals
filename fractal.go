package main

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
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
	if window {
		go showWindow()
	}

	im, err := fractal()
	if err != nil {
		return fmt.Errorf("error while generating image: %v", err)
	}
	frames = append(frames, im)

	outputPath := checkOutputPath(outLoc, "png")
	imageFile, err := os.Create(outputPath)
	if err != nil {
		return errors.New("Failed")
	}

	defer imageFile.Close()
	png.Encode(imageFile, im)

	_, err = getConfirm("quit?")
	return err
}

func Anim() error {
	if window {
		go showWindow()
	}
	imgs := make([]*image.Paletted, frameCount)
	// step := scale
	for i := 0; i < frameCount; i++ {
		im, err := fractal()
		if err != nil {
			return fmt.Errorf("error while generating frame: %v", err)
		}
		frames = append(frames, im)

		imgs[i] = im
		scale = scale + depth
	}

	delay := make([]int, frameCount)
	for i := 0; i < frameCount; i++ {
		delay[i] = 0
	}

	anim := gif.GIF{Delay: delay, Image: imgs}

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
