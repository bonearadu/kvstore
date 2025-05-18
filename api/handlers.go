package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bonearadu/kvstore/kv_store"
)

// Handler handles HTTP requests for the key-value store
type Handler struct {
	store kv_store.KeyValueStore
	mux   *http.ServeMux
}

// NewHandler creates a new Handler with the given store
func NewHandler(store kv_store.KeyValueStore) *Handler {
	h := &Handler{
		store: store,
		mux:   http.NewServeMux(),
	}

	// Register routes
	h.registerRoutes()

	return h
}

// registerRoutes sets up all the routes for the API
func (h *Handler) registerRoutes() {
	// GET /keys - List all entries (key-value pairs)
	h.mux.HandleFunc("GET /keys", h.handleListEntries)

	// GET /keys/{key} - Get a specific key
	h.mux.HandleFunc("GET /keys/", h.handleGetKey)

	// PUT /keys/{key} - Create or update a key
	h.mux.HandleFunc("PUT /keys/", h.handlePutKey)

	// DELETE /keys/{key} - Delete a specific key
	h.mux.HandleFunc("DELETE /keys/", h.handleDeleteKey)
}

// ServeHTTP delegates to the internal mux
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// extractKey extracts the key from the URL path
func extractKey(r *http.Request) string {
	// The pattern is registered as "/keys/"
	// So we need to extract everything after "/keys/"
	path := r.URL.Path
	parts := strings.Split(path, "/keys/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// handleGetKey handles GET requests for a specific key
func (h *Handler) handleGetKey(w http.ResponseWriter, r *http.Request) {
	// Extract key from path
	key := extractKey(r)

	// Get the value from the store
	value, err := h.store.Get(key)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// handlePutKey handles PUT requests to create or update a key
func (h *Handler) handlePutKey(w http.ResponseWriter, r *http.Request) {
	// Extract key from path
	key := extractKey(r)

	// Read the value from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Check if the key already exists to determine the status code
	_, err = h.store.Get(key)
	isNewKey := err != nil

	// Store the key-value pair
	err = h.store.Put(key, string(body))
	if err != nil {
		http.Error(w, "Failed to store value", http.StatusInternalServerError)
		return
	}

	// Set appropriate status code
	if isNewKey {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// handleDeleteKey handles DELETE requests for a specific key
func (h *Handler) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
	// Extract key from path
	key := extractKey(r)

	// Delete the key from the store
	err := h.store.Delete(key)
	if err != nil {
		http.Error(w, "Failed to delete key", http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusOK)
}

// handleListEntries handles GET requests to list all key-value pairs
func (h *Handler) handleListEntries(w http.ResponseWriter, r *http.Request) {
	// Get all entries from the store
	entries, err := h.store.Entries()
	if err != nil {
		http.Error(w, "Failed to list entries", http.StatusInternalServerError)
		return
	}

	// Marshal entries map to JSON
	jsonData, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, "Failed to marshal entries", http.StatusInternalServerError)
		return
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
