package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/model"
	"to-do-list/app/internal/test"
)

func TestRegisterSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	defer dbTest.CleanAllDb()
	body := &auth.RequestRegister{
		Name:     test.UserNameTest,
		Password: test.OrigUserPasswordTest,
		Email:    test.UserEmailTest,
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
	respData, errResp := io.ReadAll(writer.Body)
	if errResp != nil {
		t.Fatal(errResp)
	}
	var payload auth.ResponseAuth
	errUnmarshal := json.Unmarshal(respData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.JWT == "" {
		t.Fatal("jwt empty")
	}
	t.Log(payload.JWT)
}
func TestLoginSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	defer dbTest.CleanAllDb()
	body := &auth.RequestRegister{
		Password: test.OrigUserPasswordTest,
		Email:    test.UserEmailTest,
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
	respData, errResp := io.ReadAll(writer.Body)
	if errResp != nil {
		t.Fatal(errResp)
	}
	var payload auth.ResponseAuth
	errUnmarshal := json.Unmarshal(respData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.JWT == "" {
		t.Fatal("jwt empty")
	}
	t.Log(payload.JWT)
}
func TestConfirmSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateTempUserTest(&model.TempUser{
		Name:     test.UserNameTest,
		Email:    test.UserEmailTest,
		Password: test.HashUserPasswordTest,
		UserId:   test.UserIDTest,
	}, test.IdHashTest)
	dbTest.CreateSessionTest(&model.Session{
		SessionId:     test.SessionIdTest,
		TemporaryCode: test.TempCodeTest,
	}, test.IdHashTest)
	userToken := dbTest.CreateTemporaryJWTTest(float64(test.IdHashTest), test.SessionIdTest)
	defer dbTest.CleanAllDb()
	body := &auth.RequestConfirm{
		TempCode: test.TempCodeTest,
	}
	data, errJs := json.Marshal(&body)
	if errJs != nil {
		t.Fatal(errJs)
	}
	request, errReq := http.NewRequest(http.MethodPost, "/auth/confirm", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("X-User-Token", "Bearer "+userToken)
	query := request.URL.Query()
	query.Set("action", "register")
	request.URL.RawQuery = query.Encode()
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	respData, errResp := io.ReadAll(writer.Body)
	if errResp != nil {
		t.Fatal(errResp)
	}
	var payload auth.ResponseConfirm
	errUnmarshal := json.Unmarshal(respData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.JWT == "" {
		t.Fatal("jwt empty")
	}
	t.Log(payload.JWT)
}
