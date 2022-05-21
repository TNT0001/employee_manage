package service

import (
	"errors"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

type PermissionsServicesInterface interface {
	CheckBatchPermissions(attributes []map[string]interface{}) error
}

type PermissionsService struct {
	DB                    *gorm.DB
	permissionsRepository repository.PermissionsRepositoryInterface
}

func NewPermissionsRepository(permissionsRepository repository.PermissionsRepositoryInterface, DB *gorm.DB) PermissionsServicesInterface {
	return &PermissionsService{
		DB,
		permissionsRepository,
	}
}

func (p *PermissionsService) CheckBatchPermissions(attributes []map[string]interface{}) error {
	seen := make(map[string]bool)
	nameList := make([]string, 0)
	msgErr := make(map[string]interface{})
	for _, attribute := range attributes {
		name, ok := attribute["name"].(string)
		if !ok {
			return wraperror.NewValidationError(map[string]interface{}{
				"name": utils.ErrorInputFail,
			}, errors.New("name field is not string"))
		}
		if seen[name] {
			msgErr[name] = utils.ErrorPermissionsNameDuplicate
		} else {
			seen[name] = true
			nameList = append(nameList, name)
		}
	}
	if len(msgErr) > 0 {
		return wraperror.NewValidationError(msgErr, errors.New("duplicate permission name"))
	}

	for _, name := range nameList {
		_, err := p.permissionsRepository.TakeByConditions(map[string]interface{}{
			"name": name,
		})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			msgErr[name] = utils.ErrorPermissionsNameDuplicate
		}
	}

	if len(msgErr) > 0 {
		return wraperror.NewValidationError(msgErr, errors.New("duplicate permission name"))
	}

	return nil
}
