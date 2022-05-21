package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"strconv"
	"tungnt/emmployee_manage/pkg/dto"
	"tungnt/emmployee_manage/pkg/share/utils"
)

func (h *Handler) PatchUpdateRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.SetBadErrorResponse(c, map[string]interface{}{
			"role_id": utils.ErrorInputFail,
		})
	}
	input, err := h.GetInputAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	errResult, err := h.Validator.Validate(PatchUpdateRoleschema, input)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}
	if errResult != nil {
		h.SetJSONValidationErrorResponse(c, errResult)
		return
	}
	result := graphql.Do(graphql.Params{
		Schema: h.Graphql,
		VariableValues: map[string]interface{}{
			"role_id": roleID,
		},
		RootObject: input,
		RequestString: `
			mutation ($role_id: Int!) {
				patch_update_role (role_id: $role_id) {
					id
					name
					permissions {
						total
						list {
							id
							name
						}
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
