package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/guregu/null.v3"
)

type UserRepository struct {
	DB    *sql.DB
	Redis redis.Conn
}

func NewUserRepository(db *sql.DB, redis redis.Conn) user.Repository {
	return &UserRepository{
		DB:    db,
		Redis: redis,
	}
}

func (m *UserRepository) Store(ctx context.Context, user *models.User) error {
	// query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	query := `insert into user (username, password, nickname, profile_image) values (?, ?, ?, ?)`
	stmt, err := m.DB.PrepareContext(ctx, query)

	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, user.Username, user.Password, user.Nickname, user.ProfileImage)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (m *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	user_temp, err := redis.String(m.Redis.Do("GET", strconv.Itoa(int(id))))

	if err != nil {
		// get from db
		log.Println("GET BY ID REDIS ERR NIL")
		query := `select * from user where id= ?`
		err := m.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
		if err != nil {
			log.Println("err after query:", err.Error())
			return &models.User{}, err
		}
		// set to cache
		err = m.SetUserToCache(ctx, *user)
		if err != nil {
			log.Println("err after set to cache:", err.Error())
			return user, err
		}
		log.Println("sukses err nillll. user:", user)
		return user, nil
	}
	// geting from redis
	err = json.Unmarshal([]byte(user_temp), &user)
	if err != nil {
		log.Println("error unmarshalling")
		return &models.User{}, err
	}
	log.Println("GET BY ID SUKSES")
	return user, nil
}

func (m *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `select * from user where username= ?`
	err := m.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (m *UserRepository) UpdateNickname(ctx context.Context, id int64, nickname string) error {
	query := `update user set nickname = ? where id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("prepared failed:", err.Error())
		return err
	}
	_, err = stmt.ExecContext(ctx, null.StringFrom(nickname), id)
	if err != nil {
		log.Println("exec failed", err.Error())
		return err
	}
	// SET to redis if exist
	user, err := m.GetByID(ctx, id)
	if err != nil {
		log.Println("err:", err.Error())
		return err
	}
	err = m.SetUserToCacheIfExist(ctx, *user)
	if err != nil {
		log.Println("err:", err.Error())
		return err
	}
	return nil
}

func (m *UserRepository) UpdateProfileImage(ctx context.Context, id int64, profile_image string) error {
	query := `update user set profile_image = ? where id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, id, null.StringFrom(profile_image))
	if err != nil {
		return err
	}
	// SET to redis if exist
	user, err := m.GetByID(ctx, id)
	if err != nil {
		log.Println("err:", err.Error())
		return err
	}
	err = m.SetUserToCacheIfExist(ctx, *user)
	if err != nil {
		log.Println("err:", err.Error())
		return err
	}
	return nil

}

func (m *UserRepository) SetUserToCache(ctx context.Context, user models.User) error {
	// serialize user object
	json, err := json.Marshal(user)
	if err != nil {
		log.Println("SetUserToCache error marshail", err.Error())
		return err
	}
	// SET object
	_, err = m.Redis.Do("SET", strconv.Itoa(int(user.ID)), string(json))
	if err != nil {
		log.Println("SetUserToCache error set to redis", err.Error())
		return err
	}
	return nil
}

func (m *UserRepository) SetUserToCacheIfExist(ctx context.Context, user models.User) error {
	key := strconv.Itoa(int(user.ID))
	_, err := redis.String(m.Redis.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return err
	}
	if err != redis.ErrNil {
		// USER EXIST in cache. update user information in cache
		_, err = m.Redis.Do("DEL", key)
		if err != nil {
			return err
		}
	}
	err = m.SetUserToCache(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
