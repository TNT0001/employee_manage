package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
)

func (h *Handler) GetRoles(c *gin.Context) {
	result := graphql.Do(graphql.Params{
		Schema: h.Graphql,
		RequestString: `
			query {
				get_roles {
					total
					list {
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
