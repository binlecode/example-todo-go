package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const TOKEN_EXP_TIME_MIN = 15

// set a global secret key for signing the jwt
var jwtKey = []byte(getEnv("SECRET_KEY", "this-should-be-a-long-secret"))

// define a list of users fixture data
var hashedPassword1, _ = GenerateHashedPassword("password1")
var hashedPassword2, _ = GenerateHashedPassword("password2")
var users = []User{
	User{"alice", string(hashedPassword1), "admin,editor"},
	User{"bob", string(hashedPassword2), ""},
}

// getUser by username from users list
func getUser(username string) (User, bool) {
	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    string `json:"roles,omitempty"`
}

// Credentials is a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims is a struct to be encoded to a JWT.
type Claims struct {
	Username string `json:"username"`
	// custom claims
	Roles string `json:"roles,omitempty"`
	// embed jwt standard claims
	jwt.RegisteredClaims
}

// https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
// https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/#jwt-authentication-in-golang

// TokenMiddleware is a middleware to check the token in the request header
// and add the claims to the context.
// If the token is not valid, return 401.
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check header for token, return 401 if not found or not valid
		token := r.Header.Get("Authorization")

		// validate token
		claims, err := ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Failed authorization"))
			return
		}

		// add claims to context
		log.Info("claims: ", claims)
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		// call next handler function
		next.ServeHTTP(w, r)
	})
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// get the claims from the context
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed authorization"))
		return
	}

	// return user info
	respondWithJSON(w, http.StatusOK, claims)
}

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// if the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the credentials are valid
	tokenString, err := Authenticate(creds.Username, creds.Password)
	if err != nil {
		// if there is an error in the credentials, return an HTTP error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// return token
	//w.Write([]byte(tokenString))
	// return token in JSON
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func Authenticate(username, password string) (string, error) {
	// check if the username exists
	user, ok := getUser(username)
	if !ok {
		return "", errors.New("user not found")
	}
	//if !userExists(username) {
	//	//return "", jwt.ErrTokenInvalidSubject
	//	return "", errors.New("user not found")
	//}

	// check if the password is correct
	if !IsPasswordValid(user.Username, password) {
		return "", errors.New("invalid password")
	}

	// generate a token
	tokenString, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func userExists(username string) bool {
	// find user from users list
	for u := range users {
		if users[u].Username == username {
			return true
		}
	}
	return false
}

// IsPasswordValid checks if the password is valid for the given username
// Return false if user does not exist or password is wrong.
// Use the bcrypt package to compare the hashed password.
func IsPasswordValid(username, password string) bool {
	user, ok := getUser(username)
	if !ok {
		return false
	}
	// Use the bcrypt package to compare the hashed password
	expectedPassword := user.Password
	return bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password)) == nil
}

// GenerateHashedPassword returns the bcrypt hash of the password at the given
func GenerateHashedPassword(password string) ([]byte, error) {
	// use default cost value 10
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

// GenerateToken generates a jwt token for the given username
func GenerateToken(user User) (string, error) {
	// TODO: implement user custom claims

	// Set token expiration time
	// In JWT, the expiry time is expressed as unix milliseconds
	expTime := time.Now().Add(TOKEN_EXP_TIME_MIN * time.Minute)
	// create the claims
	claims := Claims{
		user.Username,
		user.Roles,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			Issuer:    "test",
		},
	}

	// create the token with signing algorithm, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// create the JWT string
	return token.SignedString(jwtKey)
}

// ValidateToken validates the token string using the jwt package
func ValidateToken(tokenString string) (*Claims, error) {
	// parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// check the signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}

		// return the secret key
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	// validate the token and return the custom claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
