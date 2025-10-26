package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

const reset = "\x1b[0m"

var (
	xStep, yStep *int
	palette      *string
	grayOnly     *bool
	useBG        *bool
)

func init() {
	xStep = flag.Int("x", 4, "x axis step")
	yStep = flag.Int("y", 8, "y axis step")
	palette = flag.String("palette", " .:-=+*#%@", "ascii palette")
	grayOnly = flag.Bool("gray", false, "ignore color")
	useBG = flag.Bool("bg", false, "use truecolor as background instead of foreground")
	flag.Parse()
}

func truecolor(c color.Color) string {
	if os.Getenv("NO_COLOR") == "1" {
		return ""
	}
	r, g, b, _ := c.RGBA()
	if !*useBG {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
	} else {
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)
	}
}

func main() {
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	pLen := len(*palette) - 1

	for y := bounds.Min.Y; y < bounds.Max.Y; y += *yStep {
		for x := bounds.Min.X; x < bounds.Max.X; x += *xStep {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			gray := (r*30 + g*59 + b*11) / 100 >> 8
			idx := int(gray) * pLen / 255
			ch := (*palette)[idx]

			if *grayOnly {
				fmt.Printf("%c", ch)
			} else {
				fmt.Printf("%s%c%s", truecolor(c), ch, reset)
			}
		}
		fmt.Println()
	}
}
