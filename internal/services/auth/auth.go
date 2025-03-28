package auth

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/storage"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    storage.UserSaver
	usrProvider storage.UserProvider
	appProvider storage.AppProvider
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	userSaver storage.UserSaver,
	userProvider storage.UserProvider,
	appProvider storage.AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		log:         log,
		tokenTTL:    tokenTTL,
	}
}

// errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Login выполняет аутентификацию пользователя
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		a.log.Error("failed to get user", slog.String("email", email), slog.Any("err", err))
		return "", ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid password", slog.String("email", email))
		return "", ErrInvalidCredentials
	}

	// Проверяем существование приложения
	_, err = a.appProvider.App(ctx, appID)
	if err != nil {
		a.log.Warn("app not found", slog.Int("appID", appID))
		return "", err
	}

	// Генерируем токен (здесь пока заглушка)
	token := "mocked-jwt-token"

	a.log.Info("user logged in", slog.String("email", email), slog.Int64("userID", user.ID))
	return token, nil
}

// RegisterNewUser создает нового пользователя
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	// Хешируем пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to hash password", slog.Any("err", err))
		return 0, err
	}

	// Сохраняем пользователя
	userID, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		a.log.Error("failed to save user", slog.String("email", email), slog.Any("err", err))
		return 0, err
	}

	a.log.Info("user registered", slog.String("email", email), slog.Int64("userID", userID))
	return userID, nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		a.log.Error("failed to check admin status", slog.Int64("userID", userID), slog.Any("err", err))
		return false, err
	}

	return isAdmin, nil
}
