package author

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Author represents the author model.
type Author struct {
	Id      int64  `json:"id" example:"123"`
	Name    string `json:"name" example:"Ilya"`
	Surname string `json:"surname" example:"Sokolov"`
} // @name Author

// CreateAuthorDTO is used to create author.
type CreateAuthorDTO struct {
	Name    string `json:"name" example:"Ilya"`
	Surname string `json:"surname" example:"Sokolov"`
} // @name CreateAuthorInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (a *CreateAuthorDTO) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Name,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
		validation.Field(
			&a.Surname,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
	)
}

// UpdateAuthorDTO is used to update author record.
type UpdateAuthorDTO struct {
	Id      int64  `json:"-"`
	Name    string `json:"name" example:"Ilya"`
	Surname string `json:"surname" example:"Sokolov"`
} // @name UpdateAuthorInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (a *UpdateAuthorDTO) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Name, is.Alpha, validation.Length(1, 30), validation.Required),
		validation.Field(&a.Surname, is.Alpha, validation.Length(1, 30), validation.Required),
	)
}

// UpdateAuthorPartiallyDTO is used to update author record.
type UpdateAuthorPartiallyDTO struct {
	Id      int64   `json:"-"`
	Name    *string `json:"name" example:"Ilya"`
	Surname *string `json:"surname" example:"Sokolov"`
} // @name UpdateAuthorPartiallyInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (a *UpdateAuthorPartiallyDTO) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Name, is.Alpha, validation.Length(1, 30)),
		validation.Field(&a.Surname, is.Alpha, validation.Length(1, 30)),
	)
}
