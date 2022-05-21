package output

import (
	"github.com/Nerzal/gocloak/v11"
	"github.com/graphql-go/graphql"
)

func NewUserTokenType(types map[string]*graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "user_token",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"access_token": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).AccessToken, nil
					},
				},
				"id_token": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).IDToken, nil
					},
				},
				"expires_in": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).ExpiresIn, nil
					},
				},
				"refresh_expires_in": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).RefreshExpiresIn, nil
					},
				},
				"token_type": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).TokenType, nil
					},
				},
				"not_before_policy": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).NotBeforePolicy, nil
					},
				},
				"session_state": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).SessionState, nil
					},
				},
				"scope": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(*gocloak.JWT).Scope, nil
					},
				},
			}
		}),
	})
}
