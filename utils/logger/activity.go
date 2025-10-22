package logger

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// LogActivity สร้าง log record อัตโนมัติสำหรับ activity ต่างๆ
func LogActivity(teacherID uint, action models.LogAction, detail string, schoolID *uint) error {
	log := models.Log{
		TeacherID: teacherID,
		Action:    action,
		Detail:    detail,
		CreatedAt: time.Now().Unix(),
		SchoolID:  schoolID,
	}

	if err := configs.DB.Create(&log).Error; err != nil {
		LogError(err, "Failed to create activity log", logrus.Fields{
			"teacher_id": fmt.Sprintf("%d", teacherID),
			"action":     string(action),
			"detail":     detail,
		})
		return err
	}

	LogInfo("Activity logged successfully", logrus.Fields{
		"log_id":     fmt.Sprintf("%d", log.ID),
		"teacher_id": fmt.Sprintf("%d", teacherID),
		"action":     string(action),
	})

	return nil
}

// LogActivityWithContext สร้าง log พร้อม context เพิ่มเติม
func LogActivityWithContext(teacherID uint, action models.LogAction, detail string, schoolID *uint, context map[string]interface{}) error {
	// เพิ่ม context เข้าไปใน detail
	contextDetail := detail
	if len(context) > 0 {
		contextDetail += " | Context: "
		for key, value := range context {
			contextDetail += key + "=" + formatValue(value) + " "
		}
	}

	return LogActivity(teacherID, action, contextDetail, schoolID)
}

// formatValue แปลงค่าต่างๆ เป็น string
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case uint:
		return fmt.Sprintf("%d", v)
	case int, int32, int64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
