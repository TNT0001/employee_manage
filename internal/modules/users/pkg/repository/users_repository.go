package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/users/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type UsersRepositoryInterface interface {
	GetUsers(conditions map[string]interface{}) ([]entity.Users, error)
	CountUsers(conditions map[string]interface{}) (int64, error)
	TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Users, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.Users, error)
	CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.Users, error)
	DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, role entity.Users, attributes map[string]interface{}) (entity.Users, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
}

type UserRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewUserRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) UsersRepositoryInterface {
	return &UserRepository{
		*baseRepo,
		DB,
	}
}

func (u *UserRepository) GetUsers(conditions map[string]interface{}) ([]entity.Users, error) {
	var users []entity.Users
	err := u.DB.Where(conditions).Find(&users).Error
	return users, err
}

func (u *UserRepository) CountUsers(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.Users{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *UserRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.Users, error) {
	var users []entity.Users
	for _, attribute := range attributes {
		var temp entity.Users
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		users = append(users, temp)
	}

	err := u.DB.Create(&users).Error
	return users, err
}

func (u *UserRepository) TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Users, error) {
	var user entity.Users
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Where(conditions).Take(&user).Error
	return user, err
}

func (u *UserRepository) CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.Users, error) {
	var user entity.Users
	err := utils.MapToStruct(attributes, &user)
	if err != nil {
		return entity.Users{}, err
	}

	err = tx.Model(entity.Users{}).Create(&user).Error
	return user, err
}

func (u *UserRepository) DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error {
	return tx.Where(conditions).Delete(&entity.Users{}).Error
}

func (u *UserRepository) UpdateWithTx(tx *gorm.DB, user entity.Users, attributes map[string]interface{}) (entity.Users, error) {
	err := utils.MapToStruct(attributes, &user)
	if err != nil {
		return entity.Users{}, err
	}

	err = tx.Omit(clause.Associations).Model(user).Updates(attributes).Error
	return user, err
}

func (u *UserRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var userIDList []int
	err := u.DB.Model(entity.Users{}).Where(conditions).Pluck("id", &userIDList).Error
	return userIDList, err
}

func (u *UserRepository) TakeByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) (entity.Users, error) {
	var user entity.Users
	err := tx.Where(conditions).Take(&user).Error
	return user, err
}
