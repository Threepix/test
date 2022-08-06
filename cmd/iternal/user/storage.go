package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FINDOne(ctx context.Context, id string) (User, error)
	FINDAll(ctx context.Context) (u []User, err error)
	UPDATEOne(ctx context.Context, user User) error
	DELETEOne(ctx context.Context, id string) error
}
