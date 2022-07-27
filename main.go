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

type task struct {
	ID        string `gorm:"NOT NULL" json:"id"`
	Item      string `json:"item"`
	Completed bool   `jon:"completed"`
}

const (
	errInvalidInput = "Bad request, invalid input"
	errNotFound     = "Task not found"
    LISTEN_PORT     = ":5000"
	DB_PATH         = "todo.db"
)

func main() {
	if OpenDB("DB_FILE") != nil {
		fmt.Println("Error: Failed to connect to the database")
		return
	}
	router := mux.NewRouter()
    fmt.Println("here")
    router.HandleFunc("/",GetTasks).Methods(("GET"))
	router.HandleFunc("/todo", GetTasks).Methods("GET")
	router.HandleFunc("/todo", AddTask).Methods("POST")
	router.HandleFunc("/todo/{id}", GetTask).Methods("GET")
	router.HandleFunc("/todo/{id}", CompleteTask).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteTask).Methods("DELETE")
	router.HandleFunc("/todo/{id}/{item}", UpdateTask).Methods("Post")
    err := http.ListenAndServe(LISTEN_PORT, router)
    fmt.Println(err.Error())
}

func OpenDB(path string) error {
	var err error
	db.AutoMigrate(&task{})
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	return err
}

func GetTasks(writer http.ResponseWriter, request *http.Request) {
	var tasks []task
	if db.Find(&tasks) != nil {
		http.Error(writer, "No tasks found", http.StatusNotFound)
	}
	jsonData, _ := json.Marshal(tasks)
	writer.Write(jsonData)
}

func AddTask(writer http.ResponseWriter, request *http.Request) {
	var newTask task
	if json.NewDecoder(request.Body).Decode(&newTask) != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	if db.First(&task{}, newTask.ID).Error == nil {
		http.Error(writer, "Bad request, task already exist", http.StatusConflict)
		return
	}
	db.Create(&newTask)
	jsonData, _ := json.Marshal(newTask)
	writer.Write(jsonData)
}

func GetTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var requestedTask task
	if db.First(&requestedTask, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusNotFound)
		return
	}
	jsonData, _ := json.Marshal(requestedTask)
	writer.Write(jsonData)

}

func DeleteTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	if db.Delete(&task{}, id).Error != nil {
		http.Error(writer, errNotFound, http.StatusNotFound)
		return
	}
	writer.Write([]byte("Task deleted successfully"))
}

func CompleteTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, errNotFound, http.StatusBadRequest)
		return
	}
	checkedTask := task{}
	if db.First(&checkedTask, id).Error == nil {
		http.Error(writer, errNotFound, http.StatusConflict)
		return
	}
	checkedTask.Completed = !checkedTask.Completed
	db.Save(checkedTask)

}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	item, _ := params["item"]
	if err != nil {
		http.Error(writer, errInvalidInput, http.StatusBadRequest)
		return
	}
	updatedTask := task{}
	if db.First(&updatedTask, id).Error == nil {
		http.Error(writer, "Bad request, task already exist", http.StatusConflict)
		return
	}
	updatedTask.Item = item
	db.Save(updatedTask)

}
