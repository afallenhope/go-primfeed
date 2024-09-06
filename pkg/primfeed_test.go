package primfeed

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetToken(t *testing.T) {
	// Arrange
	pf := NewPrimfeed(APIURL)
	mockToken := "0123456789abcdef"

	// Act
	pf.SetToken(mockToken)

	// Assert
	assert.Len(t, mockToken, 16)
	assert.Equal(t, "0123456789abcdef", mockToken)
}

func TestLogin(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" && r.Method == "POST" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"user":"testuser", "token":"0123456789abcdef"}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	pf := NewPrimfeed(mockServer.URL)

	// Act
	loginResponse, err := pf.Login("username", "password", nil)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "testuser", loginResponse.User)
	assert.Equal(t, "0123456789abcdef", loginResponse.Token)
}

func TestGetUserFollowers(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/entity/testuser/followers" && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `[{"id": "123", "name": "Other Test User", "handle": "othertestuser"}]`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()
	pf := NewPrimfeed(mockServer.URL)

	// Act
	followersResponse, err := pf.GetUserFollowers("testuser")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, followersResponse, 1)
	assert.Equal(t, "123", followersResponse[0].ID)
	assert.Equal(t, "Other Test User", followersResponse[0].Name)
	assert.Equal(t, "othertestuser", followersResponse[0].Handle)
}

func TestGetUserFollows(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/entity/othertestuser/followed" && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `[{"id": "1", "name": "Test User", "handle": "testuser"}]`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()
	pf := NewPrimfeed(mockServer.URL)

	// Act
	followsResponse, err := pf.GetUserFollows("othertestuser")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, followsResponse, 1)
	assert.Equal(t, "1", followsResponse[0].ID)
	assert.Equal(t, "Test User", followsResponse[0].Name)
	assert.Equal(t, "testuser", followsResponse[0].Handle)
}
