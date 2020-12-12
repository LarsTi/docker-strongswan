package main
import (
        "log"
        "context"
	"time"
        "github.com/strongswan/govici/vici"
)
func monitorConns(s *vici.Session){
        if err := s.Subscribe("ike-updown", "child-updown"); err != nil {
                log.Panicln(err)
                return
        }
        for {
                e, err := s.NextEvent(context.Background())
                if err != nil {
                        log.Println(err)
                        log.Panicln("Assuming vici went down, shutting down this application")
                        break
                }
		k := e.Message.Keys()
		if k == nil {
			continue
		}
		log.Printf("[%s] %s | %s", e.Name, k, e.Message)
		for _,v := range k {
			if(v == "up"){
				//ist ein internes fragment
				continue
			}
			if e.Name == "ike-updown" {
				if e.Message.Get("up") != nil && e.Message.Get("up") == "yes" {
					//The IKE went up
					log.Printf("IKE %s is up\n", v)
				}else{
					log.Printf("IKE %s went down, requeuing it\n", v)
					ike_start := ike_to_start {
						name: v,
						isIke: true,
						message_send: false,
						last_try: time.Now(),
					}
					ch_ike_to_start <- ike_start
				}
			} else if e.Name == "child-updown" {
				if e.Message.Get("up") != nil && e.Message.Get("up") == "yes"{
					log.Printf("A child of IKE %s is up\n", v)
				}else{
					log.Printf("A child of IKE %s went down, terminating IKE", v)
					terminate(s,v, true)
				}
			}
		}

	}
}

