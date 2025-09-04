package test

import (
	"gofiber-boilerplate/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Username: "testuser",
		Password: "testpassword",
		Email:    "test@example.com",
		Country:  "ID",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotEmpty(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.NotEmpty(t, responseBody.Data.Token)
}

func TestRegisterInvalid(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		Username: "tes",
		Password: "tes",
		Email:    "test",
		Country:  "ID",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Empty(t, responseBody.Data)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.RegisterUserRequest{
		Username: "testuser",
		Password: "testpassword",
		Email:    "test@example.com",
		Country:  "ID",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Empty(t, responseBody.Data)
}

func TestLogin(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		Identifier: "testuser",
		Password:   "testpassword",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, responseBody.Data.ID)
	assert.Equal(t, "testuser", responseBody.Data.Username)
	assert.NotEmpty(t, responseBody.Data.Token)
}

func TestLoginInvalid(t *testing.T) {
	ClearAll()
	TestRegister(t)
	requestBody := model.LoginUserRequest{
		Identifier: "testuser",
		Password:   "wrongpassword",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Empty(t, responseBody.Data)
}
