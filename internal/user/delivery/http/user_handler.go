package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/famkampm/nentrytask/pkg/auth"
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
	handler.Router.POST("/login", middlewares.SetMiddlewareJSON(handler.Login))
	handler.Router.GET("/profile/:id", middlewares.SetMiddlewareAuthentication(handler.GetUser))
	handler.Router.PUT("/profile/nickname/:id", middlewares.SetMiddlewareAuthentication(handler.EditNickname))

}

func (u *UserHandler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}

func (u *UserHandler) Store(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	user := &models.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	err = helper.Validate("register", user.Username, user.Password)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}

	_, err = u.UserUsecase.GetByUsername(context.TODO(), user.Username)
	if err != nil && err != sql.ErrNoRows {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	if err == nil {
		formatedError := helper.FormatError("username already exist")
		responses.ERROR(w, http.StatusBadRequest, formatedError)
		return
	}
	hashedPassword, err := helper.HashingPassword(user.Password)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	user.Password = hashedPassword
	err = u.UserUsecase.Store(context.TODO(), user)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
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

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	user := &models.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	err = helper.Validate("login", user.Username, user.Password)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	hashed_user, err := u.UserUsecase.GetByUsername(context.TODO(), user.Username)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	err = helper.VerifyPassword(hashed_user.Password, user.Password)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusForbidden, formatedError)
		return
	}
	token_string, err := auth.CreateTokenFromID(hashed_user.ID)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	responses.JSON(w, http.StatusOK, token_string)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	user, err := u.UserUsecase.GetByID(context.TODO(), int64(user_id))
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	userProfile := &models.UserProfile{
		ID:           int64(user_id),
		Username:     user.Username,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user_id))
	responses.JSON(w, http.StatusCreated, userProfile)
}

func (u *UserHandler) EditNickname(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := &models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	err = u.UserUsecase.UpdateNickname(context.TODO(), int64(user_id), user.Nickname.String)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	responses.JSON(w, http.StatusOK, "Nickname Updated")

}
