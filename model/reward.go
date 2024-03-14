package model

import "time"

type Reward struct {
	Id                 uint      `gorm:"primarykey"`
	Period             uint32    `gorm:"not null;unique"`
	StartTime          time.Time `gorm:"autoCreateTime"`
	Pool               uint64    `gorm:"not null"`
	MachineNum         uint32    `gorm:"not null"`
	UnitPeriodicReward uint64    `gorm:"not null"`
	TaskNum            uint32    `gorm:"not null"`
	UnitTaskReward     uint64    `gorm:"not null"`
}
