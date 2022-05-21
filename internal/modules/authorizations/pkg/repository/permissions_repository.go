package repository

import (
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type PermissionsRepositoryInterface interface {
	GetByConditions(conditions map[string]interface{}) ([]entity.Permissions, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(conditions map[string]interface{}) (entity.Permissions, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.Permissions, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
}

type PermissionsRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewPermissionsRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) PermissionsRepositoryInterface {
	return &PermissionsRepository{
		*baseRepo,
		DB,
	}
}

func (u *PermissionsRepository) GetByConditions(conditions map[string]interface{}) ([]entity.Permissions, error) {
	var permissions []entity.Permissions
	err := u.DB.Where(conditions).Find(&permissions).Error
	return permissions, err
}

func (u *PermissionsRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.Permissions{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *PermissionsRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.Permissions, error) {
	var permissions []entity.Permissions
	for _, attribute := range attributes {
		var temp entity.Permissions
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, temp)
	}

	err := u.DB.Create(&permissions).Error
	return permissions, err
}

func (u *PermissionsRepository) TakeByConditions(conditions map[string]interface{}) (entity.Permissions, error) {
	var permissions entity.Permissions
	err := u.DB.Where(conditions).Take(&permissions).Error
	return permissions, err
}

func (u *PermissionsRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var permissionIDList []int
	err := u.DB.Model(entity.Permissions{}).
		Where(conditions).
		Pluck("id", &permissionIDList).
		Error
	return permissionIDList, err
}
