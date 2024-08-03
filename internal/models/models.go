package models

type Task struct {
	ID      uint   `json:"id" gorm:"column:id;primaryKey"`
	Date    string `json:"date" gorm:"column:date;type:string;index"`
	Title   string `json:"title" gorm:"column:title;not null"`
	Comment string `json:"comment" gorm:"column:comment"`
	Repeat  string `json:"repeat" gorm:"column:repeat;type:varchar(128)"`
}

func (Task) TableName() string {
	return "scheduler"
}
