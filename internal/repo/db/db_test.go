package db

//func TestRegister(t *testing.T) {
//	mockDB, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	defer mockDB.Close()
//
//	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
//	assert.NoError(t, err)
//
//	if err = db.AutoMigrate(&md.Service{}); err != nil {
//		zap.L().Fatal("failed to migrate the database", zap.Error(err))
//	}
//
//	repo := Repository{conn: db}
//
//	t.Run("Success case", func(t *testing.T) {
//
//	})
//}
