package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bonearadu/kvstore/kv_store"
)

// MockStore is a mock implementation of the KeyValueStore interface for testing
type MockStore struct {
	GetFunc     func(key string) (string, error)
	PutFunc     func(key string, value string) error
	DeleteFunc  func(key string) error
	EntriesFunc func() ([]kv_store.Entry, error)
}

func (m *MockStore) Get(key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return "", nil
}

func (m *MockStore) Put(key string, value string) error {
	if m.PutFunc != nil {
		return m.PutFunc(key, value)
	}
	return nil
}

func (m *MockStore) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}

func (m *MockStore) Entries() ([]kv_store.Entry, error) {
	if m.EntriesFunc != nil {
		return m.EntriesFunc()
	}
	return []kv_store.Entry{}, nil
}

// TestExtractKey tests the extractKey function
func TestExtractKey(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple key",
			path:     "/keys/mykey",
			expected: "mykey",
		},
		{
			name:     "key with slashes",
			path:     "/keys/path/to/mykey",
			expected: "path/to/mykey",
		},
		{
			name:     "empty key",
			path:     "/keys/",
			expected: "",
		},
		{
			name:     "no key",
			path:     "/keys",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			got := extractKey(req)
			if got != tt.expected {
				t.Errorf("extractKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestRouting tests that requests are routed to the correct handlers
func TestRouting(t *testing.T) {
	// Create a mock store
	mockStore := &MockStore{}

	// Create a handler with the mock store
	handler := NewHandler(mockStore)

	// Test cases
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /keys",
			method:         "GET",
			path:           "/keys",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /keys/mykey",
			method:         "GET",
			path:           "/keys/mykey",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT /keys/mykey",
			method:         "PUT",
			path:           "/keys/mykey",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE /keys/mykey",
			method:         "DELETE",
			path:           "/keys/mykey",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /keys (not implemented)",
			method:         "POST",
			path:           "/keys",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the given method and path
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
		})
	}
}

// TestHandleGetKey tests the handleGetKey function
func TestHandleGetKey(t *testing.T) {
	// Create a mock store
	mockStore := &MockStore{
		GetFunc: func(key string) (string, error) {
			// This is just a stub - the actual handler is not implemented yet
			return "", nil
		},
	}

	// Create a handler with the mock store
	handler := NewHandler(mockStore)

	// Create a request
	req := httptest.NewRequest("GET", "/keys/mykey", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

// TestHandlePutKey tests the handlePutKey function
func TestHandlePutKey(t *testing.T) {
	// Create a mock store
	mockStore := &MockStore{
		PutFunc: func(key string, value string) error {
			// This is just a stub - the actual handler is not implemented yet
			return nil
		},
	}

	// Create a handler with the mock store
	handler := NewHandler(mockStore)

	// Create a request
	req := httptest.NewRequest("PUT", "/keys/mykey", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

// TestHandleDeleteKey tests the handleDeleteKey function
func TestHandleDeleteKey(t *testing.T) {
	// Create a mock store
	mockStore := &MockStore{
		DeleteFunc: func(key string) error {
			// This is just a stub - the actual handler is not implemented yet
			return nil
		},
	}

	// Create a handler with the mock store
	handler := NewHandler(mockStore)

	// Create a request
	req := httptest.NewRequest("DELETE", "/keys/mykey", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

// TestHandleListEntries tests the handleListEntries function
func TestHandleListEntries(t *testing.T) {
	// Create a mock store
	mockStore := &MockStore{
		EntriesFunc: func() ([]kv_store.Entry, error) {
			// Return some sample entries
			return []kv_store.Entry{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			}, nil
		},
	}

	// Create a handler with the mock store
	handler := NewHandler(mockStore)

	// Create a request
	req := httptest.NewRequest("GET", "/keys", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}

	// Parse the response body
	var responseList []kv_store.Entry
	err := json.Unmarshal(rr.Body.Bytes(), &responseList)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Check the response content
	expectedList := []kv_store.Entry{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	if len(responseList) != len(expectedList) {
		t.Errorf("handler returned wrong number of entries: got %v want %v",
			len(responseList), len(expectedList))
	}

	// Create a map of expected entries for easier lookup
	expectedMap := make(map[string]string)
	for _, entry := range expectedList {
		expectedMap[entry.Key] = entry.Value
	}

	// Check that each response entry matches an expected entry
	for _, entry := range responseList {
		expectedValue, exists := expectedMap[entry.Key]
		if !exists {
			t.Errorf("Unexpected key in response: %s", entry.Key)
		} else if entry.Value != expectedValue {
			t.Errorf("For key %s, expected value %s, got %s", entry.Key, expectedValue, entry.Value)
		}
	}
}
