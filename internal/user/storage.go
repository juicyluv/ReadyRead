package user

// Storage descibes user storage functionality.
type Storage interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindById(id int64) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *UpdateUserDTO) error
	UpdatePartially(user *UpdateUserPartiallyDTO) error
	Delete(id int64) error
}
