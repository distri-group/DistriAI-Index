package common

import (
	"distriai-index-solana/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// The function initializes the database connection using the provided configuration settings and migrates the database schema.
func initDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Conf.Database.Username,
		Conf.Database.Password,
		Conf.Database.Host,
		Conf.Database.Port,
		Conf.Database.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 100,
	})
	if err != nil {
		panic("database init err! " + err.Error())
	}

	// Migrate the schema
	model.AutoMigrate(db)
	Db = db
}
