package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/users/pkg/entity"
)

type UsersList struct {
	Total int64
	List  []entity.Users
}

func NewUserListType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "user_list",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"total": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(UsersList).Total, nil
					},
				},
				"list": &graphql.Field{
					Type: &graphql.List{
						OfType: types["user"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(UsersList).List, nil
					},
				},
			}
		}),
	})
}
