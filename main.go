package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

type Color struct {
	Color                  string
	Count                  int
	Ratio                  float64
	RatioWithoutBcakground float64
	Red                    uint8
	Green                  uint8
	Blue                   uint8
}

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		filepath := argsWithoutProg[0]
		colorMap := anaysisImage(filepath)

		total := 0
		max := 0

		for _, c := range colorMap {
			total = total + c.Count
			if c.Count > max {
				max = c.Count
			}
		}

		totalWithoutMax := total - max

		results := make([]Color, len(colorMap))

		for _, c := range colorMap {
			results = append(results, Color{
				Color: c.Color,
				Count: c.Count,
				Ratio: float64(c.Count) / float64(total),
				RatioWithoutBcakground: float64(c.Count) / float64(totalWithoutMax),
				Red:   c.Red,
				Green: c.Green,
				Blue:  c.Blue,
			})
		}

		for _, c := range results {
			if c.RatioWithoutBcakground > 0.01 {
				fmt.Printf("%s %.2f%%\n", c.Color, c.Ratio*100.0)
			}
		}
		drawImage(results)
	}
}

func anaysisImage(filepath string) map[string]Color {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := m.Bounds()

	colorMap := make(map[string]Color)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			color := fmt.Sprintf("%02X%02X%02X", r>>8, g>>8, b>>8)
			if colorMap[color].Count > 0 {
				colorMap[color] = Color{
					Color: color,
					Ratio: 0,
					RatioWithoutBcakground: 0,
					Count: colorMap[color].Count + 1,
					Red:   uint8(r >> 8),
					Green: uint8(g >> 8),
					Blue:  uint8(b >> 8),
				}
				continue
			}
			colorMap[color] = Color{
				Color: color,
				Ratio: 0,
				RatioWithoutBcakground: 0,
				Count: 1,
			}
		}
	}
	return colorMap
}

func drawImage(colors []Color) {
	width := 1000
	height := 40

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	startX := 0
	endX := 0
	for _, c := range colors {
		endX = startX + int(c.Ratio*float64(width))
		drawRect(m, color.RGBA{c.Red, c.Green, c.Blue, 255}, startX, 0, endX, height)
		startX = endX
	}
	saveImage(m, "output.png")
}

func drawRect(img *image.RGBA, col color.RGBA, x1, y1, x2, y2 int) {
	for ; x1 <= x2; x1++ {
		j := y1
		for ; j <= y2; j++ {
			img.Set(x1, j, col)
		}
	}
}

func saveImage(m *image.RGBA, filepath string) {
	myfile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	png.Encode(myfile, m)
}
