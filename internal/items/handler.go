package items

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/bhumi46/payment-svc/internal/httpx"
)

// Handler serves /api/items.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Register adds the items routes to mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/items", h.list)
	mux.HandleFunc("GET /api/items/{id}", h.get)
	mux.HandleFunc("POST /api/items", h.create)
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, h.store.List())
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	it, err := h.store.Get(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, it)
}

type createRequest struct {
	Name string `json:"name"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httpx.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, h.store.Create(Item{Name: req.Name}))
}
