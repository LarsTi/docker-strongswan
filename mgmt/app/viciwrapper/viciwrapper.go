package viciwrapper
import (
	"github.com/strongswan/govici/vici"
)
type ViciWrapper struct {
	ViciStruct		*viciStruct
	initialized		bool
	secretsInSystem		[]string
}
var ch_ike_to_check = make(chan string, 100)
var me *ViciWrapper
var saNameSuffix		string
var ikesInSystem		[]string

func GetWrapper() (*ViciWrapper, error) {
	if me != nil {
		return me, nil
	}
	//Singleton not yet created:
	saNameSuffix = "-net"
	me = &ViciWrapper{}
	me.ViciStruct = &viciStruct{}
	me.ViciStruct.startCommand()
	s, err := vici.NewSession()
	me.ViciStruct.endCommand(err)
	if err != nil {
		return &ViciWrapper{}, err
	}
	me.ViciStruct.session = s

	return me, nil
}
func (w *ViciWrapper) ReadSecret(pathToFile string) error {
	return loadSharedSecret(w.ViciStruct, pathToFile)
}
func (w *ViciWrapper) ReadConnection(pathToFile string) error {
	found := false
	for _, loaded := range ikesInSystem {
		if loaded == pathToFile {
			found = true
			break
		}
	}
	if found == false {
		ikesInSystem = append(ikesInSystem, pathToFile)
	}
	_, err := loadConn(w.ViciStruct, pathToFile)
	return err
}
func (w *ViciWrapper) ListIkes()([]LoadedIKE, error){
	return listSAs(w.ViciStruct)
}
func (w *ViciWrapper) WatchIkes(){
	watchIkes(w.ViciStruct)
}
func (w *ViciWrapper) MonitorConns(){
	monitorConns(w.ViciStruct)
}
func (w *ViciWrapper) GetIkesInSystem() int {
	return len(ikesInSystem)
}

