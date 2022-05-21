package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
)

func (h *Handler) PostCreateTeam(c *gin.Context) {
	input, err := h.GetInputAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	errResult, err := h.Validator.Validate(PostCreateTeamschema, input)
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
		RequestString: `
			mutation {
				post_create_team {
					id
					team_name
					country_code
					division_name
					kind
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
