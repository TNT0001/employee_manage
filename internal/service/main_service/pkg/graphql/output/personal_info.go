package output

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/users/pkg/entity"
)

func NewPersonalInfoType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "personal_info",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).ID, nil
					},
				},
				"user_id": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).UserID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).Name, nil
					},
				},
				"sur_name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).SurName, nil
					},
				},
				"full_name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).FullName, nil
					},
				},
				"email": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).Email, nil
					},
				},
				"facebook": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).FaceBook, nil
					},
				},
				"linken": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).Linken, nil
					},
				},
				"phone_number": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).PhoneNumber, nil
					},
				},
				"address": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).Address, nil
					},
				},
				"on_probationary_period": &graphql.Field{
					Type: graphql.Boolean,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).OnProbationaryPeriod, nil
					},
				},
				"start_probationary_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).StartProbationaryDate, nil
					},
				},
				"end_probationary_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.PersonalInfos).EndProbationaryDate, nil
					},
				},
			}
		}),
	})
}
