package main

import (
	"fmt"
	"log"

	"github.com/alix2013/gomiddlewarechain"
	"github.com/alix2013/gomiddlewarechain/middleware"

	"net/http"
)

func handlefunc1(w http.ResponseWriter, r *http.Request) {
	gomiddlewarechain.RequestSetContextValue(r, "handlefunc1", "handlefunc1value")
	fmt.Fprintf(w, "this is handlefunc1\n")

}

func handlefunc2(w http.ResponseWriter, r *http.Request) {

	gomiddlewarechain.RequestSetContextValue(r, "handlefunc2", "handlefunc2value")
	fmt.Fprintf(w, "this is handlefunc2\n")

}

func handlefunc3(w http.ResponseWriter, r *http.Request) {
	v1 := gomiddlewarechain.RequestGetContextValue(r, "handlefunc1").(string)
	v2 := gomiddlewarechain.RequestGetContextValue(r, "handlefunc2").(string)

	fmt.Fprintf(w, "this is handlefunc3, key value from handlefunc1 %s, from handlefunc2 %s", v1, v2)
	gomiddlewarechain.CancelRunNextHandler(r)
}

func handlefunc4(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "this is handlefunc4\n")

}

func verifyUserPass(username, password string) bool {
	if username == "alix" && password == "password" {
		return true
	}
	return false
}

func main() {

	http.Handle("/", gomiddlewarechain.ChainHandler(middleware.BasicAuthFunc("Realm", verifyUserPass), http.FileServer(http.Dir("/tmp"))))
	http.HandleFunc("/test", gomiddlewarechain.ChainHandlerFunc(middleware.AllowedHttpMethods("POST", "GET"), handlefunc1, handlefunc2, handlefunc3, handlefunc4))
	http.HandleFunc("/admin", gomiddlewarechain.ChainHandlerFunc(middleware.BasicAuthFunc("Realm", verifyUserPass), handlefunc1, handlefunc2, handlefunc3, handlefunc4))

	//http.Handle("/file", gomiddlewarechain.ChainHandler(middleware.BasicAuthFunc("Realm", verifyUserPass), http.StripPrefix("/file", http.FileServer(http.Dir("/tmp")))))

	log.Fatal(http.ListenAndServe(":8000", nil))

}
