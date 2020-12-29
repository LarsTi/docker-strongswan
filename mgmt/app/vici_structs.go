package main
import (
	"time"
	"github.com/strongswan/govici/vici"
)
type viciStruct struct {
	session		*vici.Session
	counterCommands int64
	counterErrors	int64
	lastCommand	time.Time
	execDuraLast	time.Duration
	execDuraAvgMs	int64
}
func (v viciStruct) startCommand(){
	v.lastCommand = time.Now()
}
func (v viciStruct) endCommand(hasError error ){
	v.execDuraLast = time.Since(v.lastCommand)
	if hasError != nil {
		v.counterErrors ++
	}
	v.execDuraAvgMs = ( v.execDuraAvgMs * v.counterCommands + v.execDuraLast.Milliseconds() ) / ( v.counterCommands + 1)
	v.counterCommands ++
}
type sharedSecret struct{
	Id		string			`vici:"id"`
	Typ		string			`vici:"type"`
	Data		string			`vici:"data"`
	Owners		[]string		`vici:"owners"`
}
type loadConnection struct{
	Name		string
	LocalAddrs	[]string		`vici:"local_addrs"`
	RemoteAddrs	[]string		`vici:"remote_addrs"`
	Local		AuthOpts		`vici:"local"`
	Remote		AuthOpts		`vici:"remote"`
	ChildName	string
	Children	map[string]ChildSA	`vici:"children"`
	Version		int			`vici:"version"`
	Proposals	[]string		`vici:"proposals"`
	DpdDelay	string			`vici:"dpd_delay"`
	DpdTimeout	string			`vici:"dpd_timeout"`
	Mobike		string			`vici:"mobike"`
	Encap		string			`vici:"encap"`
}
type AuthOpts struct{
	Auth		string			`vici:"auth"`
	ID		string			`vici:"id"`
}
type ChildSA struct {
	LocalTS		[]string		`vici:"local_ts"`
	RemoteTS	[]string		`vici:"remote_ts"`
	Proposals	[]string		`vici:"esp_proposals"`
}
type loadedIKE struct {
	Name		string
	Version		int			`vici:"version"`
	State		string			`vici:"state"`
	LocalHost	string			`vici:"local-host"`
	RemoteHost	string			`vici:"remote-host"`
	Initiator	string			`vici:"initiator"`
	NatRemote	string			`vici:"nat-remote"`
	NatFake		string			`vici:"nat-fake"`
	EncrAlg		string			`vici:"encr-alg"`
	EncrKey		int			`vici:"encr-keysize"`
	IntegAlg	string			`vici:"integ-alg"`
	IntegKey	string			`vici:"integ-keysize"`
	DHGroup		string			`vici:"dh-group"`
	EstablishSec	int64			`vici:"established"`
	RekeySec	int64			`vici:"rekey-time"`
	ReauthSec	int64			`vici:"reauth-time"`
	Children	map[string]loadedChilds `vici:"child-sas"`
}
type loadedChilds struct {
	Name		string			`vici:"name"`
	State		string			`vici:"state"`
	Mode		string			`vici:"mode"`
	Protocol	string			`vici:"protocol"`
	Encap		string			`vici:"encap"`
	EncrAlg		string			`vici:"encr-alg"`
	EncrKey		int			`vici:"encr-keysize"`
	IntegAlg	string			`vici:"integ-alg"`
	IntegKey	string			`vici:"integ-keysize"`
	DHGroup		string			`vici:"dh-group"`
	BytesIn		int64			`vici:"bytes-in"`
	PacketsIn	int64			`vici:"bytes-out"`
	LastInSec	int64			`vici:"use-in"`
	BytesOut	int64			`vici:"bytes-out"`
	PacketsOut	int64			`vici:"bytes-out"`
	LastOutSec	int64			`vici:"use-out"`
	EstablishSec	int64			`vici:"install-time"`
	RekeySec	int64			`vici:"rekey-time"`
	LifeTimeSec	int64			`vici:"life-time"`
	LocalTs		[]string		`vici:"local-ts"`
	RemoteTS	[]string		`vici:"remote-ts"`
}
