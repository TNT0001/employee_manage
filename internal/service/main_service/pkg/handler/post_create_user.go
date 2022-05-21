package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"tungnt/emmployee_manage/pkg/dto"
)

func (h *Handler) PostCreateUser(c *gin.Context) {
	input, err := h.GetInputAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	errResult, err := h.Validator.Validate(PostCreateUserschema, input)
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
				post_create_user {
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
