package lib

import (
	"embed"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Point struct {
	x int
	y int
}

type OGImage struct {
	Height      int
	Width       int
	MaxStartPtX int
	MaxEndPtX   int
}

const dev = false

//go:embed fonts/Inter-Regular.ttf
var embedFS embed.FS
var fontFile = "fonts/Inter-Regular.ttf"

func DrawImage(title string, subTitle string, fontSize int, subFontSize int, color string, backgroundImageURL string) image.Image {
	titleFontFace := loadFont(fontSize)
	subtitleFontFace := loadFont(subFontSize)

	const height = 627
	const width = 1200

	const centerV = height / 2
	const centerH = width / 2

	dc := gg.NewContext(width, height)

	if len(backgroundImageURL) > 0 {
		img := loadImageFromURL(backgroundImageURL)
		backgroundImage := imaging.Fill(img, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)
		dc.DrawImage(backgroundImage, 0, 0)
	} else {
		dc.DrawRectangle(0, 0, width, height)
		dc.SetHexColor("#fff")
		dc.Fill()
		dc.Clear()
	}

	dc.SetFontFace(titleFontFace)
	dc.SetHexColor(color)
	_, titleMHeight := dc.MeasureString(title)
	dc.DrawStringAnchored(title, centerH, centerV, 0.5, 0.5)
	dc.SetFontFace(subtitleFontFace)
	dc.DrawStringAnchored(subTitle, centerH, centerV+titleMHeight+20, 0.5, 0.5)
	return dc.Image()
}

func WriteImage(w io.Writer, img image.Image) error {
	if err := png.Encode(w, img); err != nil {
		return err
	}
	return nil
}

func loadFont(fontSize int) font.Face {
	fontBytes, err := embedFS.ReadFile(fontFile)

	if err != nil {
		log.Fatal("failed to find font")
	}

	ttf, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal("failed to parse font")
	}
	fontFace := truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     300,
		Hinting: font.HintingNone,
	})
	return fontFace
}

func loadImageFromURL(url string) image.Image {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		// TODO: handle errors
	}
	defer res.Body.Close()
	imageRef, _, err := image.Decode(res.Body)
	if err != nil {
		//TODO: handle error
	}
	return imageRef
}
