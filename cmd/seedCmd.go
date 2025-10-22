package cmd

import (
	"fmt"
	"log"

	"easy-attend-service/configs"

	"github.com/spf13/cobra"
)

var seedDataCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed initial data to database",
	Run: func(cmd *cobra.Command, args []string) {
		configs.ConnectDatabase()

		// Seed Genders
		genderData := []string{"ชาย", "หญิง", "อื่นๆ"}
		for _, gender := range genderData {
			sqlQuery := `
				INSERT INTO genders (name) 
				VALUES ($1) 
				ON CONFLICT (name) DO NOTHING;
			`
			if err := configs.DB.Exec(sqlQuery, gender).Error; err != nil {
				log.Printf("Warning: Failed to seed gender '%s': %v", gender, err)
			}
		}
		fmt.Println("✅ Genders seeded successfully")

		// Seed Prefixes
		prefixData := []string{"นาย", "นาง", "นางสาว", "เด็กชาย", "เด็กหญิง", "ดร.", "ผศ.", "รศ.", "ศ."}
		for _, prefix := range prefixData {
			sqlQuery := `
				INSERT INTO prefixes (name) 
				VALUES ($1) 
				ON CONFLICT (name) DO NOTHING;
			`
			if err := configs.DB.Exec(sqlQuery, prefix).Error; err != nil {
				log.Printf("Warning: Failed to seed prefix '%s': %v", prefix, err)
			}
		}
		fmt.Println("✅ Prefixes seeded successfully")

		fmt.Println("🎉 All seed data inserted successfully!")
	},
}

func init() {
	rootCmd.AddCommand(seedDataCmd)
}
