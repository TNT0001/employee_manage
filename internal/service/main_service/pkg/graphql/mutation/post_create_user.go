package mutation

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v11"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	authorRepo "tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	jobsRepository "tungnt/emmployee_manage/internal/modules/jobs/pkg/repository"
	organizationsRepository "tungnt/emmployee_manage/internal/modules/organizations/pkg/repository"
	usersEntity "tungnt/emmployee_manage/internal/modules/users/pkg/entity"
	usersRepository "tungnt/emmployee_manage/internal/modules/users/pkg/repository"
	"tungnt/emmployee_manage/internal/modules/users/pkg/service"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func PostCreateUserMutation(
	outputTypes map[string]*graphql.Object,
	db *gorm.DB,
	client *keycloak.Client,
	teamRepo organizationsRepository.TeamsRepositoryInterface,
	usersService service.UsersServicesInterface,
	userRepo usersRepository.UsersRepositoryInterface,
	userRoleRepo authorRepo.UserRolesRepositoryInterface,
	jobsRepo jobsRepository.JobsRepositoryInterface,
	personalInfosRepo usersRepository.PersonalInfosRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_create_user",
		Type: outputTypes["user"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			input := utils.GetSubMap(params.Source, "user")
			// create keycloak account
			keycloakInput := utils.GetSubMap(input, "keycloak")
			keycloakToken, err := client.LoginAdmin(
				context.Background(),
				utils.AdminUserName,
				utils.AdminPassword,
				utils.MasterRealm,
			)
			if err != nil {
				return nil, err
			}
			users, err := client.GetUsers(
				context.Background(),
				keycloakToken.AccessToken,
				utils.EmployeeManageRealm,
				gocloak.GetUsersParams{
					Email: utils.GetStringPointer(keycloakInput["email"].(string)),
					Exact: utils.GetBoolPointer(true),
				},
			)
			if err != nil {
				return nil, err
			}
			var keycloakID string
			if len(users) == 0 {
				userName := keycloakInput["user_name"].(string)
				keycloakID, err = client.CreateUser(
					context.Background(),
					keycloakToken.AccessToken,
					utils.EmployeeManageRealm,
					gocloak.User{
						Email:         utils.GetStringPointer(keycloakInput["email"].(string)),
						Username:      &userName,
						Enabled:       utils.GetBoolPointer(true),
						EmailVerified: utils.GetBoolPointer(true),
					},
				)
				if err != nil {
					return nil, err
				}

				err = client.SetPassword(
					context.Background(),
					keycloakToken.AccessToken,
					keycloakID,
					utils.EmployeeManageRealm,
					keycloakInput["password"].(string),
					false,
				)
				if err != nil {
					return nil, err
				}
			} else {
				keycloakID = utils.DerefString(users[0].ID)
			}
			// check create user
			// check user with keycloak id exist
			_, err = userRepo.TakeByConditions(false, map[string]interface{}{
				"keycloak_user_id": keycloakID,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if err == nil {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"keycloak_user_id": utils.ErrorInputFail,
				}, errors.New(fmt.Sprintf("users with keycloak_user_id: %s, already exist", keycloakID)))
			}
			userAttr := utils.GetOnlyScalar(input)
			userAttr["keycloak_user_id"] = keycloakID
			teamID := utils.GetSubInteger(userAttr, "team_id")
			if teamID != nil {
				team, err := teamRepo.TakeByConditions(false, map[string]interface{}{
					"id": teamID,
				})
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, wraperror.NewValidationError(map[string]interface{}{
						"team_id": utils.ErrorInputFail,
					}, errors.New("team_id input fail"))
				}
				// add user to group
				groupName := utils.OrganizationPrefixName + "_" + team.CountryCode
				groupList, err := client.GetGroups(
					context.Background(),
					keycloakToken.AccessToken,
					utils.EmployeeManageRealm,
					gocloak.GetGroupsParams{
						Search: &groupName,
					},
				)
				if err != nil {
					return nil, err
				}

				var groupID string
				if len(groupList) == 0 {
					groupID, err = client.CreateGroup(
						context.Background(),
						keycloakToken.AccessToken,
						utils.EmployeeManageRealm,
						gocloak.Group{
							Name: &groupName,
						},
					)
					if err != nil {
						return nil, err
					}
				} else {
					groupID = utils.DerefString(groupList[0].ID)
				}

				err = client.AddUserToGroup(
					context.Background(),
					keycloakToken.AccessToken,
					utils.EmployeeManageRealm,
					keycloakID,
					groupID,
				)
				if err != nil {
					return nil, err
				}
			} else {
				delete(userAttr, "team_id")
			}

			// check create jobs
			jobsInput := utils.GetSubSliceMap(input, "jobs")
			var jobsAttr []map[string]interface{}
			if len(jobsInput) > 0 {
				jobsAttr, err = usersService.CheckJobAssign(jobsInput, "user")
				if err != nil {
					return nil, err
				}
			}

			// check role
			userRolesInput := utils.GetSubSliceMap(input, "user_roles")
			var userRoleAttr []map[string]interface{}
			if len(userRolesInput) > 0 {
				userRoleAttr, err = usersService.CheckUserRoles(userRolesInput, "user")
				if err != nil {
					return nil, err
				}
			}

			var user usersEntity.Users
			err = utils.Transaction(db, func(tx *gorm.DB) error {
				// create user
				user, err = userRepo.CreateWithTx(tx, userAttr)
				if err != nil {
					return err
				}

				// create user_roles
				if len(userRoleAttr) > 0 {
					for _, userRoleMap := range userRoleAttr {
						userRoleMap["user_id"] = user.ID
					}
					_, err = userRoleRepo.BatchCreateWithTx(tx, userRoleAttr)
					if err != nil {
						return err
					}
				}
				// create job assign
				if len(jobsAttr) > 0 {
					for _, jobMap := range jobsAttr {
						jobMap["user_id"] = user.ID
					}
					_, err := jobsRepo.BatchCreateWithTx(tx, jobsAttr)
					if err != nil {
						return err
					}
				}

				personalInfoAttr := utils.GetOnlyScalar(input, "personal_info")
				if len(personalInfoAttr) > 0 {
					personalInfoAttr["user_id"] = user.ID
					_, err := personalInfosRepo.CreateWithTx(tx, personalInfoAttr)
					if err != nil {
						return err
					}
				}
				return nil
			})
			return user, err
		},
	}
}
