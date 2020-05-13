package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
)

func getScreenshot(bounds image.Rectangle, display int) *image.RGBA {
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	return img
}

func getCenter() image.Rectangle {
	screen := screenshot.GetDisplayBounds(0)

	centerX := screen.Max.X / 2
	centerY := screen.Max.Y / 2

	screen.Min.X = centerX - 7
	screen.Min.Y = centerY - 7
	screen.Max.X = centerX + 7
	screen.Max.Y = centerY + 7

	return screen
}

func isDiff(pixel1 color.Color, pixel2 color.Color) bool {
	r1, g1, b1, _ := pixel1.RGBA()
	r2, g2, b2, _ := pixel2.RGBA()

	return (r1 != r2 || g1 != g2 || b1 != b2)
}

func watch(toggle chan bool) {
	// get the capture dimenions
	bounds := getCenter()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	last := getScreenshot(bounds, 0)
	now := getScreenshot(bounds, 0)

	isWatching := false

	// watch loop
	for {
		if isWatching {
			now = getScreenshot(bounds, 0)

			diff := 0
			for x := 0; x < width; x++ {
				for y := 0; x < height; x++ {
					pixel1 := last.At(x, y)
					pixel2 := now.At(x, y)

					if isDiff(pixel1, pixel2) {
						diff++
					}
				}
			}

			// has enough changed
			if diff >= 8 {
				fmt.Println("> Movement")
				robotgo.MouseClick("left", true)
				time.Sleep(2700 * time.Millisecond)
				now = getScreenshot(bounds, 0)
			}

			last = now
			time.Sleep(50 * time.Millisecond)
		}

		// toggle watching
		select {
		case <-toggle:
			isWatching = !isWatching
			fmt.Println("> Toggled: ", isWatching)
			last = getScreenshot(bounds, 0)
		default:
			continue
		}
	}
}

func main() {
	fmt.Println("Now watching...")

	// start watching
	toggle := make(chan bool)
	go watch(toggle)

	// setup toggle
	EvChan := hook.Start()
	defer hook.End()

	for event := range EvChan {
		if event.Kind == hook.MouseDown {
			if event.Button == 4 {
				toggle <- true
			}
		}
	}
}
