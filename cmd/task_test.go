package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-list/app/internal/model"
	"to-do-list/app/internal/task"
	"to-do-list/app/internal/test"
)

var caseCreateTaskTest = []struct {
	NameTest      string
	RequestCreate *task.RequestCreateTaskForm
}{
	{NameTest: "AllData", RequestCreate: &task.RequestCreateTaskForm{Header: test.HeaderTest, Task: test.TaskTest, Deadline: test.DateTest}},
	{NameTest: "NoDeadline", RequestCreate: &task.RequestCreateTaskForm{Header: test.HeaderTest, Task: test.TaskTest}},
	{NameTest: "NoHeader", RequestCreate: &task.RequestCreateTaskForm{Task: test.TaskTest, Deadline: test.DateTest}},
	{NameTest: "OnlyTask", RequestCreate: &task.RequestCreateTaskForm{Task: test.TaskTest}},
}

func TestCreateTaskSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	defer dbTest.CleanAllDb()
	for _, tests := range caseCreateTaskTest {
		dbTest.CreateUserTest(&test.BaseUserTest)
		data, errMarshal := json.Marshal(&tests.RequestCreate)
		if errMarshal != nil {
			t.Fatalf("%s: %v", tests.NameTest, errMarshal)
		}
		request, errReq := http.NewRequest(http.MethodPost, "/user/task", bytes.NewBuffer(data))
		if errReq != nil {
			t.Fatalf("%s: %v", tests.NameTest, errReq)
		}
		request.Header.Set("Authorization", test.IdJWTTest)
		writer := httptest.NewRecorder()
		App().ServeHTTP(writer, request)
		if writer.Code != http.StatusCreated {
			t.Fatalf("%s: expected %d got %d", tests.NameTest, http.StatusCreated, writer.Code)
		}
		bodyData, errRead := io.ReadAll(writer.Body)
		if errRead != nil {
			t.Fatalf("%s: %v", tests.NameTest, errRead)
		}
		var payload model.TaskForm
		errUnmarshal := json.Unmarshal(bodyData, &payload)
		if errUnmarshal != nil {
			t.Fatalf("%s: %v", tests.NameTest, errUnmarshal)
		}
		if payload.TaskId == 0 {
			t.Fatalf("%s: taskId empty", tests.NameTest)
		}
		t.Log(payload)
		dbTest.CleanAllDb()
	}
}

var caseUpdateTaskTest = []struct {
	NameTest      string
	RequestUpdate *task.RequestUpdateTaskForm
}{
	{NameTest: "AllUpdate", RequestUpdate: &task.RequestUpdateTaskForm{Header: "GET UP!", Task: "i need ro work", Deadline: "2027-05-12"}},
	{NameTest: "OnlyTask", RequestUpdate: &task.RequestUpdateTaskForm{Task: "i need ro work"}},
	{NameTest: "OnlyHeader", RequestUpdate: &task.RequestUpdateTaskForm{Header: "GET UP!"}},
	{NameTest: "OnlyDeadline", RequestUpdate: &task.RequestUpdateTaskForm{Deadline: "2027-05-12"}},
	{NameTest: "HeaderAndTask", RequestUpdate: &task.RequestUpdateTaskForm{Header: "GET UP!", Task: "i need ro work"}},
	{NameTest: "TaskAndDeadline", RequestUpdate: &task.RequestUpdateTaskForm{Task: "i need ro work", Deadline: "2027-05-12"}},
	{NameTest: "HeaderAndDeadline", RequestUpdate: &task.RequestUpdateTaskForm{Header: "GET UP!", Deadline: "2027-05-12"}},
	{NameTest: "DoneTask", RequestUpdate: &task.RequestUpdateTaskForm{StatusDone: true}},
}

func TestUpdateTaskSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	defer dbTest.CleanAllDb()
	for _, tests := range caseUpdateTaskTest {
		dbTest.CreateUserTest(&test.BaseUserTest)
		dbTest.CreateTaskTest(&test.BaseTaskTest)
		data, errMarshal := json.Marshal(&tests.RequestUpdate)
		if errMarshal != nil {
			t.Fatalf("%s: %v", tests.NameTest, errMarshal)
		}
		request, errReq := http.NewRequest(http.MethodPatch, "/user/task/"+fmt.Sprint(test.TaskIdTest), bytes.NewBuffer(data))
		if errReq != nil {
			t.Fatalf("%s: %v", tests.NameTest, errReq)
		}
		request.Header.Set("Authorization", test.IdJWTTest)
		writer := httptest.NewRecorder()
		App().ServeHTTP(writer, request)
		if writer.Code != http.StatusOK {
			t.Fatalf("%s: expected %d got %d", tests.NameTest, http.StatusOK, writer.Code)
		}
		bodyData, errRead := io.ReadAll(writer.Body)
		if errRead != nil {
			t.Fatalf("%s: %v", tests.NameTest, errRead)
		}
		var payload model.TaskForm
		errUnmarshal := json.Unmarshal(bodyData, &payload)
		if errUnmarshal != nil {
			t.Fatalf("%s: %v", tests.NameTest, errUnmarshal)
		}
		if payload.TaskId == 0 {
			t.Fatalf("%s: taskId empty", tests.NameTest)
		}
		t.Log(payload)
		dbTest.CleanAllDb()
	}
}
func TestDeleteTaskSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	dbTest.CreateTaskTest(&test.BaseTaskTest)
	defer dbTest.CleanAllDb()
	request, errReq := http.NewRequest(http.MethodDelete, "/user/task/"+fmt.Sprint(test.TaskIdTest), nil)
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
func TestGetTaskSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	dbTest.CreateTaskTest(&test.BaseTaskTest)
	defer dbTest.CleanAllDb()
	request, errReq := http.NewRequest(http.MethodGet, "/user/task/"+fmt.Sprint(test.TaskIdTest), nil)
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
	var payload model.TaskForm
	errUnmarshal := json.Unmarshal(bodyData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.TaskId == 0 {
		t.Fatal("response empty")
	}
	t.Log(payload)
}
func TestGetAllMyTasksSuccess(t *testing.T) {
	dbTest := test.OpenAllTestDb()
	dbTest.CreateUserTest(&test.BaseUserTest)
	dbTest.CreateTaskTest(&test.BaseTaskTest)
	dbTest.CreateTaskTest(&model.TaskForm{
		Task:       "GOO",
		StatusDone: true,
		TaskId:     9876123,
		UserId:     test.UserIDTest,
	})
	defer dbTest.CleanAllDb()
	request, errReq := http.NewRequest(http.MethodGet, "/user/my-tasks", nil)
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
	var payload task.ResponseAllTasksPeriod
	errUnmarshal := json.Unmarshal(bodyData, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if len(payload.DoneTasks) == 0 && len(payload.ActiveTasks) == 0 {
		t.Fatal("response empty")
	}
	t.Log(payload)
}
