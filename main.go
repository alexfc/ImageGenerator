package main

import (
	conf "ImageGenerator/internal/config"
	gen "ImageGenerator/internal/generator"
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

//var ctx = context.Background()

func main() {
	//data := []byte{'[[40089262,"0308000421","0308000421","HOSEENG SPEED CONTR",120],[40089263,"0308100820","0308100820","HOSE WATER",110]]'}

	//os.Exit(0)

	imagick.Initialize()
	defer imagick.Terminate()

	conf, err := readConfig()

	if err != nil {
		panic("config not valid")
	}

	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
	//	Password: conf.Redis.Pass, // no password set
	//	DB:       conf.Redis.DB,   // use default DB
	//})

	//manufacturers, err := apiclient.GetManufacturers()
	//
	//for _, manufacturer := range manufacturers {
	//	for true {
	//		items, _ := apiclient.GetItems(manufacturer, 0, 10000)
	//		//pipe := rdb.Pipeline()
	//		for _, item := range items {
	//			fmt.Println(item)
	//			os.Exit(0)
	//			//key := item[1]
	//			//pipe.Set(ctx, )
	//		}
	//		//pipe.Exec(ctx)
	//
	//		if len(items) < 10000 {
	//			break
	//		}
	//	}
	//}
	//
	//rdb.Set(ctx, "helloWorld", 1, 0)
	//val := rdb.Get(ctx, "helloWorld")

	//fmt.Println(val)

	router := mux.NewRouter()

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
