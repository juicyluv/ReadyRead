package author

// Storage descibes author storage functionality.
type Storage interface {
	Create(user *Author) (*Author, error)
	FindById(id int64) (*Author, error)
	Update(user *UpdateAuthorDTO) error
	UpdatePartially(user *UpdateAuthorPartiallyDTO) error
	Delete(id int64) error
}
