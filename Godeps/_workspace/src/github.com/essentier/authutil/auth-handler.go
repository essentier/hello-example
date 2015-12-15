package authutil

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	//"github.com/gorilla/context"
	//"github.com/go-errors/errors"
	jwt "github.com/dgrijalva/jwt-go"
)

// type userIdKey int

// const user_ID_KEY userIdKey = 1

type AuthHandler struct {
	//JwtService *JWTService
}

// func GetUserID(r *http.Request) (string, error) {
// 	if rv := context.Get(r, user_ID_KEY); rv != nil {
// 		return rv.(string), nil
// 	}
// 	return "", errors.New("User id is not in context.")
// }

func (h *AuthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Printf("AuthHandler.ServeHTTP")
	log.Printf("http request %+v", r)
	token, err := parseFromRequest(r, jwtKeyFunc)

	//if err == nil && token.Valid && !h.JwtService.IsInBlacklist(r.Header.Get("Authorization")) {
	if err != nil || token == nil || !token.Valid {
		log.Printf("returning status unauthorized")
		rw.Header().Add("WWW-Authenticate", "basic realm=\"demo site\"")
		//rw.Header().Add("WWW-Authenticate", "Bearer realm=\"Git Access\"")
		//rw.Header().Add("WWW-Authenticate", "Bearer realm=\"JWT\"")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	//context.Set(r, user_ID_KEY, sub)
	log.Printf("token here: %+v", token)
	sub, exists := token.Claims["sub"]
	if !exists {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Missing subject in token."))
		return
	}

	query := r.URL.Query()
	query.Add("userId", sub.(string))
	r.URL.RawQuery = query.Encode()
	next(rw, r)
}

func parseFromRequest(req *http.Request, keyFunc jwt.Keyfunc) (token *jwt.Token, err error) {
	token, err = jwt.ParseFromRequest(req, keyFunc)
	if err != jwt.ErrNoTokenInRequest {
		return token, err
	}

	// Look for an Authorization header
	if ah := req.Header.Get("Authorization"); ah != "" {
		// See if there is a bearer token disguised as Basic data
		if len(ah) > 5 && strings.ToUpper(ah[0:5]) == "BASIC" {
			log.Printf("base64 encoded token %q\n", ah[6:])
			decodedData, err := base64.StdEncoding.DecodeString(ah[6:])
			if err != nil {
				return nil, jwt.ErrNoTokenInRequest
			}
			tokenString := strings.TrimSuffix(string(decodedData), ":")
			log.Printf("base64 decoded token %q\n", decodedData)

			return ParseToken(tokenString)
		}
	}

	return nil, jwt.ErrNoTokenInRequest
}
