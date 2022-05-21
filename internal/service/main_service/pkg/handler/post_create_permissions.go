package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
)

func (h *Handler) PostCreatePermissions(c *gin.Context) {
	input, err := h.GetInputAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	errResult, err := h.Validator.Validate(PostCreatePermissionschema, input)
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
				post_create_permission {
					list {
						id
						name
					}
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
