package genre

// Storage descibes a genre storage functionality.
type Storage interface {
	Create(genre *Genre) (*Genre, error)
	FindById(id int16) (*Genre, error)
	Update(genre *UpdateGenreDTO) error
	Delete(id int16) error
}
