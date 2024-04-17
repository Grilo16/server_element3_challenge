package middleware

import (
	"log"
	"time"

	"github.com/Grilo16/server_element3_challenge/user"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: int(user.Id),
					
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			idFloat64 := claims[jwt.IdentityKey].(float64)
			id := int(idFloat64)
			return &user.User{
				Id: id,
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals UserAuth
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password
			user, err := userService.GetUserByEmail(email)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
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