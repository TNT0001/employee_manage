package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	authorizations "tungnt/emmployee_manage/internal/modules/authorizations/pkg/container"
	jobs "tungnt/emmployee_manage/internal/modules/jobs/pkg/container"
	"tungnt/emmployee_manage/internal/modules/users/pkg/service"
)

type UserServiceContainer struct {
	UsersService service.UsersServicesInterface
}

func NewUserServiceContainer(
	jobsContainer jobs.JobsRepositoryContainer,
	authorizationsContainer authorizations.AuthorizationRepositoryContainer,
	logger *logrus.Logger,
	DB *gorm.DB,
) UserServiceContainer {
	userServiceRepository := service.NewUsersService(
		DB,
		jobsContainer.AssignTypesRepository,
		authorizationsContainer.RolesRepository,
	)

	return UserServiceContainer{userServiceRepository}
}
