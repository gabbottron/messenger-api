package api

import (
	//jwt "github.com/gabbottron/gin-jwt"
	"github.com/gabbottron/messenger-api/pkg/jwt"
	"github.com/gabbottron/messenger-api/pkg/datatypes"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

// Special characters for password validation
var special_chars []string = []string{"]", "^", "/", ",", "!", " ", "&", "%", "$", "#", "\"", "+", "*", "'", "\\", "[", ".", "(", ")", "-", "@", "_", "{", "}", "|", "~", "`", "<", ">", "?", "=", ";", ":"}

// Glyph groups for password validation
var re_uppercase *regexp.Regexp = regexp.MustCompile("[A-Z]")
var re_lowercase *regexp.Regexp = regexp.MustCompile("[a-z]")
var re_numbers *regexp.Regexp = regexp.MustCompile("[0-9]")

const (
	// Password strength settings
	passMustIncludeUpper     int = 1
	passMustIncludeLower     int = 1
	passMustIncludeNumbers   int = 1
	passMustIncludeSpecial   int = 1
	passMaximumRepeatedChars int = 2
	passLengthMustBe         int = 8
)

// --- Utility Functions ---
func ValidatePassword(plain_password string) bool {
	if len(plain_password) < passLengthMustBe {
		return false
	}

	if len(re_lowercase.FindAllStringIndex(plain_password, -1)) < passMustIncludeLower {
		return false
	}
	if len(re_uppercase.FindAllStringIndex(plain_password, -1)) < passMustIncludeUpper {
		return false
	}
	if len(re_numbers.FindAllStringIndex(plain_password, -1)) < passMustIncludeNumbers {
		return false
	}

	// Test for repeated characters
	repeat_count := 1
	last_char := ""
	for _, r := range plain_password {
		c := string(r)
		if c == last_char {
			repeat_count++
			if repeat_count > passMaximumRepeatedChars {
				return false
			}
		} else {
			repeat_count = 1
		}
		last_char = c
	}

	var special_count int
	for _, element := range special_chars {
		special_count += strings.Count(plain_password, element)
		if special_count >= passMustIncludeSpecial {
			return true
		}
	}

	return false
}

// Get the user ID from the claim in the JWT
func GetIDFromClaim(c *gin.Context) int {
	var claims_userid int
	claims := jwt.ExtractClaims(c)

	if claims[USER_ID_KEY] != nil {
		assert_float, ok := claims[USER_ID_KEY].(float64)
		if !ok {
			return 0
		}
		claims_userid = int(assert_float)
	}

	return claims_userid
}

// Get the session token from the claim in the JWT
func GetSessionTokenFromClaim(c *gin.Context) string {
	var claims_sessionid string
	claims := jwt.ExtractClaims(c)

	if claims[USER_SESSION_TOKEN] != nil {
		claims_sessionid = claims[USER_SESSION_TOKEN].(string)
	}

	return claims_sessionid
}

// --- HTTP Responses ---

func ReplyBadRequest(c *gin.Context, msg string) {
	response := datatypes.HttpResponseJSON{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: msg,
	}

	c.JSON(http.StatusBadRequest, response)
}

func ReplyForbidden(c *gin.Context, msg string) {
	response := datatypes.HttpResponseJSON{
		Status:  http.StatusText(http.StatusForbidden),
		Message: msg,
	}

	c.JSON(http.StatusForbidden, response)
}

func ReplyNotFound(c *gin.Context, msg string) {
	response := datatypes.HttpResponseJSON{
		Status:  http.StatusText(http.StatusNotFound),
		Message: msg,
	}

	c.JSON(http.StatusNotFound, response)
}

// USES:
// - The payload was correctly formed, but could not be
//   marshalled into the JSON struct
func ReplyUnprocessableEntity(c *gin.Context, msg string) {
	response := datatypes.HttpResponseJSON{
		Status:  http.StatusText(http.StatusUnprocessableEntity),
		Message: msg,
	}

	c.JSON(http.StatusUnprocessableEntity, response)
}

func ReplyInternalServerError(c *gin.Context, msg string) {
	response := datatypes.HttpResponseJSON{
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: msg,
	}

	c.JSON(http.StatusInternalServerError, response)
}
