// package repository_test

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/famkampm/nentrytask/internal/models"
// 	"github.com/famkampm/nentrytask/internal/user/repository"
// 	"github.com/gomodule/redigo/redis"
// 	"github.com/rafaeljusto/redigomock"

// 	// "github.com/rafaeljusto/redigomock"
// 	"github.com/stretchr/testify/assert"
// 	"gopkg.in/guregu/null.v3"
// )

// func TestStoreSuccess(t *testing.T) {
// 	// Creates sqlmock database connection and a mock to manage expectations.
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	// Closes the database and prevents new queries from starting.
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	prep := mock.ExpectPrepare("insert into user")
// 	prep.ExpectExec().WithArgs("user1", "pass1", "nick1", "prof1").
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	u := repository.NewUserRepository(db, conn)
// 	user := &models.User{
// 		Username:     "user1",
// 		Password:     "pass1",
// 		Nickname:     null.StringFrom("nick1"),
// 		ProfileImage: null.StringFrom("prof1"),
// 	}
// 	err = u.Store(context.TODO(), user)
// 	assert.NoError(t, err)
// }

// func TestStoreFailedPrepare(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()

// 	mock.ExpectPrepare("insert into user").WillReturnError(fmt.Errorf("some error"))
// 	u := repository.NewUserRepository(db, conn)
// 	user := &models.User{}
// 	err = u.Store(context.TODO(), user)
// 	assert.NotNil(t, err)
// }

// //
// func TestStoreFailed(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()

// 	prep := mock.ExpectPrepare("insert into user")
// 	prep.ExpectExec().WithArgs("user1", "pass1", "nick1", "prof1").
// 		WillReturnError(fmt.Errorf("some error"))
// 	u := repository.NewUserRepository(db, conn)
// 	user := &models.User{
// 		Username:     "user1",
// 		Password:     "pass1",
// 		Nickname:     null.StringFrom("nick1"),
// 		ProfileImage: null.StringFrom("prof1"),
// 	}
// 	err = u.Store(context.TODO(), user)
// 	assert.NotNil(t, err)
// }

// func TestGetByIDSuccess(t *testing.T) {
// 	db, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	cmd := conn.Command("GET", "1").Expect("{\"username\" : \"user1\", \"id\" : 1, \"nickname\" : \"nick1\", \"profile_image\" : \"prof1\"}")

// 	u := repository.NewUserRepository(db, conn)
// 	user, err := u.GetByID(context.TODO(), 1)
// 	if err != nil {
// 		log.Println("err:", err.Error())
// 		t.Fatal(err)
// 	}

// 	if conn.Stats(cmd) != 1 {
// 		t.Fatal("Command was not called!")

// 	}
// 	assert.Equal(t, user.ID, int64(1))
// 	assert.Equal(t, user.Username, "user1")
// 	assert.Equal(t, user.Nickname, null.StringFrom("nick1"))
// 	assert.Equal(t, user.ProfileImage, null.StringFrom("prof1"))
// }

// func TestGetByIDSuccessToDBAndSetToRedis(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	cmd := conn.Command("GET", "1").ExpectError(redis.ErrNil)

// 	rows := sqlmock.NewRows([]string{"id", "username", "password", "nickname", "profile_image"}).
// 		AddRow(1, "user1", "pass1", "nick1", "prof1")
// 	mock.ExpectQuery("select \\* from user where id= \\?").WithArgs(1).WillReturnRows(rows)

// 	cmd2 := conn.Command("SET", "1", "{\"id\":1,\"username\":\"user1\",\"password\":\"pass1\",\"nickname\":\"nick1\",\"profile_image\":\"prof1\"}").Expect("OK!")
// 	//   conn.Commmand("HMSET", "person:1", "name", person.Name, "age", person.Age, "updatedat", redigomock.NewAnyInt(), "phone", person.Phone).Expect("OK!")

// 	u := repository.NewUserRepository(db, conn)
// 	user, err := u.GetByID(context.TODO(), 1)
// 	if conn.Stats(cmd) != 1 {
// 		t.Fatal("Command was not called!")
// 	}
// 	if conn.Stats(cmd2) != 1 {
// 		t.Fatal("Command2 was not called!")
// 	}
// 	assert.Nil(t, err)
// 	assert.Equal(t, user.ID, int64(1))
// 	assert.Equal(t, user.Username, "user1")
// 	assert.Equal(t, user.Nickname, null.StringFrom("nick1"))
// 	assert.Equal(t, user.ProfileImage, null.StringFrom("prof1"))
// }

// func TestGetByIDSuccessToDBAndFailedSetToRedis(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	cmd := conn.Command("GET", "1").ExpectError(redis.ErrNil)

// 	rows := sqlmock.NewRows([]string{"id", "username", "password", "nickname", "profile_image"}).
// 		AddRow(1, "user1", "pass1", "nick1", "prof1")
// 	mock.ExpectQuery("select \\* from user where id= \\?").WithArgs(1).WillReturnRows(rows)
// 	u := repository.NewUserRepository(db, conn)
// 	user, err := u.GetByID(context.TODO(), 1)
// 	if conn.Stats(cmd) != 1 {
// 		t.Fatal("Command was not called!")
// 	}
// 	assert.NotNil(t, err)
// 	assert.Equal(t, user.ID, int64(1))
// 	assert.Equal(t, user.Username, "user1")
// 	assert.Equal(t, user.Nickname, null.StringFrom("nick1"))
// 	assert.Equal(t, user.ProfileImage, null.StringFrom("prof1"))
// }

// func TestGetByIDFailedUnmarshal(t *testing.T) {
// 	db, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	cmd := conn.Command("GET", "1").Expect("123")

// 	u := repository.NewUserRepository(db, conn)
// 	_, err = u.GetByID(context.TODO(), 1)
// 	if conn.Stats(cmd) != 1 {
// 		t.Fatal("Command was not called!")
// 	}
// 	assert.NotNil(t, err)
// }

// func TestGetByUsernameSuccess(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	// before we actually execute our api function, we need to expect required DB actions
// 	rows := sqlmock.NewRows([]string{"id", "username", "password", "nickname", "profile_image"}).
// 		AddRow(1, "user1", "pass1", "nick1", "prof1")

// 	mock.ExpectQuery("select \\* from user where username= \\?").WithArgs("user1").
// 		WillReturnRows(rows)
// 	u := repository.NewUserRepository(db, conn)
// 	user, err := u.GetByUsername(context.TODO(), "user1")
// 	assert.Nil(t, err)
// 	assert.Equal(t, user.ID, int64(1))
// 	assert.Equal(t, user.Username, "user1")
// 	assert.Equal(t, user.Nickname, null.StringFrom("nick1"))
// 	assert.Equal(t, user.ProfileImage, null.StringFrom("prof1"))
// }

// func TestGetByUsernameFailed(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	conn := redigomock.NewConn()
// 	mock.ExpectQuery("select \\* from user where username= \\?").WithArgs("user1").
// 		WillReturnError(fmt.Errorf("some error"))
// 	u := repository.NewUserRepository(db, conn)
// 	_, err = u.GetByUsername(context.TODO(), "user1")
// 	assert.NotNil(t, err)
// }
