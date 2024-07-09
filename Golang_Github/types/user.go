package types

import (
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minUsernameLen = 2
	minPasswordLen = 4
	bcrypt_cost    = 12
)

type CreateUserParams struct {
	TeleID   int64  `json:"teleID"`
	Username string `json:"user"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TeleID            int64              `json:"teleID"`
	Username          string             `bson:"user" json:"user"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encrypted_password" json:"encrypted_password"`
}

func MaskPassword(password string) string {
	if len(password) <= 2 {
		return password
	}
	return string(password[0]) + strings.Repeat("*", len(password)-2) + string(password[len(password)-1])
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt_cost)
	if err != nil {
		return nil, err
	}
	return &User{
		TeleID:            params.TeleID,
		Username:          params.Username,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func (p CreateUserParams) Check() []string {
	errors := []string{}

	if p.TeleID == 0 {
		errors = append(errors, "teleID is required")
	}

	if len(p.Username) < minUsernameLen {
		errors = append(errors, fmt.Sprintf("username lenght should be at least %d characters", minUsernameLen))
	}

	if len(p.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password lenght should be at least %d characters", minPasswordLen))
	}
	if !isEmailValid(p.Email) {
		errors = append(errors, "invalid email")
	}

	return errors
}

// func IsValidPassword(encpw, pw string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
// }

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
