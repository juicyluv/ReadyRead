package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model.
type User struct {
	Id           int64  `json:"id" example:"123"`
	Email        string `json:"email" example:"admin@example.com"`
	Username     string `json:"username" example:"admin"`
	Password     string `json:"-"`
	Verified     bool   `json:"verified" example:"true"`
	RegisteredAt string `json:"registeredAt" example:"2022/02/24"`
} // @name User

// HashPassword will encrypt current user password.
// Returns an error on failure.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword compares hashed user password with given raw password.
// If it doesn't match, returns false.
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// CreateUserDTO is used to create user.
type CreateUserDTO struct {
	Email          string  `json:"email"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	Address        *string `json:"address,omitempty"`
	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	RepeatPassword string  `json:"repeatPassword"`
} // @name CreateUserInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *CreateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Email,
			is.Email,
			validation.Required,
		),
		validation.Field(
			&u.Username,
			is.Alphanumeric,
			validation.Length(3, 20),
			validation.Required,
		),
		validation.Field(
			&u.Address,
			is.ASCII,
			validation.Length(3, 100),
		),
		validation.Field(
			&u.PhoneNumber,
			is.Alphanumeric,
			validation.Length(5, 12),
		),
		validation.Field(
			&u.Password,
			is.Alphanumeric,
			validation.Length(6, 24),
			validation.Required,
		),
		validation.Field(
			&u.RepeatPassword,
			is.Alphanumeric,
			validation.Length(6, 24),
			validation.Required,
		),
	)
}

// UpdateUserDTO is used to update user record.
type UpdateUserDTO struct {
	Id          int64   `json:"-"`
	Email       *string `json:"email"`
	Username    *string `json:"username"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
} // @name UpdateUserInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *UpdateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Username, validation.Length(3, 20), is.Alphanumeric),
		validation.Field(&u.OldPassword, is.Alphanumeric, validation.Required),
		validation.Field(&u.NewPassword, validation.Length(6, 24), is.Alphanumeric),
	)
}
