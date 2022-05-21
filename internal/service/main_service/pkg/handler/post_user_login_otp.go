package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func (h *Handler) PostUserLoginOTP(c *gin.Context) {
	realm := c.Param("realm")
	if realm == "" {
		h.SetBadErrorResponse(c, wraperror.NewValidationError(map[string]interface{}{
			"realm": utils.ErrorInputFail,
		}, errors.New("realm params invalid")))
		return
	}

	input, err := h.GetInputAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	errResult, err := h.Validator.Validate(PostUserLoginOTPSchema, input)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}
	if errResult != nil {
		h.SetJSONValidationErrorResponse(c, errResult)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:     h.Graphql,
		RootObject: input,
		VariableValues: map[string]interface{}{
			"realm": realm,
		},
		RequestString: `
			mutation ($realm: String!) {
				post_user_login_otp (realm: $realm) {
					access_token
					id_token
					expires_in
					refresh_expires_in
					token_type
					not_before_policy
					session_state
					scope
				}
			}
		`,
	})

	if result.HasErrors() {
		h.SetGenericErrorResponse(c, result.Errors[0])
		return
	}

	c.JSON(http.StatusOK, dto.BaseSuccessResponse{
		Data: result.Data,
	})

	return
}
