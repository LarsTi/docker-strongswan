package webwrapper
type sharedSecret struct {
	Path		string	`json:path,omitempty`
	Typ		string	`json:typ,omitempty`
	Data		string	`json:data,omitempty`
	Owners		string	`json:owner,omitempty`
}
type connection struct {
	Path		string	`json:path,omitempty`
	LocalAddrs	string	`json:local,omitempty`
	RemoteAddrs	string	`json:remote,omitempty`
	Version		int	`json:version,omitempty`
	Proposals	string	`json:proposals,omitempty`
	DpdDelay	string	`json:dpdDelay,omitempty`
	DpdTimeout	string	`json:dpdTimeout,omitempty`
	LocalTS		string	`json:localTS,omitempty`
	RemoteTS	string	`json:remoteTS,omitempty`
	ChildProposals	string	`json:childProposals,omitempty`
}
type retMessage struct {
	Path		string	`json:path,omitempty`
	Error		string	`json:err,omitempty`
	Success		string	`json:success,omitempty`
}
