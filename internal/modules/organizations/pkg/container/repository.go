package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/repository"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
)

type OrganizationsContainer struct {
	TeamsRepository repository.TeamsRepositoryInterface
}

func NewOrganizationsContainer(logger *logrus.Logger, DB *gorm.DB) OrganizationsContainer {
	baseRepository := baseRepo.NewBaseRepository(logger)

	teamsRepository := repository.NewTeamsRepository(baseRepository, DB)

	return OrganizationsContainer{teamsRepository}
}
