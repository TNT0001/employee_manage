package mutation

import (
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/handler"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func PostUserLoginMutation(
	outputTypes map[string]*graphql.Object,
	client *keycloak.Client,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_user_login",
		Type: outputTypes["user_token"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			realm := params.Args["realm"].(string)
			input := utils.GetSubMap(params.Source, "user_credential")
			userName, ok := input["user_name"].(string)
			if !ok {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"user_name": utils.ErrorInputFail,
				}, errors.New("user_name input fail"))
			}
			password, ok := input["password"].(string)
			if !ok {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"password": utils.ErrorInputFail,
				}, errors.New("password input fail"))
			}

			keycloakToken, err := client.Login(
				context.Background(),
				handler.ClientID,
				handler.ClientSecret,
				realm,
				userName,
				password,
			)
			if err != nil {
				return nil, err
			}

			return keycloakToken, nil
		},
	}
}
