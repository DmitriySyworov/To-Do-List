package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-list/app/internal/model"
	"to-do-list/app/internal/test"
	"to-do-list/app/internal/user"
)

var caseUpdateMyUserTest = []struct {
	NameTest      string
	RequestUpdate *user.RequestUpdateUser
}{
	{NameTest: "Update_All_Without_Name", RequestUpdate: &user.RequestUpdateUser{OriginalPassword: test.OrigUserPasswordTest, Email: "new1@gmail.com", NewPassword: "N)A(sunaijxc0879w"}},
	{NameTest: "Update_Password", RequestUpdate: &user.RequestUpdateUser{OriginalPassword: test.OrigUserPasswordTest, NewPassword: "PZ<PXPOMX9"}},
	{NameTest: "Update_Email", RequestUpdate: &user.RequestUpdateUser{OriginalPassword: test.OrigUserPasswordTest, Email: "new1234@gmail.com"}},
	{NameTest: "Update_Name", RequestUpdate: &user.RequestUpdateUser{Name: "Aboba"}},
}

func TestUpdateMyUserSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	defer dbTest.CleanAllDb()
	for _, tests := range caseUpdateMyUserTest {
		dbTest.CreateUserTest(&test.BaseUserTest)
		data, errMarshal := json.Marshal(&tests.RequestUpdate)
		if errMarshal != nil {
			t.Fatalf("%s :  %v", tests.NameTest, errMarshal)
		}
		request, errReq := http.NewRequest(http.MethodPatch, "/my-user", bytes.NewBuffer(data))
		if errReq != nil {
			t.Fatal(errReq)
		}
		request.Header.Set("Authorization", test.IdJWTTest)
		writer := httptest.NewRecorder()
		App().ServeHTTP(writer, request)
		if writer.Code != http.StatusOK {
			t.Fatalf("%s : expected %d got %d", tests.NameTest, http.StatusOK, writer.Code)
		}
		bodyData, errRead := io.ReadAll(writer.Body)
		if errRead != nil {
			t.Fatal(errRead)
		}
		var payload model.User
		errUnmarshal := json.Unmarshal(bodyData, &payload)
		if errUnmarshal != nil {
			t.Fatal(errUnmarshal)
		}
		if payload.UserId == 0 {
			t.Fatal("response empty")
		}
		t.Log(payload)
		dbTest.CleanAllDb()
	}
}
func TestGetMyUserSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	defer dbTest.CleanAllDb()
	request, errReq := http.NewRequest(http.MethodGet, "/my-user", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.IdJWTTest)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	bodyData, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload model.User
	errUnmarshal := json.Unmarshal(bodyData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.UserId == 0 {
		t.Fatal("response empty")
	}
	t.Log(payload)
}
func TestDeleteMyUserSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	defer dbTest.CleanAllDb()
	data, errMarshal := json.Marshal(&user.RequestDeleteUser{
		Password: test.OrigUserPasswordTest,
	})
	if errMarshal != nil {
		t.Fatal(errMarshal)
	}
	request, errReq := http.NewRequest(http.MethodDelete, "/my-user", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.IdJWTTest)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusNoContent {
		t.Fatalf("expected %d got %d", http.StatusNoContent, writer.Code)
	}
}
