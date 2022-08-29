package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var idMutex sync.Mutex

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

	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func AddTask(writer http.ResponseWriter, request *http.Request) {
	var newTask Task
	if json.NewDecoder(request.Body).Decode(&newTask) != nil || validateData(newTask) != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}

	idMutex.Lock()
	newTask.ID = nextId
	nextId++
	idMutex.Unlock()


	db.Create(&newTask)
	jsonData, _ := json.MarshalIndent(newTask, "", "  ")

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

	jsonData, _ := json.MarshalIndent(updatedTask, "", "  ")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(jsonData)

}

func validateData(task Task) error {
	if task.Item == "" || task.ID < 0 {
		return fmt.Errorf("Invalid task")
	}
	return nil
}
