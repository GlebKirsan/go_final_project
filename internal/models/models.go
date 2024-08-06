package models

import (
	"fmt"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) String() string {
	return fmt.Sprintf("Task{id=%d, date=%s, title=%s, comment=%s, repeat=%s}",
		t.ID, t.Date, t.Title, t.Comment, t.Repeat)
}
