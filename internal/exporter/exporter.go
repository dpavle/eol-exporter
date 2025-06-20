package exporter

import (
	"fmt"
	"net/http"
	"os"
	"eol-exporter/internal/api"
	"strconv"
	"time"
	"regexp"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zcalusic/sysinfo"
)

var (
	ProductReleaseInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "product_release_info",
			Help: "Full information about a product release cycle.",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
			"isLts",
			"isEol",
			"isEoas",
			"isEoes",
			"isMaintained",
			"latest_name",
			"latest_link",
		},
	)
	ProductDetailsInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "product_details_info",
			Help: "Full details of a product.",
		}, []string{
			"host",
			"name",
			"product",
			"label",
			//"aliases",
			"category",
			//"tags",
			"versionCommand",
			//"identifiers",
			//"labels",
			//"links",
			//"releases",
		},
	)
	EolDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eol_date",
			Help: "End of life date for the product release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
	EoasDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eoas_date",
			Help: "End of active support date for the release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
	EoesDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eoes_date",
			Help: "End of extended support date for the release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
)

func RefreshMetrics(reg *prometheus.Registry, productCycleData api.ProductCycleData, productDetailsData api.ProductDetailsData) {

	hostname, _ := os.Hostname()
	ProductReleaseInfo.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		productCycleData.Result.Name,
		productCycleData.Product,
		productCycleData.Result.Codename,
		productCycleData.Result.Label,
		strconv.FormatBool(productCycleData.Result.IsLts),
		strconv.FormatBool(productCycleData.Result.IsEol),
		strconv.FormatBool(productCycleData.Result.IsEoas),
		strconv.FormatBool(productCycleData.Result.IsEoes),
		strconv.FormatBool(productCycleData.Result.IsMaintained),
		productCycleData.Result.Latest.Name,
		productCycleData.Result.Latest.Link,

	).Set(1)

	ProductDetailsInfo.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		productDetailsData.Result.Name,
		productDetailsData.Product,
		productDetailsData.Result.Label,
		productDetailsData.Result.Category,
		productDetailsData.Result.VersionCommand,
	).Set(1)

	if productCycleData.Result.EolFrom != nil {
		EolDateUnixTS.WithLabelValues(
			fmt.Sprintf("%s", hostname),
			productCycleData.Result.Name,
			productCycleData.Product,
			productCycleData.Result.Codename,
			productCycleData.Result.Label,
		).Set(float64(productCycleData.Result.EolFrom.Unix()))
	}

	if productCycleData.Result.EoasFrom != nil {
		EoasDateUnixTS.WithLabelValues(
			fmt.Sprintf("%s", hostname),
			productCycleData.Result.Name,
			productCycleData.Product,
			productCycleData.Result.Codename,
			productCycleData.Result.Label,
		).Set(float64(productCycleData.Result.EoasFrom.Unix()))
	}
	if productCycleData.Result.EoesFrom != nil {
		EoesDateUnixTS.WithLabelValues(
			fmt.Sprintf("%s", hostname),
			productCycleData.Result.Name,
			productCycleData.Product,
			productCycleData.Result.Codename,
			productCycleData.Result.Label,
		).Set(float64(productCycleData.Result.EoesFrom.Unix()))
	}

}

func UpdateMetrics(reg *prometheus.Registry, productCycleData api.ProductCycleData, productDetailsData api.ProductDetailsData) {
	ticker := time.NewTicker(24 * time.Hour)

	for {
		RefreshMetrics(reg, productCycleData, productDetailsData)
		<-ticker.C
	}
}

func StartExporter() error {

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	reg := prometheus.NewRegistry()

	reg.MustRegister(ProductReleaseInfo)
	reg.MustRegister(ProductDetailsInfo)
	reg.MustRegister(EolDateUnixTS)
	reg.MustRegister(EoasDateUnixTS)
	reg.MustRegister(EoesDateUnixTS)

	var si sysinfo.SysInfo
	si.GetSysInfo()

	// OS
	product := si.OS.Vendor
	version := si.OS.Version

	var productCycleData api.ProductCycleData
	var productDetailsData api.ProductDetailsData
	var productCycleApiErr error
	var productDetailsApiErr error

	productCycleData, productCycleApiErr = api.FetchProductCycleData(&httpClient, "", product, version)
	productDetailsData, productDetailsApiErr = api.FetchProductDetailsData(&httpClient, "", product)
	if productCycleApiErr != nil {
		log.Println(productCycleApiErr)
	}
	if productDetailsApiErr != nil {
		log.Println(productDetailsApiErr)
	}
	go UpdateMetrics(reg, productCycleData, productDetailsData)

	// Kernel
	pattern := regexp.MustCompile(`^[0-9]+.[0-9]+`)
	version = pattern.FindString(si.Kernel.Release)

	productCycleData, productCycleApiErr = api.FetchProductCycleData(&httpClient, "", "linux", version)
	productDetailsData, productDetailsApiErr = api.FetchProductDetailsData(&httpClient, "", "linux")
	if productCycleApiErr != nil {
		log.Println(productCycleApiErr)
	}
	if productDetailsApiErr != nil {
		log.Println(productDetailsApiErr)
	}
	go UpdateMetrics(reg, productCycleData, productDetailsData)

	productCycleData, productCycleApiErr = api.FetchProductCycleData(&httpClient, "", "undefined", "undefined")
	productDetailsData, productDetailsApiErr = api.FetchProductDetailsData(&httpClient, "", "undefined")
	if productCycleApiErr != nil {
		log.Println(productCycleApiErr)
	}
	if productDetailsApiErr != nil {
		log.Println(productDetailsApiErr)
	}

	handler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)
	http.Handle("/metrics", handler)

	fmt.Println("Starting HTTP server on :2112")
	return http.ListenAndServe(":2112", nil)
}
