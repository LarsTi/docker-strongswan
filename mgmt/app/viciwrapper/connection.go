package viciwrapper

import (
        "fmt"
	"../filewrapper"
        "github.com/strongswan/govici/vici"
)
func (v *ViciWrapper) connectionFromFile(path string) (loadConnection, error){
	ret := loadConnection{}
	if path == "" {
		return ret, fmt.Errorf("[connection] Empty path not allowed")
	}
	
	ret.DpdDelay = "2s"
	ret.Mobike = "no"
	ret.Encap = "yes"
	ret.Name = path
	ret.ChildName = path + v.saNameSuffix

	ret.LocalAddrs = filewrapper.GetStringArrayFromPath(path, "LocalAddrs")
	if (len(ret.LocalAddrs) == 0 || ret.LocalAddrs[0] == ""){
		return ret, fmt.Errorf("[%s] LocalAddrs not found in config file\n", path)
	}
	ret.RemoteAddrs = filewrapper.GetStringArrayFromPath(path, "RemoteAddrs")
	if (len(ret.RemoteAddrs) == 0 || ret.RemoteAddrs[0] == ""){
		return ret, fmt.Errorf("[%s] RemoteAddrs not found in config file\n", path)
	}
	ret.Version = filewrapper.GetIntValueFromPath(path, "Version")
	if ret.Version == 0 {
		return ret, fmt.Errorf("[%s] Version not found in config file\n", path)
	}
	ret.Proposals = filewrapper.GetStringArrayFromPath(path, "proposals")
	if (len(ret.Proposals) == 0 || ret.Proposals[0] == "") {
		return ret, fmt.Errorf("[%s] proposals not found in config file\n", path)
	}
	ret.Local = AuthOpts{ Auth: "psk", ID: filewrapper.GetStringValueFromPath("me.secret", "RemoteAddrs"), }
	if ret.Local.ID == "" {
		return ret, fmt.Errorf("[%s] RemoteAddrs not found in config file\n", "me.secret")
	}
	ret.Remote = AuthOpts{ Auth: "psk", ID: filewrapper.GetStringValueFromPath(path, "RemoteAddrs"), }
	if ret.Remote.ID == "" {
		return ret, fmt.Errorf("[%s] RemoteAddrs not found in config file\n", path)
	}
	ret.Children = make(map[string]ChildSA)
	child := ChildSA{}
	child.LocalTS = filewrapper.GetStringArrayFromPath(path, "LocalTrafficSelectors")
	if len(child.LocalTS) == 0 || child.LocalTS[0] == "" {
		return ret, fmt.Errorf("[%s] LocalTrafficSelectors not found in config file\n", path)
	}
	child.RemoteTS = filewrapper.GetStringArrayFromPath(path, "RemoteTrafficSelectors")
	if len(child.RemoteTS) == 0 || child.RemoteTS[0] == "" {
		return ret, fmt.Errorf("[%s] RemoteTrafficSelectors not found in config file\n", path)
	}
	child.Proposals = filewrapper.GetStringArrayFromPath(path, "ESPProposals")
	if len(child.Proposals) == 0 || child.Proposals[0] == "" {
		return ret, fmt.Errorf("[%s] ESPProposals not found inf config file\n", path)
	}
	ret.Children[ret.ChildName] = child
	
	return ret, nil
}
func (c loadConnection) unloadConnection(v *ViciWrapper) error {
	m := vici.NewMessage()
        if err := m.Set("name", c.Name); err != nil {
                return fmt.Errorf("[unload-conn] %s\n", err)
        }
        v.startCommand()
	_, err := v.session.CommandRequest("unload-conn", m)
	v.endCommand(err)
        if err != nil{
                return fmt.Errorf("[unload-conn] %s\n", err)
        }
        return nil
}
func (c loadConnection) loadConnection(v *ViciWrapper) error {
	msg, err := vici.MarshalMessage(c)
	if err != nil {
		return fmt.Errorf("[load-conn] %s\n", err)
	}
	m := vici.NewMessage()
	m.Set(c.Name, msg)
	v.startCommand()
	_, e := v.session.CommandRequest("load-conn", m)
	v.endCommand(e)
	if e != nil {
		return fmt.Errorf("[load-conn] %s\n", e)
	}
	return nil
}
func (c loadConnection) reload(v *ViciWrapper) error {
	c.unloadConnection(v)
	return c.loadConnection(v)
}
func (c loadConnection) initiateConnection(v *ViciWrapper) error {
	m := vici.NewMessage()
	if err := m.Set("child", c.ChildName); err != nil{
		return fmt.Errorf("[initiate] %s\n", err)
	}
	if err := m.Set("ike", c.Name); err != nil {
		return fmt.Errorf("[initiate] %s\n", err)
	}
	v.startCommand()
	_, err := v.session.CommandRequest("initiate", m)
	v.endCommand(err)
	if err != nil {
		return fmt.Errorf("[initiate] %s\n", err)
	}
	return nil
}
func (c loadConnection) terminate(v *ViciWrapper) error {
	m := vici.NewMessage()
	if err := m.Set("ike", c.Name); err != nil {
		return fmt.Errorf("[terminate] %s\n", err)
	}
	if err := m.Set("force", true); err != nil {
		return fmt.Errorf("[terminate] %s\n", err)
	}
	if err := m.Set("timeout", 1000); err != nil {
		return fmt.Errorf("[terminate] %s\n", err)
	}
	v.startCommand()
	_, err := v.session.CommandRequest("terminate", m)
	v.endCommand(err)
	if err != nil {
		return fmt.Errorf("[terminate] %s\n", err)
	}
	return nil
}
func (w *ViciWrapper) loadConn(path string) (loadConnection, error){
	found := false
	for _, loaded := range w.ikesInSystem {
		if loaded == path {
			found = true
			break
		}
	}
	if found == false {
		w.ikesInSystem = append(w.ikesInSystem, path)
	}

	c, e := w.connectionFromFile(path)
	if e != nil {
		return c, e
	}
	err := c.reload(w)
	if err != nil {
		return c, err
	}
	err = c.initiateConnection(w)
	if err != nil {
		return c, err
	}
	return c, nil
}
func (v *ViciWrapper) listSAs()([]LoadedIKE, error){
	var retVar []LoadedIKE
	v.startCommand()
	msgs, err := v.session.StreamedCommandRequest("list-sas", "list-sa", nil)
	v.endCommand(err)
	if err != nil {
		return retVar, err
	}
	for _,m := range msgs.Messages() {
		if e := m.Err(); e != nil{
			//ignoring this error
			continue
		}
		for _, k := range m.Keys() {
			inbound := m.Get(k).(*vici.Message)
			var ike LoadedIKE
			if e := vici.UnmarshalMessage(inbound, &ike); e != nil {
				//ignoring this marshal/unmarshal errro!
				continue
			}
			ike.Name = k
			retVar = append(retVar, ike)

		}
	}
	return retVar, nil
}
