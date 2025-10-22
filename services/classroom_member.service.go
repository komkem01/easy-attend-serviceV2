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

type ClassroomMemberService struct{}

func NewClassroomMemberService() *ClassroomMemberService {
	return &ClassroomMemberService{}
}

func (s *ClassroomMemberService) GetAllClassroomMembers() ([]models.ClassroomMember, error) {
	logger.LogInfo("Fetching all classroom members", logrus.Fields{})

	var members []models.ClassroomMember
	if err := configs.DB.Find(&members).Error; err != nil {
		logger.LogError(err, "Failed to fetch classroom members", logrus.Fields{})
		return nil, errors.New("failed to fetch classroom members")
	}

	logger.LogInfo("Successfully fetched classroom members", logrus.Fields{
		"count": len(members),
	})

	return members, nil
}

func (s *ClassroomMemberService) GetClassroomMembersByClassroomID(classroomID uint) ([]models.ClassroomMember, error) {
	logger.LogInfo("Fetching classroom members by classroom ID", logrus.Fields{
		"classroom_id": classroomID,
	})

	var members []models.ClassroomMember
	if err := configs.DB.Where("classroom_id = ?", classroomID).Find(&members).Error; err != nil {
		logger.LogError(err, "Failed to fetch classroom members", logrus.Fields{
			"classroom_id": classroomID,
		})
		return nil, errors.New("failed to fetch classroom members")
	}

	return members, nil
}

func (s *ClassroomMemberService) CreateClassroomMember(req *requests.ClassroomMemberCreateRequest) (*models.ClassroomMember, error) {
	logger.LogInfo("Creating new classroom member", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
	})

	// Validate that either teacher_id or student_id is provided, but not both
	if (req.TeacherID == nil && req.StudentID == nil) || (req.TeacherID != nil && req.StudentID != nil) {
		return nil, errors.New("either teacher_id or student_id must be provided, but not both")
	}

	// Check if member already exists in classroom
	var existingMember models.ClassroomMember
	query := configs.DB.Where("classroom_id = ?", req.ClassroomID)

	if req.TeacherID != nil {
		query = query.Where("teacher_id = ?", *req.TeacherID)
	}
	if req.StudentID != nil {
		query = query.Where("student_id = ?", *req.StudentID)
	}

	if err := query.First(&existingMember).Error; err == nil {
		logger.LogWarning("Classroom member already exists", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
		})
		return nil, errors.New("member already exists in this classroom")
	}

	// Create new classroom member
	member := models.ClassroomMember{
		ClassroomID: &req.ClassroomID,
		TeacherID:   req.TeacherID,
		StudentID:   req.StudentID,
	}

	if err := configs.DB.Create(&member).Error; err != nil {
		logger.LogError(err, "Failed to create classroom member", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
		})
		return nil, errors.New("failed to create classroom member")
	}

	logger.LogInfo("Classroom member created successfully", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", member.ClassroomID),
	})

	return &member, nil
}

func (s *ClassroomMemberService) UpdateClassroomMember(classroomID uint, memberID uint, req *requests.ClassroomMemberUpdateRequest) (*models.ClassroomMember, error) {
	logger.LogInfo("Updating classroom member", logrus.Fields{
		"classroom_id": classroomID,
		"member_id":    memberID,
	})

	// Find existing member
	var member models.ClassroomMember
	query := configs.DB.Where("classroom_id = ?", classroomID)

	// Find member by teacher_id or student_id
	if err := query.Where("teacher_id = ? OR student_id = ?", memberID, memberID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Classroom member not found for update", logrus.Fields{
				"classroom_id": classroomID,
				"member_id":    memberID,
			})
			return nil, errors.New("classroom member not found")
		}
		logger.LogError(err, "Failed to find classroom member for update", logrus.Fields{
			"classroom_id": classroomID,
			"member_id":    memberID,
		})
		return nil, errors.New("failed to find classroom member")
	}

	// Update fields
	member.TeacherID = req.TeacherID
	member.StudentID = req.StudentID

	if err := configs.DB.Save(&member).Error; err != nil {
		logger.LogError(err, "Failed to update classroom member", logrus.Fields{
			"classroom_id": classroomID,
			"member_id":    memberID,
		})
		return nil, errors.New("failed to update classroom member")
	}

	logger.LogInfo("Classroom member updated successfully", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", member.ClassroomID),
	})

	return &member, nil
}

func (s *ClassroomMemberService) DeleteClassroomMember(classroomID uint, memberID uint) error {
	logger.LogInfo("Deleting classroom member", logrus.Fields{
		"classroom_id": classroomID,
		"member_id":    memberID,
	})

	// Delete the member
	result := configs.DB.Where("classroom_id = ? AND (teacher_id = ? OR student_id = ?)", classroomID, memberID, memberID).Delete(&models.ClassroomMember{})

	if result.Error != nil {
		logger.LogError(result.Error, "Failed to delete classroom member", logrus.Fields{
			"classroom_id": classroomID,
			"member_id":    memberID,
		})
		return errors.New("failed to delete classroom member")
	}

	if result.RowsAffected == 0 {
		logger.LogWarning("Classroom member not found for deletion", logrus.Fields{
			"classroom_id": classroomID,
			"member_id":    memberID,
		})
		return errors.New("classroom member not found")
	}

	logger.LogInfo("Classroom member deleted successfully", logrus.Fields{
		"classroom_id": classroomID,
		"member_id":    memberID,
	})

	return nil
}
