package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceService struct{}

func NewAttendanceService() *AttendanceService {
	return &AttendanceService{}
}

func (s *AttendanceService) GetAllAttendances() ([]models.Attendance, error) {
	logger.LogInfo("Fetching all attendances", logrus.Fields{})

	var attendances []models.Attendance
	if err := configs.DB.Find(&attendances).Error; err != nil {
		logger.LogError(err, "Failed to fetch attendances", logrus.Fields{})
		return nil, errors.New("failed to fetch attendances")
	}

	logger.LogInfo("Successfully fetched attendances", logrus.Fields{
		"count": len(attendances),
	})

	return attendances, nil
}

func (s *AttendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	logger.LogInfo("Fetching attendance by ID", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", id),
	})

	var attendance models.Attendance
	if err := configs.DB.Where("id = ?", id).First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Attendance not found", logrus.Fields{
				"attendance_id": fmt.Sprintf("%d", id),
			})
			return nil, errors.New("attendance not found")
		}
		logger.LogError(err, "Failed to fetch attendance", logrus.Fields{
			"attendance_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to fetch attendance")
	}

	return &attendance, nil
}

func (s *AttendanceService) GetAttendancesByClassroom(classroomID uint) ([]models.Attendance, error) {
	logger.LogInfo("Fetching attendances by classroom", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", classroomID),
	})

	var attendances []models.Attendance
	if err := configs.DB.Where("classroom_id = ?", classroomID).Find(&attendances).Error; err != nil {
		logger.LogError(err, "Failed to fetch attendances by classroom", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", classroomID),
		})
		return nil, errors.New("failed to fetch attendances")
	}

	return attendances, nil
}

func (s *AttendanceService) GetAttendancesByStudent(studentID uint) ([]models.Attendance, error) {
	logger.LogInfo("Fetching attendances by student", logrus.Fields{
		"student_id": fmt.Sprintf("%d", studentID),
	})

	var attendances []models.Attendance
	if err := configs.DB.Where("student_id = ?", studentID).Find(&attendances).Error; err != nil {
		logger.LogError(err, "Failed to fetch attendances by student", logrus.Fields{
			"student_id": fmt.Sprintf("%d", studentID),
		})
		return nil, errors.New("failed to fetch attendances")
	}

	return attendances, nil
}

func (s *AttendanceService) CreateAttendance(req *requests.AttendanceCreateRequest) (*models.Attendance, error) {
	logger.LogInfo("Creating new attendance", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
		"teacher_id":   fmt.Sprintf("%d", req.TeacherID),
		"student_id":   fmt.Sprintf("%d", req.StudentID),
		"session_date": req.SessionDate,
		"status":       string(req.Status),
	})

	// Check if attendance already exists for this student on this date in this classroom
	var existingAttendance models.Attendance
	if err := configs.DB.Where("classroom_id = ? AND student_id = ? AND session_date = ?",
		req.ClassroomID, req.StudentID, req.SessionDate).First(&existingAttendance).Error; err == nil {
		logger.LogWarning("Attendance already exists", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
			"student_id":   fmt.Sprintf("%d", req.StudentID),
			"session_date": req.SessionDate,
		})
		return nil, errors.New("attendance for this student on this date already exists")
	}

	// Create new attendance
	attendance := models.Attendance{
		ClassroomID: &req.ClassroomID,
		TeacherID:   &req.TeacherID,
		StudentID:   &req.StudentID,
		SessionDate: req.SessionDate,
		Status:      req.Status,
		CheckedAt:   req.CheckedAt,
		Remark:      req.Remark,
	}

	if err := configs.DB.Create(&attendance).Error; err != nil {
		logger.LogError(err, "Failed to create attendance", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
			"student_id":   fmt.Sprintf("%d", req.StudentID),
		})
		return nil, errors.New("failed to create attendance")
	}

	logger.LogInfo("Attendance created successfully", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", attendance.ID),
		"student_id":    fmt.Sprintf("%d", attendance.StudentID),
		"status":        string(attendance.Status),
	})

	// Log activity automatically - get school ID from classroom
	var classroom models.Classroom
	if err := configs.DB.Where("id = ?", req.ClassroomID).First(&classroom).Error; err == nil {
		logger.LogActivity(req.TeacherID, models.LogActionAttendance,
			fmt.Sprintf("บันทึกการเข้าเรียน: %s (วันที่: %s)", string(req.Status), req.SessionDate),
			classroom.SchoolID)
	}

	return &attendance, nil
}

func (s *AttendanceService) UpdateAttendance(id uint, req *requests.AttendanceUpdateRequest) (*models.Attendance, error) {
	logger.LogInfo("Updating attendance", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", id),
		"status":        string(req.Status),
	})

	var attendance models.Attendance
	if err := configs.DB.Where("id = ?", id).First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Attendance not found for update", logrus.Fields{
				"attendance_id": fmt.Sprintf("%d", id),
			})
			return nil, errors.New("attendance not found")
		}
		logger.LogError(err, "Failed to find attendance for update", logrus.Fields{
			"attendance_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to find attendance")
	}

	// Update fields
	attendance.Status = req.Status
	attendance.CheckedAt = req.CheckedAt
	attendance.Remark = req.Remark

	if err := configs.DB.Save(&attendance).Error; err != nil {
		logger.LogError(err, "Failed to update attendance", logrus.Fields{
			"attendance_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to update attendance")
	}

	logger.LogInfo("Attendance updated successfully", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", attendance.ID),
		"status":        string(attendance.Status),
	})

	return &attendance, nil
}

func (s *AttendanceService) DeleteAttendance(id uint) error {
	logger.LogInfo("Deleting attendance", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", id),
	})

	var attendance models.Attendance
	if err := configs.DB.Where("id = ?", id).First(&attendance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Attendance not found for deletion", logrus.Fields{
				"attendance_id": fmt.Sprintf("%d", id),
			})
			return errors.New("attendance not found")
		}
		logger.LogError(err, "Failed to find attendance for deletion", logrus.Fields{
			"attendance_id": fmt.Sprintf("%d", id),
		})
		return errors.New("failed to find attendance")
	}

	if err := configs.DB.Delete(&attendance).Error; err != nil {
		logger.LogError(err, "Failed to delete attendance", logrus.Fields{
			"attendance_id": fmt.Sprintf("%d", id),
		})
		return errors.New("failed to delete attendance")
	}

	logger.LogInfo("Attendance deleted successfully", logrus.Fields{
		"attendance_id": fmt.Sprintf("%d", id),
	})

	return nil
}
