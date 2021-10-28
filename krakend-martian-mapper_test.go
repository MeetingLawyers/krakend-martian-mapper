package mapper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestMap(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "map_fields": { "user_id": "token" }}`
	requestBody := map[string]string{"user_id": "John Doe", "occupation": "gardener"}
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	json_data, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer(json_data))

	// Act
	modifier.RequestModifier().ModifyRequest(req)

	// Assert
	bodyBytes, _ := io.ReadAll(req.Body)
	expectedBody := `{"occupation":"gardener","token":"John Doe"}`
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected output <%s> different than obtained <%s>", expectedBody, string(bodyBytes))
	}
	expectedQuery := "irrelevant=false&token=12345"
	if string(req.URL.RawQuery) != expectedQuery {
		t.Errorf("Expected query <%s> different than obtained <%s>", expectedQuery, string(req.URL.RawQuery))
	}
}

func TestCopy(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "copy_fields": { "user_id": "token" }}`
	requestBody := map[string]string{"user_id": "John Doe", "occupation": "gardener"}
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	json_data, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer(json_data))

	// Act
	modifier.RequestModifier().ModifyRequest(req)

	// Assert
	bodyBytes, _ := io.ReadAll(req.Body)
	expectedBody := `{"occupation":"gardener","token":"John Doe","user_id":"John Doe"}`
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected output <%s> different than obtained <%s>", expectedBody, string(bodyBytes))
	}
	expectedQuery := "irrelevant=false&token=12345&user_id=12345"
	if string(req.URL.RawQuery) != expectedQuery {
		t.Errorf("Expected query <%s> different than obtained <%s>", expectedQuery, string(req.URL.RawQuery))
	}
}

func TestBothMapOverrides(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "copy_fields": { "user_id": "token", "occupation": "job" }, "map_fields": { "user_id": "token" }}`
	requestBody := map[string]string{"user_id": "John Doe", "occupation": "gardener"}
	url := "http://example.com?irrelevant=false&user_id=12345&occupation=writer"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	json_data, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer(json_data))

	// Act
	modifier.RequestModifier().ModifyRequest(req)

	// Assert
	bodyBytes, _ := io.ReadAll(req.Body)
	expectedBody := `{"job":"gardener","occupation":"gardener","token":"John Doe"}`
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected output <%s> different than obtained <%s>", expectedBody, string(bodyBytes))
	}
	expectedQuery := "irrelevant=false&job=writer&occupation=writer&token=12345"
	if string(req.URL.RawQuery) != expectedQuery {
		t.Errorf("Expected query <%s> different than obtained <%s>", expectedQuery, string(req.URL.RawQuery))
	}
}