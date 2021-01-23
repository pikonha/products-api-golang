package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := strings.TrimRight(response.Body.String(), "\n"); body != "[]" {
		t.Errorf("Expected an empty table. Got %s\n", body)
	}
}

func TestProductNotFound(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products/10", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if _, ok := m["error"]; !ok {
		t.Errorf("Expected error key on response, got: %#v", m)
	}
}


func TestInsertProduct(t *testing.T) {
	clearTable()

	productName := "chair"
	productPrice := 20.50
	productJSON := fmt.Sprintf(`{"name": %s, "price": %f}`, productName, productPrice)

	jsonBody := []byte(productJSON)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	response := executeRequest(req)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != productName {
		t.Errorf("Expected: %s, Got: %s", m["name"], productName)
	}

	if m["price"] != productPrice {
		t.Errorf("Expected: %f, Got: %f", m["price"], productPrice)
	}

	// This happens because json unmarshal converts int to float when using interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected: '1', Got: %v", m["id"])
	}
}
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d\n", expected, actual)
	}
}