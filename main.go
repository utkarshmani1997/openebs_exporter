package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/utkarshmani1997/openebs_exporter/collector"
)

var (
	listenAddress        = flag.String("web.listen", ":9500", "Address on which to expose metrics and web interface.")
	metricsPath          = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	openEBSControllerURL = flag.String("controller.addr", "http://localhost:9501/", "Address of the OpenEBS controller monitoring.")
	namespace            = flag.String("namespace", "OpenEBS Controller", "Namespace for the OpenEBS Volume metrics.")
	homepage             = `<html>
			                <head><title>OpenEBS Exporter</title></head>
			                <body>
			                <h1>OpenEBS Exporter</h1>
			                <p><a href="` + *metricsPath + `">Metrics</a></p>
			                </body>
			                </html>`
)

func main() {
	flag.Parse()
	controllerURL, err := url.Parse(*openEBSControllerURL)

	if err != nil {
		log.Fatal(err)
	}

	exporter := collector.NewExporter(controllerURL)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	if *metricsPath == "" || *metricsPath == "/" {

		http.Handle(*metricsPath, promhttp.Handler())

	} else {

		http.Handle(*metricsPath, promhttp.Handler())
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			w.Write([]byte(homepage))
		})
	}

	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
