package gomiddlewarechain

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"context"
)

// set request context key tag whether break chain
const cancel_run_next_context_key = "CANCEL_RUN_NEXT_CONTEXT_KEY"

// if OS env set MIDDLEWARE_DEBUG, panic error show in response page
var MIDDLEWARE_DEBUG = false

func init() {
	if os.Getenv("MIDDLEWARE_DEBUG") != "" {
		MIDDLEWARE_DEBUG = true
	}
}

// chain http.HandlerFunc
// example:  http.HandleFunc("/yoururl", gomiddlewarechain.ChainHandlerFunc( handlerFunc1, handlerFun2 ))
func ChainHandlerFunc(handlers ...http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		beginChain(r)
		for _, handler := range handlers {
			if !isCancelRunNext(r) {
				start := time.Now()

				defer func() {
					if err := recover(); err != nil {
						errorStack := fmt.Sprintf("Panic Error:%s, Stack trace: %s\n", err, debug.Stack())
						log.Println("[Middleware] Recover Panic error:", errorStack)
						if MIDDLEWARE_DEBUG {
							http.Error(w, errorStack, http.StatusInternalServerError)
						} else {
							http.Error(w, "[Middleware] Internal Server Error, please check system logs", http.StatusInternalServerError)
						}

						CancelRunNextHandler(r)
					}
				}()

				handler(w, r)

				elapsed := time.Since(start)
				handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
				log.Println("[Middlewarechain]==>Performance:", handlerName, r.Method, r.RequestURI, elapsed)

			}
		}
		//endChain(r)
	})
}

// chain http.Handler
// example: http.ListenAndServe(":8000", gomiddlewarechain.ChainHandler(basicAuthMiddleware, mux)
func ChainHandler(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginChain(r)
		for _, handler := range handlers {
			if !isCancelRunNext(r) {
				start := time.Now()

				//panic check
				defer func() {
					if err := recover(); err != nil {
						errorStack := fmt.Sprintf("Panic Error:%s, Stack trace: %s\n", err, debug.Stack())
						log.Println("[Middleware] Recover Panic error:", errorStack)
						if MIDDLEWARE_DEBUG {
							http.Error(w, errorStack, http.StatusInternalServerError)
						} else {
							http.Error(w, "[Middleware] Internal Server Error, please check system logs", http.StatusInternalServerError)
						}
						CancelRunNextHandler(r)
					}
				}()

				handler.ServeHTTP(w, r)
				elapsed := time.Since(start)
				//handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

				handlerName := reflect.TypeOf(handler)
				if handlerName.Name() == "http.HandlerFunc" {
					handlerFuncName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
					log.Println("[Middlewarechain]==>Performance:", handlerFuncName, r.Method, r.RequestURI, elapsed)
				}

				log.Println("[Middlewarechain]==>Performance:", handlerName, r.Method, r.RequestURI, elapsed)

			}
		}
		//endChain(r)
	})
}

// add context key value to request
func RequestSetContextValue(r *http.Request, contextKey, contextValue interface{}) {
	newRequest := r.WithContext(context.WithValue(r.Context(), contextKey, contextValue))
	*r = *newRequest
	return

}

// get context value from request by key
func RequestGetContextValue(r *http.Request, contextKey interface{}) interface{} {
	return r.Context().Value(contextKey)
}

// set break  tag to request context
func beginChain(r *http.Request) {
	RequestSetContextValue(r, cancel_run_next_context_key, false)
}

// get cancle tag from request context
func isCancelRunNext(r *http.Request) bool {
	rv := RequestGetContextValue(r, cancel_run_next_context_key)
	if rv != nil {
		return rv.(bool)
	} else {
		return false
	}

}

// set cancel tag cancel_run_next_context_key to true
func CancelRunNextHandler(r *http.Request) {
	RequestSetContextValue(r, cancel_run_next_context_key, true)

}

/*
// for gorrila context
func endChain(r *http.Request) {
	context.Clear(r)
}

func beginChain(r *http.Request) {
	context.Set(r, cancel_run_next_context_key, false)
}

func endChain(r *http.Request) {
	context.Clear(r)
}

func isCancelRunNext(r *http.Request) bool {
	rv := context.Get(r, cancel_run_next_context_key)
	if rv != nil {
		return rv.(bool)
	} else {
		return false
	}

}

func CancelRunNextHandler(r *http.Request) {
	context.Set(r, cancel_run_next_context_key, true)
}
*/
