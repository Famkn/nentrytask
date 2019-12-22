package helper

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashingPassword(plain_password string) (string, error) {
	hashedPassword, err := Hash(plain_password)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Validate(action, username, password string) error {
	if strings.ToLower(action) == "login" || strings.ToLower(action) == "register" {
		if username == "" {
			return errors.New("Required Username")
		}
		if password == "" {
			return errors.New("Required Password")
		}
	}
	return nil
}

func FormatError(err string) error {

	if strings.Contains(err, "username") {
		return errors.New("Username Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}

	if strings.Contains(err, "no rows") {
		return errors.New("User not exist")
	}

	return errors.New("Incorrect Details")
}

func RemovePicture(profile_image string) error {
	if profile_image == "" {
		return nil
	}
	path := os.Getenv("IMAGE_PATH") + profile_image
	// log.Println("PATH TO REMOVE PICT:", path)
	err := os.Remove(path)
	if err != nil {
		log.Println("os remove err", err.Error())
		return err
	}
	return nil
}

func RandToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
