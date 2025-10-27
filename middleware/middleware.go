package middleware

import (
	"net/http"
	"strings"
	"taskapi/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuthMiddleware(userService *services.UserService) gin.HandlerFunc {
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

		/** 
			Using the database as the source of truth instead of adding auth to the claims
			this is to avoid update mismatch btwn JWT & DB
		**/

		// get user from DB
		user, err := userService.GetUserByID(userID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}	

		// enforce SMS if 2FA enabled
		if user.Enabled2FA && !user.IsSmsVerified && !strings.HasSuffix(c.FullPath(), "/verify-sms") && !strings.HasSuffix(c.FullPath(), "/request-sms") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access this path. 2FA is enabled, please verify SMS code."})
			c.Abort()
			return
		}

		// set the user_id from services/user.go, in the context so that it can be accessed in the handler
		c.Set("user_id", userID)
		c.Next() // continue to the next request
	}
}
