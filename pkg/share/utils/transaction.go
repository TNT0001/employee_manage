package utils

import "gorm.io/gorm"

func Transaction(db *gorm.DB, f func(db *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var commit bool
	defer func() {
		if !commit {
			tx.Rollback()
		}
	}()

	err := f(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	commit = true
	return nil
}
