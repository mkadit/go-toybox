package users

import "golang.org/x/crypto/bcrypt"

type Users struct {
	secretKey string
}

func New(secretKey string) *Users {
	return &Users{
		secretKey: secretKey,
	}
}

// HashPassword Hash a given password
func (u *Users) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash Check a password by comparing a hash
func (u *Users) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
