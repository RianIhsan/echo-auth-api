package models

type Book struct {
	Id          uint64 `json:"id" gorm:"primary_key;autoIncrement"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishYear uint   `json:"publish_year"`
	ISBN        string `json:"isbn" gorm:"unique"`
	Genre       string `json:"genre"`
	UserID      uint
}
