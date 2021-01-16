package webwrapper

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"../filewrapper"
	"../viciwrapper"
)
func getSecretFromFile(pathToFile string) sharedSecret{
	return sharedSecret{
		Path: strings.TrimSuffix(pathToFile, ".secret"),
		Typ: "PSK",
		Data: filewrapper.GetStringValueFromPath(pathToFile, "PSK"),
		Owners: filewrapper.GetStringValueFromPath(pathToFile, "RemoteAddrs"),
	}
}
func getSecrets(w http.ResponseWriter, r *http.Request){
	var ret []sharedSecret
	for _, secretPath := range filewrapper.GetFilesForSecrets(){
		ret = append(ret, getSecretFromFile(secretPath))
	}
	json.NewEncoder(w).Encode(ret)
}
func getSecret (w http.ResponseWriter, r *http.Request){
	path := getSecretPath(r)
	ret := getSecretFromFile(path)
	json.NewEncoder(w).Encode(ret)
}
func unloadSecret (w http.ResponseWriter, r *http.Request) {
	path := getSecretPath(r)
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check == "" {
		log.Printf("[webapi][%s] Found no LocalAddrs for change file, error\n", path)
		http.Error(w,"File does not exist", http.StatusNotFound)
		return
	}
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = wrapper.UnloadSecret(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func loadSecret (w http.ResponseWriter, r *http.Request) {
	path := getSecretPath(r)
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check == "" {
		log.Printf("[webapi][%s] Found no LocalAddrs for change file, error\n", path)
		http.Error(w,"File does not exist", http.StatusNotFound)
		return
	}
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = wrapper.ReadSecret(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func changeSecret (w http.ResponseWriter, r *http.Request) {
	path := getSecretPath(r)

	secret := sharedSecret{}
	json.NewDecoder(r.Body).Decode(&secret)
	log.Println(secret)
	//path = me.secret, secrets.path = me
	if ! strings.HasPrefix(path, secret.Path) {
		log.Printf("[webapi][%s] wrong Path set in JSON")
		http.Error(w, "wrong path", http.StatusBadRequest)
		return
	}
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check == "" {
		log.Printf("[webapi][%s] Found no LocalAddrs for change file, error\n", path)
		http.Error(w,"File does not exist", http.StatusNotFound)
		return
	}
	anyChange := false
	changed, err := changeIfNeeded(path, "RemoteAddrs", secret.Owners)
	if err != nil{
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}else if changed == true {
		anyChange = true
	}
	changed, err = changeIfNeeded(path, "PSK", secret.Data)
	if err != nil{
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}else if changed == true {
		anyChange = true
	}
	if anyChange == false {
		log.Printf("[webapi][%s] Nothing Changed\n", path)
		json.NewEncoder(w).Encode(secret)
		return
	}
	
	//vici richtig ziehen
	wrapper, errW := viciwrapper.GetWrapper()
	if errW != nil {
		log.Printf("[webapi] %s\n", errW)
		http.Error(w, errW.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.ReadSecret(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(secret)
}
func deleteSecret (w http.ResponseWriter, r *http.Request) {
	path := getSecretPath(r)
	//vici richtig ziehen
	wrapper, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.UnloadSecret(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = filewrapper.DeleteFile(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}

}
func createSecret (w http.ResponseWriter, r *http.Request){
	path := getSecretPath(r)
	
	secret := sharedSecret{}
	json.NewDecoder(r.Body).Decode(&secret)
	
	check := filewrapper.GetStringValueFromPath(path, "RemoteAddrs")
	if check != "" {
		log.Printf("[webapi][%s] Found LocalAddrs for new create file, error\n", path)
		http.Error(w,"File exists", http.StatusBadRequest)
		return
	}
	_, err := changeIfNeeded(path, "RemoteAddrs", secret.Owners)
	if err != nil{
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = changeIfNeeded(path, "PSK", secret.Owners)
	if err != nil{
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	//vici richtig ziehen
	wrapper, errW := viciwrapper.GetWrapper()
	if errW != nil {
		log.Printf("[webapi] %s\n", errW)
		http.Error(w, errW.Error(), http.StatusBadRequest)
		return
	}
	err = wrapper.ReadSecret(path)
	if err != nil {
		log.Printf("[webapi] %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(secret)
	
}
func getSecretPath(r *http.Request) string {
	params := mux.Vars(r)
	return strings.Join([]string{params["path"], "secret",}, ".")
}
