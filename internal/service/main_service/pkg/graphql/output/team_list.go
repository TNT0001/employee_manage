package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/entity"
)

type TeamList struct {
	Total int64
	List  []entity.Teams
}

func NewTeamListType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "team_list",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"total": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(TeamList).Total, nil
					},
				},
				"list": &graphql.Field{
					Type: &graphql.List{
						OfType: types["team"],
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(TeamList).List, nil
					},
				},
			}
		}),
	})
}
