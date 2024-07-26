package models

import (
	"time"
)

type Task struct {
	ID      uint      `json:"id" gorm:"column:id;primaryKey"`
	Date    time.Time `json:"date" gorm:"column:date;type:date;index"`
	Title   string    `json:"title" gorm:"column:title"`
	Comment string    `json:"comment" gorm:"column:comment"`
	Repeat  string    `json:"repeat" gorm:"column:repeat;type:varchar(128)"`
}

func (Task) TableName() string {
	return "scheduler"
}
