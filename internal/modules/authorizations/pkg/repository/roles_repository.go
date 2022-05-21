package repository

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type RolesRepositoryInterface interface {
	GetByConditions(conditions map[string]interface{}) ([]entity.Roles, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Roles, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.Roles, error)
	GetByConditionsWithPreload(conditions map[string]interface{}) ([]entity.Roles, error)
	CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.Roles, error)
	TakePreloadByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) (entity.Roles, error)
	DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, role entity.Roles, attributes map[string]interface{}) (entity.Roles, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
	TakeByConditionsWithPreload(isFromMaster bool, conditions map[string]interface{}) (entity.Roles, error)
}

type RolesRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewRolesRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) RolesRepositoryInterface {
	return &RolesRepository{
		*baseRepo,
		DB,
	}
}

func (u *RolesRepository) GetByConditions(conditions map[string]interface{}) ([]entity.Roles, error) {
	var roles []entity.Roles
	err := u.DB.Where(conditions).Find(&roles).Error
	return roles, err
}

func (u *RolesRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.Roles{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *RolesRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.Roles, error) {
	var roles []entity.Roles
	for _, attribute := range attributes {
		var temp entity.Roles
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		roles = append(roles, temp)
	}

	err := u.DB.Create(&roles).Error
	return roles, err
}

func (u *RolesRepository) TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Roles, error) {
	var roles entity.Roles
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Where(conditions).Take(&roles).Error
	return roles, err
}

func (u *RolesRepository) GetByConditionsWithPreload(conditions map[string]interface{}) ([]entity.Roles, error) {
	var roles []entity.Roles
	err := u.DB.Preload("RolesPermissions.Permission").Where(conditions).Find(&roles).Error
	return roles, err
}

func (u *RolesRepository) CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.Roles, error) {
	var role entity.Roles
	err := utils.MapToStruct(attributes, &role)
	if err != nil {
		return entity.Roles{}, err
	}

	err = tx.Model(entity.Roles{}).Create(&role).Error
	return role, err
}

func (u *RolesRepository) TakePreloadByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) (entity.Roles, error) {
	var roles entity.Roles
	err := tx.Model(entity.Roles{}).Preload("RolesPermissions.Permission").Where(conditions).Take(&roles).Error
	return roles, err
}

func (u *RolesRepository) DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error {
	return tx.Where(conditions).Delete(&entity.Roles{}).Error
}

func (u *RolesRepository) UpdateWithTx(tx *gorm.DB, role entity.Roles, attributes map[string]interface{}) (entity.Roles, error) {
	err := utils.MapToStruct(attributes, &role)
	if err != nil {
		return entity.Roles{}, err
	}

	err = tx.Model(role).Updates(attributes).Error
	return role, err
}

func (u *RolesRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var roleIDList []int
	err := u.DB.Model(entity.Roles{}).Where(conditions).Pluck("id", &roleIDList).Error
	return roleIDList, err
}

func (u *RolesRepository) TakeByConditionsWithPreload(isFromMaster bool, conditions map[string]interface{}) (entity.Roles, error) {
	var roles entity.Roles
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Preload("RolesPermissions.Permission").Where(conditions).Take(&roles).Error
	return roles, err
}
