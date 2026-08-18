package player

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type fakeRepository struct {
	createPlayer Player
	createError  error
	findPlayer   Player
	findError    error
}

func (repository fakeRepository) Create(context.Context, string) (Player, error) {
	return repository.createPlayer, repository.createError
}

func (repository fakeRepository) FindByID(context.Context, string) (Player, error) {
	return repository.findPlayer, repository.findError
}

func TestCreate(t *testing.T) {
	expectedPlayer := Player{
		ID:        "21faea44-59df-4ac7-b75b-c82d51a3de00",
		Name:      "sumom",
		CreatedAt: time.Date(2026, time.August, 18, 0, 0, 0, 0, time.UTC),
	}
	handler := NewHandler(fakeRepository{createPlayer: expectedPlayer})

	request := httptest.NewRequest(http.MethodPost, "/players", strings.NewReader(`{"name":" sumom "}`))
	recorder := httptest.NewRecorder()

	handler.Create(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("ステータスコード = %d, want %d", response.StatusCode, http.StatusCreated)
	}

	var body Player
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("応答JSONの読み取りに失敗しました: %v", err)
	}

	if body != expectedPlayer {
		t.Fatalf("応答 = %#v, want %#v", body, expectedPlayer)
	}
}

func TestCreateReturnsBadRequestForEmptyName(t *testing.T) {
	handler := NewHandler(fakeRepository{})

	request := httptest.NewRequest(http.MethodPost, "/players", strings.NewReader(`{"name":"  "}`))
	recorder := httptest.NewRecorder()

	handler.Create(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("ステータスコード = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetReturnsNotFound(t *testing.T) {
	handler := NewHandler(fakeRepository{findError: ErrNotFound})
	router := chi.NewRouter()
	router.Get("/players/{id}", handler.Get)

	request := httptest.NewRequest(http.MethodGet, "/players/21faea44-59df-4ac7-b75b-c82d51a3de00", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("ステータスコード = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestGetReturnsBadRequestForInvalidID(t *testing.T) {
	handler := NewHandler(fakeRepository{})
	router := chi.NewRouter()
	router.Get("/players/{id}", handler.Get)

	request := httptest.NewRequest(http.MethodGet, "/players/not-a-uuid", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("ステータスコード = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetReturnsInternalServerError(t *testing.T) {
	handler := NewHandler(fakeRepository{findError: errors.New("database error")})
	router := chi.NewRouter()
	router.Get("/players/{id}", handler.Get)

	request := httptest.NewRequest(http.MethodGet, "/players/21faea44-59df-4ac7-b75b-c82d51a3de00", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("ステータスコード = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
