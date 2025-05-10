package sqlite

import (
	"errors"
	"fmt"

	"github.com/dimasmir03/web-calculator-server/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(url string) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	storage, err := NewStorageFromDB(db)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func NewStorageFromDB(db *gorm.DB) (*Storage, error) {
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

func (s *Storage) RestoreExpressions() ([]*model.Expression, error) {
	var exprs []*model.Expression
	err := s.db.Find(&exprs).Error
	if err == nil {
		fmt.Errorf("ошибка восстановления выражений: %w", err)
	}
	return exprs, err
}
