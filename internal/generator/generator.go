package generator

import (
	"ImageGenerator/internal/config"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
	"strings"
)

type GenerationParams struct {
	Style  string
	Lang   string
	Size   string
	Vendor string
	OEM    string
}

func logError(err interface{}) {
	if err != nil {
		fmt.Println(err)
	}
}

func wordWrapAnnotation(mw *imagick.MagickWand, dw *imagick.DrawingWand, text string, maxWidth float64) ([]string, float64) {
	words := strings.Split(text, " ")
	var lines []string

	i := 0
	var lineHeight float64 = 0

	for len(words) > 0 {
		i += 1
		metrics := mw.QueryFontMetrics(dw, strings.Join(words[:i], " "))
		lineHeight = math.Max(metrics.TextHeight, lineHeight)

		if metrics.TextWidth > maxWidth || len(words) <= i {
			lines = append(lines, strings.Join(words[:i], " "))
			words = words[i:]
			i = 0
		}
	}

	return lines, lineHeight
}

func distortImage(style config.Style, name string, oem string, c chan *imagick.MagickWand) {

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	pw.SetColor("none")
	err := mw.NewImage(style.LabelWidth, style.LabelHeight+45, pw)

	if err != nil {
		panic(err)
	}

	pw.SetColor("black")
	err = dw.SetFont("fonts/roboto_condensed.ttf")

	if err != nil {
		panic(err)
	}

	dw.SetFontSize(style.TitleFontSize)
	dw.SetFillColor(pw)
	dw.SetTextAntialias(true)

	lines, lineHeight := wordWrapAnnotation(mw, dw, name, float64(style.LabelWidth))

	startY := 19

	if len(lines) == 1 {
		startY = 35
	}

	for key, line := range lines {
		if key > 2 {
			break
		}
		dw.Annotation(0, float64(startY)+float64(key)*lineHeight, line)
	}

	dw.SetFontSize(style.OemFontSize)
	dw.SetTextAlignment(imagick.ALIGN_CENTER)

	dw.Annotation(float64(style.LabelWidth/2), 92, oem)

	dw.Line(0, 60, float64(style.LabelWidth), 60)
	err = mw.DrawImage(dw)
	logError(err)

	args := []float64{0, 0, 0, 0, float64(mw.GetImageWidth()), 0, float64(mw.GetImageWidth()), 61, 0, float64(mw.GetImageHeight()), 0, float64(mw.GetImageHeight()), float64(mw.GetImageWidth()), float64(mw.GetImageHeight()), float64(mw.GetImageWidth()), float64(mw.GetImageHeight()) + 61}
	err = mw.DistortImage(imagick.DISTORTION_PERSPECTIVE, args, false)
	logError(err)

	dw.Destroy()
	pw.Destroy()

	c <- mw
}

func GenerateImage(filename string, style config.Style, name string, oem string) ([]byte, error) {
	ch := make(chan *imagick.MagickWand)

	go distortImage(style, name, oem, ch)

	mw := imagick.NewMagickWand()
	err := mw.ReadImage(filename)
	logError(err)

	err = mw.CompositeImage(<-ch, imagick.COMPOSITE_OP_OVER, true, int(style.X), int(style.Y))
	logError(err)

	mw.SetImageDepth(uint(4))

	result := mw.GetImageBlob()
	mw.Destroy()

	return result, nil
}
