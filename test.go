package main

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
)

func TestPing(t *testing.T) {
	env, _ := godotenv.Read(filepath.Join(".env"))
	app := setupApp(env)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	app.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
