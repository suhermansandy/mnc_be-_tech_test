package model

import "github.com/jinzhu/gorm"

// ChangeLog is object for change log mdm
type ChangeLog struct {
	gorm.Model
	TableName string `json:"table_name" gorm:"type:varchar(100);"`
	RowID     *uint  `json:"row_id"`
	IsSync    bool   `json:"is_sync" gorm:"default:false"`
}
