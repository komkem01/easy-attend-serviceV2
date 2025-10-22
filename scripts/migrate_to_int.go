package main

import (
	"easy-attend-service/configs"
	"easy-attend-service/database/migrations"
	"fmt"
	"log"
)

func main() {
	// เชื่อมต่อ database
	configs.ConnectDatabase()

	// รัน migration
	if err := migrations.CreateIntIDTables(configs.DB); err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Migration to int ID completed successfully!")
}
