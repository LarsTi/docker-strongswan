package viciwrapper
import (
        "log"
        "context"
	"time"
	"fmt"
)
func (v *ViciWrapper) monitorConns(){
	return
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
func (v *ViciWrapper) watchIkes() {
	log.Printf("[watch] Start watching for %d ikes\n", len(v.ikesInSystem))
	ticker := time.NewTicker(20 * time.Second)

	for {
		select {
			case ikeName := <-v.checkChannel:
				ike,child,err := v.findIke(ikeName)
				if ike == 1 && child == 1 && err == nil {
					log.Printf("[watch] Ike %s is connected, check completed\n", ikeName)
					continue
				}else if err != nil{
					v.reinitiateConn(ikeName)
				}else if (ike > 1 || child > 1){
					log.Printf("[watch] IKE %s (or child of IKE) is multiple times connected (Ike: %d, Child: %d\n", ikeName, ike, child)
					log.Printf("[watch] Terminating IKE %s and rechecking it!\n", ikeName)
					c, errC := connectionFromFile(ikeName)
					if errC != nil {
						log.Printf("[watch] [%s] %s\n", ikeName, errC)
						continue
					}
					errT := c.terminate(v)
					if errT != nil {
						log.Printf("[watch] [%s] %s\n", ikeName, errC)
						continue
					}
					log.Printf("[watch] IKE %s was terminated, so it should be able to be checked again. Requeuing it!\n", ikeName)
					v.checkChannel <- ikeName
				}
			case <- ticker.C:
				for _,ikeName := range v.ikesInSystem {
					v.checkChannel <- ikeName
				}
		}
	}
}

func (v *ViciWrapper) reinitiateConn(ikeName string) (bool, error){
	log.Printf("Received %s to restart\n", ikeName)
	log.Printf("[%s] waiting 5 seconds, so we are sure!\n", ikeName)
	time.Sleep(5 * time.Second)
	log.Printf("[%s] checking if it was restarted\n", ikeName)
	
	ikeCnt,saCnt, errIke := v.findIke(ikeName)
	if errIke != nil {
		return false, fmt.Errorf("[%s] %s", ikeName, errIke)
	}
	foundIke := false
	foundSA := false
	if ikeCnt == 1 {
		foundIke = true
	}
	if saCnt == 1 {
		foundSA = true
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
		v.checkChannel <- ikeName
		return false, nil
	}else if foundIke == false && foundSA == false {
		conn, errConn := connectionFromFile(ikeName)
		if errConn != nil {
			return false, fmt.Errorf("[%s] could not load Connection %s", ikeName, errConn)
		}
		log.Printf("[%s] trying to connect\n", ikeName)
		if errInitiate := conn.initiateConnection(v); errInitiate != nil {
			log.Printf("[%s] could not connect: %s\n", ikeName, errInitiate)
			v.checkChannel <- ikeName
			return false, nil
		}
	}else if foundIke == true && foundSA == true {
		log.Printf("[%s] is up, no changes", ikeName)
		return true, nil
	}else {
		return false, fmt.Errorf("[%s] is in an invalid and unrecoverable state!\n", ikeName)
	}

	ikeCnt,saCnt, errIke = v.findIke(ikeName)
	if errIke != nil {
		return false, fmt.Errorf("[%s] %s", ikeName, errIke)
	}
	foundIke = false
	foundSA = false
	if ikeCnt == 1 {
		foundIke = true
	}
	if saCnt == 1 {
		foundSA = true
	}
	if foundIke == true && foundSA == true {
		log.Printf("[%s] is now connected!\n", ikeName)
		return true, nil
	}else{
		log.Printf("[%s] could not be reconnected!\n", ikeName)
		return false, nil
	}

}
func (v *ViciWrapper) findIke(ikeName string)(int, int, error){
	ikes, err := v.listSAs()
	if err != nil {
		log.Fatalf("[%s] %s", ikeName, err)
		return 0,0, err
	}
	ikeCnt := 0
	saCnt := 0
	for _, ike := range ikes {
		if ike.Name == ikeName {
			//log.Printf("[%s] IKE is connected\n", ikeName)
			ikeCnt ++
		}else {
			continue
		}
		if (len(ike.Children) > 0) {
			saCnt ++
			//log.Printf("[%s] SA is connected, nothing to do\n", ikeName)
		}
	}

	return ikeCnt, saCnt, nil

}
