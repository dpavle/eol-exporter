package exporter

import (
	"fmt"
	"net/http"
	"os"
	"os-eol-exporter/internal/api"
	"strconv"
	"time"
	"github.com/themakers/osinfo"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	osEolInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "os_eol_info",
			Help: "Information about the end of life for the host OS.",
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
	osEolDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "os_eol_date",
			Help: "OS end-of-life date as seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
	osEoasDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "os_eoas_date",
			Help: "OS end-of-life date as seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
	osEoesDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "os_eoes_date",
			Help: "OS end-of-life date as seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
		},
	)
)

func RefreshMetrics(reg *prometheus.Registry, apiData api.EOLData) {

	hostname, _ := os.Hostname()
	osEolInfo.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		apiData.Result.Name,
		apiData.Product,
		apiData.Result.Codename,
		apiData.Result.Label,
		strconv.FormatBool(apiData.Result.IsLts),
		strconv.FormatBool(apiData.Result.IsEol),
		strconv.FormatBool(apiData.Result.IsEoas),
		strconv.FormatBool(apiData.Result.IsEoes),
		strconv.FormatBool(apiData.Result.IsMaintained),
		apiData.Result.Latest.Name,
		apiData.Result.Latest.Link,

	).Set(1)

	osEolDateUnixTS.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		apiData.Result.Name,
		apiData.Product,
		apiData.Result.Codename,
		apiData.Result.Label,
	).Set(float64(apiData.Result.EolFrom.Unix()))

	if apiData.Result.EoasFrom != nil {
		osEoasDateUnixTS.WithLabelValues(
			fmt.Sprintf("%s", hostname),
			apiData.Result.Name,
			apiData.Product,
			apiData.Result.Codename,
			apiData.Result.Label,
		).Set(float64(apiData.Result.EoasFrom.Unix()))
	}
	if apiData.Result.EoesFrom != nil {
		osEoesDateUnixTS.WithLabelValues(
			fmt.Sprintf("%s", hostname),
			apiData.Result.Name,
			apiData.Product,
			apiData.Result.Codename,
			apiData.Result.Label,
		).Set(float64(apiData.Result.EoesFrom.Unix()))
	}

}

func UpdateMetrics(reg *prometheus.Registry, data api.EOLData) {
	ticker := time.NewTicker(24 * time.Hour)

	for {
		RefreshMetrics(reg, data)
		<-ticker.C
	}
}

func StartExporter() error {
	reg := prometheus.NewRegistry()

	reg.MustRegister(osEolInfo)
	reg.MustRegister(osEolDateUnixTS)
	reg.MustRegister(osEoasDateUnixTS)
	reg.MustRegister(osEoesDateUnixTS)

	os_info := osinfo.GetInfo()
	product := os_info.LinuxRelease.ID
	version := os_info.LinuxRelease.VersionID

	data, _ := api.FetchProductCycleEOLData(product, version)
	go UpdateMetrics(reg, data)

	handler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)
	http.Handle("/metrics", handler)

	fmt.Println("Starting HTTP server on :2112")
	return http.ListenAndServe(":2112", nil)
}
