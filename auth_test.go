package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test the auth package

func TestGenerateToken(t *testing.T) {
	username := "test"
	tokenString, err := GenerateToken(username)
	assert.NotNil(t, tokenString)
	assert.Nil(t, err)
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
