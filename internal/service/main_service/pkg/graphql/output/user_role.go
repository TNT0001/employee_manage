package output

import (
	"errors"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	authorizationsRepository "tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
)

func NewUserRoleType(
	types map[string]*graphql.Object,
	roleRepository authorizationsRepository.RolesRepositoryInterface,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "user_role",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.UserRoles).ID, nil
					},
				},
				"user_id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.UserRoles).UserID, nil
					},
				},
				"role": &graphql.Field{
					Type: types["role"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						role, err := roleRepository.TakeByConditionsWithPreload(false, map[string]interface{}{
							"id": params.Source.(entity.UserRoles).RoleID,
						})
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return nil, nil
						}

						return role, err
					},
				},
			}
		}),
	})
}
