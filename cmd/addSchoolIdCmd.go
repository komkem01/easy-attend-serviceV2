package cmd

import (
	"fmt"
	"log"

	"easy-attend-service/configs"

	"github.com/spf13/cobra"
)

var addSchoolIDCmd = &cobra.Command{
	Use:   "add-school-id",
	Short: "Add school_id column to students table",
	Run: func(cmd *cobra.Command, args []string) {
		configs.ConnectDatabase()

		// Add school_id column to students table
		sqlQuery := `
			ALTER TABLE students 
			ADD COLUMN school_id UUID;
		`

		if err := configs.DB.Exec(sqlQuery).Error; err != nil {
			log.Fatalf("Failed to add school_id column: %v", err)
		}

		fmt.Println("Successfully added school_id column to students table")

		// Add foreign key constraint
		constraintQuery := `
			ALTER TABLE students 
			ADD CONSTRAINT fk_students_school 
			FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;
		`

		if err := configs.DB.Exec(constraintQuery).Error; err != nil {
			log.Printf("Warning: Failed to add foreign key constraint: %v", err)
		} else {
			fmt.Println("Successfully added foreign key constraint")
		}
	},
}

func init() {
	rootCmd.AddCommand(addSchoolIDCmd)
}
