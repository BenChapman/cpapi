package user

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u UserRepo) GetNinjas(token string) ([]User, error) {
	ninjas := []User{}
	if u.WithToken(token).Where("user_type = ?", ROLE_NINJA).Find(&ninjas).RecordNotFound() {
		return []User{}, errors.New("Could not find user")
	}
	return ninjas, nil
}

func (u UserRepo) GetNinja(token string, id int) (User, error) {
	ninja := User{}
	if u.WithToken(token).Find(&ninja, id).RecordNotFound() {
		return ninja, errors.New(fmt.Sprintf("Could not find ninja %d", id))
	}
	return ninja, nil
}

func (u UserRepo) WithToken(token string) *gorm.DB {
	return u.db.Where("token = ?", token)
}
