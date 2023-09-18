package models

type User struct {
	Id           uint   `json:"id" gorm:"primarykey"`
	Name         string `json:"name" gorm:"not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `json:"password" gorm:"not null"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Books        []Book `gorm:"foreignkey:UserID"`
}
