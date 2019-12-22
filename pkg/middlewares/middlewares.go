package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/famkampm/nentrytask/pkg/auth"
	"github.com/famkampm/nentrytask/pkg/responses"
	"github.com/julienschmidt/httprouter"
)

func SetMiddlewareAuthentication(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := auth.TokenValidFromID(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		user_id, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("Unauthorized"))
			return
		}
		extracted_user_id, err := auth.ExtractTokenID(r)
		if err != nil || extracted_user_id != int64(user_id) {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r, ps)
	}
}

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
