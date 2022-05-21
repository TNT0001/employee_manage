package query

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
)

func GetRolesQuery(
	outputTypes map[string]*graphql.Object,
	rolesRepository repository.RolesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "get_roles",
		Type: outputTypes["role_list"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			roles, err := rolesRepository.GetByConditionsWithPreload(nil)
			if err != nil {
				return nil, err
			}
			total, err := rolesRepository.CountByConditions(nil)
			if err != nil {
				return nil, err
			}

			return output.RoleList{
				Total: total,
				List:  roles,
			}, nil
		},
	}
}
