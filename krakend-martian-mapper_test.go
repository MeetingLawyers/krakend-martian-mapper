package mapper

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestEmptyBody(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "map_fields": { "user_id": "token" }}`
	requestBody := ""
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer( []byte(requestBody) ))

	// Act
	modifier.RequestModifier().ModifyRequest(req)

	// Assert
	bodyBytes, _ := io.ReadAll(req.Body)
	expectedBody := ``
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected output <%s> different than obtained <%s>", expectedBody, string(bodyBytes))
	}
	expectedQuery := "irrelevant=false&token=12345"
	if string(req.URL.RawQuery) != expectedQuery {
		t.Errorf("Expected query <%s> different than obtained <%s>", expectedQuery, string(req.URL.RawQuery))
	}
}

func TestNotValidJsonBody(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "map_fields": { "user_id": "token" }}`
	requestBody := `{"invalidjson": {[}]}`
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer( []byte(requestBody) ))

	// Act
	modifier.RequestModifier().ModifyRequest(req)

	// Assert
	bodyBytes, _ := io.ReadAll(req.Body)
	expectedBody := ``
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected output <%s> different than obtained <%s>", expectedBody, string(bodyBytes))
	}
	expectedQuery := "irrelevant=false&token=12345"
	if string(req.URL.RawQuery) != expectedQuery {
		t.Errorf("Expected query <%s> different than obtained <%s>", expectedQuery, string(req.URL.RawQuery))
	}
}

func TestMap(t *testing.T) {
	// Arrange
	cfg := `{"source": "request", "map_fields": { "user_id": "token" }}`
	requestBody := `{"user_id": "John Doe", "occupation": "gardener"}`
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer( []byte(requestBody) ))

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
	requestBody := `{"user_id": "John Doe", "occupation": "gardener"}`
	url := "http://example.com?irrelevant=false&user_id=12345"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer( []byte(requestBody) ))

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
	requestBody := `{"user_id": "John Doe", "occupation": "gardener"}`
	url := "http://example.com?irrelevant=false&user_id=12345&occupation=writer"
	requestType := "POST"

	modifier, _ := MapperFromJSON([]byte(cfg))
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer( []byte(requestBody) ))

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