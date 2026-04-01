package middleware

import (
	"log"
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

		// 2FA properties checks
		enabled2FA, _ := claims["enabled_2fa"].(bool)
        otpVerified, _ := claims["is_otp_verified"].(bool)

		// allowed path/routes
        allowed := []string{"/enable-2fa", "/verify-otp"}
        path := c.FullPath() // canonical route pattern

		// If 2FA is enabled but OTP is not yet verified,
        // restrict access to all routes except those explicitly allowed.
        if enabled2FA && !otpVerified {
            allowedRoute := false
            for _, route := range allowed {
                if strings.Contains(path, route) {
                    allowedRoute = true
                    break
                }
            }

			// if route is not allowed, block access
            if !allowedRoute {
                c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                    "error": "2FA required. Please verify OTP.",
                })
                return
            }
        }

		// set the user_id from services/user.go, in the context so that it can be accessed in the handler
		c.Set("user_id", userID)
		log.Println("Middleware userID:", userID)

		c.Next() // continue to the next request
	}
}
