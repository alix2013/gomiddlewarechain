package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alix2013/gomiddlewarechain"
	"github.com/alix2013/gomiddlewarechain/middleware"

	"net/http"
)

var signkey = []byte("yoursecretkey")
var contextTokenKey = "user"

func publicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is public page\n")

}

func authUser(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	//code auth user and password
	return true
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("get user password from POST FORM", username, password)
	//validate user credentials
	if authUser(username, password) {
		//genereate token for map
		claimMap := make(map[string]interface{})
		claimMap["username"] = username
		claimMap["exp"] = time.Now().Add(time.Minute * 20).Unix()
		//tokenString := middleware.GenerateHS256JWTForMap(claimMap, signkey)
		tokenString, err := middleware.GenerateDefaultJWTForMap(claimMap, signkey)
		if err != nil {

			log.Println("generate token failed", err)
			http.Error(w, "Generate token failed", 500)

		} else {
			w.WriteHeader(http.StatusOK)
			retJSON := map[string]string{"token": tokenString}
			json.NewEncoder(w).Encode(retJSON)
		}
	} else {
		http.Error(w, "Invalid credentials\n", http.StatusForbidden)
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if userMap, ok := (gomiddlewarechain.RequestGetContextValue(r, contextTokenKey)).(map[string]interface{}); ok {
		username := userMap["username"]
		fmt.Fprint(w, "welcome:", username, "\nyou see the page with valid token!\n")
	} else {
		log.Println("No context for key:", contextTokenKey)
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/login", gomiddlewarechain.ChainHandlerFunc(middleware.AllowedHttpMethod("POST"), LoginHandler))
	mux.HandleFunc("/protect", gomiddlewarechain.ChainHandlerFunc(middleware.AllowedHttpMethod("GET"), ProtectedHandler))
	mux.HandleFunc("/public", gomiddlewarechain.ChainHandlerFunc(middleware.AllowedHttpMethod("GET"), publicHandler))

	jwtmiddleware := middleware.NewDefaultJWTMiddleware(
		middleware.JWTDefaultMiddlewareOptions{
			SigningKey: signkey, ContextKeyForClaim: contextTokenKey, ExcludeURI: []string{"/public", "/login"}})

	//signkey, contextTokenKey, []string{"/login", "/public"})

	log.Fatal(http.ListenAndServe(":8000", gomiddlewarechain.ChainHandler(jwtmiddleware, mux)))

}
