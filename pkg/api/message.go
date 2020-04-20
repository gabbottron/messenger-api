package api

import (
	//jwt "github.com/gabbottron/gin-jwt"
	"github.com/gabbottron/messenger-api/pkg/jwt"
	"github.com/gabbottron/messenger-api/pkg/datastore"
	"github.com/gabbottron/messenger-api/pkg/datatypes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	MessageRoute  string = "/message"
	MessagesRoute string = "/messages"
)

func HandleMessageCreateRequest(c *gin.Context) {
	// The message object
	var messageData datatypes.MessageJSON

	// Attempt to bind the body to message obj
	err := c.ShouldBindJSON(&messageData)
	if err != nil {
		log.Printf("Create message -> Could not bind request body to struct!: %s", err.Error())
		ReplyUnprocessableEntity(c, "Could not bind request body to struct!")
		return
	}

	// Get the validated userID
	claim_userid := GetIDFromClaim(c)

	// Set the sender ID based on the validated ID in the claim
	messageData.SenderID = claim_userid

	// Attempt to insert the message
	err = datastore.InsertMessageRecord(&messageData)
	if err != nil {
		log.Printf("Create message -> datastore.InsertMessageRecord() returned an error:  %s", err.Error())

		// The error type is specific and known (defined in datastore)
		if de, ok := err.(*datastore.DatastoreError); ok {
			switch de.Code {
			case datastore.ErrorConstraintViolation:
				ReplyUnprocessableEntity(c, de.Text)
			case datastore.ErrorForbidden:
				ReplyForbidden(c, de.Text)
			}
			return
		}

		ReplyInternalServerError(c, "There was a database error inserting the new message record!")
		return
	}

	c.JSON(http.StatusOK, messageData)
	return
}

func HandleMessagesRequest(c *gin.Context) {
	// Get the userID from the claims
	claim_userid := GetIDFromClaim(c)

	var recipient int
	var start int
	var limit = 100

	qstrings := c.Request.URL.Query()
	if len(qstrings) > 1 {
		if val, ok := qstrings["recipient"]; ok {
			i, err := strconv.Atoi(strings.ToLower(strings.TrimSpace(val[0])))
			if err == nil {
				recipient = i
			} else {
				// Recipient is REQUIRED
				ReplyUnprocessableEntity(c, "Could not understand recipient value!")
				return
			}
		} else {
			// Recipient is REQUIRED
			ReplyUnprocessableEntity(c, "Missing required querystring param (recipient)!")
			return
		}

		if val, ok := qstrings["start"]; ok {
			i, err := strconv.Atoi(strings.ToLower(strings.TrimSpace(val[0])))
			if err == nil {
				start = i
			} else {
				ReplyUnprocessableEntity(c, "Could not understand start value!")
				return
			}
		} else {
			// Start is REQUIRED
			ReplyUnprocessableEntity(c, "Missing required querystring param (start)!")
			return
		}

		if val, ok := qstrings["limit"]; ok {
			i, err := strconv.Atoi(strings.ToLower(strings.TrimSpace(val[0])))
			if err == nil {
				// Control for sensible values, negative values passed to
				// db as a limit will cause an exception
				if i > 0 && i < 10000 {
					limit = i
				}
			}
		}
		/* DEBUG
		fmt.Println("Keys passed in from querystring:")
		for k, v := range qstrings {
			fmt.Printf("key[%s] value[%s]\n", k, v)
		} */
	} else {
		// There are 2 required querystring params
		ReplyUnprocessableEntity(c, "This route requires querystring params (recipient, start)!")
		return
	}

	// Now check the DB for the tags
	messages, err := datastore.GetMessages(claim_userid, recipient, start, limit)
	if err != nil {
		log.Printf("ERROR in HandleMessagesRequest -> %s", err.Error())
		ReplyNotFound(c, "Messages not found!")
		return
	}

	c.JSON(http.StatusOK, messages)
	return
}

func AddMessageRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {
	// Add protected routes
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// CRUD
		auth.GET(MessagesRoute, HandleMessagesRequest)
		auth.POST(MessageRoute, HandleMessageCreateRequest)
		auth.POST(MessagesRoute, HandleMessageCreateRequest)
	}

	return r
}
