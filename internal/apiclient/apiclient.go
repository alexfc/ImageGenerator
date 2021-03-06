package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	API_URL                  = "65.21.55.165"
	API_METHOD_MANUFACTURERS = "get_manufacturers"
	API_METHOD_ITEMS         = "get_all_items"
)

type Manufacturer struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Enabled    bool   `json:"is_enabled"`
	PartsCount int    `json:"parts_count"`
}

type Item struct {
	Id              int
	Number          string
	FormattedNumber string
	Name            string
	Weight          int
}

func getMethodUrl(method string) string {
	return "http://" + API_URL + "/" + method
}

func GetManufacturers() (map[int]Manufacturer, error) {

	var result map[int]Manufacturer
	url := getMethodUrl(API_METHOD_MANUFACTURERS)

	response, err := http.Get(url)

	if err != nil {
		return map[int]Manufacturer{}, err
	}

	body, err := ioutil.ReadAll(response.Body)

	var jsonObj map[string]map[string]json.RawMessage

	json.Unmarshal(body, &jsonObj)
	json.Unmarshal(jsonObj["result"]["manufacturers"], &result)

	return result, nil
}

func GetItems(manufacturer Manufacturer, offset int, limit int) ([]Item, error) {
	var data [][]interface{}
	var items []Item

	query := "?manufacturer=" + manufacturer.Slug + "&offset_id=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	url := getMethodUrl(API_METHOD_ITEMS) + query

	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)

	var jsonObj map[string]map[string]json.RawMessage

	json.Unmarshal(body, &jsonObj)
	json.Unmarshal(jsonObj["result"]["rows"], &data)

	for _, item := range data {
		name, ok := item[3].(string)

		if !ok {
			name = ""
		}

		itemData := Item{
			Id:              int(item[0].(float64)),
			Number:          item[1].(string),
			FormattedNumber: item[2].(string),
			Name:            name,
		}
		items = append(items, itemData)
	}

	return items, nil
}
