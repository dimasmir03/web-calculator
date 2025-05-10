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
		return nil, fmt.Errorf("error opening gorm db connection: %w", err)
	}
	storage, err := NewStorageFromDB(db)
	if err != nil {
		return nil, fmt.Errorf("error creating storage from db: %w", err)
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
	err := s.db.AutoMigrate(arr...)
	if err != nil {
		return fmt.Errorf("error migrating db: %w", err)
	}
	return nil
}

func (s *Storage) CreateExpression(expression *model.Expression) error {
	err := s.db.Create(expression).Error
	if err != nil {
		return fmt.Errorf("error creating expression: %w", err)
	}
	return nil
}

func (s *Storage) GetExpressionByID(id string) (*model.Expression, error) {
	var expr model.Expression
	err := s.db.Where("id = ?", id).First(&expr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting expression by id: %w", err)
	}
	return &expr, nil
}

func (s *Storage) UpdateExpression(expr *model.Expression) error {
	expr1, err := s.GetExpressionByID(expr.ID)
	if err != nil {
		return fmt.Errorf("error getting expression by id: %w", err)
	}
	expr1.Result = expr.Result
	expr1.Status = expr.Status
	err = s.db.Save(expr1).Error
	if err != nil {
		return fmt.Errorf("error updating expression: %w", err)
	}
	return nil
}

func (s *Storage) CreateUser(user *model.User) error {
	err := s.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
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
