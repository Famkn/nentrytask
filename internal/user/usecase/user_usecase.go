package usecase

import (
	"context"
	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/famkampm/nentrytask/pkg/helper"
	"log"
)

type userUsecase struct {
	userRepoMysql user.Repository
	userRepoRedis user.Repository
}

func NewUserUsecase(mysql user.Repository, redis user.Repository) user.Usecase {
	return &userUsecase{
		userRepoMysql: mysql,
		userRepoRedis: redis,
	}
}

func (u *userUsecase) Store(ctx context.Context, user *models.User) error {
	// THE PRIORITY IS TO STORE TO MYSQL FIRST. AND THEN STORE TO REDIS IF NECESSARY
	err := u.userRepoMysql.Store(ctx, user)
	if err != nil {
		log.Println("errror storing to mysql from user usecase.err:", err.Error())
		return err
	}
	// IF ERROR WHEN STORING TO REDIS, IT DOESN'T REALLY MATTER. SO NO ERROR. JUST LOG
	err = u.userRepoRedis.Store(ctx, user)
	if err != nil {
		log.Println("errror storing to redis from user usecase.err:", err.Error())
	}
	return nil
}

func (u *userUsecase) GetByID(ctx context.Context, id int64) (*models.User, error) {
	// GET FROM REDIS FIRST. IF NOT EXIST. GET TO DB
	user, err := u.userRepoRedis.GetByID(ctx, id)
	if err != nil {
		log.Println("usecase GET BY ID FROM REDIS err:", err.Error())
	} else {
		return user, nil
	}
	user, err = u.userRepoMysql.GetByID(ctx, id)
	if err != nil {
		log.Println("usecase get by id from mysql err:", err.Error())
		return &models.User{}, err
	}
	return user, nil
}

func (u *userUsecase) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.userRepoMysql.GetByUsername(ctx, username)
}

func (u *userUsecase) UpdateNickname(ctx context.Context, id int64, nickname string) error {
	// UPDATE MUST APPLY TO BOTH REDIS AND MYSQL
	err := u.userRepoRedis.UpdateNickname(ctx, id, nickname)
	if err != nil {
		log.Println("usecase failed to update nickname redis repo:", err.Error())
		return err
	}
	err = u.userRepoMysql.UpdateNickname(ctx, id, nickname)
	if err != nil {
		log.Println("usecase failed to update nickname mysql repo:", err.Error())
		return err
	}
	return nil
}

func (u *userUsecase) UpdateProfileImage(ctx context.Context, id int64, profile_image string) error {
	// UPDATE MUST APPLY TO BOTH REDIS AND MYSQL
	err := u.userRepoRedis.UpdateProfileImage(ctx, id, profile_image)
	if err != nil {
		log.Println("usecase failed to update profile image redis repo:", err.Error())
		return err
	}
	err = u.userRepoMysql.UpdateProfileImage(ctx, id, profile_image)
	if err != nil {
		log.Println("usecase failed to update profile image mysql repo:", err.Error())
		return err
	}
	return nil
}

func (u *userUsecase) ValidateUserPassword(ctx context.Context, username, password string) (*models.User, error) {
	user, err := u.GetByUsername(ctx, username)
	if err != nil {
		return &models.User{}, err
	}
	err = helper.VerifyPassword(user.Password, password)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil

}
