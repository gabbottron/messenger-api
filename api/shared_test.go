package api

import (
	"github.com/stretchr/testify/assert"
)

// Empty JSON payload
var EmptyJSON string = "{}"

func (suite *ApiTestSuite) TestValidatePasswords() {
	assert.False(suite.T(), ValidatePassword(""), "Empty: Should be false.")
	assert.False(suite.T(), ValidatePassword("aA@6bb"), "Too short: Should be false.")
	assert.False(suite.T(), ValidatePassword("aAfadskfjdf6bb"), "No Specials: Should be false.")
	assert.False(suite.T(), ValidatePassword("ABCDE@@6786......"), "No Lowercase: Should be false.")
	assert.False(suite.T(), ValidatePassword("abcdefgh@@6789......"), "No Uppercase: Should be false.")
	assert.False(suite.T(), ValidatePassword("ABCDE@@defkl.+-^%."), "No Numbers: Should be false.")
	assert.False(suite.T(), ValidatePassword("AAAAAA@@defkl.+-^%."), "Too many repeated characters: Should be false.")

	assert.True(suite.T(), ValidatePassword("aA@6. krA t f bb"), "Meets Specs: Should be true")
}
