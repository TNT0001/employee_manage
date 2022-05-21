package query

import (
	"github.com/graphql-go/graphql"
	authRepo "tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
)

func GetPermissionsQuery(
	outputTypes map[string]*graphql.Object,
	permissionsRepository authRepo.PermissionsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "get_permissions",
		Type: outputTypes["permission_list"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			permissionList, err := permissionsRepository.GetByConditions(nil)
			if err != nil {
				return nil, err
			}
			total, err := permissionsRepository.CountByConditions(nil)
			if err != nil {
				return nil, err
			}

			return output.PermissionList{
				Total: total,
				List:  permissionList,
			}, nil
		},
	}
}
