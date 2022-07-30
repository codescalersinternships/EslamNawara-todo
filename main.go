package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID        int64  `gorm:"NOT NULL" json:"id"`
	Item      string `json:"item"`
	Completed bool   `jon:"completed"`
}

const (
	errInvalidInput = "Bad request, invalid input"
	errNotFound     = "Task not found"
	LISTEN_PORT     = ":5000"
	DB_PATH         = "gorm.db"
)

func main() {

	if OpenDB("DB_FILE") != nil {
		fmt.Println("Error: Failed to connect to the database")
		return
	}
	router := mux.NewRouter()
	SetRoutes(router)
	http.ListenAndServe(LISTEN_PORT, router)
}

func OpenDB(path string) error {
	var err error
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	db.AutoMigrate(&Task{})
	return err
}

func SetRoutes(router *mux.Router) {
	router.HandleFunc("/todo", GetTasks).Methods("GET")
	router.HandleFunc("/todo", AddTask).Methods("POST")
	router.HandleFunc("/todo", UpdateTask).Methods("PATCH")
	router.HandleFunc("/todo/{id}", GetTask).Methods("GET")
	router.HandleFunc("/todo/{id}", CompleteTask).Methods("PATCH")
	router.HandleFunc("/todo/{id}", DeleteTask).Methods("DELETE")
}

func GetTasks(writer http.ResponseWriter, request *http.Request) {
	var tasks []Task
	db.Find(&tasks)
	jsonData, _ := json.MarshalIndent(tasks, "", "  ")
	request.Header.Set("Content-Type", "application/json")
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
	request.Header.Set("Content-Type", "application/json")
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
	request.Header.Set("Content-Type", "application/json")
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
	request.Header.Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write([]byte("Task deleted successfully"))
}

func CompleteTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}

	checkedTask := Task{}
	if db.First(&checkedTask, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusConflict)
		return
	}
	checkedTask.Completed = !checkedTask.Completed
	db.Save(checkedTask)
	request.Header.Set("Content-Type", "application/json")
	jsonData, _ := json.MarshalIndent(checkedTask, "", "  ")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(jsonData)

}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task Task
	if json.NewDecoder(request.Body).Decode(&task) != nil || validateData(task) != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	id := task.ID
	item := task.Item

	updatedTask := Task{}
	if db.First(&updatedTask, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusConflict)
		return
	}
	updatedTask.Item = item
	db.Save(updatedTask)
	request.Header.Set("Content-Type", "application/json")
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
