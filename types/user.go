package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("firstName length should be at least %d charachters", minFirstNameLength))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("lastName length should be at least %d charachters", minLastNameLength))
	}
	if len(params.Password) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password length should be at least %d charachters", minPasswordLength))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, "email is invalid")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName        string             `bson:"firstName" json:"firstName"`
	LastName         string             `bson:"lastName" json:"lastName"`
	Email            string             `bson:"email" json:"email"`
	EncryptedPasswod string             `bson:"encryptedPasswod" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:        params.FirstName,
		LastName:         params.LastName,
		Email:            params.Email,
		EncryptedPasswod: string(hashedPassword),
	}, nil
}
