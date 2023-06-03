package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/manifoldco/promptui"
	"golang.org/x/image/font/gofont/goregular"
)

type rectangle struct {
	Min, Max point
}

type point struct {
	X, Y float64
}

func (r rectangle) Dx() float64 {
	return math.Abs(r.Min.X - r.Max.X)
}
func (r rectangle) Dy() float64 {
	return math.Abs(r.Min.Y - r.Max.Y)
}

func checkOutputPath(s, ext string) string {
	defaultName := fmt.Sprintf("output.%s", ext)
	if filepath.Ext(s) != ext {
		return strings.Join([]string{".", filepath.Dir(s), defaultName}, "/")
	}
	return s
}

func fractal(zoomLevel float64) *image.Paletted {
	window := rectangle{bottomLeft, point{bottomLeft.X + 3.0/zoomLevel, bottomLeft.Y + 3.0/zoomLevel}}
	im := image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)

	hist := make([]int, maxIterations)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			adjustedX := float64(window.Min.X) + float64(x)/float64(width)*float64(window.Dx())
			adjustedY := float64(window.Min.Y) + float64(y)/float64(height)*float64(window.Dy())
			hist[Mandelbrot(adjustedX, adjustedY, maxIterations)-1] += 1
		}
	}
	vals := make([]float64, maxIterations)
	total := 0
	for _, val := range hist {
		total += val
	}
	vals[0] = float64(hist[0]) / float64(total)
	for v := 1; v < len(vals); v++ {
		vals[v] = vals[v-1] + float64(hist[v])/float64(total)
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			adjustedX := float64(window.Min.X) + float64(x)/float64(width)*float64(window.Dx())
			adjustedY := float64(window.Min.Y) + float64(y)/float64(height)*float64(window.Dy())
			val := 255 * vals[Mandelbrot(adjustedX, adjustedY, maxIterations)-1]
			r := uint16(val * 225 * 1.0)
			g := uint16(val * 225 * 1.50)
			b := uint16(val * 225 * 1.50) //uint8(math.Min(val * 255 + 40, 255))
			a := uint16(255)
			if val >= 255*vals[len(vals)-1] {
				r = 0
				g = 0
				b = 0
				a = 0
			}

			// col := color.Gray16{Y: uint16(val * 255)}
			col := color.RGBA64{r, g, b, a}
			im.Set(x, y, col)

		}
	}

	return im
}

// Returns the number of iterations
func Mandelbrot(x, y float64, maxIteration int) int {
	x0 := x
	y0 := y
	x = 0.0
	y = 0.0

	iteration := 0

	for x*x+y*y < 6 && iteration < maxIteration {
		xtemp := x*x - y*y + x0
		y = 2*x*y + y0
		x = xtemp
		iteration += 1
	}
	return iteration
}

func getConfirm(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func showWindow(ready chan bool, finishedFrameCount *int) {
	cfg := pixelgl.WindowConfig{
		Title:  "Now with a UI!",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}

	f, err := os.Open("./output/loading.gif")
	if err != nil {
		panic(err)
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}

	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 25,
	})
	txt := text.New(win.GetPos(), text.NewAtlas(face, text.ASCII))

	b := false
	for !win.Closed() {
		select {
		case <-ready:
			b = true
		default:

		}
		if b {
			updateWindow(win, frames, nil, "", 16)
		} else {
			updateWindow(win, g.Image, txt, fmt.Sprintf("%d/%d Frames Finished.", *finishedFrameCount, frameCount), 0)
		}
	}
}

func updateWindow(win *pixelgl.Window, f []*image.Paletted, txt *text.Text, s string, sleep int) {
	for _, img := range f {
		win.Clear(color.Black)
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())

		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if txt != nil {
			txt.Clear()
			_, _ = txt.WriteString(s)
			txt.Draw(win, pixel.IM)
		}

		win.Update()
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}

}
