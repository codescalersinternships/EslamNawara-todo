package main

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID        int64  `gorm:"NOT NULL" json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

func OpenDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	db.AutoMigrate(&Task{})
	return err
}

func SetRoutes(router *mux.Router) {
	router.HandleFunc("/todo", GetTasks).Methods("GET")
	router.HandleFunc("/todo", AddTask).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateTask).Methods("PATCH")
	router.HandleFunc("/todo/{id}", GetTask).Methods("GET")
	router.HandleFunc("/todo/{id}", DeleteTask).Methods("DELETE")
}
