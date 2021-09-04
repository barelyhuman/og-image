package main

import (
	"flag"
	"log"
	"os"

	lib "github.com/barelyhuman/og-image/lib"
)

func main() {
	title := flag.String("title", "", "Title")
	subTitle := flag.String("desc", "", "Description or Sub Title")
	fontSize := flag.Int("size-one", 16, "Font Size for title")
	subFontSize := flag.Int("size-two", 12, "Font Size For description")
	color := flag.String("color", "#000", "Font Color")
	backgroundImageURL := flag.String("background-url", "", "URL for the background")
	outFile := flag.String("out", "./og-image.png", "File to output to, will export a png")

	flag.Parse()

	file, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	img := lib.DrawImage(*title, *subTitle, *fontSize, *subFontSize, *color, *backgroundImageURL)
	lib.WriteImage(file, img)
}
