package genre

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Genre represents the genre model.
type Genre struct {
	Id    int16  `json:"id" example:"123"`
	Genre string `json:"genre" example:"fantasy"`
} // @name Genre

// CreateGenreDTO is used to create genre.
type CreateGenreDTO struct {
	Genre string `json:"genre" example:"fantasy"`
} // @name CreateGenreInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (g *CreateGenreDTO) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.Genre,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
	)
}

// UpdateGenreDTO is used to update genre record.
type UpdateGenreDTO struct {
	Id    int16  `json:"-"`
	Genre string `json:"genre" example:"fantasy"`
} // @name UpdateGenreInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (g *UpdateGenreDTO) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.Genre,
			is.Alpha,
			validation.Length(1, 30),
			validation.Required,
		),
	)
}
