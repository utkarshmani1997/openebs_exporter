package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/utkarshmani1997/openebs_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Version of nsq_exporter. Set at build time.
const Version = "0.0.0.dev"

var (
	listenAddress        = flag.String("web.listen", ":9500", "Address on which to expose metrics and web interface.")
	metricsPath          = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	OpenEBSControllerURL = flag.String("controller.addr", "http://localhost:9501/v1/stats", "Address of the controller.")
	namespace            = flag.String("namespace", "OpenEBS", "Namespace for the OpenEBS metrics.")

	enabledCollectors = flag.String("collect", "stats.OpenEBS", "Comma-separated list of collectors to use.")
	statsRegistry     = map[string]func(namespace string) collector.StatsCollector{
		"OpenEBS": collector.VolumeStats,
	}
)

func main() {
	flag.Parse()

	ex, err := createOpenEBSExecutor()
	if err != nil {
		log.Fatalf("error creating nsq executor: %v", err)
	}
	prometheus.MustRegister(ex)

	http.Handle(*metricsPath, promhttp.Handler())
	if *metricsPath != "" && *metricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>OpenEBS Exporter</title></head>
			<body>
			<h1>OpenEBS Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	log.Print("listening to ", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createOpenEBSExecutor() (*collector.Executor, error) {

	openEBSControllerURL, err := normalizeURL(*OpenEBSControllerURL)
	if err != nil {
		return nil, err
	}

	ex := collector.NewExecutor(*namespace, openEBSControllerURL)

	for _, param := range strings.Split(*enabledCollectors, ",") {
		param = strings.TrimSpace(param)
		parts := strings.SplitN(param, ".", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid collector name: %s", param)
		}
		if parts[0] != "stats" {
			return nil, fmt.Errorf("invalid collector prefix: %s", parts[0])
		}

		name := parts[1]
		c, has := statsRegistry[name]
		if !has {
			return nil, fmt.Errorf("unknown stats collector: %s", name)
		}
		ex.Use(c(*namespace))
	}
	return ex, nil
}

func normalizeURL(ustr string) (string, error) {
	ustr = strings.ToLower(ustr)
	if !strings.HasPrefix(ustr, "https://") && !strings.HasPrefix(ustr, "http://") {
		ustr = "http://" + ustr
	}

	u, err := url.Parse(ustr)
	if err != nil {
		return "", err
	}
	fmt.Println(u, u.String())
	if u.Path == "" {
		u.Path = "v1/stats"
	}
	return u.String(), nil
}
