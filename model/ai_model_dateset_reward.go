package model

import "time"

type AiModelDatasetRewardPeriod struct {
	Id                 uint      `gorm:"primarykey"`
	Period             uint32    `gorm:"not null;unique"`
	StartTime          time.Time `gorm:"autoCreateTime"`
	Pool               uint64    `gorm:"not null"`
	AiModelNum         uint32    `gorm:"not null"`
	DatasetNum         uint32    `gorm:"not null"`
	UnitPeriodicReward uint64    `gorm:"not null"`
}

type AiModelDatasetReward struct {
	Id             uint   `gorm:"primarykey"`
	Period         uint32 `gorm:"not null;index:idx_rewards_period_owner"`
	Owner          string `gorm:"size:44;not null;index:idx_rewards_period_owner"`
	AiModelNum     uint32 `gorm:"not null"`
	DatasetNum     uint32 `gorm:"not null"`
	PeriodicReward uint64 `gorm:"not null"`
	Reported       bool   `gorm:"not null"`
}
