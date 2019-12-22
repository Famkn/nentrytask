package repository

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/guregu/null.v3"
)

type redisUserRepository struct {
	// Redis redis.Conn
	RedisPool *redis.Pool
}

// func NewRedisUserRepository(conn redis.Conn) user.Repository {
// 	return &redisUserRepository{
// 		Redis: conn,
// 	}
// }

func NewRedisUserRepository(redisPool *redis.Pool) user.Repository {
	return &redisUserRepository{
		RedisPool: redisPool,
	}
}

// func NewRedisUer
// *redis.Pool

func (r *redisUserRepository) Store(ctx context.Context, user *models.User) error {
	// serialize user object
	json, err := json.Marshal(user)
	if err != nil {
		// log.Println("SetUserToCache error marshail", err.Error())
		return err
	}
	// SET object
	// log.Println("json dari user:", string(json))
	conn := r.RedisPool.Get()
	defer conn.Close()
	// _, err = r.Redis.Do("SET", strconv.Itoa(int(user.ID)), string(json))
	_, err = conn.Do("SET", strconv.Itoa(int(user.ID)), string(json))
	if err != nil {
		// log.Println("SetUserToCache error set to redis", err.Error())
		return err
	}
	return nil
}

func (r *redisUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	conn := r.RedisPool.Get()
	defer conn.Close()

	// user_temp, err := redis.String(r.Redis.Do("GET", strconv.Itoa(int(id))))
	user_temp, err := redis.String(conn.Do("GET", strconv.Itoa(int(id))))
	// log.Println("user_temp:", user_temp)

	if err != nil {
		log.Println("getbyid err1:", err.Error())
		return &models.User{}, err
	}
	user := &models.User{}
	err = json.Unmarshal([]byte(user_temp), &user)
	if err != nil {
		// log.Println("getbyid err2:", err.Error())
		return &models.User{}, err
	}
	return user, nil

}

func (r *redisUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return &models.User{}, nil
}

func (r *redisUserRepository) UpdateNickname(ctx context.Context, id int64, nickname string) error {
	// log.Println("update nickname caled")
	user, err := r.GetByID(ctx, id)
	if err != nil {
		log.Println("update nickname getbyid err:", err.Error())
		return err
	}
	user.Nickname = null.StringFrom(nickname)
	err = r.Store(ctx, user)
	if err != nil {
		log.Println("update nickname store err:", err.Error())
		return err
	}
	return nil

}

func (r *redisUserRepository) UpdateProfileImage(ctx context.Context, id int64, profile_image string) error {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.ProfileImage = null.StringFrom(profile_image)
	err = r.Store(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
