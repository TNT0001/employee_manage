package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	authorizationsRepository "tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	jobsRepository "tungnt/emmployee_manage/internal/modules/jobs/pkg/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

type UsersServicesInterface interface {
	CheckJobAssign(jobs []map[string]interface{}, msgFormat string) ([]map[string]interface{}, error)
	CheckUserRoles(userRoles []map[string]interface{}, msgFormat string) ([]map[string]interface{}, error)
}

type UsersService struct {
	db              *gorm.DB
	assignTypesRepo jobsRepository.AssignTypesRepositoryInterface
	rolesRepo       authorizationsRepository.RolesRepositoryInterface
}

func NewUsersService(
	db *gorm.DB,
	assignTypesRepository jobsRepository.AssignTypesRepositoryInterface,
	rolesRepository authorizationsRepository.RolesRepositoryInterface,
) UsersServicesInterface {
	return &UsersService{
		db,
		assignTypesRepository,
		rolesRepository,
	}
}

func (u *UsersService) CheckJobAssign(jobs []map[string]interface{}, msgFormat string) ([]map[string]interface{}, error) {
	var total int
	var assignTypeIDInputList []int
	seen := make(map[int]bool)
	msgErr1 := make(map[string]interface{})
	jobsAttr := make([]map[string]interface{}, 0)
	for i, job := range jobs {
		temp := make(map[string]interface{})
		assignTypeID := utils.GetSubInteger(job, "assign_type_id")
		if assignTypeID == nil {
			msgErr1[fmt.Sprintf(msgFormat, fmt.Sprintf(".job[%d].assign_type_id", i))] = utils.ErrorInputFail
		} else {
			if seen[*assignTypeID] {
				msgErr1[fmt.Sprintf(msgFormat, fmt.Sprintf(".job[%d].assign_type_id", i))] = utils.ErrorInputFail
			}
			assignTypeIDInputList = append(assignTypeIDInputList, *assignTypeID)
			temp["assign_type_id"] = assignTypeID
		}

		assignPercent := utils.GetSubInteger(job, "assign_percent")
		if assignPercent == nil {
			msgErr1[fmt.Sprintf(msgFormat, fmt.Sprintf(".job[%d].assign_percent", i))] = utils.ErrorInputFail
		} else {
			if *assignPercent < utils.MinimumAssignPercent {
				msgErr1[fmt.Sprintf(msgFormat, fmt.Sprintf(".job[%d].assign_percent", i))] = utils.ErrorInputFail
			}
			temp["assign_percent"] = *assignPercent
			total += *assignPercent
		}
		temp["project_name"] = job["project_name"].(string)
		jobsAttr = append(jobsAttr, temp)
	}
	if total > 100 {
		msgErr1[fmt.Sprintf(msgFormat, ".job.assign_percent")] = utils.ErrorTotalAssignPercentLagerThan100
	}

	var finalMsg map[string]interface{}
	if len(assignTypeIDInputList) > 0 {
		assignIDList, err := u.assignTypesRepo.PluckIDByConditions(map[string]interface{}{
			"id": assignTypeIDInputList,
		})
		if err != nil {
			return jobsAttr, err
		}

		msgErr2, ok := utils.CompareIntSlice(assignTypeIDInputList, assignIDList, fmt.Sprintf(msgFormat, ".job[%d].assign_type_id"))
		if !ok {
			finalMsg = utils.MergeMap(msgErr1, msgErr2)
		} else {
			finalMsg = msgErr1
		}
	}

	if len(finalMsg) > 0 {
		return jobsAttr, wraperror.NewValidationError(finalMsg, errors.New("job assign errors"))
	}

	return jobsAttr, nil
}

func (u *UsersService) CheckUserRoles(userRoles []map[string]interface{}, msgFormat string) ([]map[string]interface{}, error) {
	var roleIDInputList []int
	msgErr1 := make(map[string]interface{})
	userRolesAttr := make([]map[string]interface{}, 0)
	seen := make(map[int]bool)
	for i, userRole := range userRoles {
		roleID := utils.GetSubInteger(userRole, "role_id")
		if roleID == nil {
			msgErr1[fmt.Sprintf(msgFormat, fmt.Sprintf(".user_roles[%d].role_id", i))] = utils.ErrorInputFail
		} else {
			roleIDInputList = append(roleIDInputList, utils.DerefInt(roleID))
			if !seen[*roleID] {
				userRolesAttr = append(userRolesAttr, map[string]interface{}{
					"role_id": roleID,
				})
				seen[*roleID] = true
			}
		}
	}

	roleIDList, err := u.rolesRepo.PluckIDByConditions(map[string]interface{}{
		"id": roleIDInputList,
	})
	if err != nil {
		return nil, err
	}

	msgErr2, ok := utils.CompareIntSlice(roleIDList, roleIDInputList, fmt.Sprintf(msgFormat, ".user_roles[%d].role_id"))
	var finalMsg map[string]interface{}
	if !ok {
		finalMsg = utils.MergeMap(msgErr1, msgErr2)
	}

	if len(finalMsg) > 0 {
		return nil, wraperror.NewValidationError(finalMsg, errors.New("user_roles errors"))
	}

	return userRolesAttr, nil
}
