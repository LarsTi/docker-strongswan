package main
import (
	"log"
	"time"
	"github.com/strongswan/govici/vici"
)
func connectionIniator(s *vici.Session){
	for{
		select{
		case ike_start := <-ch_ike_to_start:
			if ike_start.name == "" || ike_start.name == "(unnamed)" {
				log.Printf("Ignoring unnamed IKE\n")
				continue
			}
			if time.Since(ike_start.last_try) < (2 * time.Second) {
				if ike_start.message_send == false {
				log.Printf("IKE %s had last try at %d:%d:%d, waiting 2 seconds before restarting\n", ike_start.name,
					ike_start.last_try.Hour(), ike_start.last_try.Minute() ,ike_start.last_try.Second())
				}
				ike := ike_to_start {
					name: ike_start.name,
					isIke: ike_start.isIke,
					message_send: true,
					last_try: ike_start.last_try,
				}
				ch_ike_to_start <- ike
				continue
			}
			//terminate(s, ike_start.name, true)
			//IKE neu starten (erst IKE, dann Child)
			//initiate(s, ike_start.name, true)
			log.Printf("Got Ike to Start %s\n", ike_start.name)
			//b := initiate(s, ike_start.name, false)
			//if b == true {
			//	log.Printf("[%s] went up in connection Initiator\n", ike_start.name)
			//}else{
			//	log.Printf("[%s]Initiate was issued but did not work!\n", ike_start.name)
			//	ike := ike_to_start {
			//		name: ike_start.name,
			//		isIke: ike_start.isIke,
			//		message_send: false,
			//		last_try: time.Now(),
			//	}
			//	ch_ike_to_start <- ike
			//}

		}
	}
}
