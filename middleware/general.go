package middleware

import (
	"net/http"
	"strings"

	"github.com/alix2013/gomiddlewarechain"
)

func AllowedHttpMethods(method ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range method {
			if strings.ToLower(r.Method) == strings.ToLower(m) {
				return
			}
		}

		http.Error(w, "Bad Request(MethodNotAllowed)", http.StatusMethodNotAllowed)
		gomiddlewarechain.CancelRunNextHandler(r)
	}
}

func AllowedHttpMethodsWithErrorMsg(errorMsg string, errorStatusCode int, method ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range method {
			if strings.ToLower(r.Method) == strings.ToLower(m) {
				return
			}
		}

		http.Error(w, errorMsg, errorStatusCode)
		gomiddlewarechain.CancelRunNextHandler(r)
	}
}

/*
type RecoverMiddleware struct {
}

func NewRecoverMiddleware() http.Handler {
	return RecoverMiddleware{}
}

func (rec RecoverMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprintf("Internal Server Error:%v", err), http.StatusInternalServerError)
			//gomiddlewarechain.CancelRunNextHandler(r)
		}
	}()
}
func RecoverFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, fmt.Sprintf("Internal Server Error:%v", err))
				//http.Error(w, fmt.Sprintf("Internal Server Error:%v", err), http.StatusInternalServerError)
				//gomiddlewarechain.CancelRunNextHandler(r)
			}
		}()
	})
}
*/
