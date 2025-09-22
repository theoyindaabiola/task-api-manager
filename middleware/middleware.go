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
		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer") {
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

		// check 2FA/TOTP claims
		enabled2FA, ok := claims["enabled_2fa"].(bool)
		if !ok {
			// claim is missing
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing 2FA claim"})
			return
		}
		totpVerified, _ := claims["is_totp_verified"].(bool)

		// enforce TOTP if 2FA enabled
		if enabled2FA && !totpVerified && !strings.HasSuffix(c.FullPath(), "/verify-totp") {
            c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access this path. 2FA is enabled, please verify TOTP."})
            c.Abort()
            return
		}

		// set the user_id from services/user.go, in the context so that it can be accessed in the handler
		c.Set("user_id", userID)
		c.Next() // continue to the next request
	}
}
