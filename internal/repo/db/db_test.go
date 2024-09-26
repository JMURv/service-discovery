package db

import (
	"context"
	md "github.com/JMURv/service-discovery/pkg/model"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
	"testing"
)

func TestRegister(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	if err = db.AutoMigrate(&md.Service{}); err != nil {
		zap.L().Fatal("failed to migrate the database", zap.Error(err))
	}

	repo := Repository{conn: db}

	t.Run("Success case", func(t *testing.T) {
		mock.ExpectQuery("SELECT * FROM")
		err := repo.Register(context.Background(), "test-name", "test-addr")
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
