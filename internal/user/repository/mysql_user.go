package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/internal/user"
	"gopkg.in/guregu/null.v3"
)

type mysqlUserRepository struct {
	DB *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) user.Repository {
	return &mysqlUserRepository{
		DB: db,
	}

}

func (m *mysqlUserRepository) Store(ctx context.Context, user *models.User) error {
	// query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	query := `insert into user (username, password, nickname, profile_image) values (?, ?, ?, ?)`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
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

func (m *mysqlUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `select * from user where id= ?`
	err := m.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `select * from user where username= ?`
	err := m.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (m *mysqlUserRepository) UpdateNickname(ctx context.Context, id int64, nickname string) error {
	query := `update user set nickname = ? where id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("prepared failed:", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, null.StringFrom(nickname), id)
	if err != nil {
		log.Println("exec failed", err.Error())
		return err
	}
	return nil
}

func (m *mysqlUserRepository) UpdateProfileImage(ctx context.Context, id int64, profile_image string) error {
	query := `update user set profile_image = ? where id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {

		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, null.StringFrom(profile_image), id)
	if err != nil {
		return err
	}
	return nil
}
