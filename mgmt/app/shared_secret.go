package main

import (
        "log"
        "github.com/strongswan/govici/vici"
)
type secret struct {    
        Id              string          `vici:"id"`
        Typ             string          `vici:"type"`
        Data            string          `vici:"data"`
        Owners          []string        `vici:"owners"`
}
func countSecrets(s *vici.Session) int {
        loaded, e := s.CommandRequest("get-shared", nil)
        if e != nil {
                log.Panicln(e)
                return 0
        }
	
	return len(loaded.Get("keys").([]string))
}
func isSecretLoaded(s *vici.Session, secretId string) bool{
        loaded, e := s.CommandRequest("get-shared", nil)
        if e != nil {
                log.Println("Can not read loaded secretes from charon")
                log.Panicln(e)
                return false
        }
        for _, value := range loaded.Get("keys").([]string){
                if value == secretId {
                        //log.Printf("Secret %s is loaded\n", secretId)
                        return true
                }
        }
        //log.Printf("Secret %s is not loaded\n", secretId)
        return false
}
func unloadSecret(s *vici.Session, secretId string) bool{
        log.Printf("Unloading Secret %s\n", secretId)
        m := vici.NewMessage()
        if err := m.Set("id", secretId); err != nil {
                log.Println("Could not create Unload-Shared-Secret Message")
                log.Panicln(err)
                return false
        }
        _, e := s.CommandRequest("unload-shared", m)
        if e != nil {
                log.Println("Could not unload Secret")
                log.Panicln(e)
                return false
        }
        //fmt.Printf("Unloaded Secret %s\n", secretId)
        return true

}
func loadSharedSecret(s *vici.Session, path string) bool{
        psk := secret{
                Id: getStringValueFromPath(path, "RemoteAddrs"),
                Typ: "IKE",
                Data: getStringValueFromPath(path, "PSK"),
                Owners: getStringArrayFromPath(path, "RemoteAddrs"),
        }
	if psk.Data == "" {
		log.Printf("Secret in file %s is no PSK\n", path)
		return false
	}
        log.Printf("Loading SharedSecret for %s from path %s\n", psk.Id, path)
        if isSecretLoaded(s, psk.Id) {
                log.Println("Secret existed, reloading it now!")
                unloadSecret(s, psk.Id)
        }
        m,err := vici.MarshalMessage(psk)
        if err != nil {
                log.Panicf("[%s]%s",path, err)
                return false
        }

        _, err2 := s.CommandRequest("load-shared", m)
        if err2 != nil {
                log.Printf("Could not load PSK %s\n", psk.Id)
                log.Panicln(err2)
                return false
        }
        return isSecretLoaded(s, psk.Id)
}
