package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/jobs/pkg/repository"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
)

type JobsRepositoryContainer struct {
	AssignTypesRepository repository.AssignTypesRepositoryInterface
	JobsRepository        repository.JobsRepositoryInterface
}

func NewJobsRepositoryContainer(logger *logrus.Logger, DB *gorm.DB) JobsRepositoryContainer {
	baseRepository := baseRepo.NewBaseRepository(logger)

	assignTypesRepository := repository.NewAssignTypesRepository(baseRepository, DB)
	jobsRepository := repository.NewJobsRepository(baseRepository, DB)

	return JobsRepositoryContainer{
		assignTypesRepository,
		jobsRepository,
	}
}
