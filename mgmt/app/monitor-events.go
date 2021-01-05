package main
import (
        "log"
        "context"
	"time"
	//"strings"
        //"github.com/strongswan/govici/vici"
)
func monitorConns(v *viciStruct){
	v.startCommand()
        if err := v.session.Subscribe("ike-updown", "child-updown"); err != nil {
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
				go reinitiateConn(v, value, e.Name, 0)
			}
		}

	}
}
func reinitiateConn(v *viciStruct, ikeName string, event string, cnt int) {
	log.Printf("Received %s to restart\n", ikeName)
	if event == "ike-updown" {
		log.Println("Ignoring ike, only responding to child-updown")
		return
	}
	if cnt > 10 {
		log.Printf("Please check ike %s, it was tried 10 times before giving up!", ikeName)
		return
	}
	log.Printf("[%s] waiting 5 seconds, so we are sure!\n", ikeName)
	time.Sleep(5 * time.Second)
	log.Printf("[%s] checking if it was restarted\n", ikeName)
	
	foundIke, foundSA, errIke := findIke(v, ikeName)
	if errIke != nil {
		log.Fatalf("[%s] %s\n", ikeName, errIke)
		return
	}
	
	if foundIke == true && foundSA == false {
		c, errC := connectionFromFile(ikeName)
		if errC != nil {
			log.Panicf("[%s] could not load Connection: %s\n", ikeName, errC)
			return
		}
		if errTerminate := c.terminate(v); errTerminate != nil {
			log.Printf("[%s] could not terminate Connection: %s\n", ikeName, errTerminate)
		}
		log.Printf("[%s] should now be terminated, restarting process to initiate\n", ikeName)
		reinitiateConn(v, ikeName, event, (cnt + 1))
		log.Printf("[%s] should now be connected!\n", ikeName)
		return
	}else if foundIke == false && foundSA == false {
		conn, errConn := connectionFromFile(ikeName)
		if errConn != nil {
			log.Panicf("[%s] could not load Connection %s\n", ikeName, errConn)
			return
		}
		log.Printf("[%s] trying to connect\n", ikeName)
		if errInitiate := conn.initiateConnection(v); errInitiate != nil {
			log.Printf("[%s] could not connect: %s\n", ikeName, errInitiate)
			reinitiateConn(v, ikeName, event, (cnt + 1) )
			log.Printf("[%s] should now be connected!\n", ikeName)
			return
		}
	}else if foundIke == true && foundSA == true {
		log.Printf("[%s] is up, no changes", ikeName)
		return
	}else {
		log.Panicf("[%s] is in an invalid and unrecoverable state!\n", ikeName)
		return
	}

	foundIke, foundSA, errIke = findIke(v, ikeName)
	if errIke != nil {
		log.Fatalf("[%s] %s\n", ikeName, errIke)
		return
	}
	if foundIke == true && foundSA == true {
		log.Printf("[%s] is now connected!\n", ikeName)
	}else{
		log.Printf("[%s] could not be reconnected!\n", ikeName)
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
			log.Printf("[%s]IKE is connected\n", ikeName)
		}else {
			continue
		}
		if (len(ike.Children) > 0) {
			log.Printf("[%s] SA is connected, nothing to do\n", ikeName)
			return true, true, nil
		}
		return true, false, nil
	}
	return false, false, nil

}
