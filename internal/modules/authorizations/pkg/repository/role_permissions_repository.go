package repository

import (
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type RolePermissionsRepositoryInterface interface {
	GetByConditions(conditions map[string]interface{}) ([]entity.RolePermissions, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(conditions map[string]interface{}) (entity.RolePermissions, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.RolePermissions, error)
	BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.RolePermissions, error)
	DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error
}

type RolePermissionsRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewRolePermissionsRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) RolePermissionsRepositoryInterface {
	return &RolePermissionsRepository{
		*baseRepo,
		DB,
	}
}

func (u *RolePermissionsRepository) GetByConditions(conditions map[string]interface{}) ([]entity.RolePermissions, error) {
	var rolePermissions []entity.RolePermissions
	err := u.DB.Where(conditions).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (u *RolePermissionsRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.RolePermissions{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *RolePermissionsRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.RolePermissions, error) {
	var rolePermissions []entity.RolePermissions
	for _, attribute := range attributes {
		var temp entity.RolePermissions
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, temp)
	}

	err := u.DB.Create(&rolePermissions).Error
	return rolePermissions, err
}

func (u *RolePermissionsRepository) TakeByConditions(conditions map[string]interface{}) (entity.RolePermissions, error) {
	var rolePermission entity.RolePermissions
	err := u.DB.Where(conditions).Take(&rolePermission).Error
	return rolePermission, err
}

func (u *RolePermissionsRepository) BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.RolePermissions, error) {
	var rolePermissions []entity.RolePermissions
	for _, attribute := range attributes {
		var temp entity.RolePermissions
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, temp)
	}

	err := tx.Model(entity.RolePermissions{}).Create(&rolePermissions).Error
	return rolePermissions, err
}

func (u *RolePermissionsRepository) DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error {
	return tx.Where(conditions).Delete(&entity.RolePermissions{}).Error
}
