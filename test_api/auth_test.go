package test_api

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	johnID       int
	johnEmail    = fmt.Sprintf("john.doe%d@example.com", time.Now().Unix())
	johnPassword = "password123"
)

func TestSignUpAPI(t *testing.T) {
	type SignUpResponse struct {
		ID           int    `json:"id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var endpoint = serverAddr + "/auth/sign-up"

	client := &http.Client{}

	t.Run("valid sign-up request", func(t *testing.T) {
		signUpPayload := map[string]string{
			"name":        "John",
			"second_name": "Doe",
			"email":       johnEmail,
			"password":    johnPassword,
		}

		body, _ := json.Marshal(signUpPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response SignUpResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.NotZero(t, response.ID)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)

		johnID = response.ID
	})

	testMissingFields := func(t *testing.T, name, secondName, email, password string) {
		invalidPayload := map[string]string{
			"name":        name,
			"second_name": secondName,
			"email":       email,
			"password":    password,
		}

		body, _ := json.Marshal(invalidPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	t.Run("invalid sign-up request - missing name", func(t *testing.T) {
		testMissingFields(t, "", "Doe", johnEmail, johnPassword)
	})

	t.Run("invalid sign-up request - missing second name", func(t *testing.T) {
		testMissingFields(t, "John", "", johnEmail, johnPassword)
	})

	t.Run("invalid sign-up request - missing email", func(t *testing.T) {
		testMissingFields(t, "John", "Doe", "", johnPassword)
	})

	t.Run("invalid sign-up request - missing password", func(t *testing.T) {
		testMissingFields(t, "John", "Doe", johnEmail, "")
	})

	t.Run("johnEmail already in use", func(t *testing.T) {
		signUpPayload := map[string]string{
			"name":        "John",
			"second_name": "Doe",
			"email":       johnEmail,
			"password":    "password123",
		}

		body, _ := json.Marshal(signUpPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSignInAPI(t *testing.T) {
	type SignInResponse struct {
		ID           int    `json:"id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var endpoint = serverAddr + "/auth/sign-in"

	client := &http.Client{}

	t.Run("valid sign-in request", func(t *testing.T) {
		signUpPayload := map[string]string{
			"email":    johnEmail,
			"password": "password123",
		}

		body, _ := json.Marshal(signUpPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response SignInResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.NotZero(t, response.ID)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
	})

	testMissingFields := func(t *testing.T, email, password string) {
		invalidPayload := map[string]string{
			"email":    email,
			"password": password,
		}

		body, _ := json.Marshal(invalidPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	t.Run("invalid sign-in request - missing password", func(t *testing.T) {
		testMissingFields(t, "john.doe@example.com", "")
	})

	t.Run("invalid sign-in request - missing email", func(t *testing.T) {
		testMissingFields(t, "", "password123")
	})

	t.Run("invalid johnPassword", func(t *testing.T) {
		signUpPayload := map[string]string{
			"name":        "John",
			"second_name": "Doe",
			"email":       johnEmail,
			"password":    "invalid",
		}

		body, _ := json.Marshal(signUpPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestRefreshTokensAPI(t *testing.T) {
	type RefreshTokensResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var endpoint = serverAddr + "/auth/refresh-tokens"

	client := &http.Client{}

	t.Run("valid refresh tokens request", func(t *testing.T) {
		refreshTokensPayload := map[string]interface{}{
			"id":       johnID,
			"password": johnPassword,
		}

		body, _ := json.Marshal(refreshTokensPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response RefreshTokensResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
	})

	t.Run("invalid refresh tokens request - missing password", func(t *testing.T) {
		invalidPayload := map[string]interface{}{
			"id": johnID,
		}

		body, _ := json.Marshal(invalidPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid refresh tokens request - missing id", func(t *testing.T) {
		invalidPayload := map[string]interface{}{
			"id":       -1,
			"password": "password123",
		}

		body, _ := json.Marshal(invalidPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unauthorized refresh tokens request", func(t *testing.T) {
		refreshTokensPayload := map[string]interface{}{
			"id":       johnID,
			"password": "wrongpassword",
		}

		body, _ := json.Marshal(refreshTokensPayload)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
