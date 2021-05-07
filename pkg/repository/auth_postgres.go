package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/p44elovod/auth-with-gopg/models"
)

type AuthPostgres struct {
	db *pg.DB
}

func NewAuthPostgres(db *pg.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	_, err := r.db.Model(&user).
		Returning("id").
		Insert()

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *AuthPostgres) GetUserID(username, password string) (int, error) {
	var user models.User
	err := r.db.Model(&user).
		Column("id").
		Where("username = ?", username).
		Where("password_hash = ?", password).
		Limit(1).
		Select()

	return user.ID, err
}

func (r *AuthPostgres) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).
		Where("id = ?", id).
		Select()

	return user, err
}
