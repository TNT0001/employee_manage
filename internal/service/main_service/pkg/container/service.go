package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	authorizations "tungnt/emmployee_manage/internal/modules/authorizations/pkg/container"
	users "tungnt/emmployee_manage/internal/modules/users/pkg/container"
)

type ServicesContainer struct {
	AuthorizationsContainer authorizations.AuthorizationServiceContainer
	UserContainer           users.UserServiceContainer
}

func NewServicesContainer(logger *logrus.Logger, DB *gorm.DB, container RepositoryContainer) ServicesContainer {
	authorizationsContainer := authorizations.NewAuthorizationServiceContainer(
		logger,
		DB,
		container.AuthorizationsContainer,
	)

	userContainer := users.NewUserServiceContainer(
		container.JobsContainer,
		container.AuthorizationsContainer,
		logger,
		DB,
	)

	return ServicesContainer{
		authorizationsContainer,
		userContainer,
	}
}
