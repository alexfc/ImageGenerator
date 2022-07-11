package main

import (
	"ImageGenerator/internal/apiclient"
	conf "ImageGenerator/internal/config"
	gen "ImageGenerator/internal/generator"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strings"
)

func readConfig() (*conf.Config, error) {
	configFile, err := os.ReadFile("config.yml")

	if err != nil {
		panic("can't read config file")
	}

	config := new(conf.Config)

	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		panic("config is wrong")
	}

	return config, err
}

func getFilename(params gen.GenerationParams) (string, error) {
	filepath := ""

	for _, part := range []string{"templates", params.Style, params.Lang, params.Size} {
		filepath += part + string(os.PathSeparator)
	}

	files, err := os.ReadDir(filepath)

	if err != nil {
		return "", err
	}

	filename := ""

	for _, file := range files {
		if strings.HasPrefix(file.Name(), params.Vendor) {
			filename = file.Name()
		}
	}

	return filepath + string(os.PathSeparator) + filename, nil
}

func parseParams(r *http.Request) gen.GenerationParams {
	vars := mux.Vars(r)

	return gen.GenerationParams{
		Style:  vars["style"],
		Lang:   vars["lang"],
		Size:   vars["zoom"],
		Vendor: vars["manufacturer"],
		OEM:    vars["oem"],
	}
}

func cacheWarmer(rdb *redis.Client) error {
	manufacturers, err := apiclient.GetManufacturers()
	offset := 0

	if err != nil {
		return err
	}

	for _, manufacturer := range manufacturers {
		offset = 0
		for true {
			items, _ := apiclient.GetItems(manufacturer, offset, 10000)
			pipe := rdb.Pipeline()
			for _, item := range items {
				val := item.FormattedNumber + "##" + item.Name
				pipe.Set(ctx, item.Number, val, 0)
				offset = item.Id
			}
			pipe.Exec(ctx)

			if len(items) < 10000 {
				break
			}
		}
	}

	return nil
}

func cacheClear(rdb *redis.Client) {
	rdb.FlushDB(ctx)
}

var ctx = context.Background()

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	conf, err := readConfig()

	if err != nil {
		panic("config not valid")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Pass,
		DB:       conf.Redis.DB,
	})

	router := mux.NewRouter()

	router.HandleFunc("/cache-warm", func(w http.ResponseWriter, r *http.Request) {
		cacheWarmer(rdb)
		w.Write([]byte("Done"))
	})

	router.HandleFunc("/cache-clear", func(w http.ResponseWriter, r *http.Request) {
		cacheClear(rdb)
		w.Write([]byte("Done"))
	})

	router.HandleFunc("/{style:[a-z-]+}/{lang:[a-z]+}/{zoom:[0-9x]+}/{manufacturer:[a-z]+}/{oem:[a-zA-Z0-9-]+}.png", func(w http.ResponseWriter, r *http.Request) {
		params := parseParams(r)
		style := conf.GetStyle(params.Style)
		filename, err := getFilename(params)

		if err != nil {
			//handle error
		}

		image, err := gen.GenerateImage(filename, style)

		if err != nil {
			panic(err)
		}

		w.Write(image)
	})

	//args := os.Args[1:]
	//port := args[0]

	http.Handle("/", router)

	http.ListenAndServe(":5555", nil)
}
