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
	"path/filepath"
	"sort"
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

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	log.SetOutput(f)

	if len(argsWithoutProg) > 0 {
		filename := argsWithoutProg[0]
		generateColors(filename)
	} else {
		fmt.Println("usage: go run main.go sample.jpg")
	}
}

func generateColors(filename string) {
	if _, err := os.Stat(filename); err != nil {
		log.Fatal(err)
	}
	colors := anaysisImage(filename)

	sort.Sort(sort.Reverse(ByColor(colors)))

	drawImage(colors, getSaveToFilename(filename))
}

func anaysisImage(filepath string) []Color {
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

	i := 0
	for _, c := range colorMap {
		results[i] = Color{
			Color: c.Color,
			Count: c.Count,
			Ratio: float64(c.Count) / float64(total),
			RatioWithoutBcakground: float64(c.Count) / float64(totalWithoutMax),
			Red:   c.Red,
			Green: c.Green,
			Blue:  c.Blue,
		}
		i++
	}

	return results
}

type ByColor []Color

func (a ByColor) Len() int           { return len(a) }
func (a ByColor) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByColor) Less(i, j int) bool { return a[i].Count < a[j].Count }

func drawImage(colors []Color, saveTo string) {
	width := 1000
	height := 40

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	startX := 0
	endX := 0
	for _, c := range colors {
		log.Print(c.Color, c.Count, c.Ratio)
		if int(c.Ratio*float64(width)) > 1 && !isWhiteOrGray(c) {
			endX = startX + int(c.Ratio*float64(width))
			drawRect(m, color.RGBA{c.Red, c.Green, c.Blue, 255}, startX, 0, endX, height)
			startX = endX
		}
	}

	if startX < width {
		drawRect(m, color.RGBA{255, 255, 255, 255}, startX, width, 0, height)
	}

	saveImage(m, saveTo)
}

func isWhiteOrGray(c Color) bool {
	return c.Red != 0 && c.Red == c.Green && c.Red == c.Blue
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

func getSaveToFilename(filename string) string {
	ext := filepath.Ext(filename)
	return filepath.Join(filename[0:len(filename)-len(ext)] + "_color" + ext)
}
