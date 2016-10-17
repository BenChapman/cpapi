package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserType int

const (
	ROLE_ADMIN UserType = iota
	ROLE_NINJA
	ROLE_PARENT
	ROLE_MENTOR
)

type User struct {
	gorm.Model
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	UserType    UserType  `json:"-"`
	Token       string    `json:"-"`
}
