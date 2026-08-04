package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Veter-ok/MarkTogether/internal/document"
)

type DocumentHandler struct {
	store document.Store
	mutex *sync.RWMutex
}

func NewDocumentHandler(store document.Store) *DocumentHandler {
	return &DocumentHandler{
		store: store,
		mutex: &sync.RWMutex{},
	}
}

func (docHandler *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	doc := document.NewDocument(req.Name)
	if err := docHandler.store.Save(doc); err != nil {
		http.Error(w, "Failed to save document: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":   doc.ID,
		"name": doc.Title,
	})
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.store.Get(id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if doc == nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}
