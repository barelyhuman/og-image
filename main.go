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
	fontSize := flag.Int("size", 45, "Font Size")
	outFile := flag.String("out", "./og-image.png", "File to output to, will export a png")

	flag.Parse()

	file, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	img := lib.DrawImage(*title, *subTitle, *fontSize)
	lib.WriteImage(file, &img)
}
