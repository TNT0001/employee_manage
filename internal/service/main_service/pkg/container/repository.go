package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	authorizations "tungnt/emmployee_manage/internal/modules/authorizations/pkg/container"
	jobs "tungnt/emmployee_manage/internal/modules/jobs/pkg/container"
	organizations "tungnt/emmployee_manage/internal/modules/organizations/pkg/container"
	users "tungnt/emmployee_manage/internal/modules/users/pkg/container"
)

type RepositoryContainer struct {
	UsersContainer          users.UserContainer
	OrganizationContainer   organizations.OrganizationsContainer
	AuthorizationsContainer authorizations.AuthorizationRepositoryContainer
	JobsContainer           jobs.JobsRepositoryContainer
}

func NewRepositoryContainer(logger *logrus.Logger, DB *gorm.DB) RepositoryContainer {
	usersContainer := users.NewUserContainer(logger, DB)
	organizationsContainer := organizations.NewOrganizationsContainer(logger, DB)
	authorizationsContainer := authorizations.NewAuthorizationRepositoryContainer(logger, DB)
	jobsContainer := jobs.NewJobsRepositoryContainer(logger, DB)

	return RepositoryContainer{
		usersContainer,
		organizationsContainer,
		authorizationsContainer,
		jobsContainer,
	}
}
