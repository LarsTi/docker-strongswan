package viciwrapper

import (
        "fmt"
	"../filewrapper"
        "github.com/strongswan/govici/vici"
)
func (v *ViciWrapper) connectionFromFile(path string) (loadConnection, error){
	ret := loadConnection{
		Name: path,
		LocalAddrs: filewrapper.GetStringArrayFromPath(path, "LocalAddrs"),
		RemoteAddrs: filewrapper.GetStringArrayFromPath(path, "RemoteAddrs"),
		Local: AuthOpts { Auth: "psk", ID: filewrapper.GetStringValueFromPath("me.secret", "RemoteAddrs"), },
		Remote: AuthOpts { Auth: "psk", ID: filewrapper.GetStringValueFromPath(path, "RemoteAddrs"), },
		ChildName: path + v.saNameSuffix,
		Children: make(map[string]ChildSA),
		Version: filewrapper.GetIntValueFromPath(path, "Version"),
		Proposals: filewrapper.GetStringArrayFromPath(path, "proposals"),
		DpdDelay: "2s",
		Mobike: "no",
		Encap: "yes",
	}
	ret.Children[ret.ChildName] = ChildSA{
		LocalTS: filewrapper.GetStringArrayFromPath(path, "LocalTrafficSelectors"),
		RemoteTS: filewrapper.GetStringArrayFromPath(path, "RemoteTrafficSelectors"),
		Proposals: filewrapper.GetStringArrayFromPath(path, "ESPProposals"),
	}

	//TODO: check if everything is set!
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
