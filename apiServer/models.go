package apiserver

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type task struct {
    ID        string `gorm:"NOT NULL" json:"id"`
	Item      string `json:"item"`
	Completed bool   `jon:"completed"`
}
var db,ERR = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})  
