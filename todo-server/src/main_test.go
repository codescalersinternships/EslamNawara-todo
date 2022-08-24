package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

const (
	task1        = `{"id": 1, "item": "clean the room", "completed": false}`
	task2        = `{"id": 2, "item": "do the assignment", "completed": false}`
	task3        = `{"id": 3, "item": "work out", "completed": false}`
	task4        = `{"id": 4, "item": "make the bed", "completed": false}`
	invalidTask  = `{"id": "5", "item": "make the bed", "completed": false}`
)

func TestGetTasks(t *testing.T) {
	t.Run("Empty database", func(t *testing.T) {
		defer os.Remove(DB_PATH)
		OpenDB()
		request := httptest.NewRequest(http.MethodGet, "localhost:5000/todo", nil)
		response := httptest.NewRecorder()
		GetTasks(response, request)
		got := response.Body.String()
		want := "[]"
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("Not empty db", func(t *testing.T) {
		defer os.Remove(DB_PATH)
		OpenDB()
		fillDB()
		request := httptest.NewRequest(http.MethodGet, "localhost:5000/todo", nil)
		response := httptest.NewRecorder()
		GetTasks(response, request)
		got := make([]Task, 0)
		json.NewDecoder(response.Body).Decode(&got)
		want := []Task{{ID: 1, Item: "clean the room", Completed: false},
			{ID: 2, Item: "do the assignment", Completed: false},
			{ID: 3, Item: "work out", Completed: false},
			{ID: 4, Item: "make the bed", Completed: false}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}
func TestAddTask(t *testing.T) {
	t.Run("Get existing task ", func(t *testing.T) {
		defer os.Remove(DB_PATH)
		OpenDB()
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer([]byte(task1)))
		response := httptest.NewRecorder()
		AddTask(response, request)
		got := Task{}
		json.NewDecoder(response.Body).Decode(&got)
		want := Task{ID: 1, Item: "clean the room", Completed: false}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("Get invalid task", func(t *testing.T) {
		defer os.Remove(DB_PATH)
		OpenDB()
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo/2", bytes.NewBuffer([]byte(invalidTask)))
		response := httptest.NewRecorder()
		AddTask(response, request)
		got := response.Body.String()
		want := "Bad request, invalid input\n"
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

}

func fillDB() {
	request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer([]byte(task1)))
	response := httptest.NewRecorder()
	AddTask(response, request)

	request = httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer([]byte(task2)))
	response = httptest.NewRecorder()
	AddTask(response, request)

	request = httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer([]byte(task3)))
	response = httptest.NewRecorder()
	AddTask(response, request)

	request = httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer([]byte(task4)))
	response = httptest.NewRecorder()
	AddTask(response, request)
}
