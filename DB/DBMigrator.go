package DB

import (
	"RateLimitAPI/Models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Migrator struct {
	db *gorm.DB
}

func New() (*Migrator, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Migrator{
		db,
	}, nil
}

func (m *Migrator) AutoMigrateModels() {
	err := m.db.AutoMigrate(&Models.UserModel{})
	if err != nil {
		log.Print("migration failed")
	}
}

// UserRepositoryInterface Implementation

func (m *Migrator) UsernameExist(username string) (bool, error) {
	var count int64 = 0
	tx := m.db.Model(&Models.UserModel{}).Where("username = ?", username).Count(&count)
	if tx.Error != nil {
		return false, &MigratorError{
			CouldNotFind,
		}
	} else if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *Migrator) FindByUsername(username string) (*Models.UserModel, error) {
	var product Models.UserModel
	tx := m.db.First(&product, "username = ?", username)
	if tx.Error != nil {
		return nil, &MigratorError{
			CouldNotFind,
		}
	} else {
		return &product, nil
	}
}

func (m *Migrator) Save(user Models.UserModel) error {
	tx := m.db.Create(&user)
	if tx.Error != nil {
		return &MigratorError{
			CouldNotCreate,
		}
	} else {
		return nil
	}
}
