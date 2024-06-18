package model

import "time"

type Dataset struct {
	Id         uint      `gorm:"primarykey"`
	Owner      string    `gorm:"size:44;not null;uniqueIndex:idx_dataset_owner_name"`
	Name       string    `gorm:"size:50;not null;uniqueIndex:idx_dataset_owner_name"`
	Scale      uint8     `gorm:"not null"`
	License    uint8     `gorm:"not null"`
	Type1      uint8     `gorm:"not null"`
	Type2      uint8     `gorm:"not null"`
	Tags       string    `gorm:"size:128;not null"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoCreateTime"`
}

type DatasetHeat struct {
	Id        uint   `gorm:"primarykey"`
	Owner     string `gorm:"size:44;not null;uniqueIndex:idx_dataset_heats_owner_name"`
	Name      string `gorm:"size:50;not null;uniqueIndex:idx_dataset_heats_owner_name"`
	Likes     uint   `gorm:"not null"`
	Downloads uint   `gorm:"not null"`
	Clicks    uint   `gorm:"not null"`
	Status    uint8  `gorm:"not null;comment:'0-private,1-review,2-public'"`
	Reason    string `gorm:"size:128;not null;comment:'why review fails'"`
	Size      uint32 `gorm:"not null"`
}

type DatasetLike struct {
	Id      uint   `gorm:"primarykey"`
	Account string `gorm:"size:44;not null"`
	Owner   string `gorm:"size:44;not null;"`
	Name    string `gorm:"size:50;not null;"`
}
