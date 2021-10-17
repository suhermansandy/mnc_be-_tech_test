package model

import "time"

// Default model database
type Default struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by"`
	DeletedBy string     `json:"deleted_by"`
	IsUpload  string     `gorm:"default:'N';type:varchar(1)"`
}
