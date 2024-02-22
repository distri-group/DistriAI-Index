package model

type Reward struct {
	Id         uint   `gorm:"primarykey"`
	Period     uint32 `gorm:"not null;unique"`
	Pool       uint64 `gorm:"not null"`
	MachineNum uint32 `gorm:"not null"`
}
