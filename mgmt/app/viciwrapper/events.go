package viciwrapper
import (
        "log"
        "context"
	"time"
	"fmt"
	"../filewrapper"
)
func (v *ViciWrapper) monitorConns(){
	lastEventTime := make(map [string]time.Time)

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
		log.Printf("[%s] %s\n", e.Name, k)
		for _,value := range k {
			if(value == "up"){
				//ist ein internes fragment
				continue
			}
			if value == "" || value == "(unnamed)" {
				//ignoring unnamed SAs
				continue
			}
			lastEvent, ok := lastEventTime[value]
			if (ok && time.Since(lastEvent) > 20 * time.Second){
				v.checkChannel <- value
			}
			lastEventTime[value] = time.Now()
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
				conn, errConn := v.connectionFromFile(ikeName)
				if errConn != nil{
					v.UnloadConnection(ikeName)
					v.ReadConnection(ikeName)
					continue
				}
				ike, err := v.findIke(ikeName)
				if err != nil {
					log.Println(err)
					v.terminateChannel <- conn
					continue
				}
				ikeExpected := v.ikesInSystem[ikeName]
				if ikeExpected.numberRemoteTS != ike.numberRemoteTS {
					log.Printf("[%s] Remote Traffic Selectors: expected %d, found\n", ikeName, ikeExpected.numberRemoteTS, ike.numberRemoteTS)
				}else if ikeExpected.numberLocalTS != ike.numberLocalTS {
					log.Printf("[%s] Local Traffic Selectors: expected %d, found\n", ikeName, ikeExpected.numberLocalTS, ike.numberLocalTS)
				}else if ikeExpected.numberChildren != ike.numberChildren {
					log.Printf("[%s] Children: expected %d, found\n", ikeName, ikeExpected.numberChildren, ike.numberChildren)
				}else{
					log.Printf("[%s] looks good!\n", ikeName)
					continue
				}
				if ikeExpected.numberRemoteTS > ike.numberRemoteTS ||
					ikeExpected.numberLocalTS > ike.numberLocalTS ||
					ikeExpected.numberChildren > ike.numberChildren {
					v.initiateChannel <- conn
				}else{
					v.terminateChannel <- conn
				}
			case <- ticker.C:
				for _,ike := range v.ikesInSystem {
					if ike.initiator == false {
						continue
					}
					v.checkChannel <- ike.ikeName
				}
		}
	}
}
func (v *ViciWrapper) checkIke(ikeName string) (bool, error){
	conn, errC := v.connectionFromFile(ikeName)
	if errC != nil {
		return false, errC
	}
	if conn.Name == ikeName {
		return true, nil
	}
	ikes, errS := v.listSAs()
	if errS != nil {
		return false, errS
	}
	for _, ike := range ikes {
		if ike.Name == ikeName {
			return true, nil
		}
	}
	return false, nil
}
func (v *ViciWrapper) findIke(ikeName string)(ikeInSystem, error){
	retVal := ikeInSystem{
		ikeName: ikeName,
		initiator: filewrapper.GetBoolValueFromPath(ikeName, "Initiator"),
		numberRemoteTS: 0,
		numberLocalTS: 0,
		numberChildren: 0,
	}
	ikes, err := v.listSAs()
	if err != nil {
		log.Fatalf("[%s] %s", ikeName, err)
		return retVal, err
	}
	ikeCnt := 0
	for _, ike := range ikes {
		if ike.Name == ikeName {
			ikeCnt ++
		}else {
			continue
		}
		for _, child := range ike.Children {
			retVal.numberChildren += 1
			retVal.numberRemoteTS += len(child.RemoteTS)
			retVal.numberLocalTS += len(child.LocalTS)
		}
	}
	if ikeCnt != 1 {
		return retVal, fmt.Errorf("[%s] there are %d ikes connected, 1 expected!", ikeName, ikeCnt)
	}
	return retVal, nil

}
