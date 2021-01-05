package main
import (
        "log"
        "github.com/strongswan/govici/vici"
        "time"
)
//up to 100 ikes queued
var ch_ike_to_start = make(chan ike_to_start, 100)
type ike_to_start struct {
        name			string
        isIke			bool
	message_send		bool
        last_try		time.Time
}
var saNameSuffix string
func main() {
	saNameSuffix = "-net"
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
		_, err := loadConn(v, f)
		if err != nil {
			log.Printf("[%s] connection not loaded: %s\n", f, err)
		}else{
			log.Printf("[%s] connection loaded successful\n", f)
		}
        }
        go monitorConns(v)
	go runPrometheus(v)
	for {
		time.Sleep(1 * time.Second)
		listSAs(v)
	}
}

