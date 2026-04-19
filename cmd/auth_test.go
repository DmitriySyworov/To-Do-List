package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/test"
)

func TestRegisterSuccess(t *testing.T) {
	body := &auth.RequestRegister{
		Name:     test.UserTestName,
		Password: test.OrigTestUserPassword,
		Email:    test.UserTestEmail,
	}
	data, errJs := json.Marshal(&body)
	if errJs != nil {
		t.Fatal(errJs)
	}
	request, errReq := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
}
func TestLoginSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateTestUser(&test.BaseTestUser)
	defer dbTest.CleanAllDb()
	body := &auth.RequestRegister{
		Password: test.OrigTestUserPassword,
		Email:    test.UserTestEmail,
	}
	data, errJs := json.Marshal(&body)
	if errJs != nil {
		t.Fatal(errJs)
	}
	request, errReq := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
}
