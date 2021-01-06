package main
import (
        "log"
        "github.com/strongswan/govici/vici"
        "time"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
//up to 100 ikes queued
var ch_ike_to_check = make(chan string, 100)
var ikesInSystem []string

var saNameSuffix string
func main() {
	saNameSuffix = "-net"
	
	//Initializing vici
	start := time.Now()
        s, err := vici.NewSession()
	end := time.Now()
        if err != nil {
                log.Println("Could not load vici, going down")
                log.Panicln(err)
                return
        }
	v := &viciStruct {
		session: s,
		counterCommands: 1,
		lastCommand: start,
		execDuraLast: end.Sub(start),
		execDuraAvgMs: end.Sub(start).Milliseconds(),
	}
	
	log.Println("Vici loaded, starting operations")
	
	//Initializing Connectiosn
	a := getFiles()
        for _, f := range a {
		e := loadSharedSecret(v, f)
		if e != nil {
			log.Printf("[%s] Shared Secret not loaded: %s\n", f, e)
		}else{
			log.Printf("[%s] Shared Secret loaded\n", f)
		}
        }
        for _, f := range a {
                if f == "me" {
                        log.Println("me is not a valid Connection, but a SharedSecret!")
                        continue
                }
		
		ikesInSystem = append(ikesInSystem, f)

		_, err := loadConn(v, f)
		if err != nil {
			log.Printf("[%s] connection not loaded: %s\n", f, err)
		}else{
			log.Printf("[%s] connection loaded successful\n", f)
		}
        }

	//Initializing Collectors for Prometheus:
	strongswanCollector := NewStrongswanCollector(v)
	strongswanCollector.init()
	http.Handle("/metrics", promhttp.Handler())

	//Starting monitoring Threads:
        go monitorConns(v)
	go watchIkes(v)

	//Running Prometheus (blocking):
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

