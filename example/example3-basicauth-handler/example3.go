package main

import (
	"fmt"
	"log"

	"github.com/alix2013/gomiddlewarechain"
	"github.com/alix2013/gomiddlewarechain/middleware"

	"net/http"
)

//verify user/password from cache, file, database etc...
//this just demo
func verifyUserPass(username, password string) bool {
	if username == "alix" && password == "password" {
		return true
	}
	return false
}

func adminHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this admin page, authentication required\n")
}

func userHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this user page, authentication required\n")
}

func publicHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this public page\n")
}

func main() {
	//http.HandleFunc("/ba1", gomiddlewarechain.ChainHandleFunc(http.HandlerFunc(middleware.NewBasicAuth("Realm", verifyUserPass))))

	basicAuthMiddleware := middleware.NewBasicAuthMiddleware(
		middleware.BasicAuthOptions{
			Realm:         "authentication required",
			UserValidator: verifyUserPass,
			ExcludeURI:    []string{"/public", "/"}})

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/tmp")))
	mux.HandleFunc("/admin", adminHandleFunc)
	mux.HandleFunc("/user", userHandleFunc)
	mux.HandleFunc("/public", publicHandleFunc)

	log.Println("Server listening 8000")
	log.Fatal(http.ListenAndServe(":8000", gomiddlewarechain.ChainHandler(basicAuthMiddleware, mux)))

}
