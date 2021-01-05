package main

import (
        "fmt"
        "github.com/strongswan/govici/vici"
)
func countSecrets(v *viciStruct) (int, error) {
        v.startCommand()
	loaded, e := v.session.CommandRequest("get-shared", nil)
        v.endCommand(e)
	if e != nil {
                return 0, fmt.Errorf("[get-shared] %s\n", e)
        }
	
	return len(loaded.Get("keys").([]string)), nil
}
func isSecretLoaded(v *viciStruct, secretId string) (bool, error){
	v.startCommand()
        loaded, e := v.session.CommandRequest("get-shared", nil)
	v.endCommand(e)
        if e != nil {
                return false, fmt.Errorf("[get-shared] %s\n", e)
        }
        for _, value := range loaded.Get("keys").([]string){
                if value == secretId {
                        return true, nil
                }
        }
        return false, nil
}
func unloadSecret(v *viciStruct, secretId string) error{
        m := vici.NewMessage()
        if err := m.Set("id", secretId); err != nil {
                return fmt.Errorf("[unload-shared] %s\n", err)
        }
        v.startCommand()
	_, e := v.session.CommandRequest("unload-shared", m)
	v.endCommand(e)
        if e != nil {
                return fmt.Errorf("[unload-shared] %s\n", e)
        }
        return nil

}
func loadSharedSecret(v *viciStruct, path string) error{
        psk := sharedSecret{
                Id: getStringValueFromPath(path, "RemoteAddrs"),
                Typ: "IKE",
                Data: getStringValueFromPath(path, "PSK"),
                Owners: getStringArrayFromPath(path, "RemoteAddrs"),
        }
	if psk.Data == "" {
		return fmt.Errorf("Secret in file %s is no PSK\n", path)
	}
	isLoaded, err := isSecretLoaded(v, psk.Id)
	if err != nil {
		return err
	}else if isLoaded {
                unloadSecret(v, psk.Id)
        }
        m, err := vici.MarshalMessage(psk)
        if err != nil {
                return fmt.Errorf("[%s] %s\n",path, err)
        }
	v.startCommand()
        _, err2 := v.session.CommandRequest("load-shared", m)
	v.endCommand(err2)
        if err2 != nil {
                return fmt.Errorf("[%s] %s\n", path, err2)
        }
        return nil
}
