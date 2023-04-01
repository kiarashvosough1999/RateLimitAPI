//  Copyright 2023 KiarashVosough and other contributors
//
//  Permission is hereby granted, free of charge, to any person obtaining
//  a copy of this software and associated documentation files (the
//  Software"), to deal in the Software without restriction, including
//  without limitation the rights to use, copy, modify, merge, publish,
//  distribute, sublicense, and/or sell copies of the Software, and to
//  permit persons to whom the Software is furnished to do so, subject to
//  the following conditions:
//
//  The above copyright notice and this permission notice shall be
//  included in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
//  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
//  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
//  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
