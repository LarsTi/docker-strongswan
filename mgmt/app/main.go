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
func main() {
        s, err := vici.NewSession()
        if err != nil {
                log.Println("Could not load vici, going down")
                log.Panicln(err)
                return
        }
	log.Println("Vici loaded, starting operations")
	a := getFiles()
        for _, f := range a {
                loadSharedSecret(s, f)
        }
        for _, f := range a {
                if f == "me" {
                        log.Println("me is not a valid Connection, but a SharedSecret!")
                        continue
                }
                loadConn(s, f)
        }
	log.Printf("Found %d Secrets", countSecrets(s))
	go connectionIniator(s)
        go monitorConns(s)
	for {
		time.Sleep(1 * time.Second)
	}
}

