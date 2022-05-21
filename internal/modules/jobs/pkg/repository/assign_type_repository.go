package repository

import (
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/jobs/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type AssignTypesRepositoryInterface interface {
	GetByConditions(conditions map[string]interface{}) ([]entity.AssignTypes, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(conditions map[string]interface{}) (entity.AssignTypes, error)
	Create(input map[string]interface{}) (entity.AssignTypes, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
}

type AssignTypesRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewAssignTypesRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) AssignTypesRepositoryInterface {
	return &AssignTypesRepository{
		*baseRepo,
		DB,
	}
}

func (u *AssignTypesRepository) GetByConditions(conditions map[string]interface{}) ([]entity.AssignTypes, error) {
	var assignTypes []entity.AssignTypes
	err := u.DB.Where(conditions).Find(&assignTypes).Error
	return assignTypes, err
}

func (u *AssignTypesRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.AssignTypes{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *AssignTypesRepository) TakeByConditions(conditions map[string]interface{}) (entity.AssignTypes, error) {
	var assignType entity.AssignTypes
	err := u.DB.Where(conditions).Take(&assignType).Error
	return assignType, err
}

func (u *AssignTypesRepository) Create(input map[string]interface{}) (entity.AssignTypes, error) {
	var assignType entity.AssignTypes
	err := utils.MapToStruct(input, &assignType)
	if err != nil {
		return entity.AssignTypes{}, err
	}
	err = u.DB.Create(&assignType).Error
	return assignType, err
}

func (u *AssignTypesRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var assignTypeIDList []int
	err := u.DB.Model(entity.AssignTypes{}).Where(conditions).Pluck("id", &assignTypeIDList).Error
	return assignTypeIDList, err
}
