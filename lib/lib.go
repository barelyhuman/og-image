package lib

import (
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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

func (ogImage *OGImage) drawText(text string, img *image.RGBA, fontFace font.Face, point fixed.Point26_6, startPoint int) {

	fittingText := getTruncatedText(fontFace, text, startPoint, ogImage.MaxEndPtX)

	if len(fittingText) < len(text) {
		fittingText = fittingText[0:len(fittingText)-3] + "..."
	}

	fontDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: fontFace,
		Dot:  point,
	}

	fontDrawer.DrawString(fittingText)
}

func (ogImage *OGImage) drawTitle(text string, img *image.RGBA, fontSize int) {
	fontBytes, err := embedFS.ReadFile(fontFile)
	const offset = 2

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

	xToNegate, err := findStartPointX(text, fontFace)

	if err != nil {
		log.Fatal("couldn't find start point")
	}

	middleLine := Point{
		x: ogImage.Width / 2,
		y: ogImage.Height / 2,
	}

	dotStart := (middleLine.x - xToNegate)
	if dotStart < ogImage.MaxStartPtX {
		dotStart = ogImage.MaxStartPtX
	}

	titlePositionY := int(ogImage.Height - int(float32(ogImage.Height)/offset))
	titlePoint := freetype.Pt(dotStart, titlePositionY)
	ogImage.drawText(text, img, fontFace, titlePoint, dotStart)
}

func (ogImage *OGImage) drawSubTitle(text string, img *image.RGBA, fontSize int) {
	fontBytes, err := embedFS.ReadFile(fontFile)
	const offset = 2.5

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

	xToNegate, err := findStartPointX(text, fontFace)

	if err != nil {
		log.Fatal("couldn't find start point")
	}

	middleLine := Point{
		x: ogImage.Width / 2,
		y: ogImage.Height / 2,
	}

	dotStart := (middleLine.x - xToNegate)
	if dotStart < ogImage.MaxStartPtX {
		dotStart = ogImage.MaxStartPtX
	}

	titlePositionY := int(ogImage.Height - int(float32(ogImage.Height)/offset))
	titlePoint := freetype.Pt(dotStart, titlePositionY)
	ogImage.drawText(text, img, fontFace, titlePoint, dotStart)
}

func (ogImage *OGImage) drawCircle(img draw.Image, x0, y0, r float64, c color.Color) {
	x, y, dx, dy := float64(r-1), float64(0), float64(1), float64(1)
	err := dx - (r * float64(2))

	for x > y {
		img.Set(int(x0+x), int(y0+y), c)
		img.Set(int(x0+y), int(y0+x), c)
		img.Set(int(x0-y), int(y0+x), c)
		img.Set(int(x0-x), int(y0+y), c)
		img.Set(int(x0-x), int(y0-y), c)
		img.Set(int(x0-y), int(y0-x), c)
		img.Set(int(x0+y), int(y0-x), c)
		img.Set(int(x0+x), int(y0-y), c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func (ogImage *OGImage) drawBoundaries(img *image.RGBA) {

	xMid := ogImage.Width / 2
	yMid := ogImage.Height / 2

	for y := 0; y <= ogImage.Height; y++ {
		img.Set(xMid, y, color.Black)
	}

	for x := 0; x <= ogImage.Width; x++ {
		img.Set(x, yMid, color.Black)
	}

}

// DrawImage - draw the og image and return a image.Image to be used with
// an io.Writer to write it to the needed destination
func DrawImage(title string, subtitle string, fontSize int) image.Image {
	ogImage := OGImage{}

	ogImage.Width = 1200
	ogImage.Height = 627

	ogImage.MaxStartPtX = ogImage.Width / 8
	ogImage.MaxEndPtX = ogImage.Width - ogImage.Width/8

	imgRGBA := image.NewRGBA(image.Rect(0, 0, ogImage.Width, ogImage.Height))
	white := color.RGBA{255, 255, 255, 255}
	var img image.Image = imgRGBA

	draw.Draw(imgRGBA, imgRGBA.Bounds(), &image.Uniform{white}, image.Point{X: 0, Y: 0}, draw.Src)

	for y := 0; y <= ogImage.Height; y++ {
		for x := 0; x <= ogImage.Width; x++ {
			if x%40 == 0 && y%40 == 0 {
				ogImage.drawCircle(imgRGBA, float64(x), float64(y), 1.1, color.RGBA{153, 153, 153, 255})
			}
		}
	}

	if dev {
		ogImage.drawBoundaries(imgRGBA)
	}

	ogImage.drawTitle(title, imgRGBA, fontSize)
	ogImage.drawSubTitle(subtitle, imgRGBA, int(float64(fontSize)/2))

	return img
}

// WriteImage - a simple writer that encodes the generated img into a png
// you can use other writers for advanced use cases
func WriteImage(w io.Writer, img *image.Image) error {
	if err := png.Encode(w, *img); err != nil {
		return err
	}
	return nil
}

func findStartPointX(text string, face font.Face) (int, error) {

	if len(text) <= 0 {
		return 0, nil
	}

	aWidth := font.MeasureString(face, text)
	return int(float64(aWidth.Round()) * float64(0.5)), nil
}

func getTruncatedText(face font.Face, text string, startPoint int, maxPoint int) string {
	dotAdvanceEnd := font.MeasureString(face, text).Round()

	if startPoint+dotAdvanceEnd <= maxPoint {
		return text
	}
	return getTruncatedText(face, text[0:len(text)-1], startPoint, maxPoint)
}
