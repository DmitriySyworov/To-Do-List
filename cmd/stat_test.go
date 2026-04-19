package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-list/app/internal/stat"
	"to-do-list/app/internal/test"
)

func TestGetMyStat(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateTestUser(&test.BaseTestUser)
	dbTest.CreateStat(test.UserTestID, test.UserTestName)
	defer dbTest.CleanAllDb()
	request, errReq := http.NewRequest(http.MethodGet, "/user/my-stat", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.JWTTestToken)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}

	data, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload stat.ResponseMyStat
	errJs := json.Unmarshal(data, &payload)
	if errJs != nil {
		t.Fatal(errJs)
	}
	if payload.DoneTask == "" && payload.DeleteTask == "" && payload.ActiveTask == "" {
		t.Fatal("empty response")
	}
}
