package repository

import (
	"tungnt/emmployee_manage/pkg/share/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BaseRepository struct.
type BaseRepository struct {
	Logger *logrus.Logger
}

// NewBaseRepository returns NewBaseRepository instance.
func NewBaseRepository(logger *logrus.Logger) *BaseRepository {
	return &BaseRepository{Logger: logger}
}

// Paginate scope
func Paginate(pageData map[string]int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := 1
		if valPage, ok := pageData["Page"]; ok && valPage > 0 {
			page = valPage
		}

		pageSize := utils.DefaultPerPage
		if valPageSize, ok := pageData["PerPage"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
