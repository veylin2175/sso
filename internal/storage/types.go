package storage

import (
	"context"
	"sso/internal/domain/models"
)

type Provider interface {
	UserSaver
	UserProvider
	AppProvider
	Pinger
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}
