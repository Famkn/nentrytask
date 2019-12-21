package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user/mocks"
	"github.com/famkampm/nentrytask/internal/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v3"
)

func TestStoreSuccessUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := &models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}

	mockUserRepoMysql.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
	mockUserRepoRedis.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.Store(context.TODO(), mockUser)
	assert.NoError(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)

}

func TestStoreFailedMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := &models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}

	mockUserRepoMysql.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("Unexpected")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.Store(context.TODO(), mockUser)
	assert.Error(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestStoreFailedRedisUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := &models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	mockUserRepoMysql.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
	mockUserRepoRedis.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("Unexpected")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.Store(context.TODO(), mockUser)
	assert.NoError(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)

}

func TestGetByIDSuccessRedisUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	mockUserRepoRedis.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockUser, nil).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	user, err := u.GetByID(context.TODO(), mockUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestGetByIDFailedRedisSuccessMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := models.User{
		ID:           int64(1),
		Username:     "user1",
		Password:     "pass1",
		Nickname:     null.StringFrom("nick1"),
		ProfileImage: null.StringFrom("prof1"),
	}
	mockUserRepoRedis.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&models.User{}, errors.New("Unexpected")).Once()
	mockUserRepoMysql.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockUser, nil).Once()

	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	user, err := u.GetByID(context.TODO(), mockUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestGetByIDFailedUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUser := models.User{
		ID: int64(1),
	}
	mockUserRepoRedis.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&models.User{}, errors.New("Unexpected")).Once()
	mockUserRepoMysql.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&models.User{}, errors.New("Unexpected")).Once()

	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	user, err := u.GetByID(context.TODO(), mockUser.ID)
	assert.Error(t, err)
	assert.NotNil(t, user)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestGetByUsernameSuccessUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)

	mockUserRepoMysql.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&models.User{}, nil).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	user, err := u.GetByUsername(context.TODO(), "user1")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestGetByUsernameFailedUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)

	mockUserRepoMysql.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&models.User{}, errors.New("some error")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	user, err := u.GetByUsername(context.TODO(), "user1")
	assert.Error(t, err)
	assert.NotNil(t, user)

	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestUpdateNicknameSuccessUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateNickname", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	mockUserRepoMysql.On("UpdateNickname", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.NoError(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)

}

func TestUpdateNicknameFailedRedisSuccessMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateNickname", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(errors.New("some error")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.Error(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestUpdateNicknameSuccRedisSuccessFailedMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateNickname", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	mockUserRepoMysql.On("UpdateNickname", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(errors.New("some error")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateNickname(context.TODO(), int64(1), "nick1")
	assert.Error(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}
func TestUpdateProfileImageSuccessUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateProfileImage", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	mockUserRepoMysql.On("UpdateProfileImage", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.NoError(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestUpdateProfileImageFailedRedisSuccessMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateProfileImage", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(errors.New("some error")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.Error(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}

func TestUpdateProfileImageSuccessRedisFailedMysqlUsecase(t *testing.T) {
	mockUserRepoMysql := new(mocks.Repository)
	mockUserRepoRedis := new(mocks.Repository)
	mockUserRepoRedis.On("UpdateProfileImage", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(nil).Once()
	mockUserRepoMysql.On("UpdateProfileImage", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(errors.New("some error")).Once()
	u := usecase.NewUserUsecase(mockUserRepoMysql, mockUserRepoRedis)
	err := u.UpdateProfileImage(context.TODO(), int64(1), "prof1")
	assert.Error(t, err)
	mockUserRepoMysql.AssertExpectations(t)
	mockUserRepoRedis.AssertExpectations(t)
}
