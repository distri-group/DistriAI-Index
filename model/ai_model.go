package model

import "time"

type AiModel struct {
	Id         uint      `gorm:"primarykey"`
	Owner      string    `gorm:"size:44;not null;uniqueIndex:idx_ai_models_owner_name"`
	Name       string    `gorm:"size:50;not null;uniqueIndex:idx_ai_models_owner_name"`
	Framework  uint8     `gorm:"not null"`
	License    uint8     `gorm:"not null"`
	Type1      uint8     `gorm:"not null"`
	Type2      uint8     `gorm:"not null"`
	Tags       string    `gorm:"size:128;not null"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoCreateTime"`
}

type AiModelHeat struct {
	Id        uint   `gorm:"primarykey"`
	Owner     string `gorm:"size:44;not null;uniqueIndex:idx_ai_model_heats_owner_name"`
	Name      string `gorm:"size:50;not null;uniqueIndex:idx_ai_model_heats_owner_name"`
	Likes     uint32 `gorm:"not null"`
	Downloads uint32 `gorm:"not null"`
	Clicks    uint32 `gorm:"not null"`
	Review    uint8  `gorm:"not null;comment:'0-pending,1-accept,2-reject'"`
	Size      uint32 `gorm:"not null"`
}

type AiModelLike struct {
	Id      uint   `gorm:"primarykey"`
	Account string `gorm:"size:44;not null"`
	Owner   string `gorm:"size:44;not null"`
	Name    string `gorm:"size:50;not null"`
}
