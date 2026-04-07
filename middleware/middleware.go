package middleware

import (
	"net/http"
	"strings"
	"taskapi/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract header from the request
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// check tokenString if it's present and trim off the "Bearer" prefix, it's not needed
		// extract the token from the tokenString
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		// check the token validity by calling the ParseToken function from the services package
		token, err := services.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// uses AbortWithStatusJSON to stop further processing of the request if the token is invalid
		// checks claims are valid, extract claim
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		// 2FA properties checks
		enabled2FA, _ := claims["enabled_2fa"].(bool)
		isVerified, _ := claims["is_2fa_verified"].(bool)
		twoFactorType, _ := claims["two_factor_type"].(string)

		if enabled2FA && !isVerified {
			path := c.Request.URL.Path
			var allowedRoutes []string

			switch twoFactorType {
				case "email":
					allowedRoutes = []string{
						"/verify-email-2fa",
						"/enable-email-2fa",
						"/disable-email-2fa",
					}
				case "totp":
					allowedRoutes = []string{
						"/verify-totp",
						"/enable-totp",
						"/disable-totp",
					}
				default:
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"error": "Invalid 2FA configuration",
					})
					return
			}

			allowed := false
			for _, route := range allowedRoutes {
				if path == route {
					allowed = true
					break
				}
			}

			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "2FA verification required",
				})
				return
			}
		}

		// set the user_id from services/user.go, in the context so that it can be accessed in the handler
		c.Set("user_id", userID)
		c.Next() // continue to the next request
	}
}
