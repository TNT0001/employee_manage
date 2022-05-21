package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/users/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type PersonalInfosRepositoryInterface interface {
	GetUsers(conditions map[string]interface{}) ([]entity.PersonalInfos, error)
	CountUsers(conditions map[string]interface{}) (int64, error)
	TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.PersonalInfos, error)
	BatchCreate(attributes []map[string]interface{}) ([]entity.PersonalInfos, error)
	CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.PersonalInfos, error)
	DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, role entity.PersonalInfos, attributes map[string]interface{}) (entity.PersonalInfos, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
}

type PersonalInfosRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewPersonalInfosRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) PersonalInfosRepositoryInterface {
	return &PersonalInfosRepository{
		*baseRepo,
		DB,
	}
}

func (u *PersonalInfosRepository) GetUsers(conditions map[string]interface{}) ([]entity.PersonalInfos, error) {
	var personalInfos []entity.PersonalInfos
	err := u.DB.Where(conditions).Find(&personalInfos).Error
	return personalInfos, err
}

func (u *PersonalInfosRepository) CountUsers(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.PersonalInfos{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *PersonalInfosRepository) BatchCreate(attributes []map[string]interface{}) ([]entity.PersonalInfos, error) {
	var personalInfos []entity.PersonalInfos
	for _, attribute := range attributes {
		var temp entity.PersonalInfos
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		personalInfos = append(personalInfos, temp)
	}

	err := u.DB.Create(&personalInfos).Error
	return personalInfos, err
}

func (u *PersonalInfosRepository) TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.PersonalInfos, error) {
	var personalInfo entity.PersonalInfos
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Where(conditions).Take(&personalInfo).Error
	return personalInfo, err
}

func (u *PersonalInfosRepository) CreateWithTx(tx *gorm.DB, attributes map[string]interface{}) (entity.PersonalInfos, error) {
	var personalInfo entity.PersonalInfos
	err := utils.MapToStruct(attributes, &personalInfo)
	if err != nil {
		return entity.PersonalInfos{}, err
	}

	err = tx.Model(entity.PersonalInfos{}).Create(&personalInfo).Error
	return personalInfo, err
}

func (u *PersonalInfosRepository) DeleteByConditionsWithTx(tx *gorm.DB, conditions map[string]interface{}) error {
	return tx.Where(conditions).Delete(&entity.PersonalInfos{}).Error
}

func (u *PersonalInfosRepository) UpdateWithTx(tx *gorm.DB, personalInfo entity.PersonalInfos, attributes map[string]interface{}) (entity.PersonalInfos, error) {
	err := utils.MapToStruct(attributes, &personalInfo)
	if err != nil {
		return entity.PersonalInfos{}, err
	}

	err = tx.Omit(clause.Associations).Model(personalInfo).Updates(attributes).Error
	return personalInfo, err
}

func (u *PersonalInfosRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var personalInfoIDList []int
	err := u.DB.Model(entity.PersonalInfos{}).Where(conditions).Pluck("id", &personalInfoIDList).Error
	return personalInfoIDList, err
}
