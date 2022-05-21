package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
)

type RoleList struct {
	Total int64
	List  []entity.Roles
}

func NewRoleListType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "role_list",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"total": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(RoleList).Total, nil
					},
				},
				"list": &graphql.Field{
					Type: &graphql.List{
						OfType: types["role"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(RoleList).List, nil
					},
				},
			}
		}),
	})
}
