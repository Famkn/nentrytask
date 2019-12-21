package middlewares

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// func SetMiddlewareAuthentication(next httprouter.Handle) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		err := auth.TokenValid(r)
// 		if err != nil {
// 			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
// 			return
// 		}
// 		next(w, r, ps)
// 	}
// }

func MiddlewareTestHttpRouter(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		next(w, r, ps)
	}
}

func SetMiddlewareJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r, ps)
	}
}
