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

type EOLData struct {
	Result Result `json:"result"`
}


type Latest struct {
	Name string     `json:"name"`
	Date CustomTime `json:"date"`
	Link string     `json:"link"`
}

type Result struct {
	Name         string     `json:"name"`
	Codename     string     `json:"codename"`
	Label        string     `json:"label"`
	ReleaseDate  CustomTime `json:"releaseDate"`
	IsLts        bool       `json:"isLts"`
	LtsFrom      *string    `json:"ltsFrom"`
	IsEoas       bool       `json:"isEoas"`
	EoasFrom     *CustomTime `json:"eoasFrom"`
	IsEol        bool       `json:"isEol"`
	EolFrom      CustomTime `json:"eolFrom"`
	IsEoes       bool       `json:"isEoes"`
	EoesFrom     *string    `json:"eoesFrom"`
	IsMaintained bool       `json:"isMaintained"`
	Latest       Latest     `json:"latest"`
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
	fmt.Println(date)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func FetchProductCycleEOLData(product string, release string) (data EOLData, err error) {

	baseUrl := "https://endoflife.date/api/v1"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/products/%s/releases/%s", baseUrl, product, release), nil)

	if err != nil {
		log.Fatal(err)
		return data, err
	}

	res, getErr := httpClient.Do(req)

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

	return data, nil
}
