package api

import (
	//jwt "github.com/gabbottron/gin-jwt"
	"github.com/gabbottron/messenger-api/pkg/jwt"
	"github.com/gabbottron/messenger-api/pkg/datastore"
	"github.com/gabbottron/messenger-api/pkg/datatypes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	UserRoute  string = "/user"
	UsersRoute string = "/users"
)

func HandleUserCreateRequest(c *gin.Context) {
	var userData datatypes.MsgrUserJSON

	// Attempt to bind the body to user obj
	err := c.ShouldBindJSON(&userData)
	if err != nil {
		log.Printf("Create user -> Could not bind request body to struct!: %s", err.Error())
		ReplyUnprocessableEntity(c, "Could not bind request body to struct!")
		return
	}

	if userData.Password == nil || userData.Username == nil {
		ReplyUnprocessableEntity(c, "Username and password are required fields!")
		return
	}

	// Validate the password
	if !ValidatePassword(*userData.Password) {
		log.Printf("Create user -> Password failed validation!")
		ReplyUnprocessableEntity(c, "Password did not meet security criteria!")
		return
	}

	// Now hash and salt the password
	hash, err := datastore.HashAndSalt(*userData.Password)
	if err != nil {
		log.Printf("Create user -> Password hashing failed!")
		ReplyInternalServerError(c, "There was an internal server error processing your request!")
		return
	}

	// Set the password to the hash
	*userData.Password = hash

	err = datastore.InsertUserRecord(&userData)
	if err != nil {
		log.Printf("Create user -> datastore.InsertUserRecord() returned an error:  %s", err.Error())

		// The error type is specific and known (defined in datastore)
		if de, ok := err.(*datastore.DatastoreError); ok {
			switch de.Code {
			case datastore.ErrorConstraintViolation:
				ReplyUnprocessableEntity(c, de.Text)
			}
			return
		}

		ReplyInternalServerError(c, "There was a database error inserting the new user record!")
		return
	}

	// no reason to send the password back, even hashed...
	*userData.Password = "**OBFUSCATED**"

	c.JSON(http.StatusOK, userData)
	return
}

func AddUserRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {
	r.POST(UserRoute, HandleUserCreateRequest)
	r.POST(UsersRoute, HandleUserCreateRequest)

	return r
}
