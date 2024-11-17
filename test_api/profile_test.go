package test_api

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestGetSmallProfileAndUpdateSmallProfileAPI(t *testing.T) {
	type SmallProfileResponse struct {
		Name        string `json:"name"`
		SecondName  string `json:"second_name"`
		Avatar      string `json:"avatar"`
		Birthday    string `json:"birthday"`
		Description string `json:"description"`
		Gender      int32  `json:"gender"`
		Email       string `json:"email"`
	}

	var endpoint = serverAddr + "/profile/"

	client := &http.Client{}

	johnAccessToken := getJohnAccessToken()

	t.Run("update small profile", func(t *testing.T) {
		updateSmallProfilePayload := map[string]interface{}{
			"name":        johnUpdatedName,
			"second_name": johnUpdatedSecondName,
			// empty avatar
			"birthday": johnUpdatedBirthday,
			// empty description
			"gender": johnUpdatedGender,
		}

		body, _ := json.Marshal(updateSmallProfilePayload)
		req, err := http.NewRequest(http.MethodPatch, endpoint+"small", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", johnAccessToken)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("valid get small profile request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, endpoint+strconv.Itoa(johnID), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response SmallProfileResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, johnUpdatedName, response.Name)
		assert.Equal(t, johnUpdatedSecondName, response.SecondName)
		assert.Equal(t, johnUpdatedBirthday, response.Birthday)
		assert.Equal(t, johnUpdatedGender, response.Gender)

		assert.Equal(t, "", response.Avatar)
		assert.Equal(t, "", response.Description)
	})

	t.Run("update small profile description only", func(t *testing.T) {
		updateSmallProfilePayload := map[string]interface{}{
			"description": johnUpdatedDescription,
		}

		body, _ := json.Marshal(updateSmallProfilePayload)
		req, err := http.NewRequest(http.MethodPatch, endpoint+"small", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", johnAccessToken)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("valid get small profile request after update description", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, endpoint+strconv.Itoa(johnID), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response SmallProfileResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, johnUpdatedName, response.Name)
		assert.Equal(t, johnUpdatedSecondName, response.SecondName)
		assert.Equal(t, johnUpdatedBirthday, response.Birthday)
		assert.Equal(t, johnUpdatedGender, response.Gender)
		assert.Equal(t, johnEmail, response.Email)
		assert.Equal(t, johnUpdatedDescription, response.Description)

		assert.Equal(t, "", response.Avatar)
	})

	t.Run("get small profile with invalid id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, endpoint+"invalid", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get small profile with not existing id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, endpoint+"1113111222", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
