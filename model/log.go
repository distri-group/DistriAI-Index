package model

import "time"

type Log struct {
	Id        uint   `gorm:"primarykey"`
	OrderUuid string `gorm:"size:34;not null"`
	Content   string `gorm:"type:longtext"`
	CreatedAt time.Time
}
