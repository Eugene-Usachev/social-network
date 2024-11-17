package test_api

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"time"
)

var (
	johnID                 int
	johnEmail              = fmt.Sprintf("john.doe%d@example.com", time.Now().Unix())
	johnPassword           = "password123"
	johnUpdatedName        = "JohnU"
	johnUpdatedSecondName  = "DoeU"
	johnUpdatedBirthday    = "2000-01-01"
	johnUpdatedDescription = "Updated description"
	johnUpdatedGender      = int32(1)
)

type refreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getJohnAccessToken() string {
	client := &http.Client{}
	refreshTokensPayload := map[string]interface{}{
		"id":       johnID,
		"password": johnPassword,
	}

	body, _ := json.Marshal(refreshTokensPayload)
	req, err := http.NewRequest(http.MethodPost, serverAddr+"/auth/refresh-tokens", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		log.Fatalf("Failed to refresh tokens: %v", resp.StatusCode)
	}

	var response refreshTokensResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return response.AccessToken
}
