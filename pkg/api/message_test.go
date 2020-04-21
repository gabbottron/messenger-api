package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var Message_Create_Text_Good string = `{"sender": %s, "recipient": %s, "content": {"type": "text", "text": %s}}`
var Message_Create_Image_Good string = `{"sender": %s, "recipient": %s, "content": {"type": "image", "imageuri": %s, "imagetype": %s}}`
var Message_Create_Video_Good string = `{"sender": %s, "recipient": %s, "content": {"type": "video", "videouri": %s, "encoding": %s, "length": %s}}`

// Create ----------------------
func (suite *ApiTestSuite) TestCreateMessageWithNoRequestBody() {
	code, token := TryLogin(suite, Login_Geoff_Good)
	suite.Token = token
	assert.Equal(suite.T(), 200, code, "Login should work.")

	// Must reset the recorder before making another request
	suite.ResetRecorder()

	assert.Equal(suite.T(), 422, tryCreateMessage(suite, ""), "Create should return 422")
}

func (suite *ApiTestSuite) TestCreateTextMessage() {
	code, token := TryLogin(suite, Login_Geoff_Good)
	suite.Token = token
	assert.Equal(suite.T(), 200, code, "Login should work.")

	// Must reset the recorder before making another request
	suite.ResetRecorder()

	payload := fmt.Sprintf(Message_Create_Text_Good, "1", "2", "\"Hey man. How goes?\"")

	assert.Equal(suite.T(), 200, tryCreateMessage(suite, payload), "Create should return 200")
}

func (suite *ApiTestSuite) TestCreateImageMessage() {
	code, token := TryLogin(suite, Login_Geoff_Good)
	suite.Token = token
	assert.Equal(suite.T(), 200, code, "Login should work.")

	// Must reset the recorder before making another request
	suite.ResetRecorder()

	payload := fmt.Sprintf(Message_Create_Image_Good, "1", "2", "\"http://www.amazon.com/img/someimage.jpg\"", "\"jpeg\"")

	assert.Equal(suite.T(), 200, tryCreateMessage(suite, payload), "Create should return 200")
}

func (suite *ApiTestSuite) TestCreateVideoMessage() {
	code, token := TryLogin(suite, Login_Geoff_Good)
	suite.Token = token
	assert.Equal(suite.T(), 200, code, "Login should work.")

	// Must reset the recorder before making another request
	suite.ResetRecorder()

	payload := fmt.Sprintf(Message_Create_Video_Good, "1", "2", "\"http://www.amazon.com/vid/somevid.mpg\"", "\"mpeg\"", "521")

	assert.Equal(suite.T(), 200, tryCreateMessage(suite, payload), "Create should return 200")
}

func tryCreateMessage(suite *ApiTestSuite, json_payload string) int {
	// This is actually setting the router up to handle this and only this request type
	// This request type requires authentication, so we use router grouping and middleware
	auth := suite.Router.Group("/auth")
	auth.Use(suite.AuthMiddleware.MiddlewareFunc())
	{
		auth.POST(MessageRoute, HandleMessageCreateRequest)
	}

	// Serve the single request to the router (gin.Engine)
	suite.Router.ServeHTTP(suite.Recorder, GetRequestAuthed(http.MethodPost, "/auth/message", json_payload, suite.Token))

	// Get the response
	resp := suite.Recorder.Result()

	// Make sure the content type is set correctly
	assert.Equal(suite.T(), ContentTypeJSON, resp.Header.Get("Content-Type"), "Content-Type in response should be JSON.")

	return resp.StatusCode
}
