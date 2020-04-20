package api

import (
	"encoding/json"
	"fmt"
	//jwt "github.com/gabbottron/gin-jwt"
	"github.com/gabbottron/messenger-api/pkg/jwt"
	"github.com/gabbottron/messenger-api/pkg/datastore"
	"github.com/gabbottron/messenger-api/pkg/datatypes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ApiTestSuite struct {
	suite.Suite
	// ResponseRecorder is an implementation of http.ResponseWriter that
	// records its mutations for later inspection in tests.
	Recorder *httptest.ResponseRecorder
	// Context maintains state between middleware
	Context *gin.Context
	// Engine is the router instance (muxer, etc)
	Router *gin.Engine

	// This is our ginJWT middleware object
	AuthMiddleware *jwt.GinJWTMiddleware

	// A JWT session token if login has occured
	// this is probably in the conetxt.. may remove it
	Token string

	// Tracking variables for initialization
	Env_loaded     bool
	Db_initialized bool
}

// Run the test suite setup routine before each test
func (suite *ApiTestSuite) SetupTest() {
	fmt.Println("In SetupTest!")
	// Setup the auth middleware
	suite.AuthMiddleware, _ = GetAuthMiddleware()

	suite.Recorder = httptest.NewRecorder()
	suite.Context, suite.Router = gin.CreateTestContext(suite.Recorder)

	suite.Token = ""

	if !suite.Env_loaded {
		//fmt.Println("Loading ENV...")
		err := godotenv.Load("../../.env")
		assert.Nil(suite.T(), err, "ENV load should return no errors")
		if err == nil {
			suite.Env_loaded = true
		}
	}

	if !suite.Db_initialized {
		//fmt.Println("Initializing the database...")
		err := datastore.InitDB()
		assert.Nil(suite.T(), err, "Datastore.InitDB should return no errors")
		if err == nil {
			suite.Db_initialized = true
		}
	}
}

func (suite *ApiTestSuite) ResetRecorder() {
	suite.Recorder = httptest.NewRecorder()
	suite.Context, suite.Router = gin.CreateTestContext(suite.Recorder)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

// This function will attempt to POST to the login route
// It will return the status code and the token (if result was 200)
func TryLogin(suite *ApiTestSuite, json_payload string) (int, string) {
	// This is actually setting the router up to handle this and only this request type
	suite.Router.POST(LoginRoute, suite.AuthMiddleware.LoginHandler)

	// Serve the single request to the router (gin.Engine)
	suite.Router.ServeHTTP(suite.Recorder, GetRequest(http.MethodPost, LoginRoute, json_payload))

	// Get the response
	resp := suite.Recorder.Result()
	resp_body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Response body decoding shoud not return an error.")

	// Make sure the content type is set correctly
	assert.Equal(suite.T(), ContentTypeJSON, resp.Header.Get("Content-Type"), "Content-Type in response should be JSON.")

	// Store the JWT here
	var token string

	if resp.StatusCode == 200 {
		data := new(datatypes.HttpLoginResponseJSON)
		err = json.Unmarshal([]byte(resp_body), data)
		assert.Nil(suite.T(), err, "Got 200 from login, should be able to unmarshal the response body.")
		if err == nil {
			token = data.Token
		}
	}

	return resp.StatusCode, token
}

// This function builds a http request object and returns it
func GetRequest(method string, route string, body string) *http.Request {
	reader := strings.NewReader(body)

	request := httptest.NewRequest(method, route, reader)

	// Set the headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request
}

// This function builds a http request object and returns it
func GetRequestAuthed(method string, route string, body string, token string) *http.Request {
	request := GetRequest(method, route, body)

	// Set the authorization header with the JWT
	request.Header.Set("Authorization", "Bearer "+token)

	return request
}
