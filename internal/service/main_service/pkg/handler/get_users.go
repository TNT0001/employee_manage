package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
)

func (h *Handler) GetUser(c *gin.Context) {
	result := graphql.Do(graphql.Params{
		Schema: h.Graphql,
		RequestString: `
			query {
				get_users {
					total
					list {
						id
						keycloak_user_id
						user_name
						join_date
						team {
							team_name
							country_code
							division_name
							kind
						}
						personal_info {
							full_name
							email
							phone_number
							on_probationary_period
						}
						job_assigns {
							project_name
							user_id
							assign_type_id
						}
						user_roles {
							role {
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
