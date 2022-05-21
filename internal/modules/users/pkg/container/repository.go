package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/users/pkg/repository"
	baseRepo "tungnt/emmployee_manage/pkg/share/repository"
)

type UserContainer struct {
	UsersRepository         repository.UsersRepositoryInterface
	PersonalInfosRepository repository.PersonalInfosRepositoryInterface
}

func NewUserContainer(logger *logrus.Logger, DB *gorm.DB) UserContainer {
	baseRepository := baseRepo.NewBaseRepository(logger)

	userRepository := repository.NewUserRepository(baseRepository, DB)
	personalInfosRepository := repository.NewPersonalInfosRepository(baseRepository, DB)

	return UserContainer{
		userRepository,
		personalInfosRepository,
	}
}
