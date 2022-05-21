package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"

	//"tungnt/emmployee_manage/internal/modules/users/pkg/repository"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
)

type AuthorizationRepositoryContainer struct {
	PermissionsRepository     repository.PermissionsRepositoryInterface
	RolesRepository           repository.RolesRepositoryInterface
	RolePermissionsRepository repository.RolePermissionsRepositoryInterface
	UserRolesRepository       repository.UserRolesRepositoryInterface
}

func NewAuthorizationRepositoryContainer(logger *logrus.Logger, DB *gorm.DB) AuthorizationRepositoryContainer {
	baseRepository := baseRepo.NewBaseRepository(logger)

	permissionsRepository := repository.NewPermissionsRepository(baseRepository, DB)
	rolesRepository := repository.NewRolesRepository(baseRepository, DB)
	rolesPermissionsRepository := repository.NewRolePermissionsRepository(baseRepository, DB)
	userRolesRepository := repository.NewUserRolesRepository(baseRepository, DB)

	return AuthorizationRepositoryContainer{
		permissionsRepository,
		rolesRepository,
		rolesPermissionsRepository,
		userRolesRepository,
	}
}
