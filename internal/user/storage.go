package user

import "context"

// Storage descibes a user storage functionality.
type Storage interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	UpdatePartially(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
