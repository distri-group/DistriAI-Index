package model

type Machine struct {
	Id                     uint   `gorm:"primarykey"`
	Owner                  string `gorm:"size:44;not null;uniqueIndex:idx_machines_owner_uuid"`
	Uuid                   string `gorm:"size:34;not null;uniqueIndex:idx_machines_owner_uuid"`
	Metadata               string `gorm:"size:2048;not null"`
	Status                 uint8  `gorm:"not null"`
	Price                  uint64 `gorm:"not null"`
	MaxDuration            uint32 `gorm:"not null"`
	Disk                   uint32 `gorm:"not null"`
	CompletedCount         uint32 `gorm:"not null"`
	FailedCount            uint32 `gorm:"not null"`
	Score                  uint32 `gorm:"not null"`
	ClaimedPeriodicRewards uint64 `gorm:"not null"`
	ClaimedTaskRewards     uint64 `gorm:"not null"`
	OrderPda               string `gorm:"size:44;not null"`
	Gpu                    string `gorm:"size:256"`
	GpuCount               uint32
	Region                 string `gorm:"size:32"`
	Tflops                 float64
	Ram                    float64
}

type MachineInfo struct {
	Id             uint   `gorm:"primarykey"`
	Owner          string `gorm:"size:44;not null;uniqueIndex:idx_machines_owner_uuid"`
	Uuid           string `gorm:"size:34;not null;uniqueIndex:idx_machines_owner_uuid"`
	CachedModels   string `gorm:"size:2048;not null"`
	CachedDatasets string `gorm:"size:2048;not null"`
}
