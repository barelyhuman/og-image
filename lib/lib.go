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

type OGImage struct {
	Height      int
	Width       int
	MaxStartPtX int
	MaxEndPtX   int
}

//go:embed fonts/Inter-Regular.ttf
var embedFS embed.FS
var fontFile = "fonts/Inter-Regular.ttf"

func DrawImage(title string, subTitle string, fontSize int, subFontSize int, color string, backgroundImageURL string, backgroundImageColor string, padding int) image.Image {
	titleFontFace := loadFont(fontSize)
	subtitleFontFace := loadFont(subFontSize)

	const height = 627
	const width = 1200

	dc := gg.NewContext(width, height)

	if len(backgroundImageColor) > 0 {
		dc.DrawRectangle(0, 0, width, height)
		dc.SetHexColor(backgroundImageColor)
		dc.Fill()
	} else if len(backgroundImageURL) > 0 {
		img := loadImageFromURL(backgroundImageURL)
		backgroundImage := imaging.Fill(img, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)
		dc.DrawImage(backgroundImage, 0, 0)
	} else {
		dc.DrawRectangle(0, 0, width, height)
		dc.SetHexColor("#fff")
		dc.Fill()
		dc.Clear()
	}

	// Calculate the padded area
	paddedWidth := width - 2*padding
	centerH := width / 2
	centerV := height / 2

	dc.SetFontFace(titleFontFace)
	dc.SetHexColor(color)

	// Draw the title, wrapped within the padded area
	titleMaxWidth := float64(paddedWidth)
	// Measure title height for subtitle placement
	_, titleMHeight := dc.MeasureMultilineString(title, titleMaxWidth)
	dc.DrawStringWrapped(title, float64(centerH), float64(centerV), 0.5, 1, titleMaxWidth, 1.5, gg.AlignCenter)

	// Draw the subtitle below the title, also wrapped
	dc.SetFontFace(subtitleFontFace)
	subtitleMaxWidth := float64(paddedWidth)
	dc.DrawStringWrapped(subTitle, float64(centerH), float64(centerV)+titleMHeight+20, 0.5, 0, subtitleMaxWidth, 1.5, gg.AlignCenter)
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
