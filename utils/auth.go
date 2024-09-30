package utils

import "golang.org/x/crypto/bcrypt"

type authUtils struct{}

type AuthUtils interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

func (u *authUtils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *authUtils) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAuthUtils() AuthUtils {
	return &authUtils{}
}
