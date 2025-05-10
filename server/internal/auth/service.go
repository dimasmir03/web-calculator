package auth

import (
	"time"

	"github.com/dimasmir03/web-calculator-server/internal/model"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db        *sqlite.Storage
	jwtSecret string
	log       *logrus.Logger
}

func NewService(db *sqlite.Storage, jwtSecret string) *Service {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
		log:       logger,
	}
}

func (s *Service) Login(login, password string) (string, error) {
	user, err := s.db.GetUserByLogin(login)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Error("Failed to get user by login")
		return "", echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.log.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Error("Password comparison failed")
		return "", echo.ErrUnauthorized
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Error("Failed to sign token")
		return "", err
	}

	s.log.WithFields(logrus.Fields{
		"login": login,
	}).Info("User logged in successfully")
	return signedToken, nil
}

func (s *Service) Register(login, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Error("Failed to generate password hash")
		return err
	}
	if err := s.db.CreateUser(&model.User{
		Login:        login,
		PasswordHash: string(hash),
	}); err != nil {
		s.log.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Error("Failed to create user")
		return err
	}

	return nil
}
