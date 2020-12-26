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
