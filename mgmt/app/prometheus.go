package main
import (
	"net/http"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
func runPrometheus(v *viciStruct){
	viciCollector := NewViciCollector(v)
	strongswanCollector := NewStrongswanCollector(v)
	viciCollector.init()
	strongswanCollector.init()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
