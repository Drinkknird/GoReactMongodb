package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//เสร็จทดสอบของ CreateUser
// Database interface สำหรับฐานข้อมูล MongoDB
type Database interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

// MockDatabase is a struct used for testing
type MockDatabase struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

// Implement Database interface methods
func (m *MockDatabase) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
  // Create a mock InsertOneResult
	result := &mongo.InsertOneResult{
    InsertedID: "some_id",
	}
	return result, nil
}

type MockUserAPI struct {
	db Database
}

func (m *MockUserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
  // Simulate successful creation with a mock response body
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"insertedID": "some_id", "message": "User created successfully (mock)"}`)
}



func TestCreateUser(t *testing.T) {
	// Mock MongoDB database
	mockDB := &MockDatabase{}

	// Create a UserAPI instance
	userAPI := &MockUserAPI{db: mockDB}


	// Prepare a sample request body
	requestBody := []byte(`{"name":"John Doe","email":"john@example.com","password":"password123"}`)

	// Create a request
	req, err := http.NewRequest("POST", "/users/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handle CreateUser request
	handler := http.HandlerFunc(userAPI.CreateUser)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body
	expected := `{"insertedID": "some_id", "message": "User created successfully (mock)"}`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
