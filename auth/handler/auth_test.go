package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createTestRouter(t *testing.T) *gin.Engine {
	router := gin.New()

	db := createTestDB(t)

	authHandler := NewAuthHandler(db)

	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)

	cleanUpTestDB(t, db)

	t.Cleanup(func() {
		closeTestDB(t, db)
	})

	return router
}

func createTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DB_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func cleanUpTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS users CASCADE")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Database dropped table ")
}

func closeTestDB(t *testing.T, db *sql.DB) {
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegister(t *testing.T) {
	r := createTestRouter(t)

	body := map[string]string{
		"username": "testuser",
		"email":    "zakhar@example.com",
		"password": "testpassword",
		"fullname": "Zakhar Sln",
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		t.Errorf("error to create json body: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost, "/api/auth/register",
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status code not 201: %v", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["user"] == nil {
		t.Fatalf("user not found in response")
	}

	t.Logf("Register test success")
}
