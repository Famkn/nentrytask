package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/famkampm/nentrytask/pkg/auth"
	"github.com/famkampm/nentrytask/pkg/helper"
	"github.com/famkampm/nentrytask/pkg/middlewares"
	"github.com/famkampm/nentrytask/pkg/responses"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/guregu/null.v3"
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
	handler.Router.PUT("/profile/image/:id", middlewares.SetMiddlewareAuthentication(handler.UpdateProfileImage))
}

func (u *UserHandler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("home")

	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}

func (u *UserHandler) Store(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	defer r.Body.Close()
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
	defer r.Body.Close()
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
	defer r.Body.Close()
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

func (u *UserHandler) UpdateProfileImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	user, err := u.UserUsecase.GetByID(context.TODO(), int64(user_id))
	if err != nil {
		formatedError := helper.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}
	err = helper.RemovePicture(user.ProfileImage.String)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	newPathImage, err := u.SaveImageToFile(r)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	// log.Println("newpathimage:", newPathImage)
	user.ProfileImage = null.StringFrom(newPathImage)
	err = u.UserUsecase.UpdateProfileImage(context.TODO(), user.ID, newPathImage)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	responses.JSON(w, http.StatusOK, "ProfileImage Updated")
}

func (u *UserHandler) SaveImageToFile(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image") //retrieve the file from form data

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "INVALID_FILE", err
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/png":
		break
	default:
		return "INVALID_FILE_TYPE", err
	}
	fileName := helper.RandToken(12)
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		return "CANT_READ_FILE_TYPE", err
	}
	// log.Println("fileName+fileEndings[0]:", fileName+fileEndings[0])
	newPath := filepath.Join(os.Getenv("IMAGE_PATH"), fileName+fileEndings[0])
	// fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		return "CANT_WRITE_FILE", err
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		return "CANT_WRITE_FILE", err
	}
	return fileName + fileEndings[0], nil

}
