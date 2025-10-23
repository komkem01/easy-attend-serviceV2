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

type ClassroomService struct{}

func NewClassroomService() *ClassroomService {
	return &ClassroomService{}
}

func (s *ClassroomService) GetClassroomByID(id uint) (*models.Classroom, error) {
	logger.LogInfo("Fetching classroom by ID", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", id),
	})

	var classroom models.Classroom
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&classroom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Classroom not found", logrus.Fields{
				"classroom_id": fmt.Sprintf("%d", id),
			})
			return nil, errors.New("classroom not found")
		}
		logger.LogError(err, "Failed to fetch classroom", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to fetch classroom")
	}

	return &classroom, nil
}

func (s *ClassroomService) CreateClassroom(req *requests.ClassroomCreateRequest) (*models.Classroom, error) {
	logger.LogInfo("Creating new classroom", logrus.Fields{
		"name":       req.Name,
		"school_id":  fmt.Sprintf("%d", req.SchoolID),
		"teacher_id": fmt.Sprintf("%d", req.TeacherID),
	})

	// Check if classroom name already exists in the same school
	var existingClassroom models.Classroom
	if err := configs.DB.Where("name = ? AND school_id = ? AND deleted_at IS NULL", req.Name, req.SchoolID).First(&existingClassroom).Error; err == nil {
		logger.LogWarning("Classroom creation failed - name already exists in school", logrus.Fields{
			"name":      req.Name,
			"school_id": fmt.Sprintf("%d", req.SchoolID),
		})
		return nil, errors.New("classroom with this name already exists in this school")
	}

	// Create new classroom
	classroom := models.Classroom{
		SchoolID:  &req.SchoolID,
		TeacherID: &req.TeacherID,
		Name:      req.Name,
		Grade:     req.Grade,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := configs.DB.Create(&classroom).Error; err != nil {
		logger.LogError(err, "Failed to create classroom", logrus.Fields{
			"name": req.Name,
		})
		return nil, errors.New("failed to create classroom")
	}

	logger.LogInfo("Classroom created successfully", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", classroom.ID),
		"name":         classroom.Name,
	})

	// Log activity automatically
	logger.LogActivity(req.TeacherID, models.LogActionCreateClassroom, fmt.Sprintf("สร้างห้องเรียนใหม่: %s", req.Name), &req.SchoolID)

	return &classroom, nil
}

func (s *ClassroomService) UpdateClassroom(id uint, req *requests.ClassroomUpdateRequest) (*models.Classroom, error) {
	logger.LogInfo("Updating classroom", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", id),
		"name":         req.Name,
	})

	var classroom models.Classroom
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&classroom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Classroom not found for update", logrus.Fields{
				"classroom_id": fmt.Sprintf("%d", id),
			})
			return nil, errors.New("classroom not found")
		}
		logger.LogError(err, "Failed to find classroom for update", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to find classroom")
	}

	// Check if name is being changed and if it already exists in the same school
	if req.Name != classroom.Name {
		var existingClassroom models.Classroom
		if err := configs.DB.Where("name = ? AND school_id = ? AND id != ? AND deleted_at IS NULL", req.Name, req.SchoolID, id).First(&existingClassroom).Error; err == nil {
			logger.LogWarning("Classroom update failed - name already exists in school", logrus.Fields{
				"classroom_id": fmt.Sprintf("%d", id),
				"name":         req.Name,
				"school_id":    fmt.Sprintf("%d", req.SchoolID),
			})
			return nil, errors.New("classroom with this name already exists in this school")
		}
	}

	// Update fields
	classroom.SchoolID = &req.SchoolID
	classroom.TeacherID = &req.TeacherID
	classroom.Name = req.Name
	classroom.Grade = req.Grade
	classroom.UpdatedAt = time.Now().Unix()

	if err := configs.DB.Save(&classroom).Error; err != nil {
		logger.LogError(err, "Failed to update classroom", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", id),
		})
		return nil, errors.New("failed to update classroom")
	}

	logger.LogInfo("Classroom updated successfully", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", classroom.ID),
		"name":         classroom.Name,
	})

	// Log activity automatically
	logger.LogActivity(req.TeacherID, models.LogActionUpdateClassroom, fmt.Sprintf("อัพเดทห้องเรียน: %s", req.Name), &req.SchoolID)

	return &classroom, nil
}

func (s *ClassroomService) DeleteClassroom(id uint) error {
	logger.LogInfo("Deleting classroom", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", id),
	})

	var classroom models.Classroom
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&classroom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Classroom not found for deletion", logrus.Fields{
				"classroom_id": fmt.Sprintf("%d", id),
			})
			return errors.New("classroom not found")
		}
		logger.LogError(err, "Failed to find classroom for deletion", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", id),
		})
		return errors.New("failed to find classroom")
	}

	// Soft delete
	deleteTime := time.Now().Unix()
	classroom.DeletedAt = &deleteTime

	if err := configs.DB.Save(&classroom).Error; err != nil {
		logger.LogError(err, "Failed to delete classroom", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", id),
		})
		return errors.New("failed to delete classroom")
	}

	logger.LogInfo("Classroom deleted successfully", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", id),
	})

	// Log activity automatically
	logger.LogActivity(*classroom.TeacherID, models.LogActionDeleteClassroom, fmt.Sprintf("ลบห้องเรียน: %s", classroom.Name), classroom.SchoolID)

	return nil
}

// GetClassroomsByTeacher gets all classrooms for a specific teacher
func (s *ClassroomService) GetClassroomsByTeacher(teacherID uint) ([]models.Classroom, error) {
	var classrooms []models.Classroom

	if err := configs.DB.Where("teacher_id = ?", teacherID).Find(&classrooms).Error; err != nil {
		return nil, errors.New("failed to get classrooms by teacher")
	}

	return classrooms, nil
}
