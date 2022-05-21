package output

import (
	"errors"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	authorizationsRepository "tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	jobsRepository "tungnt/emmployee_manage/internal/modules/jobs/pkg/repository"
	organizationsRepository "tungnt/emmployee_manage/internal/modules/organizations/pkg/repository"
	"tungnt/emmployee_manage/internal/modules/users/pkg/entity"
	usersRepository "tungnt/emmployee_manage/internal/modules/users/pkg/repository"
)

func NewUserType(
	types map[string]*graphql.Object,
	teamRepo organizationsRepository.TeamsRepositoryInterface,
	personalInfosRepository usersRepository.PersonalInfosRepositoryInterface,
	jobsRepository jobsRepository.JobsRepositoryInterface,
	userRolesRepository authorizationsRepository.UserRolesRepositoryInterface,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "user",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Users).ID, nil
					},
				},
				"keycloak_user_id": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Users).KeycloakUserID, nil
					},
				},
				"user_name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Users).UserName, nil
					},
				},
				"join_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Users).JoinDate, nil
					},
				},
				"team": &graphql.Field{
					Type: types["team"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						teamID := params.Source.(entity.Users).TeamID
						if teamID != nil {
							team, err := teamRepo.TakeByConditions(true, map[string]interface{}{
								"id": teamID,
							})
							if errors.Is(err, gorm.ErrRecordNotFound) {
								return nil, nil
							}
							return team, err
						}
						return nil, nil
					},
				},
				"personal_info": &graphql.Field{
					Type: types["personal_info"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						personalInfo, err := personalInfosRepository.TakeByConditions(true, map[string]interface{}{
							"user_id": params.Source.(entity.Users).ID,
						})
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return nil, nil
						}
						return personalInfo, err
					},
				},
				"job_assigns": &graphql.Field{
					Type: &graphql.List{
						OfType: types["job"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						jobs, err := jobsRepository.GetByConditions(true, map[string]interface{}{
							"user_id": params.Source.(entity.Users).ID,
						})
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return nil, nil
						}
						return jobs, err
					},
				},
				"user_roles": &graphql.Field{
					Type: &graphql.List{
						OfType: types["user_role"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						usesRoles, err := userRolesRepository.GetByConditions(true, map[string]interface{}{
							"user_id": params.Source.(entity.Users).ID,
						})
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return nil, nil
						}
						return usesRoles, err
					},
				},
			}
		}),
	})
}
