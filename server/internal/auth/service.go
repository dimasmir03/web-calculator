package auth

import (
	"time"

	"github.com/dimasmir03/web-calculator-server/internal/model"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db        *sqlite.Storage
	jwtSecret string
}

func NewService(db *sqlite.Storage, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Login(login, password string) (string, error) {
	user, err := s.db.GetUserByLogin(login)
	if err != nil {
		return "", echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", echo.ErrUnauthorized
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) Register(login, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.db.CreateUser(&model.User{
		Login:        login,
		PasswordHash: string(hash),
	}); err != nil {
		return err
	}

	return nil
}
