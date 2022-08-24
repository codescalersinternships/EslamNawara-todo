package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	errInvalidInput = "Bad request, invalid input"
	errNotFound     = "Task not found"
	DB_PATH         = "gorm.db"
	URL             = "http://localhost:8080"
)

func GetTasks(writer http.ResponseWriter, request *http.Request) {
	var tasks []Task
	db.Find(&tasks)
	jsonData, _ := json.MarshalIndent(tasks, "", "  ")

	writer.Header().Set("Access-Control-Allow-Origin", URL)
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func AddTask(writer http.ResponseWriter, request *http.Request) {
	var newTask Task
	if json.NewDecoder(request.Body).Decode(&newTask) != nil || validateData(newTask) != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	if db.First(&Task{}, newTask.ID).Error == nil {
		http.Error(writer, "Bad request, task already exist", http.StatusConflict)
		return
	}
	db.Create(&newTask)
	jsonData, _ := json.MarshalIndent(newTask, "", "  ")

	writer.Header().Set("Access-Control-Allow-Origin", URL)
	writer.Header().Add("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE")
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(http.StatusCreated)
	writer.Write(jsonData)
}

func GetTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	var requestedTask Task
	if db.First(&requestedTask, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusNotFound)
		return
	}
	jsonData, _ := json.MarshalIndent(requestedTask, "", "  ")
	writer.Header().Set("Access-Control-Allow-Origin", URL)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(jsonData)

}

func DeleteTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	if db.Delete(&Task{}, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write([]byte("Task deleted successfully"))
}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task Task
	if json.NewDecoder(request.Body).Decode(&task) != nil || validateData(task) != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	id := task.ID
	item := task.Item
	completed := task.Completed

	updatedTask := Task{}
	if db.First(&updatedTask, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusConflict)
		return
	}
	updatedTask.Item = item
	updatedTask.Completed = completed
	db.Save(updatedTask)

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	writer.Header().Add("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE")
	writer.Header().Set("content-type", "application/json")

	jsonData, _ := json.MarshalIndent(updatedTask, "", "  ")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(jsonData)

}

func validateData(task Task) error {
	if task.Item == "" || task.ID <= 0 {
		return fmt.Errorf("Invalid task")
	}
	return nil
}
