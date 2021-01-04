package main

import (
        "fmt"
        "github.com/strongswan/govici/vici"
)
func connectionFromFile(path string) (loadConnection, error){
	ret := loadConnection{
		Name: path,
		LocalAddrs: getStringArrayFromPath(path, "LocalAddrs"),
		RemoteAddrs: getStringArrayFromPath(path, "RemoteAddrs"),
		Local: AuthOpts { Auth: "psk", ID: getStringValueFromPath("me", "RemoteAddrs"), },
		Remote: AuthOpts { Auth: "psk", ID: getStringValueFromPath(path, "RemoteAddrs"), },
		ChildName: path + "-net",
		Children: make(map[string]ChildSA),
		Version: getIntValueFromPath(path, "Version"),
		Proposals: getStringArrayFromPath(path, "proposals"),
		DpdDelay: "2s",
		Mobike: "no",
		Encap: "yes",
	}
	ret.Children[ret.ChildName] = ChildSA{
		LocalTS: getStringArrayFromPath(path, "LocalTrafficSelectors"),
		RemoteTS: getStringArrayFromPath(path, "RemoteTrafficSelectors"),
		Proposals: getStringArrayFromPath(path, "ESPProposals"),
	}

	//TODO: check if everything is set!
	return ret, nil
}
func (c loadConnection) unloadConnection(v *viciStruct) error {
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
func (c loadConnection) loadConnection(v *viciStruct) error {
	msg, err := vici.MarshalMessage(c)
	if err != nil {
		return fmt.Errorf("[load-conn] %s\n", err)
	}
	m := vici.NewMessage()
	m.Set(c.Name, msg)
	fmt.Println(msg)
	fmt.Println(m)
	fmt.Println(c.Children[c.ChildName])
	v.startCommand()
	_, e := v.session.CommandRequest("load-conn", m)
	v.endCommand(e)
	if e != nil {
		return fmt.Errorf("[load-conn] %s\n", e)
	}
	return nil
}
func (c loadConnection) reload(v *viciStruct) error {
	c.unloadConnection(v)
	return c.loadConnection(v)
}
func (c loadConnection) initiateConnection(v *viciStruct) error {
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
func (c loadConnection) terminate(v *viciStruct) error {
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
func loadConn(v *viciStruct, path string) (loadConnection, error){
	c, e := connectionFromFile(path)
	if e != nil {
		return c, e
	}
	err := c.reload(v)
	if err != nil {
		return c, err
	}
	err = c.initiateConnection(v)
	if err != nil {
		return c, err
	}
	return c, nil
func listSAs(v *viciStruct)([]loadedIKE, error){
	var retVar []loadedIKE
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
			var ike loadedIKE
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
}
