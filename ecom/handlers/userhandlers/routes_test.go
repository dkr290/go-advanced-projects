package userhandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/stretchr/testify/assert"
)

// MockmysqlDB is a mock implementation of the database interface
type mockMysqlDB struct{}

func TestHandleRegister_InvalidEmail(t *testing.T) {
	db := &mockMysqlDB{}
	handler := NewUserHandler(db)

	testCases := []struct {
		name           string
		payload        types.RegisterUserPayload
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid Email",
			payload: types.RegisterUserPayload{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "jo@example.com",
				Password:  "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Invalid Email",
			payload: types.RegisterUserPayload{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid payload",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			payloadBytes, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.handleRegister(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check the response body
			var responseBody map[string]string
			json.Unmarshal(rr.Body.Bytes(), &responseBody)

			if tc.expectedError != "" {
				assert.Contains(t, responseBody["error"], tc.expectedError)
			} else {
				assert.Empty(t, responseBody["error"])
			}

			// Print the actual response for debugging
			t.Logf("Response Status: %d", rr.Code)
			t.Logf("Response Body: %s", rr.Body.String())
		})
	}
}

func (m *mockMysqlDB) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}
func (m *mockMysqlDB) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockMysqlDB) CreateUser(user types.User) error {
	return nil
}
