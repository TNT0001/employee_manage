package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tungnt/emmployee_manage/pkg/dto"
	"tungnt/emmployee_manage/pkg/share/jwt"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/utils"
)

func NewAuthMiddleware(client *keycloak.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, exists := jwt.GetTokenByHeader(c)
		if !exists {
			data := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": utils.ErrorInputFail,
					},
				},
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, data)
			return
		}

		realmName, err := jwt.GetRealmNameByToken(token)
		if err != nil {
			data := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": utils.ErrorInputFail,
					},
				},
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, data)
			return
		}

		publicKey, err := client.GetCerts(c, realmName)
		if err != nil || publicKey == nil {
			data := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": utils.ErrorInputFail,
					},
				},
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, data)
			return
		}

		var rsaKey string
		for _, key := range *publicKey.Keys {
			if utils.DerefString(key.Kty) == "RSA" {
				rsaKey = (*key.X5c)[0]
				break
			}
		}

		if rsaKey == "" {
			data := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": utils.ErrorInputFail,
					},
				},
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, data)
			return
		}
		mapClaim, err := jwt.Verify(token, rsaKey)
		if err != nil {
			data := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": utils.ErrorInputFail,
					},
				},
			}

			if strings.Contains(err.Error(), utils.ErrorTokenExpired) {
				c.AbortWithStatusJSON(http.StatusRequestTimeout, data)
				return
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", mapClaim)
		c.Next()
	}
}
