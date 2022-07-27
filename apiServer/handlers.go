package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getTasks(writer http.ResponseWriter, request *http.Request) {
	var tasks []task
	db.Find(&tasks)
	jsonData, _ := json.Marshal(tasks)
	writer.Write(jsonData)
}

func addTask(writer http.ResponseWriter, request *http.Request) {
	var newTask task
	err := json.NewDecoder(request.Body).Decode(&newTask)
	if err != nil {
		http.Error(writer, "Bad request, invalid input", http.StatusBadRequest)
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

func getTask(writer http.ResponseWriter, request *http.Request) {
	id := 0               //how to get the id
	err := fmt.Errorf("") //parsing id error
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var requestedTask task
	if db.First(&requestedTask, id).Error != nil {
		http.Error(writer, "Task not found", http.StatusNotFound)
		return
	}
	jsonData, _ := json.Marshal(requestedTask)
	writer.Write(jsonData)

}

func deleteTask(writer http.ResponseWriter, request *http.Request) {
	id := 0
	err := fmt.Errorf("")
	if err != nil {
		http.Error(writer, "Bad request, invalid input", http.StatusBadRequest)
		return
	}
	if db.Delete(&task{}, id).Error != nil {
		http.Error(writer, "Task not found", http.StatusNotFound)
		return
	}
}

func completeTask(writer http.ResponseWriter, request *http.Request) {
    id := 0

}

func updateTask(writer http.ResponseWriter, request *http.Request) {

}
