package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/service"
)

type AuthorizationServiceContainer struct {
	PermissionsService service.PermissionsServicesInterface
}

func NewAuthorizationServiceContainer(logger *logrus.Logger, DB *gorm.DB, authRepoContainer AuthorizationRepositoryContainer) AuthorizationServiceContainer {
	permissionsService := service.NewPermissionsRepository(authRepoContainer.PermissionsRepository, DB)

	return AuthorizationServiceContainer{permissionsService}
}
