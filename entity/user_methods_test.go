package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	expectedUser := Users{
		Uid:      "94891akc",
		Username: "Test Name",
		Email:    "test@gmail.com",
		Type:     "normal",
	}

	user, err := NewUser("94891akc", "Test Name", "test@gmail.com", "normal")

	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Uid, user.Uid)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Type, user.Type)

}

func TestNewUserValidation(t *testing.T) {
	type test struct {
		uid      string
		username string
		email    string
		userType string
		expError error
	}

	tests := []test{
		{
			uid:      "94891akc",
			username: "Test Name",
			email:    "test@gmail.com",
			userType: "normal",
			expError: nil,
		},
		{
			uid:      "",
			username: "Test Name",
			email:    "test@gmail.com",
			userType: "normal",
			expError: ErrInvalidUserUidBlank,
		},
		{
			uid:      "94891akc",
			username: "",
			email:    "test@gmail.com",
			userType: "normal",
			expError: ErrInvalidUserNameBlank,
		},
		{
			uid:      "94891akc",
			username: "Test Name",
			email:    "",
			userType: "normal",
			expError: ErrInvalidUserEmailBlank,
		},
		{
			uid:      "94891akc",
			username: "Test Name",
			email:    "test@gmail.com",
			userType: "",
			expError: ErrInvalidUserTypeBlank,
		},
	}

	for _, testCase := range tests {
		_, err := NewUser(testCase.uid, testCase.username, testCase.email,
			testCase.userType)
		assert.Equal(t, testCase.expError, err)
	}

}
