package model

type Sync struct {
	Id          uint   `gorm:"primarykey"`
	BlockNumber uint64 `gorm:"not null"`
}
