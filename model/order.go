package model

import "time"

type Order struct {
	Id          uint      `gorm:"primarykey"`
	Uuid        string    `gorm:"size:34;not null;unique"`
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
	Model1Owner string    `gorm:"size:44;not null"`
	Model1Name  string    `gorm:"size:50;not null"`
	Model2Owner string    `gorm:"size:44;not null"`
	Model2Name  string    `gorm:"size:50;not null"`
	Model3Owner string    `gorm:"size:44;not null"`
	Model3Name  string    `gorm:"size:50;not null"`
	Model4Owner string    `gorm:"size:44;not null"`
	Model4Name  string    `gorm:"size:50;not null"`
	Model5Owner string    `gorm:"size:44;not null"`
	Model5Name  string    `gorm:"size:50;not null"`
}
