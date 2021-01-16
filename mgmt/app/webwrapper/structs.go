package webwrapper
type sharedSecret struct {
	Path		string	`json:Path,omitempty`
	Typ		string	`json:Typ,omitempty`
	Data		string	`json:Data,omitempty`
	Owners		string	`json:Owners,omitempty`
}
type connection struct {
	Path		string	`json:Path,omitempty`
	LocalAddrs	string	`json:LocalAddrs,omitempty`
	RemoteAddrs	string	`json:RemoteAddrs,omitempty`
	Version		string	`json:Version,omitempty`
	Proposals	string	`json:Proposals,omitempty`
	DpdDelay	string	`json:DpdDelay,omitempty`
	DpdTimeout	string	`json:DpdTimeout,omitempty`
	LocalTS		string	`json:LocalTS,omitempty`
	RemoteTS	string	`json:RemoteTS,omitempty`
	ChildProposals	string	`json:ChildProposals,omitempty`
}
type retMessage struct {
	Path		string	`json:Path,omitempty`
	Error		string	`json:Err,omitempty`
	Success		string	`json:Success,omitempty`
}
