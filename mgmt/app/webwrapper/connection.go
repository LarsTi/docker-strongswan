package webwrapper

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
	"../filewrapper"
	"../viciwrapper"
)
func getConnectionFromFile(pathToFile string) connection {
	udp := filewrapper.GetStringValueFromPath(pathToFile, "UDPEncap")
	if udp == "" {
		udp = "yes"
	}
	initiator := filewrapper.GetStringValueFromPath(pathToFile, "Initiator")
	if initiator == "" {
		initiator = "no"
	}
	return connection {
		Path: pathToFile,
		LocalAddrs: filewrapper.GetStringValueFromPath(pathToFile, "LocalAddrs"),
		RemoteAddrs: filewrapper.GetStringValueFromPath(pathToFile, "RemoteAddrs"),
		Version: filewrapper.GetStringValueFromPath(pathToFile,"Version"),
		Proposals: filewrapper.GetStringValueFromPath(pathToFile, "proposals"),
		LocalTS: filewrapper.GetStringValueFromPath(pathToFile, "LocalTrafficSelectors"),
		RemoteTS: filewrapper.GetStringValueFromPath(pathToFile, "RemoteTrafficSelectors"),
		ChildProposals: filewrapper.GetStringValueFromPath(pathToFile, "ESPProposals"),
		UDPEncap: udp,
		Initiator: initiator,
	}
}
func getConnection(w http.ResponseWriter, r *http.Request) {
	path := getConnectionPath(r)
	ret := getConnectionFromFile(path)
	json.NewEncoder(w).Encode(ret)
}
func getConnections(w http.ResponseWriter, r *http.Request) {
	var ret []connection
	for _, connectionPath := range filewrapper.GetFilesForConnections(){
		ret = append(ret, getConnectionFromFile(connectionPath))
	}
	json.NewEncoder(w).Encode(ret)
}
func createConnection(w http.ResponseWriter, r *http.Request){
	path := getConnectionPath(r)

	newConnection := connection{}
	json.NewDecoder(r.Body).Decode(&newConnection)

	if path != newConnection.Path {
		errorPrefix(w,"Path", path)
		return
	}
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check != "" {
		log.Printf("[webapi][%s] Found RemoteAddrs for change file, error!\n", path)
		http.Error(w, "File exists", http.StatusBadRequest)
		return
	}
	anyChange := false
	changed, err := changeIfNeeded(path, "RemoteAddrs", newConnection.RemoteAddrs)
	if err != nil {
		errorPrefix(w, "RemoteAddrs", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "LocalAddrs", newConnection.LocalAddrs)
	if err != nil {
		errorPrefix(w, "LocalAddrs", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "proposals", newConnection.Proposals)
	if err != nil {
		errorPrefix(w, "proposals", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "Version", newConnection.Version)
	if err != nil {
		errorPrefix(w, "Version", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "RemoteTrafficSelectors", newConnection.RemoteTS)
	if err != nil {
		errorPrefix(w, "RemoteTrafficSelectors", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "LocalTrafficSelectors", newConnection.LocalTS)
	if err != nil {
		errorPrefix(w, "LocalTrafficSelectors", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "ESPProposals", newConnection.ChildProposals)
	if err != nil {
		errorPrefix(w, "ESPProposals", path)
		return
	}else if changed == true {
		anyChange = true
	}
	if anyChange == false {
		log.Printf("[webapi][%s] Nothing changed\n", path)
		json.NewEncoder(w).Encode(newConnection)
		return
	}

	//vici richtig ziehen
	loadConnection(w,r)
}
func changeConnection(w http.ResponseWriter, r *http.Request){
	path := getConnectionPath(r)

	newConnection := connection{}
	json.NewDecoder(r.Body).Decode(&newConnection)

	if path != newConnection.Path {
		errorPrefix(w,"Path", path)
		return
	}
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check == "" {
		log.Printf("[webapi][%s] Found no RemoteAddrs for change file, error!\n", path)
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}
	anyChange := false
	changed, err := changeIfNeeded(path, "RemoteAddrs", newConnection.RemoteAddrs)
	if err != nil {
		errorPrefix(w, "RemoteAddrs", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "LocalAddrs", newConnection.LocalAddrs)
	if err != nil {
		errorPrefix(w, "LocalAddrs", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "proposals", newConnection.Proposals)
	if err != nil {
		errorPrefix(w, "proposals", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "Version", newConnection.Version)
	if err != nil {
		errorPrefix(w, "Version", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "RemoteTrafficSelectors", newConnection.RemoteTS)
	if err != nil {
		errorPrefix(w, "RemoteTrafficSelectors", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "LocalTrafficSelectors", newConnection.LocalTS)
	if err != nil {
		errorPrefix(w, "LocalTrafficSelectors", path)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "ESPProposals", newConnection.ChildProposals)
	if err != nil {
		errorPrefix(w, "ESPProposals", path)
		return
	}else if changed == true {
		anyChange = true
	}
	if anyChange == false {
		log.Printf("[webapi][%s] Nothing changed\n", path)
		json.NewEncoder(w).Encode(newConnection)
		return
	}

	//vici richtig ziehen
	loadConnection(w,r)
}
func deleteConnection(w http.ResponseWriter, r *http.Request) {
	path := getConnectionPath(r)
	
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n",err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.UnloadConnection(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = filewrapper.DeleteFile(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func errorPrefix(w http.ResponseWriter, prefix string, path string){
	log.Printf("[webapi][%s] Prefix %s wrong\n", path, prefix)
	http.Error(w, "Error", http.StatusBadRequest)
}
func loadConnection(w http.ResponseWriter, r *http.Request) {
	path := getConnectionPath(r)
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.ReadConnection(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func unloadConnection(w http.ResponseWriter, r *http.Request) {
	path := getConnectionPath(r)
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.UnloadConnection(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func getConnectionPath(r *http.Request) string {
	params := mux.Vars(r)
	return params["path"]
}
