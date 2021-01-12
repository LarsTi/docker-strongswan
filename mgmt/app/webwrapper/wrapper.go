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

	router.HandleFunc("/api/secrets", GetSecrets).Methods("GET")
	router.HandleFunc("/api/secrets/{path}", GetSecret).Methods("GET")
	router.HandleFunc("/api/secrets/{path}", ChangeSecret).Methods("PUT")
	router.HandleFunc("/api/secrets/{path}", CreateSecret).Methods("POST")
	router.HandleFunc("/api/secrets/{path}", DeleteSecret).Methods("DELETE")
	router.HandleFunc("/api/secrets/{path}/unload", UnloadSecret).Methods("PUT")
	router.HandleFunc("/api/secrets/{path}/load", LoadSecret).Methods("PUT")
	
	router.HandleFunc("/api/connections/", GetConnections).Methods("GET")
	router.HandleFunc("/api/connections/{path}", GetConnection).Methods("GET")
	router.HandleFunc("/api/connections/{path}", ChangeConnection).Methods("PUT")
	router.HandleFunc("/api/connections/{path}", CreateConnection).Methods("POST")
	router.HandleFunc("/api/connections/{path}", DeleteConnection).Methods("DELETE")
	router.HandleFunc("/api/connections/{path}/unload", UnloadConnection).Methods("PUT")
	router.HandleFunc("/api/connections/{path}/load", LoadConnection).Methods("PUT")

	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port) , router))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[webapi-request] %s: %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	    })
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

func GetConnections(w http.ResponseWriter, r *http.Request){
	getConnections(w,r)
}
func GetConnection(w http.ResponseWriter, r *http.Request){
	getConnection(w,r)
}
func ChangeConnection(w http.ResponseWriter, r *http.Request){
	changeConnection(w,r)
}
func DeleteConnection (w http.ResponseWriter, r *http.Request) {
	deleteConnection(w,r)
}
func CreateConnection (w http.ResponseWriter, r *http.Request) {
	createConnection(w,r)
}
func UnloadConnection (w http.ResponseWriter, r *http.Request) {
	unloadConnection(w,r)
}
func LoadConnection (w http.ResponseWriter, r *http.Request) {
	loadConnection(w,r)
}

func GetSecrets (w http.ResponseWriter, r *http.Request){
	getSecrets(w,r)
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
