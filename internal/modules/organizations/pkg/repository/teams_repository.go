package repository

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type TeamsRepositoryInterface interface {
	GetByConditions(conditions map[string]interface{}) ([]entity.Teams, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Teams, error)
	Create(input map[string]interface{}) (entity.Teams, error)
}

type TeamsRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewTeamsRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) TeamsRepositoryInterface {
	return &TeamsRepository{
		*baseRepo,
		DB,
	}
}

func (u *TeamsRepository) GetByConditions(conditions map[string]interface{}) ([]entity.Teams, error) {
	var teams []entity.Teams
	err := u.DB.Where(conditions).Find(&teams).Error
	return teams, err
}

func (u *TeamsRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.Teams{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *TeamsRepository) TakeByConditions(isFromMaster bool, conditions map[string]interface{}) (entity.Teams, error) {
	var team entity.Teams
	queryBuilder := u.DB
	if isFromMaster {
		queryBuilder = queryBuilder.Clauses(dbresolver.Write)
	}
	err := queryBuilder.Where(conditions).Take(&team).Error
	return team, err
}

func (u *TeamsRepository) Create(input map[string]interface{}) (entity.Teams, error) {
	var team entity.Teams
	err := utils.MapToStruct(input, &team)
	if err != nil {
		return entity.Teams{}, err
	}
	err = u.DB.Create(&team).Error
	return team, err
}
