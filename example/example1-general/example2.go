package main

import (
	"fmt"
	"log"

	"github.com/alix2013/gomiddlewarechain"

	"net/http"
)

func publicHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this public page\n")
}

//if OS env MIDDLEWARE_DEBUG=1, panic error stack will display in broswer
// if not set env MIDDLEWARE_DEBUG, just show "[Middleware] Internal Server Error, please check system logs" http code 500
func panicFuncDemo(w http.ResponseWriter, r *http.Request) {
	//panic("Panic Error ...")
	//simulate a runtime error
	var slice []int
	slice[0] = 0
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/tmp")))
	mux.HandleFunc("/public", publicHandleFunc)
	mux.HandleFunc("/panic", panicFuncDemo)
	log.Println("Server listening 8000")
	log.Fatal(http.ListenAndServe(":8000", gomiddlewarechain.ChainHandler(mux)))

}
