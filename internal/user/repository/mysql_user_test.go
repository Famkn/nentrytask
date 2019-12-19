package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user/repository"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestStoreSuccess(t *testing.T) {
	// Creates sqlmock database connection and a mock to manage expectations.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// Closes the database and prevents new queries from starting.
	defer db.Close()

	prep := mock.ExpectPrepare("insert into user")
	prep.ExpectExec().WithArgs("user1", "pass1", "nick1", "prof1").
		WillReturnResult(sqlmock.NewResult(1, 1))

		// a := articleRepo.NewMysqlArticleRepository(db)
	u := repository.NewMysqlUserRepository(db)
	user := &models.User{
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	err = u.Store(context.TODO(), user)
	assert.NoError(t, err)
}

func TestStoreFailedPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("insert into user").WillReturnError(fmt.Errorf("some error"))
	u := repository.NewMysqlUserRepository(db)
	user := &models.User{}
	err = u.Store(context.TODO(), user)
	assert.NotNil(t, err)
}

func TestStoreFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("insert into user")
	prep.ExpectExec().WithArgs("user1", "pass1", "nick1", "prof1").
		WillReturnError(fmt.Errorf("some error"))
	u := repository.NewMysqlUserRepository(db)
	user := &models.User{
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	err = u.Store(context.TODO(), user)
	assert.NotNil(t, err)
}

func TestGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	defer db.Close()

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"id", "username", "password", "nickname", "profile_image"}).
		AddRow(1, "user1", "pass1", "nick1", "prof1")

	mock.ExpectQuery("select \\* from user where id= \\?").WithArgs(1).WillReturnRows(rows)
	u := repository.NewMysqlUserRepository(db)
	user, err := u.GetByID(context.TODO(), 1)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, int64(1))
}

func TestGetByIDFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("select \\* from user where id= \\?").WithArgs(1).WillReturnError(fmt.Errorf("some error"))
	u := repository.NewMysqlUserRepository(db)
	_, err = u.GetByID(context.TODO(), 1)
	assert.NotNil(t, err)
}

func TestGetByUsernameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"id", "username", "password", "nickname", "profile_image"}).
		AddRow(1, "user1", "pass1", "nick1", "prof1")

	mock.ExpectQuery("select \\* from user where username= \\?").WithArgs("user1").
		WillReturnRows(rows)
	u := repository.NewMysqlUserRepository(db)
	user, err := u.GetByUsername(context.TODO(), "user1")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, int64(1))
}

func TestGetByUsernameFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("select \\* from user where username= \\?").WithArgs("user1").
		WillReturnError(fmt.Errorf("some error"))

	u := repository.NewMysqlUserRepository(db)
	_, err = u.GetByUsername(context.TODO(), "user1")
	assert.NotNil(t, err)
}

func TestUpdateNicknameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	id := 1
	nickname := "nick1"

	prep := mock.ExpectPrepare("update user set nickname = \\? where id = \\?")
	prep.ExpectExec().WithArgs(nickname, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateNickname(context.TODO(), int64(id), nickname)
	assert.Nil(t, err)
}

func TestUpdateNicknameFailedPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectPrepare("update user set nickname = ? where id = ?").
		WillReturnError(fmt.Errorf("prepared error"))

	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.NotNil(t, err)
}

func TestUpdateNickNameFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	prep := mock.ExpectPrepare("update user set nickname = \\? where id = \\?")
	prep.ExpectExec().WithArgs("nickname", int64(1)).
		WillReturnError(fmt.Errorf("some error"))

	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateNickname(context.TODO(), int64(1), "nickname")
	assert.NotNil(t, err)
}

func TestUpdateProfileImageSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	id := int64(1)
	profile_image := "nick1"

	prep := mock.ExpectPrepare("update user set profile_image = \\? where id = \\?")
	prep.ExpectExec().WithArgs(id, profile_image).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateProfileImage(context.TODO(), id, profile_image)
	assert.Nil(t, err)
}

func TestUpdateProfileImageFailedPrepared(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("update user set profile_image = \\? where id = \\?").
		WillReturnError(fmt.Errorf("some error"))

	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateProfileImage(context.TODO(), int64(1), "prof")
	assert.NotNil(t, err)
}

func TestUpdateProfileImageFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	prep := mock.ExpectPrepare("update user set profile_image = \\? where id = \\?")
	prep.ExpectExec().WithArgs("prof1", int64(1)).
		WillReturnError(fmt.Errorf("some error"))
	u := repository.NewMysqlUserRepository(db)
	err = u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.NotNil(t, err)
}
