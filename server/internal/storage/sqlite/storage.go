package sqlite

import (
	"errors"

	"github.com/dimasmir03/web-calculator-server/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(url string) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open(url))
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (s *Storage) Migrate() error {
	arr := []interface{}{
		&model.Expression{},
		&model.User{},
		&model.Operation{},
	}
	return s.db.AutoMigrate(arr...)
}

func (s *Storage) CreateExpression(expression *model.Expression) error {
	return s.db.Create(expression).Error
}

func (s *Storage) GetExpressionByID(id string) (*model.Expression, error) {
	var expr model.Expression
	err := s.db.Where("id = ?", id).First(&expr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &expr, err
}

func (s *Storage) UpdateExpression(expr *model.Expression) error {
	expr1, err := s.GetExpressionByID(expr.ID)
	if err != nil {
		return err
	}
	expr1.Result = expr.Result
	expr1.Status = expr.Status
	return s.db.Save(expr1).Error
}

func (s *Storage) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *Storage) GetUserByLogin(login string) (*model.User, error) {
	var user model.User
	err := s.db.Where("login = ?", login).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
