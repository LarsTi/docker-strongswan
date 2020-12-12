package main

import (
        "log"
        "github.com/strongswan/govici/vici"
)

type connection struct {
        Name string // This field will NOT be marshaled!

        LocalAddrs []string            `vici:"local_addrs"`
        RemoteAddrs []string            `vici:"remote_addrs"`
        Local      *authOpts          `vici:"local"`
        Remote     *authOpts         `vici:"remote"`
        Children   map[string]*childSA `vici:"children"`
        Version         int             `vici:"version"`
        Proposals       []string        `vici:"proposals"`
        Dpd_delay       string          `vici:"dpd_delay"`
        Dpd_timeout     string          `vici:"dpd_timeout"`
        Mobike          string          `vici:"mobike"`
        Unique          string          `vici:"unique"`

}
type authOpts struct {
        Auth  string   `vici:"auth"`
        ID    string   `vici:"id"`
}

type childSA struct {
        LocalTrafficSelectors []string `vici:"local_ts"`
        RemoteTrafficSelectors []string `vici:"remote_ts"`
        Updown                string   `vici:"updown"`
        ESPProposals          []string `vici:"esp_proposals"`
}
func unloadConnection(s *vici.Session, path string) bool {
        m := vici.NewMessage()
        if err := m.Set("name", path); err != nil {
                log.Panicln("Error unloading: Could not set Name")
                return false
        }
        _, err := s.CommandRequest("unload-conn", m)
        if err != nil{
                log.Printf("Could not unload connection %s\n", path)
                log.Println(err)
                return false
        }else{
                log.Printf("Unloaded connection %s\n", path)
                return true
        }
        return true
}
func loadConnection(s *vici.Session, path string) bool{
        l := &authOpts{
                Auth: "psk",
                ID: getStringValueFromPath("me", "RemoteAddrs"),
        }
        r := &authOpts{
                Auth: "psk",
                ID: getStringValueFromPath(path, "RemoteAddrs"),
        }


        childname := path + "-net"

        var children = make(map [string]*childSA)
        children[childname] = &childSA{
                LocalTrafficSelectors: getStringArrayFromPath(path, "LocalTrafficSelectors"),
                RemoteTrafficSelectors: getStringArrayFromPath(path, "RemoteTrafficSelectors"),
                ESPProposals: getStringArrayFromPath(path, "ESPProposals"),
        }
        c := connection {
                Name: path,
                LocalAddrs: getStringArrayFromPath(path, "LocalAddrs"),
                RemoteAddrs: getStringArrayFromPath(path, "RemoteAddrs"),
                Local: l,
                Remote: r,
                Children: children,
                Version: getIntValueFromPath(path, "Version"),
                Proposals: getStringArrayFromPath(path, "proposals"),
                Dpd_delay: "2s",
                Mobike: "no",
                Unique: "replace",
        }
        msg, err := vici.MarshalMessage(c)
        if err != nil {
                log.Panicf("[%s]%s\n",path ,err)
        }

        m := vici.NewMessage()
        m.Set(path, msg)

        st, err := s.CommandRequest("load-conn", m)
        if err != nil{
                log.Println("Error executing Command \"load-conn\"")
                log.Println(m)
                log.Panicln(err)
                return false
        }
        log.Printf("Loaded Connection %s with success status: %s\n", path, st.Get("success"))
        return true
}
func initiate(s *vici.Session, path string, isIke bool) bool{
        m := vici.NewMessage()
        childname := path + "-net"
        if isIke == false {
                if err := m.Set("child", childname); err != nil {
                        log.Panicln(err)
                        return false
                }
        }
        if err := m.Set("ike", path); err != nil {
                log.Panicln(err)
                return false
        }

        ms, err := s.StreamedCommandRequest("initiate", "control-log", m)
        if err != nil {
                log.Panicln(err)
                return false
        }

        for _, msg := range ms.Messages() {
                if err := msg.Err(); err != nil {
                        log.Printf("Error initiating %s\n", path)
                        log.Println(err);
                        return false
                }
        }
        if isIke == false{
                log.Printf("Initiated %s\n", childname)
        }else{
                log.Printf("Initiated %s\n", path)
        }
        return true
}
func terminate(s *vici.Session, path string, isIke bool) bool{
        if isIke == true {
                terminate(s, path, false)
        }
        childname := path + "-net"
        ike := vici.NewMessage()
        if err := ike.Set("ike", path); err != nil {
                log.Panicln(err)
                return false
        }
        if isIke == false {
                if err:= ike.Set("child", childname); err != nil{
                        log.Panicln(err)
                        return false
                }
        }
        if err := ike.Set("force", true); err != nil {
                log.Panicln(err)
                return false
        }
        if err := ike.Set("timeout", 1000); err != nil {
                log.Panicln(err)
                return false
        }
        scrike, err := s.StreamedCommandRequest("terminate", "control-log", ike)
        if err != nil {
                log.Panicln(err)
                return false
        }
        for _, msg := range scrike.Messages() {
                if msg.Get("terminated") != nil {
                        if isIke == false{
                                log.Printf("Terminiertes Child: %s\n", msg.Get("terminated"))
                        }else{
                                log.Printf("Terminierte SAs: %s\n", msg.Get("terminated"))
                        }
                }
        }
        if isIke == false{
                log.Printf("Tried to terminate child %s\n", childname)
        }else {
                log.Printf("Tried to terminate ike %s\n", path)
        }
        return true
}
func loadConn(s *vici.Session, path string) bool{
        terminate(s, path, true)
        unloadConnection(s, path)
        if loadConnection(s, path) == false{
                return false
        }
        if initiate(s, path, true) == false{
                return false
        }
        return initiate(s, path, false)
}
