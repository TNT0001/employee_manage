package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
)

type PermissionList struct {
	Total int64
	List  []entity.Permissions
}

func NewPermissionListType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "permission_list",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"total": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(PermissionList).Total, nil
					},
				},
				"list": &graphql.Field{
					Type: &graphql.List{
						OfType: types["permission"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(PermissionList).List, nil
					},
				},
			}
		}),
	})
}
