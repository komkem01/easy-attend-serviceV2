package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

// GetAllLogs - อ่านข้อมูล log ทั้งหมด (เรียงตามเวลาล่าสุดก่อน)
func (s *LogService) GetAllLogs() ([]models.Log, error) {
	logger.LogInfo("Fetching all logs", logrus.Fields{})

	var logs []models.Log
	if err := configs.DB.Order("created_at DESC").Find(&logs).Error; err != nil {
		logger.LogError(err, "Failed to fetch logs", logrus.Fields{})
		return nil, errors.New("failed to fetch logs")
	}

	logger.LogInfo("Successfully fetched logs", logrus.Fields{
		"count": len(logs),
	})

	return logs, nil
}

// GetLogByID - อ่านข้อมูล log ตาม ID
func (s *LogService) GetLogByID(id uint) (*models.Log, error) {
	logger.LogInfo("Fetching log by ID", logrus.Fields{
		"log_id": fmt.Sprintf("%d", id),
	})

	var log models.Log
	if err := configs.DB.Where("id = ?", id).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Log not found", logrus.Fields{
				"log_id": fmt.Sprintf("%d", id),
			})
			return nil, errors.New("log not found")
		}
		logger.LogError(err, "Failed to fetch log", logrus.Fields{
			"log_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to fetch log")
	}

	return &log, nil
}

func (s *LogService) GetLogsByTeacher(teacherID string) ([]models.Log, error) {
	logger.LogInfo("Fetching logs by teacher", logrus.Fields{
		"teacher_id": teacherID,
	})

	var logs []models.Log
	if err := configs.DB.Where("teacher_id = ?", teacherID).Order("created_at DESC").Find(&logs).Error; err != nil {
		logger.LogError(err, "Failed to fetch logs by teacher", logrus.Fields{
			"teacher_id": teacherID,
		})
		return nil, errors.New("failed to fetch logs")
	}

	return logs, nil
}

func (s *LogService) GetLogsByAction(action models.LogAction) ([]models.Log, error) {
	logger.LogInfo("Fetching logs by action", logrus.Fields{
		"action": string(action),
	})

	var logs []models.Log
	if err := configs.DB.Where("action = ?", action).Order("created_at DESC").Find(&logs).Error; err != nil {
		logger.LogError(err, "Failed to fetch logs by action", logrus.Fields{
			"action": string(action),
		})
		return nil, errors.New("failed to fetch logs")
	}

	return logs, nil
}

func (s *LogService) CreateLog(req *requests.LogCreateRequest) (*models.Log, error) {
	logger.LogInfo("Creating new log", logrus.Fields{
		"teacher_id": fmt.Sprintf("%d", req.TeacherID),
		"action":     string(req.Action),
		"detail":     req.Detail,
	})

	// Create new log
	log := models.Log{
		TeacherID: req.TeacherID,
		Action:    req.Action,
		Detail:    req.Detail,
		CreatedAt: time.Now().Unix(),
	}

	if err := configs.DB.Create(&log).Error; err != nil {
		logger.LogError(err, "Failed to create log", logrus.Fields{
			"teacher_id": fmt.Sprintf("%d", req.TeacherID),
			"action":     string(req.Action),
		})
		return nil, errors.New("failed to create log")
	}

	logger.LogInfo("Log created successfully", logrus.Fields{
		"log_id":     fmt.Sprintf("%d", log.ID),
		"teacher_id": fmt.Sprintf("%d", log.TeacherID),
		"action":     string(log.Action),
	})

	return &log, nil
}
