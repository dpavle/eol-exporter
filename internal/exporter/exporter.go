package exporter

import (
	"fmt"
	"net/http"
	"os"
	"os-eol-exporter/internal/api"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetrics(reg *prometheus.Registry, apiData api.EOLData) {

	var (
		osEolInfo = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "os_eol_info",
				Help: "Information about the end of life for the host OS.",
			}, []string{
				"host",
				"name",
				"codename",
				"label",
				"releaseDate",
				"isLts",
			},
		)
		osEolDateUnixTS = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "os_eol_date",
				Help: "OS end-of-life date as seconds since Unix epoch (Unix Timestamp).",
			}, []string{
				"host",
				"name",
				"codename",
			},
		)
	)

	hostname, _ := os.Hostname()
	osEolInfo.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		apiData.Result.Name,
		apiData.Result.Codename,
		apiData.Result.Label,
		apiData.Result.ReleaseDate.String(),
		strconv.FormatBool(apiData.Result.IsEol),

	).Set(1)

	osEolDateUnixTS.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		apiData.Result.Name,
		apiData.Result.Codename,
	).Set(float64(apiData.Result.EolFrom.Unix()))

	reg.MustRegister(osEolInfo)
	reg.MustRegister(osEolDateUnixTS)
}

func StartExporter() error {
	reg := prometheus.NewRegistry()

	data, _ := api.FetchProductCycleEOLData("fedora", "42")
	RegisterMetrics(reg, data)

	handler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)
	http.Handle("/metrics", handler)

	fmt.Println("Starting HTTP server on :2112")
	return http.ListenAndServe(":2112", nil)
}
