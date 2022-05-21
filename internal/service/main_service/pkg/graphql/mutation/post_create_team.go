package mutation

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v11"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/repository"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func PostCreateTeamMutation(
	outputTypes map[string]*graphql.Object,
	client *keycloak.Client,
	teamsRepo repository.TeamsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_user_login",
		Type: outputTypes["team"],
		Args: map[string]*graphql.ArgumentConfig{
			"realm": {
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			input := utils.GetSubMap(params.Source, "team")
			countryCode, existCountryCode := input["country_code"].(string)
			if !existCountryCode {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"country_code": utils.ErrorInputFail,
				}, errors.New("country_code input fail"))
			}
			teamName, existTeamName := input["team_name"].(string)
			if !existTeamName {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"team_name": utils.ErrorInputFail,
				}, errors.New("team_name input fail"))
			}
			divisionName, existDividsion := input["division_name"].(string)
			if !existDividsion {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"team_name": utils.ErrorInputFail,
				}, errors.New("team_name input fail"))
			}

			conditions := map[string]interface{}{
				"country_code": countryCode,
				"team_name":    teamName,
			}
			if existDividsion {
				conditions["division_name"] = divisionName
			} else {
				conditions["division_name"] = nil
			}

			_, err := teamsRepo.TakeByConditions(false, conditions)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if err == nil {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"country_code":  utils.ErrorTeamDuplicate,
					"team_name":     utils.ErrorTeamDuplicate,
					"division_name": utils.ErrorTeamDuplicate,
				}, errors.New("teams is duplicate"))
			}

			adminToken, err := client.LoginAdmin(
				context.Background(),
				utils.AdminUserName,
				utils.AdminPassword,
				utils.MasterRealm,
			)
			if err != nil {
				return nil, err
			}

			groupName := utils.OrganizationPrefixName + "_" + countryCode
			groupList, err := client.GetGroups(
				context.Background(),
				adminToken.AccessToken,
				utils.EmployeeManageRealm,
				gocloak.GetGroupsParams{
					Search: &groupName,
				},
			)
			if err != nil {
				return nil, err
			}
			if len(groupList) == 0 {
				_, err := client.CreateGroup(context.Background(), adminToken.AccessToken, utils.EmployeeManageRealm, gocloak.Group{
					Name: &groupName,
				})
				if err != nil {
					return nil, err
				}
			}

			return teamsRepo.Create(input)
		},
	}
}
