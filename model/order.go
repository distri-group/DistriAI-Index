package model

import "time"

type Order struct {
	Id          uint      `gorm:"primarykey"`
	Uuid        string    `gorm:"size:34;not null;uniqueIndex"`
	Buyer       string    `gorm:"size:44;not null"`
	Seller      string    `gorm:"size:44;not null"`
	MachineUuid string    `gorm:"size:34;not null"`
	Price       uint64    `gorm:"not null"`
	Duration    uint32    `gorm:"not null"`
	Total       uint64    `gorm:"not null"`
	Metadata    string    `gorm:"size:2048;not null"`
	Status      uint8     `gorm:"not null"`
	OrderTime   time.Time `gorm:"autoCreateTime"`
	StartTime   time.Time `gorm:"autoCreateTime"`
	RefundTime  time.Time `gorm:"autoCreateTime"`
}
