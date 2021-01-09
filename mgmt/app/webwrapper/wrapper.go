package webwrapper

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"log"
	"fmt"
	"../filewrapper"
)
func RunWebApi(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/api/secrets/{path}", GetSecret).Methods("GET")
	router.HandleFunc("/api/secrets/{path}", ChangeSecret).Methods("PUT")
	router.HandleFunc("/api/secrets/{path}", CreateSecret).Methods("POST")
	router.HandleFunc("/api/secrets/{path}", DeleteSecret).Methods("DELETE")
	router.HandleFunc("/api/secrets/{path}/unload", UnloadSecret).Methods("PUT")
	router.HandleFunc("/api/secrets/{path}/load", LoadSecret).Methods("PUT")

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port) , router))
}
func changeIfNeeded(path string, prefix string, value string)(bool, error){
	oldValue := filewrapper.GetStringValueFromPath(path, prefix)
	if (value == "") {
		err := fmt.Errorf("[webapi][%s] - changing value for %s permitted, empty value not allowed", path, prefix)
		log.Println(err)
		return false, err
	}else if oldValue == value {
		log.Printf("[webapi][%s] - not changing value for %s (equally) \n", path, prefix)
		return false, nil
	}else{
		log.Printf("[webapi][%s] - changing value for %s\n", path, prefix)
		err := filewrapper.WriteOrReplaceLine(path, prefix, value)
		if err != nil {
			log.Printf("[webapi][%s] - %s\n", path, err)
			return false, err
		}
		return true, nil
	}
}
func GetSecret (w http.ResponseWriter, r *http.Request){
	getSecret(w,r)
}
func ChangeSecret (w http.ResponseWriter, r *http.Request) {
	changeSecret(w,r)
}
func DeleteSecret (w http.ResponseWriter, r *http.Request) {
	deleteSecret(w,r)
}
func CreateSecret (w http.ResponseWriter, r *http.Request) {
	createSecret(w,r)
}
func UnloadSecret (w http.ResponseWriter, r *http.Request) {
	unloadSecret(w,r)
}
func LoadSecret (w http.ResponseWriter, r *http.Request) {
	loadSecret(w,r)
}
