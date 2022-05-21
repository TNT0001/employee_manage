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

func PostUserLoginOTPMutation(
	outputTypes map[string]*graphql.Object,
	client *keycloak.Client,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_user_login_otp",
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
			otp, ok := input["otp"].(string)
			if !ok {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"otp": utils.ErrorInputFail,
				}, errors.New("otp input fail"))
			}

			keycloakToken, err := client.LoginOtp(
				context.Background(),
				handler.ClientID,
				handler.ClientSecret,
				realm,
				userName,
				password,
				otp,
			)
			if err != nil {
				return nil, err
			}

			return keycloakToken, nil
		},
	}
}
