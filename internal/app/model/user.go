package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// NewUser The function will return a pointer to the array of type User
func NewUser() *Users {
	return &Users{}
}

// User ...
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	FirstName         string `json:"firstname"`
	LastName          string `json:"lastname"`
	Role              string `json:"role"`
	Groups            string `json:"groups"`
	Tabel             string `json:"tabel"`
	Groupmix          string `json:"groupmix"`
}

// Users ...
type Users []User

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// BeforCreate ...
func (u *User) BeforCreate() error {
	pass := u.Password
	fmt.Println("Password:", pass)

	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		fmt.Println("Hash:    ", enc)
		u.EncryptedPassword = enc
	}
	return nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
	//	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	//	return err == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
