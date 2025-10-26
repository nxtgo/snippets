package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
)

var (
	inPath  = flag.String("in", "in.png", "input")
	outPath = flag.String("out", "out.png", "output")
	thr     = flag.Float64("t", 60, "threshold")
	sz      = flag.Int("s", 8, "sample size")
	soft    = flag.Bool("soft", true, "soft alpha")
)

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*inPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	n := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			n.Set(x, y, color.NRGBAModel.Convert(img.At(x, y)))
		}
	}

	bg := bgcol(n, w, h, *sz)
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	var wg sync.WaitGroup
	for i := range cpu {
		y0 := i * h / cpu
		y1 := (i + 1) * h / cpu

		wg.Go(func() {
			for y := y0; y < y1; y++ {
				for x := range w {
					p := n.NRGBAAt(x, y)
					d := dist(p, bg)
					if d <= *thr {
						if *soft {
							a := uint8(float64(p.A) * d / *thr)
							out.Set(x, y, color.NRGBA{p.R, p.G, p.B, a})
						} else {
							out.Set(x, y, color.NRGBA{p.R, p.G, p.B, 0})
						}
					} else {
						out.Set(x, y, p)
					}
				}
			}
		})
	}
	wg.Wait()

	of, _ := os.Create(*outPath)
	defer of.Close()

	png.Encode(of, out)
	fmt.Printf("done: %s\n", *outPath)
}

func bgcol(m *image.NRGBA, w, h, s int) color.NRGBA {
	f := func(x0, y0 int) (r, g, b, n int64) {
		for y := range s {
			for x := range s {
				p := m.NRGBAAt(x0+x, y0+y)
				r, g, b, n = r+int64(p.R), g+int64(p.G), b+int64(p.B), n+1
			}
		}
		return
	}

	s1r, s1g, s1b, s1n := f(0, 0)
	s2r, s2g, s2b, s2n := f(w-s, 0)
	s3r, s3g, s3b, s3n := f(0, h-s)
	s4r, s4g, s4b, s4n := f(w-s, h-s)
	n := s1n + s2n + s3n + s4n

	return color.NRGBA{uint8((s1r + s2r + s3r + s4r) / n), uint8((s1g + s2g + s3g + s4g) / n), uint8((s1b + s2b + s3b + s4b) / n), 255}
}

func dist(a, b color.NRGBA) float64 {
	dr := float64(a.R) - float64(b.R)
	dg := float64(a.G) - float64(b.G)
	db := float64(a.B) - float64(b.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}
