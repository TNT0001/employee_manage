package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
	"time"
)

// Decoding JWT to get payload, not verifying JWT
func Decode(JWTToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err.Error() == jwt.ErrInvalidKeyType.Error() {
		return token, nil
	}
	return nil, err
}

// Generate HS256 JWT token
func GenerateHS256JWT(secret string, keyID string, payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, err
}

// Verify JWT token
func Verify(token string, keyString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	publicKey := "-----BEGIN PUBLIC KEY-----\n" + keyString + "\n" + "-----END PUBLIC KEY-----"
	keyByte, keyErr := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if keyErr != nil {
		if os.Getenv("STAGE") == "DEV" && (errors.Is(keyErr, jwt.ErrNotRSAPublicKey) || errors.Is(keyErr, jwt.ErrKeyMustBePEMEncoded)) {
			_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return nil, nil
			})
			return claims, err
		} else {
			return claims, keyErr
		}
	}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return keyByte, nil
	})

	return claims, err
}

func GetUserIDByToken(token string) (string, error) {
	jwtToken, err := Decode(token)
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", nil
	}

	return userID, nil
}

func GetRealmNameByToken(token string) (string, error) {
	jwtToken, err := Decode(token)
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	urlPath, ok := claims["iss"].(string)
	if !ok {
		return "", nil
	}

	stringSplit := strings.Split(urlPath, "/")
	return stringSplit[len(stringSplit)-1], nil
}

func GetEmailByToken(token string) (string, error) {
	jwtToken, err := Decode(token)
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", nil
	}

	return email, nil
}

func VerifyExpiredByToken(token string) (bool, error) {
	jwtToken, err := Decode(token)
	if err != nil {
		return false, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	return claims.VerifyExpiresAt(time.Now().Unix(), true), nil
}

func GetTokenByHeader(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.Index(authHeader, "Bearer") == 0 {
		return authHeader[7:], true
	}
	return "", false
}