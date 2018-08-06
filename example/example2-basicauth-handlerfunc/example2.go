package main

import (
	"fmt"
	"log"

	"github.com/alix2013/gomiddlewarechain"
	"github.com/alix2013/gomiddlewarechain/middleware"

	"net/http"
)

func adminHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this admin page, authentication required\n")
}

func userHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this user page, authentication required\n")
}

func publicHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this public page\n")
}

//verify user/password from cache, file, database etc...
//this just demo
func verifyUserPass(username, password string) bool {
	if username == "alix" && password == "password" {
		return true
	}
	return false
}

func main() {

	http.Handle("/", gomiddlewarechain.ChainHandler(middleware.BasicAuthFunc("Realm", verifyUserPass), http.FileServer(http.Dir("/tmp"))))
	http.HandleFunc("/admin", gomiddlewarechain.ChainHandlerFunc(middleware.BasicAuthFunc("Realm", verifyUserPass), adminHandleFunc))
	http.HandleFunc("/user", gomiddlewarechain.ChainHandlerFunc(middleware.BasicAuthFunc("Realm", verifyUserPass), userHandleFunc))
	http.HandleFunc("/public", publicHandleFunc)

	log.Fatal(http.ListenAndServe(":8000", nil))

}
