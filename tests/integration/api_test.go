package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestAuth(t *testing.T) {
	t.Run("Successful authentication", func(t *testing.T) {
		payload := map[string]string{"username": "testuser", "password": "testpass"}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		if result["token"] == "" {
			t.Error("Token not received")
		}
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		payload := map[string]string{"username": "", "password": ""}
		body, _ := json.Marshal(payload)
		resp, _ := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(body))
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Authentication with wrong password", func(t *testing.T) {
		payload := map[string]string{"username": "testuser", "password": "wrongpass"}
		body, _ := json.Marshal(payload)
		resp, _ := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(body))
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Re-authentication with valid credentials", func(t *testing.T) {
		token1 := getAuthToken("testuser1", "testpass1")
		token2 := getAuthToken("testuser1", "testpass1")
		if token1 == "" || token2 == "" {
			t.Error("Token not received")
		}
		if token1 != token2 {
			t.Error("Tokens don't match")
		}
	})
}

func getAuthToken(username, password string) string {
	payload := map[string]string{"username": username, "password": password}
	body, _ := json.Marshal(payload)
	resp, _ := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(body))
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["token"]
}

func TestSendCoins(t *testing.T) {
	token := getAuthToken("testuser", "testpass")

	t.Run("Successful send", func(t *testing.T) {
		payload := map[string]interface{}{"toUser": "testuser1", "amount": 10}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Insufficient coins", func(t *testing.T) {
		payload := map[string]interface{}{"toUser": "testuser1", "amount": 9999}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Send negative amount", func(t *testing.T) {
		payload := map[string]interface{}{"toUser": "testuser1", "amount": -10}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Send coins to non-existent user", func(t *testing.T) {
		payload := map[string]interface{}{"toUser": "nonexistentuser", "amount": 10}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Trying send without token", func(t *testing.T) {
		payload := map[string]interface{}{"toUser": "testuser1", "amount": 10}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestBuyItem(t *testing.T) {
	token := getAuthToken("testuser", "testpass")

	t.Run("Successful purchase", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/buy/pen", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Insufficient coins", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/buy/sword", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Trying purchase without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/buy/pen", nil)
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestGetUserInfo(t *testing.T) {
	t.Run("Succesful info extraction", func(t *testing.T) {
		token := getAuthToken("testuser", "testpass")

		req, _ := http.NewRequest("GET", baseURL+"/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		if info["coins"] == nil || info["inventory"] == nil {
			log.Println(info["coins"], info["inventory"])
			t.Error("Invalid response")
		}
	})

	t.Run("Trying access to info without token", func(t *testing.T) {
		resp, _ := http.Get(baseURL + "/api/info")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
