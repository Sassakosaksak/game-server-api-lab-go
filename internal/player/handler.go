package player

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repository Repository
}

type createPlayerRequest struct {
	Name string `json:"name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	name := strings.TrimSpace(request.Name)
	if length := utf8.RuneCountInString(name); length == 0 || length > 30 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_name"})
		return
	}

	createdPlayer, err := handler.repository.Create(r.Context(), name)
	if err != nil {
		log.Printf("Playerの作成に失敗しました: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}

	writeJSON(w, http.StatusCreated, createdPlayer)
}

func (handler *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_player_id"})
		return
	}

	foundPlayer, err := handler.repository.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "player_not_found"})
			return
		}

		log.Printf("Playerの取得に失敗しました: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}

	writeJSON(w, http.StatusOK, foundPlayer)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("JSON応答の書き込みに失敗しました: %v", err)
	}
}
