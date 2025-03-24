package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password  string `gorm:"not null"`
	Admin     bool   `gorm:"default:false"`
	AuthToken string `gorm:"index"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword compares input password with hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
