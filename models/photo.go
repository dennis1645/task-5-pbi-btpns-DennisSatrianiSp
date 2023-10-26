package models

type Photo struct {
	ID       uint   `json:"id" gorm:"primaryKey;not null"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id" gorm:"not null"`

	Users Users `json:"user" gorm:"foreignKey:UserID"`
}
