package user

import (
	"context"
	"github.com/famkampm/nentrytask/internal/models"
)

type Repository interface {
	Store(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateNickname(ctx context.Context, id int64, nickname string) error
	UpdateProfileImage(ctx context.Context, id int64, profile_image string) error

	// RemoveProfileImage(ctx context.Context) error
	// SaveProfileImage(ctx context.Context) error

}
