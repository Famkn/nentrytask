package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"gopkg.in/guregu/null.v3"
)

type memoryUserRepository struct {
	hm map[int64]string
}

func NewMemoryUserRepository(newmap map[int64]string) user.Repository {
	return &memoryUserRepository{
		hm: newmap,
	}
}
func (m *memoryUserRepository) Store(ctx context.Context, user *models.User) error {
	user_byte, err := json.Marshal(user)
	if err != nil {
		log.Println("memory err store. marsshal:", err.Error())
		return err
	}
	m.hm[user.ID] = string(user_byte)
	return nil
}
func (m *memoryUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	ret := m.hm[id]
	if ret == "" {
		return &models.User{}, errors.New("ID NOT FOUND")
	}
	user := &models.User{}
	err := json.Unmarshal([]byte(ret), &user)
	if err != nil {
		log.Println("memory geybyid err unmarshal.", err.Error())
	}
	return user, nil
}
func (m *memoryUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return &models.User{}, nil
}
func (m *memoryUserRepository) UpdateNickname(ctx context.Context, id int64, nickname string) error {
	user, err := m.GetByID(ctx, id)
	if err != nil {
		log.Println("update nickname memory getbyid err:", err.Error())
		return err
	}
	user.Nickname = null.StringFrom(nickname)
	err = m.Store(ctx, user)
	if err != nil {
		log.Println("update nickname memory store err:", err.Error())
		return err
	}
	return nil

}
func (m *memoryUserRepository) UpdateProfileImage(ctx context.Context, id int64, profile_image string) error {
	user, err := m.GetByID(ctx, id)
	if err != nil {
		log.Println("update profile memory getbyid err:", err.Error())
		return err
	}
	user.ProfileImage = null.StringFrom(profile_image)
	err = m.Store(ctx, user)
	if err != nil {
		log.Println("update nickname memory store err:", err.Error())
		return err
	}
	return nil
}
