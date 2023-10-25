package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test the auth package

func TestGetUser(t *testing.T) {
	user, ok := getUser("alice")
	assert.True(t, ok)
	assert.Equal(t, "alice", user.Username)
	user, ok = getUser("wrong")
	assert.False(t, ok)
}

func TestGenerateToken(t *testing.T) {
	user := User{Username: "test", Roles: "admin,editor"}
	tokenString, err := GenerateToken(user)
	assert.NotNil(t, tokenString)
	assert.Nil(t, err)
}

func TestValidateToken(t *testing.T) {
	user := User{Username: "test", Roles: "admin,editor"}
	tokenString, err := GenerateToken(user)
	assert.NotNil(t, tokenString)
	claims, err := ValidateToken(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, "test", claims.Username)
	assert.Equal(t, "admin,editor", claims.Roles)
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		username string
		password string
	}{
		{"alice", "wrong"},
		{"alice", "password1"},
		{"bob", "wrong"},
		{"wrong", "password2"},
	}
	for _, test := range tests {
		tokenString, err := Authenticate(test.username, test.password)
		if test.username == "alice" && test.password == "password1" {
			assert.NotNil(t, tokenString)
			assert.Nil(t, err)
		} else {
			assert.Empty(t, tokenString)
			assert.NotNil(t, err)
		}
	}
}
