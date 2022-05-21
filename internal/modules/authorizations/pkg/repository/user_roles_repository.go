package repository

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type UserRolesRepositoryInterface interface {
	GetByConditions(isFromMaster bool, conditions map[string]interface{}) ([]entity.UserRoles, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.UserRoles, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.UserRoles, error)
	GetByConditionsWithPreload(conditions map[string]interface{}) ([]entity.UserRoles, error)
	CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.UserRoles, error)
	TakePreloadByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) (entity.UserRoles, error)
	DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error
	BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.UserRoles, error)
}

type UserRolesRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewUserRolesRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) UserRolesRepositoryInterface {
	return &UserRolesRepository{
		*baseRepo,
		DB,
	}
}

func (u *UserRolesRepository) GetByConditions(isFromMaster bool, conditions map[string]interface{}) ([]entity.UserRoles, error) {
	var roles []entity.UserRoles
	queryBuilder := u.DB
	if isFromMaster {
		queryBuilder = queryBuilder.Clauses(dbresolver.Write)
	}
	err := queryBuilder.Where(conditions).Find(&roles).Error
	return roles, err
}

func (u *UserRolesRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.UserRoles{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *UserRolesRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.UserRoles, error) {
	var userRoles []entity.UserRoles
	for _, attribute := range attributes {
		var temp entity.UserRoles
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, temp)
	}

	err := u.DB.Create(&userRoles).Error
	return userRoles, err
}

func (u *UserRolesRepository) TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.UserRoles, error) {
	var userRole entity.UserRoles
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Where(conditions).Take(&userRole).Error
	return userRole, err
}

func (u *UserRolesRepository) GetByConditionsWithPreload(conditions map[string]interface{}) ([]entity.UserRoles, error) {
	var userRoles []entity.UserRoles
	err := u.DB.Preload("RolesPermissions.Permission").Where(conditions).Find(&userRoles).Error
	return userRoles, err
}

func (u *UserRolesRepository) CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.UserRoles, error) {
	var userRole entity.UserRoles
	err := utils.MapToStruct(attributes, &userRole)
	if err != nil {
		return entity.UserRoles{}, err
	}

	err = tx.Model(entity.UserRoles{}).Create(&userRole).Error
	return userRole, err
}

func (u *UserRolesRepository) TakePreloadByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) (entity.UserRoles, error) {
	var userRole entity.UserRoles
	err := tx.Model(entity.UserRoles{}).Where(conditions).Take(&userRole).Error
	return userRole, err
}

func (u *UserRolesRepository) DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error {
	return tx.Where(conditions).Delete(&entity.UserRoles{}).Error
}

func (u *UserRolesRepository) BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.UserRoles, error) {
	var userRoles []entity.UserRoles
	for _, attribute := range attributes {
		var temp entity.UserRoles
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, temp)
	}

	err := tx.Create(&userRoles).Error
	return userRoles, err
}
