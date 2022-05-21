package repository

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"tungnt/emmployee_manage/internal/modules/jobs/pkg/entity"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
)

type JobsRepositoryInterface interface {
	GetByConditions(isFromMaster bool, conditions map[string]interface{}) ([]entity.Jobs, error)
	CountByConditions(conditions map[string]interface{}) (int64, error)
	TakeByConditions(conditions map[string]interface{}) (entity.Jobs, error)
	Create(input map[string]interface{}) (entity.Jobs, error)
	PluckIDByConditions(conditions map[string]interface{}) ([]int, error)
	BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.Jobs, error)
}

type JobsRepository struct {
	BaseRepository baseRepo.BaseRepository
	DB             *gorm.DB
}

func NewJobsRepository(baseRepo *baseRepo.BaseRepository, DB *gorm.DB) JobsRepositoryInterface {
	return &JobsRepository{
		*baseRepo,
		DB,
	}
}

func (u *JobsRepository) GetByConditions(isFromMaster bool, conditions map[string]interface{}) ([]entity.Jobs, error) {
	var jobs []entity.Jobs
	queryBuidler := u.DB
	if isFromMaster {
		queryBuidler = queryBuidler.Clauses(dbresolver.Write)
	}
	err := queryBuidler.Where(conditions).Find(&jobs).Error
	return jobs, err
}

func (u *JobsRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var total int64
	err := u.DB.Model(entity.Jobs{}).Where(conditions).Count(&total).Error
	return total, err
}

func (u *JobsRepository) TakeByConditions(conditions map[string]interface{}) (entity.Jobs, error) {
	var job entity.Jobs
	err := u.DB.Where(conditions).Take(&job).Error
	return job, err
}

func (u *JobsRepository) Create(input map[string]interface{}) (entity.Jobs, error) {
	var job entity.Jobs
	err := utils.MapToStruct(input, &job)
	if err != nil {
		return entity.Jobs{}, err
	}
	err = u.DB.Create(&job).Error
	return job, err
}

func (u *JobsRepository) PluckIDByConditions(conditions map[string]interface{}) ([]int, error) {
	var jobIDList []int
	err := u.DB.Model(entity.Jobs{}).Where(conditions).Pluck("id", &jobIDList).Error
	return jobIDList, err
}

func (u *JobsRepository) BatchCreateWithTx(tx *gorm.DB, attributes []map[string]interface{}) ([]entity.Jobs, error) {
	var jobs []entity.Jobs
	for _, attribute := range attributes {
		var temp entity.Jobs
		err := utils.MapToStruct(attribute, &temp)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, temp)
	}

	err := tx.Create(&jobs).Error
	return jobs, err
}
