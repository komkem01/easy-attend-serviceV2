package migrations

import (
	"easy-attend-service/models"

	"gorm.io/gorm"
)

// CreateIntIDTables สร้างตารางใหม่ด้วย int ID แทน UUID
func CreateIntIDTables(db *gorm.DB) error {
	// Drop existing tables first (careful in production!)
	err := db.Migrator().DropTable(
		&models.Log{},
		&models.Attendance{},
		&models.ClassroomMember{},
		&models.Classroom{},
		&models.Student{},
		&models.Teacher{},
		&models.School{},
		&models.Gender{},
		&models.Prefix{},
	)
	if err != nil {
		return err
	}

	// Create tables with new schema
	err = db.AutoMigrate(
		&models.School{},
		&models.Teacher{},
		&models.Student{},
		&models.Classroom{},
		&models.ClassroomMember{},
		&models.Attendance{},
		&models.Log{},
		&models.Gender{},
		&models.Prefix{},
	)
	if err != nil {
		return err
	}

	// Insert default data
	return insertDefaultData(db)
}

func insertDefaultData(db *gorm.DB) error {
	// Insert default genders
	genders := []models.Gender{
		{Name: "ชาย"},
		{Name: "หญิง"},
		{Name: "ไม่ระบุ"},
	}

	for _, gender := range genders {
		var existing models.Gender
		if err := db.Where("name = ?", gender.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&gender).Error; err != nil {
					return err
				}
			}
		}
	}

	// Insert default prefixes
	prefixes := []models.Prefix{
		{Name: "นาย"},
		{Name: "นาง"},
		{Name: "นางสาว"},
		{Name: "ด.ช."},
		{Name: "ด.ญ."},
		{Name: "เด็กชาย"},
		{Name: "เด็กหญิง"},
	}

	for _, prefix := range prefixes {
		var existing models.Prefix
		if err := db.Where("name = ?", prefix.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&prefix).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}
