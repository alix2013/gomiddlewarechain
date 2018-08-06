package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alix2013/gomiddlewarechain"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTDefaultMiddlewareOptions struct {
	SigningKey         []byte
	ContextKeyForClaim string
	ExcludeURI         []string
}

type JWTHS256Middleware struct {
	Options JWTDefaultMiddlewareOptions
}

func NewDefaultJWTMiddleware(options JWTDefaultMiddlewareOptions) JWTHS256Middleware {
	return JWTHS256Middleware{Options: options}
}

// implement http.ServeHTTP
func (jwths256 JWTHS256Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	found := false
	for _, exuri := range jwths256.Options.ExcludeURI {
		if exuri == r.RequestURI {
			found = true
			log.Println("[Middleware] JWTHS256Middleware match exclude URI:", r.RequestURI)
			break
		}
	}
	if !found {
		jwths256handlerfunc(jwths256.Options.SigningKey, jwths256.Options.ContextKeyForClaim, w, r)
	}

}

// genearete HS256 token string for map
func GenerateDefaultJWTForMap(mapClaim map[string]interface{}, signingKey []byte) (tokenString string, err error) {
	return GenerateHS256JWTForMap(mapClaim, signingKey)
}

// generate HS256 token string for map
func GenerateHS256JWTForMap(mapClaim map[string]interface{}, signingKey []byte) (tokenString string, err error) {

	token := jwt.New(jwt.SigningMethodHS256)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for k, v := range mapClaim {
			claims[k] = v
		}
		//claims["username"] = user.Username
		//claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
		tokenString, _ = token.SignedString(signingKey)
		return
	} else {
		return "", errors.New("[Middleware] GenerateHS256JWTForMap failed")
	}

}

// verify HS256 token string for map
// if not error, return decoded token map
func VerifyHS256TokenForMap(tokenString string, signingKey []byte) (bool, error, map[string]interface{}) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("[Middleware]Unexpected signing method: %v", token.Header["alg"])
		}
		// signingKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["username"], claims["exp"])
		retMap := make(map[string]interface{})
		for k, v := range claims {
			retMap[k] = v
		}
		return true, nil, retMap
	} else {
		//log.Println("[Middleware] VerifyHS256TokenForMap Error:", err)
		return false, err, nil
	}

}

// for  http.handlerFunc detailed implement
func jwths256handlerfunc(signingKey []byte, contextKeyforClaim string, w http.ResponseWriter, r *http.Request) {
	tokenString := extractToken(r)
	if tokenString == "" {
		gomiddlewarechain.CancelRunNextHandler(r)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

	} else {
		if ok, err, claimMap := VerifyHS256TokenForMap(tokenString, signingKey); ok {
			gomiddlewarechain.RequestSetContextValue(r, contextKeyforClaim, claimMap)

		} else {
			gomiddlewarechain.CancelRunNextHandler(r)
			log.Println("[Middleware] JWTHandlerFunc VerifyHS256TokenForMap error:", err)
			http.Error(w, "Invalid Token", http.StatusForbidden)
		}
	}

}

// http.HandlerFun for simple usecase
func DefaultJWTHandlerFunc(signingKey []byte, contextKeyforClaim string) http.HandlerFunc {
	return JWTHS256HandlerFunc(signingKey, contextKeyforClaim)
}

// HS256 HandlerFunc for simle usecase
func JWTHS256HandlerFunc(signingKey []byte, contextKeyforClaim string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwths256handlerfunc(signingKey, contextKeyforClaim, w, r)
	}
}

// exatrct token string from request
func extractToken(r *http.Request) (tokenString string) {

	tokenString = r.Header.Get("Token")
	if tokenString != "" {
		return
	}

	tokenString = r.Header.Get("token")
	if tokenString != "" {
		return
	}

	// from query params
	tokenString = r.URL.Query().Get("jwt")
	if tokenString != "" {
		return
	}

	tokenString = r.URL.Query().Get("token")
	if tokenString != "" {
		return
	}

	// from authorization header
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		tokenString = bearer[7:]
	}
	if tokenString != "" {
		return
	}

	//  from cookie
	cookie, err := r.Cookie("jwt")
	if err == nil {
		tokenString = cookie.Value
	}
	if tokenString != "" {
		return
	}
	return ""
}

/*
type JWTMapClaimMiddleware struct {
	signingMethod      string
	signingKey         []byte
	contextKeyforClaim string
	excludeURI         []string
}

// signingMethod : HS256, HS384, HS512
func GenerateJWTForMap(signingMethod string, mapClaim map[string]interface{}, signingKey []byte) (tokenString string, err error) {

	var token *jwt.Token

	switch signingMethod {
	case "HS256":
		token = jwt.New(jwt.SigningMethodHS256)
	case "HS386":
		token = jwt.New(jwt.SigningMethodHS384)
	case "HS512":
		token = jwt.New(jwt.SigningMethodHS512)
	default:
		return "", errors.New("SigningMethod unsupported")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for k, v := range mapClaim {
			claims[k] = v
		}
	} else {
		return "", errors.New("MapClaims error")
	}

	//claims["username"] = user.Username
	//claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenString, _ = token.SignedString(signingKey)
	return

}

//same as VerifyHS256TokenForMap
func VerifyJWTForMap(tokenString string, signingKey []byte) (bool, error, map[string]interface{}) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// signingKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//return true, nil
		//fmt.Println(claims["username"], claims["exp"])
		retMap := make(map[string]interface{})
		for k, v := range claims {
			retMap[k] = v
		}
		return true, nil, retMap
	} else {
		log.Println("VerifyJWTForMap Error", err)
		return false, err, nil
	}

}

func jwthhandlerfunc(signingKey []byte, contextKeyforClaim string, w http.ResponseWriter, r *http.Request) {

	tokenString := extractToken(r)
	//fmt.Println("extracted token:", tokenString)
	if tokenString == "" {
		gomiddlewarechain.CancelRunNextHandler(r)
		//w.WriteHeader(401)
		//fmt.Fprintf(w, "Permission denied")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

	} else {
		if ok, err, claimMap := VerifyJWTForMap(tokenString, signingKey); ok {
			gomiddlewarechain.RequestSetContextValue(r, contextKeyforClaim, claimMap)

		} else {
			//fmt.Fprintf(w, err.Error())
			log.Println("jwthhandlerfunc erro", err)
			gomiddlewarechain.CancelRunNextHandler(r)

			w.WriteHeader(500)
			fmt.Fprintf(w, "Token Error")
			log.Println("jwthhandlerfunc VerifyHS256TokenForMap failed", err)

		}
	}

}

// for http.handler
func (jwtmid JWTMapClaimMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	found := false
	for _, exuri := range jwtmid.excludeURI {
		if exuri == r.RequestURI {
			found = true
			log.Println("JWTMapClaimMiddleware match exclude URI:", r.RequestURI)
			break
		}
	}
	if !found {
		jwthhandlerfunc(jwtmid.signingKey, jwtmid.contextKeyforClaim, w, r)
	}

}

func JWTHSHandlerFunc(signingKey []byte, contextKeyforClaim string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwthhandlerfunc(signingKey, contextKeyforClaim, w, r)
	}
}

func NewJWTMapClaimMiddleware(signingMethod string, signingKey []byte, contextKeyforClaim string, excludeURI []string) JWTMapClaimMiddleware {
	return JWTMapClaimMiddleware{signingMethod, signingKey, contextKeyforClaim, excludeURI}
}

*/
