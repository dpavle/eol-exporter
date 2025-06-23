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
	"plugin"
	"net"
	"github.com/spf13/viper"

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
	ReleaseDateUnixTS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "release_date",
			Help: "Release date of the product release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).",
		}, []string{
			"host",
			"name",
			"product",
			"codename",
			"label",
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

func RegisterTimeSeries(reg *prometheus.Registry, productCycleData api.ProductCycleData, productDetailsData api.ProductDetailsData) {

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

	ReleaseDateUnixTS.WithLabelValues(
		fmt.Sprintf("%s", hostname),
		productCycleData.Result.Name,
		productCycleData.Product,
		productCycleData.Result.Codename,
		productCycleData.Result.Label,
	).Set(float64(productCycleData.Result.ReleaseDate.Unix()))

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

func FetchDataAndUpdateMetrics(reg *prometheus.Registry, httpClient *http.Client, product string, version string) {
	ticker := time.NewTicker(24 * time.Hour)

	for {
		var productCycleData api.ProductCycleData
		var productDetailsData api.ProductDetailsData
		var productCycleApiErr error
		var productDetailsApiErr error

		productCycleData, productCycleApiErr = api.FetchProductCycleData(httpClient, "", product, version)
		productDetailsData, productDetailsApiErr = api.FetchProductDetailsData(httpClient, "", product)
		if productCycleApiErr != nil {
			log.Println(productCycleApiErr)
		}
		if productDetailsApiErr != nil {
			log.Println(productDetailsApiErr)
		}
		RegisterTimeSeries(reg, productCycleData, productDetailsData)
		<-ticker.C
	}
}

type DataPlugin interface {
	GetProductAndVersion() (product string, version string, err error)
}

func StartExporter() error {

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	reg := prometheus.NewRegistry()

	reg.MustRegister(ProductReleaseInfo)
	reg.MustRegister(ProductDetailsInfo)
	reg.MustRegister(ReleaseDateUnixTS)
	reg.MustRegister(EolDateUnixTS)
	reg.MustRegister(EoasDateUnixTS)
	reg.MustRegister(EoesDateUnixTS)

	// OS
	var si sysinfo.SysInfo
	si.GetSysInfo()
	go FetchDataAndUpdateMetrics(reg, &httpClient, si.OS.Vendor, si.OS.Version)

	// Kernel
	kernelPattern := regexp.MustCompile(`^[0-9]+.[0-9]+`)
	go FetchDataAndUpdateMetrics(reg, &httpClient, "linux", kernelPattern.FindString(si.Kernel.Release))

	// Plugins
	var EnabledPlugins []string = viper.GetStringSlice("plugins")

	for _, EnabledPlugin := range EnabledPlugins {
		p, err := plugin.Open(fmt.Sprintf("plugins/%s/%s.so", EnabledPlugin, EnabledPlugin))
		if err != nil {
			log.Fatalf("Failed to open plugin: %s", err)
		}
		symbol, err := p.Lookup("DataPlugin")
		if err != nil {
			log.Fatalf("Failed to find plugin symbol 'DataPlugin': %s", err)
		}

		pluginInstance, ok := symbol.(DataPlugin)
		if !ok {
			log.Fatal("Symbol does not implement the DataPlugin interface!")
		}

		PluginReturnedProduct, PluginReturnedVersion, err := pluginInstance.GetProductAndVersion()
		if err != nil {
			log.Fatalf("Failed executing plugin: %s", err)
		}
		go FetchDataAndUpdateMetrics(reg, &httpClient, PluginReturnedProduct, PluginReturnedVersion)
	}

	handler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)
	http.Handle("/metrics", handler)

	port := viper.GetString("port")
	addr := viper.GetString("address")

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		return err
	}

	fmt.Printf("Starting HTTP server on %s:%s\n", addr, port)
	return http.Serve(ln, nil)
}
