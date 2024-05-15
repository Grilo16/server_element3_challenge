package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/Nerzal/gocloak/v13"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var client gocloak.GoCloak = *gocloak.NewClient("http://keycloak:8080")
const (
	realm      = "Element3"
	clientID   = "e3-challenge-server"
	clientSecret = "Q4xgljWZtLQGj50FTZqRC4rkEkBUfS0u"
)

type UserAuth struct{
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func InitializeAuthMiddleware(userService *user.UserService) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: jwt.IdentityKey,

		Authenticator: func(c *gin.Context) (interface{}, error) {

			grantType := "password"
			clientID := clientID 
			clientSecret := clientSecret 
			scope := "openid profile email offline_access"
			var loginVals UserAuth
			
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			params := gocloak.TokenOptions{
				GrantType: &grantType,
				ClientID:  &clientID,
				ClientSecret: &clientSecret,
				Scope:    &scope,
				Username: &loginVals.Email,
				Password: &loginVals.Password,
			}
			// Authenticate against Keycloak
			token, err := client.GetToken(c, realm, params)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			// Assuming token is not null and authentication is successful
			fmt.Println(token)
			return token, nil

		
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*user.User); ok {
				return true
			}

			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware
}


	// var loginVals UserAuth
			// if err := c.ShouldBind(&loginVals); err != nil {
			// 	return "", jwt.ErrMissingLoginValues
			// }
			// email := loginVals.Email
			// password := loginVals.Password
			// user, err := userService.GetUserByEmail(email)
			// if err != nil {
			// 	return nil, jwt.ErrFailedAuthentication
			// }

			// if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			// 	return nil, jwt.ErrFailedAuthentication
			// }
			// return user, nil