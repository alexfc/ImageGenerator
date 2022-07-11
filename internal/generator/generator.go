package generator

import (
	"ImageGenerator/internal/config"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
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

func distortImage(style config.Style, c chan *imagick.MagickWand) {
	name := "Block Sub-Assy, Cylinder, Toyota"
	oem := "11401-69425"

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	pw.SetColor("none")
	err := mw.NewImage(style.LabelWidth, style.LabelHeight*2, pw)

	if err != nil {
		panic(err)
	}

	pw.SetColor("black")
	err = dw.SetFont("fonts/roboto_condensed.ttf")

	if err != nil {
		panic(err)
	}

	dw.SetFontSize(style.FontSize)
	dw.SetFillColor(pw)

	dw.SetTextAntialias(true)

	dw.Annotation(0, 20, name)

	dw.SetFontSize(25)
	dw.Annotation(60, 80, oem)
	dw.Line(0, 45, 240, 45)
	err = mw.DrawImage(dw)
	logError(err)
	args := []float64{0, 0, 0, 0, float64(mw.GetImageWidth()), 0, float64(mw.GetImageWidth()), 61, 0, float64(mw.GetImageHeight()), 0, float64(mw.GetImageHeight()), float64(mw.GetImageWidth()), float64(mw.GetImageHeight()), float64(mw.GetImageWidth()), float64(mw.GetImageHeight()) + 61}
	err = mw.DistortImage(imagick.DISTORTION_PERSPECTIVE, args, false)
	logError(err)

	dw.Destroy()
	pw.Destroy()

	c <- mw
}

func GenerateImage(filename string, style config.Style) ([]byte, error) {

	mw := imagick.NewMagickWand()
	ch := make(chan *imagick.MagickWand)

	go distortImage(style, ch)

	err := mw.ReadImage(filename)
	logError(err)

	err = mw.CompositeImage(<-ch, imagick.COMPOSITE_OP_OVER, true, int(style.X), int(style.Y))
	logError(err)

	result := mw.GetImageBlob()
	mw.Destroy()

	return result, nil
}
