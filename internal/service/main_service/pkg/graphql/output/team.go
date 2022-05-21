package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/entity"
)

func NewTeamType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "team",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Teams).ID, nil
					},
				},
				"team_name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Teams).TeamName, nil
					},
				},
				"country_code": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Teams).CountryCode, nil
					},
				},
				"division_name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Teams).DivisionName, nil
					},
				},
				"kind": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Teams).Kind, nil
					},
				},
			}
		}),
	})
}
