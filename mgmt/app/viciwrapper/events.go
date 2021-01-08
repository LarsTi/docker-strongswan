package viciwrapper
import (
        "log"
        "context"
	"time"
)
func (v *ViciWrapper) monitorConns(){
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
			v.checkChannel <- value
		}

	}
}
func (v *ViciWrapper) watchIkes() {
	go v.monitorConns()
	log.Printf("[watch] Start watching for %d ikes\n", len(v.ikesInSystem))
	ticker := time.NewTicker(20 * time.Second)

	for {
		select {
			case conn := <- v.terminateChannel:
				log.Printf("[%s] received to terminate\n", conn.Name)
				if errTerminate := conn.terminate(v); errTerminate != nil {
					log.Printf("[%s] could not terminate Connection: %s\n", conn.Name, errTerminate)
					continue
				}
				v.initiateChannel <- conn
			case conn := <- v.initiateChannel:
				log.Printf("[%s] received to initiate\n", conn.Name)
				if errInitiate := conn.initiateConnection(v); errInitiate != nil {
					log.Printf("[%s] could not connect: %s\n", conn.Name, errInitiate)
					continue
				}
				v.checkChannel <- conn.Name
			case ikeName := <- v.checkChannel:
				ikeCnt,saCnt, _ := v.findIke(ikeName)
				if ikeCnt == 1 && saCnt == 1 {
					log.Printf("[%s] is correct connected and operational\n", ikeName)
					continue
				}
				conn, errConn := v.connectionFromFile(ikeName)
				if errConn != nil {
					log.Printf("[%s] can not be read correctly, ignoring\n", ikeName)
					continue
				}
				log.Printf("[%s] is wrong connected, %d Ikes found, %d children found\n", ikeName, ikeCnt, saCnt)
				if ikeCnt > 1 || saCnt > 1 {
					//irgendwas ist hier falsch, erstmal disconnecten
					v.terminateChannel <- conn
				}else {
					//ike ist connected, sa nicht, also einfach neu connecten lassen
					v.initiateChannel <- conn
				}
			case <- ticker.C:
				for _,ikeName := range v.ikesInSystem {
					v.checkChannel <- ikeName
				}
		}
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
			saCnt += len(ike.Children)
			//log.Printf("[%s] SA is connected, nothing to do\n", ikeName)
		}
	}
	return ikeCnt, saCnt, nil

}
