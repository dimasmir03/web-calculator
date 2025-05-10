package sqlite_test

import (
	"testing"

	"github.com/dimasmir03/web-calculator-server/internal/model"
	sqlite2 "github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStorage(t *testing.T) {
	// Используем SQLite в памяти для тестов
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	storage, _ := sqlite2.NewStorageFromDB(db)

	t.Run("Create and Get Expression", func(t *testing.T) {
		expr := &model.Expression{
			ID:         uuid.New().String(),
			UserID:     "test-user",
			Expression: "2+2",
			Status:     "pending",
		}

		// Тест создания
		err := storage.CreateExpression(expr)
		assert.NoError(t, err)
		assert.NotEmpty(t, expr.ID)

		// Тест получения
		found, err := storage.GetExpressionByID(expr.ID)
		assert.NoError(t, err)
		assert.Equal(t, expr.Expression, found.Expression)
	})

	t.Run("Restore Expressions", func(t *testing.T) {
		// Создаем тестовые данные
		expr1 := &model.Expression{ID: uuid.New().String(), UserID: "test-user", Expression: "2+2", Status: "pending"}
		expr2 := &model.Expression{ID: uuid.New().String(), UserID: "test-user", Expression: "3*3", Status: "pending"}
		storage.CreateExpression(expr1)
		storage.CreateExpression(expr2)

		// Тест восстановления
		exprs, err := storage.RestoreExpressions()
		assert.NoError(t, err)
		assert.Len(t, exprs, 2)
	})
}
