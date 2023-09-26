package middleware

import (
	"Mahathirrr/belajar-golang-restful-api/helper"
	"Mahathirrr/belajar-golang-restful-api/model/web"
	"net/http"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if "RAHASIA" == r.Header.Get("X-API-KEY") {
		//ok
		middleware.Handler.ServeHTTP(w, r)
	} else {
		//error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.WriteToResponseBody(w, webResponse)
	}
}
