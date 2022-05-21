package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"strconv"
	"tungnt/emmployee_manage/pkg/dto"
	"tungnt/emmployee_manage/pkg/share/utils"
)

func (h *Handler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		h.SetBadErrorResponse(c, map[string]interface{}{
			"role_id": utils.ErrorInputFail,
		})
	}

	result := graphql.Do(graphql.Params{
		Schema: h.Graphql,
		VariableValues: map[string]interface{}{
			"role_id": roleID,
		},
		RequestString: `
			mutation ($role_id: Int!){
				delete_role (role_id: $role_id)
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
