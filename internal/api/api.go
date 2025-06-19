package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"fmt"
	"strconv"
)

type ProductCycleData struct {
	Result ProductCycleResult `json:"result"`
	Product string
}

type ProductDetailsData struct {
	Result ProductDetailsResult `json:"result"`
	Product string
}

type ProductDetailsResult struct {
	Name string `json:"name"`
	Label        string     `json:"label"`
	Category string `json:"category"`
	VersionCommand string `json:"versionCommand"`
	Links ProductDetailsLinks `json:"links"`
}

type ProductDetailsLinks struct {
	HTML string `json:"html"`
	ReleasePolicy string `json:"releasePolicy"`
}

type ProductCycleLatest struct {
	Name string     `json:"name"`
	Date CustomTime `json:"date"`
	Link string     `json:"link"`
}

type ProductCycleResult struct {
	Name         string     `json:"name"`
	Codename     string     `json:"codename"`
	Label        string     `json:"label"`
	ReleaseDate  CustomTime `json:"releaseDate"`
	IsLts        bool       `json:"isLts"`
	LtsFrom      *CustomTime    `json:"ltsFrom"`
	IsEoas       bool       `json:"isEoas"`
	EoasFrom     *CustomTime `json:"eoasFrom"`
	IsEol        bool       `json:"isEol"`
	EolFrom      *CustomTime `json:"eolFrom"`
	IsEoes       bool       `json:"isEoes"`
	EoesFrom     *CustomTime `json:"eoesFrom"`
	IsMaintained bool       `json:"isMaintained"`
	Latest       ProductCycleLatest     `json:"latest"`
	Custom       *string    `json:"custom"`
}

// parsing JSON date strings as time.Time
type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	const layout = "2006-01-02"
	rawDate, _ := strconv.Unquote(string(b))

	date, err := time.Parse(layout, rawDate)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func FetchProductCycleData(client *http.Client, baseUrl string, product string, release string) (data ProductCycleData, err error) {

	if baseUrl == "" {
		baseUrl = "https://endoflife.date/api/v1"
	}
	data.Product = product

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/products/%s/releases/%s", baseUrl, product, release), nil)

	if err != nil {
		log.Fatal(err)
		return data, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
		return data, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
		return data, readErr
	}

	err = json.Unmarshal(body, &data)

	return data, err
}

func FetchProductDetailsData(client *http.Client, baseUrl string, product string) (data ProductDetailsData, err error) {

	if baseUrl == "" {
		baseUrl = "https://endoflife.date/api/v1"
	}
	data.Product = product

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/products/%s", baseUrl, product), nil)

	if err != nil {
		log.Fatal(err)
		return data, err
	}

	res, getErr := client.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
		return data, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
		return data, readErr
	}

	err = json.Unmarshal(body, &data)

	return data, err
}
