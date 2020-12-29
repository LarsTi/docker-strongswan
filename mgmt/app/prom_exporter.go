package main
import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
func runPrometheus(v *viciStruct){
	collector := NewViciCollector(v)
	collector.init()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
