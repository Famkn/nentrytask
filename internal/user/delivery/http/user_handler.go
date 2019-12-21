package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/famkampm/nentrytask/pkg/helper"
	"github.com/famkampm/nentrytask/pkg/middlewares"
	"github.com/famkampm/nentrytask/pkg/responses"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	Router      *httprouter.Router
	UserUsecase user.Usecase
}

func NewUserHandler(router *httprouter.Router, us user.Usecase) {
	handler := &UserHandler{
		Router:      router,
		UserUsecase: us,
	}
	handler.Router.GET("/", handler.Home)

	handler.Router.POST("/register", middlewares.SetMiddlewareJSON(handler.Store))

}

func (u *UserHandler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}

func (u *UserHandler) Store(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := &models.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = helper.Validate("register", user.Username, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = u.UserUsecase.GetByUsername(context.TODO(), user.Username)
	if err != nil && err != sql.ErrNoRows {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err == nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	hashedPassword, err := helper.HashingPassword(user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user.Password = hashedPassword
	err = u.UserUsecase.Store(context.TODO(), user)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user.ID))
	userProfile := &models.UserProfile{}
	userProfile.ID = user.ID
	userProfile.Username = user.Username
	userProfile.Nickname = user.Nickname
	userProfile.ProfileImage = user.ProfileImage
	responses.JSON(w, http.StatusCreated, userProfile)
}
