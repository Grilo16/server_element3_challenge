package middleware

import (
	"net/http"
	"strings"

	"github.com/Grilo16/server_element3_challenge/keycloak"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenAuthMiddleware(userService *user.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keycloakService := keycloak.NewKeycloakService()
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer Token missing"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token with keycloak
		_, err := keycloakService.RetrospectToken(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		userSub, err := parsedToken.Claims.GetSubject()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		authenticatedUser, err := userService.GetUserBySub(userSub)
		if err != nil {
			// Assume error means user not found; create new user
			userDetails := user.User{
				FirstName: claims["given_name"].(string),
				LastName:  claims["family_name"].(string),
				Email:     claims["email"].(string),
				Sub:       userSub,
			}
			authenticatedUser, err = userService.CreateNewUser(&userDetails)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.Set("sub", userSub)
		ctx.Set("user", authenticatedUser)
		ctx.Next()
	}
}
