package language

// Storage descibes a language storage functionality.
type Storage interface {
	Create(genre *Language) (*Language, error)
	FindById(id int16) (*Language, error)
	Update(genre *UpdateLanguageDTO) error
	Delete(id int16) error
}
