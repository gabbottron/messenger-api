package api

import (
	//"errors"
	"github.com/gabbottron/gin-jwt"
	"github.com/gabbottron/messenger-api/src/datastore"
	"github.com/gabbottron/messenger-api/src/datatypes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// The user identity key for JWT middleware
	KEY_USERID = "userid"

	APIHealthCheckRoute string = "/check"

	LoginRoute string = "/login"

	// Identifiers for the user in claims/context
	USER_ID_KEY        string = "id"
	USER_SESSION_TOKEN string = "sessiontoken"

	// For identifying the content type in Content-Type header
	ContentTypeJSON string = "application/json; charset=utf-8"
)

// For marshaling the POST login values
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AddHealthCheckRoutes(r *gin.Engine) *gin.Engine {
	r.GET(APIHealthCheckRoute, HandleHealthCheckRequest)
	return r
}

func AddAdditionalRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {
	r = AddUserRoutes(r, authMiddleware)
	r = AddMessageRoutes(r, authMiddleware)

	return r
}

func HandleHealthCheckRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})
}

func loginHandler(c *gin.Context) (interface{}, error) {
	var loginAuth login
	if err := c.ShouldBind(&loginAuth); err != nil {
		log.Println(err)
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := datastore.AuthenticateUser(loginAuth.Username, loginAuth.Password)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	return user, err
}

func authHandler(user interface{}, c *gin.Context) bool {
	// Note: No advanced auth checking at router level currently...
	return true
}

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm: "messengerAPI",
		// TODO: Change this before deploying!
		Key:         []byte("Hey, guys! Whoa, Big Gulps, huh? Alright... Welp, see ya later!"),
		Timeout:     time.Hour * 8,
		MaxRefresh:  time.Hour * 8,
		IdentityKey: USER_ID_KEY,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*datatypes.User); ok {
				return jwt.MapClaims{
					USER_ID_KEY:        v.UserID,
					USER_SESSION_TOKEN: uuid.New().String(),
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: loginHandler,
		Authorizator:  authHandler,
		IdentityHandler: func(c *gin.Context) interface{} {
			// extract the claims
			claims := jwt.ExtractClaims(c)

			// Populate and return the user object
			return &datatypes.User{
				UserID:       int(claims[USER_ID_KEY].(float64)),
				SessionToken: claims[USER_SESSION_TOKEN].(string),
			}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				// Decode user id and type from the token string if need be
				//"userid":   decoded_userid,
				//"usertype": decoded_usertype,
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

func SetCrossOrigin(r *gin.Engine) *gin.Engine {
	// Set the application endpoint
	APPEndpoint := os.Getenv("APP_ENDPOINT")
	if APPEndpoint == "" {
		log.Fatal("APP_ENDPOINT Was not set in .env!")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{APPEndpoint},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return r
}

// Initialize the router
func InitRouter() *gin.Engine {
	log.Println("Initializing router...")

	// Build the router
	r := gin.Default()

	// Set cross-origin support
	r = SetCrossOrigin(r)

	// max file size upload
	//r.MaxMultipartMemory = 8 << 20

	// Assign the authorization middleware
	authMiddleware, err := GetAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Add Health Check Routes
	r = AddHealthCheckRoutes(r)

	// Add login routes
	r.POST("/login", authMiddleware.LoginHandler)

	// Add protected routes
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	// Each section of the API will set up public/private
	// routes in implementation file...
	r = AddAdditionalRoutes(r, authMiddleware)

	return r
}
