package mutation

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/service"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
	"tungnt/emmployee_manage/pkg/share/utils"
)

func PostCreatePermissionsMutation(
	outputTypes map[string]*graphql.Object,
	permissionsService service.PermissionsServicesInterface,
	permissionsRepo repository.PermissionsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_user_login",
		Type: outputTypes["permission_list"],
		Args: map[string]*graphql.ArgumentConfig{
			"realm": {
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			input := utils.GetSubSliceMap(params.Source, "permissions")
			err := permissionsService.CheckBatchPermissions(input)
			if err != nil {
				return nil, err
			}

			permissions, err := permissionsRepo.BatchCreate(input)
			if err != nil {
				return nil, err
			}

			return output.PermissionList{
				Total: int64(len(permissions)),
				List:  permissions,
			}, nil
		},
	}
}
