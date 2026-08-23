package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandlerUsesEmbeddedTemplate(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	indexHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), "Business Loan Calculator") {
		t.Fatalf("response did not contain the embedded index template")
	}
}

func TestMortgageCalculation(t *testing.T) {
	form := strings.NewReader("amount=300000&rate=6&years=30")
	req := httptest.NewRequest(http.MethodPost, "/mortgage", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	mortgageHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), "$1,799") {
		t.Fatalf("response did not contain the expected monthly payment")
	}
}
