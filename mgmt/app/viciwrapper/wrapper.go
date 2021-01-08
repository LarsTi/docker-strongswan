package viciwrapper
import (
	"github.com/strongswan/govici/vici"
	"log"
)
var me *ViciWrapper

func GetWrapper() (*ViciWrapper, error) {
	if me != nil {
		return me, nil
	}
	//Singleton not yet created:
	me = &ViciWrapper{}
	me.startCommand()
	s, err := vici.NewSession()
	me.endCommand(err)
	if err != nil {
		return &ViciWrapper{}, err
	}
	me.session = s
	me.checkChannel = make(chan string, 100)
	me.terminateChannel = make(chan loadConnection, 10)
	me.initiateChannel = make(chan loadConnection, 10)
	me.saNameSuffix = "-net"
	return me, nil
}
func (w *ViciWrapper) GetViciMetrics() ViciMetrics{
	secrets, err := w.countSecrets()
	if err != nil {
		log.Println(err)
		secrets = 0
	}
	return ViciMetrics{
		CounterCommands: w.counterCommands,
		CounterErrors: w.counterErrors,
		LastCommand: w.lastCommand,
		ExecDuraLast: w.execDuraLast,
		ExecDuraAvgNs: w.execDuraAvgMs,
		LoadedSecrets: int64(secrets),
	}
}
func (w *ViciWrapper) ReadSecret(pathToFile string) error {
	return w.loadSharedSecret(pathToFile)
}
func (w *ViciWrapper) ReadConnection(pathToFile string) error {
	_, err := w.loadConn(pathToFile)
	return err
}
func (w *ViciWrapper) ListIkes()([]LoadedIKE, error){
	return w.listSAs()
}
func (w *ViciWrapper) WatchIkes(){
	w.watchIkes()
}
func (w *ViciWrapper) GetIkesInSystem() int {
	return len(w.ikesInSystem)
}

