package Routes

import (
	"first-api/Controllers"
	"first-api/Structs"
	"first-api/Models"
	"time"
	"log"
	"errors"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

var identityKey = "email"

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())



	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Second * 30,
		MaxRefresh:  time.Second * 60,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Models.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[identityKey]
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Structs.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email

			var user Models.User
			err := Models.GetUserByEmail(&user, email)

			if err != nil {
				return nil, errors.New("Invalid Email")
			} else {
				return user, nil
			}

		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time){
			c.JSON(code, gin.H{
				"code":    code,
				"message": token,
				"expire": expire,
			})
		},

		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})


	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/refresh_token",authMiddleware.RefreshHandler)

	r.GET("/http-get",func(c *gin.Context){
		resp, err := http.Get("https://itx.trainz.id/")

		if err != nil {
		    log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
		    bodyBytes, err := ioutil.ReadAll(resp.Body)
		    if err != nil {
		        log.Fatal(err)
		    }
		    bodyString := string(bodyBytes)
		    //log.Info(bodyString)
		    c.String(200, bodyString)
		}
	})
	r.Use(authMiddleware.MiddlewareFunc())
	{

		grp1 := r.Group("/user-api")
		{
			grp1.GET("user", Controllers.GetUsers)
			grp1.POST("user", Controllers.CreateUser)
			grp1.GET("user/:id", Controllers.GetUserByID)
			grp1.PUT("user/:id", Controllers.UpdateUser)
			grp1.DELETE("user/:id", Controllers.DeleteUser)
		}
	}

	return r
}
