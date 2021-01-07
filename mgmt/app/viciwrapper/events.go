package viciwrapper
import (
        "log"
        "context"
	"time"
	"fmt"
)
func monitorConns(v *viciStruct){
	v.startCommand()
        if err := v.session.Subscribe("child-updown"); err != nil {
		v.endCommand(err)
                log.Panicln(err)
                return
        }
	v.endCommand(nil)
        for {
                e, err := v.session.NextEvent(context.Background())
                if err != nil {
                        log.Println(err)
                        log.Panicln("Assuming vici went down, shutting down this application")
                        break
                }
		k := e.Message.Keys()
		if k == nil {
			continue
		}
		log.Printf("[%s] %s | %s\n", e.Name, k, e.Message)
		for _,value := range k {
			if(value == "up"){
				//ist ein internes fragment
				continue
			}
			if value == "" || value == "(unnamed)" {
				//ignoring unnamed SAs
				continue
			}
			if e.Message.Get("up") != nil && e.Message.Get("up") == "yes" {
				log.Printf("[%s] %s went up\n", e.Name, value)
			}else{
				ch_ike_to_check <- value
			}
		}

	}
}
func reinitiateConn(v *viciStruct, ikeName string) (bool, error){
	log.Printf("Received %s to restart\n", ikeName)
	log.Printf("[%s] waiting 5 seconds, so we are sure!\n", ikeName)
	time.Sleep(5 * time.Second)
	log.Printf("[%s] checking if it was restarted\n", ikeName)
	
	foundIke, foundSA, errIke := findIke(v, ikeName)
	if errIke != nil {
		return false, fmt.Errorf("[%s] %s", ikeName, errIke)
	}
	
	if foundIke == true && foundSA == false {
		c, errC := connectionFromFile(ikeName)
		if errC != nil {
			return false, fmt.Errorf("[%s] could not load Connection: %s\n", ikeName, errC)
		}
		if errTerminate := c.terminate(v); errTerminate != nil {
			log.Printf("[%s] could not terminate Connection: %s\n", ikeName, errTerminate)
		}
		log.Printf("[%s] should now be terminated, restarting process to initiate\n", ikeName)
		ch_ike_to_check <- ikeName
		return false, nil
	}else if foundIke == false && foundSA == false {
		conn, errConn := connectionFromFile(ikeName)
		if errConn != nil {
			return false, fmt.Errorf("[%s] could not load Connection %s", ikeName, errConn)
		}
		log.Printf("[%s] trying to connect\n", ikeName)
		if errInitiate := conn.initiateConnection(v); errInitiate != nil {
			log.Printf("[%s] could not connect: %s\n", ikeName, errInitiate)
			ch_ike_to_check <- ikeName
			return false, nil
		}
	}else if foundIke == true && foundSA == true {
		log.Printf("[%s] is up, no changes", ikeName)
		return true, nil
	}else {
		return false, fmt.Errorf("[%s] is in an invalid and unrecoverable state!\n", ikeName)
	}

	foundIke, foundSA, errIke = findIke(v, ikeName)
	if errIke != nil {
		return false, fmt.Errorf("[%s] %s", ikeName, errIke)
	}
	if foundIke == true && foundSA == true {
		log.Printf("[%s] is now connected!\n", ikeName)
		return true, nil
	}else{
		log.Printf("[%s] could not be reconnected!\n", ikeName)
		return false, nil
	}

}
func watchIkes(v *viciStruct) {
	log.Printf("[watch] Start watching for %d ikes\n", len(ikesInSystem))
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
			case ikeName := <-ch_ike_to_check:
				ike,child,err := findIke(v, ikeName)
				if ike == true && child == true && err == nil {
					log.Printf("[watch] Ike %s is connected, check completed\n", ikeName)
					continue
				}else{
					reinitiateConn(v, ikeName)
				}
			case <- ticker.C:
				for _,ikeName := range ikesInSystem {
					ch_ike_to_check <- ikeName
				}
			default:
				time.Sleep(1 * time.Second)
		}
	}
}
func findIke(v *viciStruct, ikeName string)(bool, bool, error){
	ikes, err := listSAs(v)
	if err != nil {
		log.Fatalf("[%s] %s", ikeName, err)
		return false, false, err
	}
	for _, ike := range ikes {
		if ike.Name == ikeName {
			//log.Printf("[%s] IKE is connected\n", ikeName)
		}else {
			continue
		}
		if (len(ike.Children) > 0) {
			//log.Printf("[%s] SA is connected, nothing to do\n", ikeName)
			return true, true, nil
		}
		return true, false, nil
	}
	return false, false, nil

}
