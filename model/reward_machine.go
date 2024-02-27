package model

type RewardMachine struct {
	Id        uint   `gorm:"primarykey"`
	Period    uint32 `gorm:"not null;index:idx_machines_period_owner_machineid"`
	Owner     string `gorm:"size:44;not null;index:idx_machines_period_owner_machineid"`
	MachineId string `gorm:"size:34;not null;index:idx_machines_period_owner_machineid"`
	TaskNum   uint32 `gorm:"not null"`
	Claimed   bool   `gorm:"not null"`
}
