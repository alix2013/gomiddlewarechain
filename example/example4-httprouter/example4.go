package main

import (
	"fmt"
	"log"

	"github.com/alix2013/gomiddlewarechain"
	"github.com/alix2013/gomiddlewarechain/middleware"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "this is index page!\n")
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "this is login page!\n")
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

	basicAuthMiddleware := middleware.NewBasicAuthMiddleware(
		middleware.BasicAuthOptions{
			Realm:         "authentication required",
			UserValidator: verifyUserPass,
			ExcludeURI:    []string{"/login"}})

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/login", Login)
	log.Fatal(http.ListenAndServe(":8000", gomiddlewarechain.ChainHandler(basicAuthMiddleware, router)))

}
