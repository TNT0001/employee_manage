package mutation

import (
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func PostAdminUserLoginMutation(
	outputTypes map[string]*graphql.Object,
	client *keycloak.Client,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_admin_user_login",
		Type: outputTypes["user_token"],
		Args: map[string]*graphql.ArgumentConfig{
			"realm": {
				Type: graphql.String,
			},
		},
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

			keycloakToken, err := client.LoginAdmin(
				context.Background(),
				userName,
				password,
				realm,
			)
			if err != nil {
				return nil, err
			}

			return keycloakToken, nil
		},
	}
}
