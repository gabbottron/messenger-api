package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
)

// Login payloads
var Login_Geoff_Good string = `{"username": "gabbott","password": "blarg"}`
var Login_Geoff_Bad_NoPass string = `{"username": "gabbott"}`
var Login_Geoff_Bad_WrongPass string = `{"username": "gabbott","password": "notmypass"}`

// User payloads
var User_Create_Minimal string = `{"username": "testuser1", "password": "blarg182@fRt"}`
var User_Create_Missing_Password string = `{"username": "testuser2"}`
var User_Create_Blank_Password string = `{"username": "testuser2", "password": ""}`
var User_Create_Missing_Username string = `{"password": "blarg182@fRt"}`
var User_Create_Username_Too_Short string = `{"username": "c"}`

// BEGIN TESTS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (suite *ApiTestSuite) TestLoginGeoffGood() {
	code, _ := TryLogin(suite, Login_Geoff_Good)
	assert.Equal(suite.T(), 200, code)
}

func (suite *ApiTestSuite) TestLoginGeoffNoPass() {
	code, _ := TryLogin(suite, Login_Geoff_Bad_NoPass)
	assert.Equal(suite.T(), 401, code)
}

func (suite *ApiTestSuite) TestLoginGeoffWrongPass() {
	code, _ := TryLogin(suite, Login_Geoff_Bad_WrongPass)
	assert.Equal(suite.T(), 401, code)
}

// Test User Creation  --------------------------------------
func (suite *ApiTestSuite) TestCreateUserWithNoRequestBody() {
	assert.Equal(suite.T(), 422, tryCreateUser(suite, ""))
}

func (suite *ApiTestSuite) TestCreateUserWithEmpty() {
	assert.Equal(suite.T(), 422, tryCreateUser(suite, EmptyJSON))
}

func (suite *ApiTestSuite) TestCreateUserWithMissingReqField() {
	assert.Equal(suite.T(), 422, tryCreateUser(suite, User_Create_Missing_Password))
}

func (suite *ApiTestSuite) TestCreateUserWithBlankReqField() {
	assert.Equal(suite.T(), 422, tryCreateUser(suite, User_Create_Blank_Password))
}

func (suite *ApiTestSuite) TestCreateUserWithMinimal() {
	assert.Equal(suite.T(), 200, tryCreateUser(suite, User_Create_Minimal))
}

func (suite *ApiTestSuite) TestCreateUserWithTakenUsername() {
	assert.Equal(suite.T(), 422, tryCreateUser(suite, User_Create_Minimal))
}

// This function will attempt to POST to the user route
func tryCreateUser(suite *ApiTestSuite, json_payload string) int {
	// This is actually setting the router up to handle this and only this request type
	suite.Router.POST(UserRoute, HandleUserCreateRequest)

	// Serve the single request to the router (gin.Engine)
	suite.Router.ServeHTTP(suite.Recorder, GetRequest(http.MethodPost, "/user", json_payload))

	// Get the response
	resp := suite.Recorder.Result()

	// Make sure the content type is set correctly
	assert.Equal(suite.T(), ContentTypeJSON, resp.Header.Get("Content-Type"), "Content-Type in response should be JSON.")

	return resp.StatusCode
}
