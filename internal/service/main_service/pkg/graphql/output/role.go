package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
)

func NewRoleType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "role",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Roles).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Roles).Name, nil
					},
				},
				"permissions": &graphql.Field{
					Type: types["permission_list"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						var permissions []entity.Permissions
						for _, rolePermission := range params.Source.(entity.Roles).RolesPermissions {
							permissions = append(permissions, rolePermission.Permission)
						}

						return PermissionList{
							Total: int64(len(permissions)),
							List:  permissions,
						}, nil
					},
				},
			}
		}),
	})
}
