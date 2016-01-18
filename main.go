package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		filepath := argsWithoutProg[0]
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

		colorMap := make(map[string]int)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				color := fmt.Sprintf("%02X%02X%02X", r>>8, g>>8, b>>8)
				if colorMap[color] > 0 {
					colorMap[color] = colorMap[color] + 1
					continue
				}
				colorMap[color] = 1
			}
		}

		for color, count := range colorMap {
			fmt.Printf("%s %d\n", color, count)
		}
	}
}
