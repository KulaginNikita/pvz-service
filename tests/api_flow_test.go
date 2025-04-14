package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func Test_FullFlow_CreateAndCloseReception(t *testing.T) {
	client := &http.Client{}

	moderatorToken := getToken(t, "moderator")
	employeeToken := getToken(t, "employee")


	var pvzID uuid.UUID
	t.Run("create_pvz", func(t *testing.T) {
		body := map[string]any{
			"city": "Казань",
		}
		resp := post(t, client, baseURL+"/pvz", body, http.StatusCreated, moderatorToken)
		defer resp.Body.Close()

		var result struct {
			Id string `json:"id"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		pvzID, _ = uuid.Parse(result.Id)
	})


	t.Run("create_reception", func(t *testing.T) {
		body := map[string]any{
			"pvzId": pvzID,
		}
		resp := post(t, client, baseURL+"/receptions", body, http.StatusCreated, employeeToken)
		defer resp.Body.Close()

		var result struct {
			Id string `json:"id"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		_,_ = uuid.Parse(result.Id)
	})

	t.Run("add_50_products", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			body := map[string]any{
				"pvzId": pvzID, // ✅ корректный pvzId
				"type":  "одежда",
			}
			resp := post(t, client, baseURL+"/products", body, http.StatusCreated, employeeToken)
			resp.Body.Close()
		}
	})

	t.Run("close_reception", func(t *testing.T) {
		url := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID.String())
		resp := post(t, client, url, nil, http.StatusOK, employeeToken)
		resp.Body.Close()
	})
}

func getToken(t *testing.T, role string) string {
	client := &http.Client{}
	body := map[string]string{
		"role": role,
	}
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(body))

	req, err := http.NewRequest(http.MethodPost, baseURL+"/dummyLogin", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var token string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&token))
	return token
}

func post(t *testing.T, client *http.Client, url string, body any, expectedStatus int, token string) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected %d, got %d", expectedStatus, resp.StatusCode)
	}

	return resp
}
