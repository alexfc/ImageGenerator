package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/gographics/imagick.v3/imagick"
	"net/http"
	"os"
	"strings"
)

type GenerationParams struct {
	Style  string
	Lang   string
	Size   string
	Vendor string
	OEM    string
}

type Config struct {
	Styles []Style
}

type Style struct {
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	ImageWidth  uint    `json:"imageWidth"`
	ImageHeight uint    `json:"imageHeight"`
	LabelWidth  uint    `json:"labelWidth"`
	LabelHeight uint    `json:"labelHeight"`
	FontSize    float64 `json:"fontSize"`
}

func (config *Config) getStyle(style string) Style {
	for _, s := range config.Styles {
		if s.Name == style {
			return s
		}
	}
	return Style{}
}

func generateImage(params GenerationParams, style Style) ([]byte, error) {
	filepath := ""

	for _, part := range []string{"templates", params.Style, params.Lang, params.Size} {
		filepath += part + string(os.PathSeparator)
	}

	files, err := os.ReadDir(filepath)

	if err != nil {
		return []byte{}, err
	}

	filename := ""

	for _, file := range files {
		if strings.HasPrefix(file.Name(), params.Vendor) {
			filename = file.Name()
		}
	}

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	cw := mw.Clone()

	err = mw.ReadImage(filepath + string(os.PathSeparator) + filename)

	if err != nil {
		panic(err)
	}

	pw.SetColor("none")
	mw.NewImage(style.ImageWidth, style.ImageHeight, pw)

	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.SetFont("fonts/roboto_condensed.ttf")
	dw.SetFontSize(style.FontSize)
	// Add a black outline to the text
	pw.SetColor("black")
	dw.SetStrokeColor(pw)
	// Turn antialias on - not sure this makes a difference
	dw.SetTextAntialias(true)
	// Now draw the text
	dw.Annotation(style.X, style.Y, style.Title)
	// Draw the image on to the mw
	cw.DrawImage(dw)
	mw.CompositeImage(cw, imagick.COMPOSITE_OP_OVER, true, 0, 0)

	return mw.GetImageBlob(), nil
}

func parseParams(r *http.Request) GenerationParams {
	vars := mux.Vars(r)

	return GenerationParams{
		Style:  vars["style"],
		Lang:   vars["lang"],
		Size:   vars["zoom"] + string('x'),
		Vendor: vars["manufacturer"],
		OEM:    vars["oem"],
	}
}

func readConfig() []byte {
	configFile, err := os.ReadFile("config.json")

	if err != nil {
		panic("can't read config file")
	}

	return configFile
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	config := new(Config)

	err := json.Unmarshal(readConfig(), &config.Styles)

	if err != nil {
		panic("config not valid")
	}

	router := mux.NewRouter()
	router.HandleFunc("/{style:[a-z-]+}/{lang:[a-z]+}/{zoom:[0-9]+}/{manufacturer:[a-z]+}/{oem:[a-zA-Z0-9-]+}.png", func(w http.ResponseWriter, r *http.Request) {
		params := parseParams(r)
		style := config.getStyle(params.Style)
		image, err := generateImage(params, style)
		if err != nil {
			panic(err)
		}
		w.Write(image)
	})

	http.Handle("/", router)

	http.ListenAndServe(":4444", nil)
}
