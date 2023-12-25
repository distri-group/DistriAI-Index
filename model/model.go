package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.Migrator().DropTable(
		&Sync{},
		&Machine{},
		&Order{},
	)
	db.AutoMigrate(
		&Mailbox{},
		&Sync{},
		&Machine{},
		&Order{},
		&Log{},
	)
}
