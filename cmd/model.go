package cmd

import (
	"easy-attend-service/configs"
	"easy-attend-service/database/migrations"
	"log"
)

func modelUp() error {
	log.Printf("Executing model up...")
	configs.AutoMigrate()
	return nil
}

func modelDown() error {
	log.Printf("Executing model down...")
	log.Printf("Manual model rollback not implemented - use database tools")
	return nil
}

func migrateToIntID() error {
	log.Printf("Migrating to int ID schema...")
	return migrations.CreateIntIDTables(configs.DB)
}
