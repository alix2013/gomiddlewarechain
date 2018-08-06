package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/alix2013/gomiddlewarechain"
)

type UserPassAuthFunc func(username, password string) bool
type BasicAuthOptions struct {
	Realm         string
	UserValidator UserPassAuthFunc
	ExcludeURI    []string
}

type BasicAuth struct {
	Options BasicAuthOptions
}

func NewBasicAuthMiddleware(options BasicAuthOptions) BasicAuth {
	return BasicAuth{Options: options}
}

func displayBasicAuthError(w http.ResponseWriter, r *http.Request, realm string) {
	gomiddlewarechain.CancelRunNextHandler(r)
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)

	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	//w.WriteHeader(401)
	//	w.Write([]byte("401 Unauthorized\n"))

}
func authentication(realm string, validator UserPassAuthFunc, w http.ResponseWriter, r *http.Request) {
	authz := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authz) != 2 || authz[0] != "Basic" {
		displayBasicAuthError(w, r, realm)
		return
	}
	userpass, _ := base64.StdEncoding.DecodeString(authz[1])
	uparray := strings.SplitN(string(userpass), ":", 2)

	if len(uparray) != 2 || !validator(uparray[0], uparray[1]) {
		displayBasicAuthError(w, r, realm)
		return
	}
}

func (auth BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, exuri := range auth.Options.ExcludeURI {
		if exuri == r.RequestURI {
			found = true
			log.Println("[Middleware] BasicAuthMiddleware match exclude URI:", r.RequestURI)
			break
		}
	}
	if !found {
		authentication(auth.Options.Realm, auth.Options.UserValidator, w, r)
	}

}

func (auth BasicAuth) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authentication(auth.Options.Realm, auth.Options.UserValidator, w, r)
	})

}

func (auth BasicAuth) Handler() http.Handler {
	return auth
}

///////////////////////////////////////////////////////////////////////////
// middleware.BasicAuthFunc for simple http.handlerFunc
func BasicAuthFunc(realm string, userpassAuth func(username, password string) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authentication(realm, userpassAuth, w, r)
	}
}
