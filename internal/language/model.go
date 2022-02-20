package language

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Language represents the language model.
type Language struct {
	Id       int16  `json:"id" example:"123"`
	Language string `json:"language" example:"ru"`
} // @name Language

// CreateLanguageDTO is used to create language.
type CreateLanguageDTO struct {
	Language string `json:"language" example:"ru"`
} // @name CreateLanguageInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (l *CreateLanguageDTO) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(
			&l.Language,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
	)
}

// UpdateLanguageDTO is used to update language record.
type UpdateLanguageDTO struct {
	Id       int16  `json:"-"`
	Language string `json:"language" example:"en"`
} // @name UpdateLanguageInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (l *UpdateLanguageDTO) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(
			&l.Language,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
	)
}
