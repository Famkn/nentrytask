package repository_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user/repository"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestStoreSuccessRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("SET", "1", "{\"id\":1,\"username\":\"user1\",\"password\":\"pass1\",\"nickname\":\"nick1\",\"profile_image\":\"prof1\"}").Expect("OK!")
	u := repository.NewRedisUserRepository(conn)
	user := &models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	err := u.Store(context.TODO(), user)
	assert.Nil(t, err)
}

func TestStoreFailedRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("SET", "1", "{}").ExpectError(fmt.Errorf("Some error"))
	u := repository.NewRedisUserRepository(conn)
	err := u.Store(context.TODO(), &models.User{})
	assert.NotNil(t, err)
}

func TestGetByIDSuccessRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("{\"username\" : \"user1\", \"id\" : 1, \"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")
	u := repository.NewRedisUserRepository(conn)
	user, err := u.GetByID(context.TODO(), int64(1))
	assert.Nil(t, err)
	assert.Equal(t, user.ID, int64(1))
	assert.Equal(t, user.Username, "user1")
	assert.Equal(t, user.Nickname, null.StringFrom("nick1"))
	assert.Equal(t, user.ProfileImage, null.StringFrom("prof1"))
}

func TestGetByIDFailedRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").ExpectError(fmt.Errorf("some error"))
	u := repository.NewRedisUserRepository(conn)
	user, err := u.GetByID(context.TODO(), int64(1))
	log.Println("USER GET APA ISINYA:", user)
	log.Println("ERROR NYA APA NII:", err.Error())
	assert.NotNil(t, err)
}

func TestGetByIDFailedUnmarshalRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("123")
	u := repository.NewRedisUserRepository(conn)
	_, err := u.GetByID(context.TODO(), int64(1))
	assert.NotNil(t, err)
}

func TestGetByUsernameSuccessRedis(t *testing.T) {
	conn := redigomock.NewConn()
	u := repository.NewRedisUserRepository(conn)
	_, err := u.GetByUsername(context.TODO(), "user1")
	assert.Nil(t, err)
}

func TestUpdateNicknameSuccessRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("{\"id\" : 1, \"username\" : \"user1\", \"password\":\"pass1\",\"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")
	conn.Command("SET", "1", "{\"id\":1,\"username\":\"user1\",\"password\":\"pass1\",\"nickname\":\"nick1\",\"profile_image\":\"prof1\"}").Expect("OK!")
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.Nil(t, err)
}

func TestUpdateNicknameFailedGetIDRedis(t *testing.T) {
	conn := redigomock.NewConn()
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.NotNil(t, err)
}

func TestUpdateNicknameFailedStoreRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("{\"id\" : 1, \"username\" : \"user1\", \"password\":\"pass1\",\"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.NotNil(t, err)
}

func TestUpdateProfileImageSuccessRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("{\"id\" : 1, \"username\" : \"user1\", \"password\":\"pass1\",\"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")
	conn.Command("SET", "1", "{\"id\":1,\"username\":\"user1\",\"password\":\"pass1\",\"nickname\":\"nick1\",\"profile_image\":\"prof1\"}").Expect("OK!")
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.Nil(t, err)
}

func TestUpdateProfileImageFailedGetIDRedis(t *testing.T) {
	conn := redigomock.NewConn()
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.NotNil(t, err)
}

func TestUpdateProfileImageFailedStoreRedis(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("GET", "1").Expect("{\"id\" : 1, \"username\" : \"user1\", \"password\":\"pass1\",\"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")
	u := repository.NewRedisUserRepository(conn)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.NotNil(t, err)
}
