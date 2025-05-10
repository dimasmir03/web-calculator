package auth_test

import (
	"testing"

	"github.com/dimasmir03/web-calculator-server/internal/auth"
	sqlite2 "github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthService(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	storage, _ := sqlite2.NewStorageFromDB(db)
	service := auth.NewService(storage, "test-secret")

	t.Run("Register and Login", func(t *testing.T) {
		err := service.Register("testuser", "password123")
		assert.NoError(t, err)

		token, err := service.Login("testuser", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid Login", func(t *testing.T) {
		_, err := service.Login("testuser", "wrongpass")
		assert.Error(t, err)
	})
}
