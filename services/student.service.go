package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudentService struct{}

func NewStudentService() *StudentService {
	return &StudentService{}
}

// generateStudentNo creates a new student number in the format STD001, STD002, etc.
func (s *StudentService) generateStudentNo() (string, error) {
	var lastStudent models.Student

	// Find the student with the highest student_no that starts with "STD"
	err := configs.DB.Where("student_no LIKE 'STD%'").
		Order("CAST(SUBSTRING(student_no, 4) AS INTEGER) DESC").
		First(&lastStudent).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// If no students exist or error finding them, start with STD001
	if err == gorm.ErrRecordNotFound {
		return "STD001", nil
	}

	// Extract the number part from the last student number
	lastNo := lastStudent.StudentNo
	if len(lastNo) < 4 || !strings.HasPrefix(lastNo, "STD") {
		return "STD001", nil
	}

	numPart := lastNo[3:] // Remove "STD" prefix
	lastNum, err := strconv.Atoi(numPart)
	if err != nil {
		return "STD001", nil
	}

	// Increment and format as STD001, STD002, etc.
	nextNum := lastNum + 1
	return fmt.Sprintf("STD%03d", nextNum), nil
}

func (s *StudentService) GetAllStudents(page, limit int) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	// Count total records
	if err := configs.DB.Model(&models.Student{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count students")
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get students with pagination
	if err := configs.DB.Offset(offset).Limit(limit).Find(&students).Error; err != nil {
		return nil, 0, errors.New("failed to get students")
	}

	return students, total, nil
}

func (s *StudentService) GetStudentByID(id uint) (*models.Student, error) {
	var student models.Student
	if err := configs.DB.Where("id = ?", id).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("failed to get student")
	}
	return &student, nil
}

func (s *StudentService) CreateStudent(req *requests.StudentCreateRequest) (*models.Student, error) {
	// Generate student number if not provided
	studentNo := req.StudentNo
	if studentNo == "" {
		generatedNo, err := s.generateStudentNo()
		if err != nil {
			logger.LogError(err, "Failed to generate student number", logrus.Fields{})
			return nil, errors.New("failed to generate student number")
		}
		studentNo = generatedNo
	}

	logger.LogInfo("Creating new student", logrus.Fields{
		"student_no":  studentNo,
		"school_name": req.SchoolName,
		"first_name":  req.Firstname,
		"last_name":   req.Lastname,
	})

	// Check if student already exists by student number
	var existingStudent models.Student
	if err := configs.DB.Where("student_no = ?", studentNo).First(&existingStudent).Error; err == nil {
		logger.LogWarning("Student creation failed - student number already exists", logrus.Fields{
			"student_no": studentNo,
		})
		return nil, errors.New("student with this student number already exists")
	}

	// Find or create school by name
	var school models.School
	if err := configs.DB.Where("name = ?", req.SchoolName).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogInfo("Creating new school", logrus.Fields{
				"school_name": req.SchoolName,
			})
			// Create new school
			school = models.School{
				Name: req.SchoolName,
			}
			if err := configs.DB.Create(&school).Error; err != nil {
				logger.LogError(err, "Failed to create school", logrus.Fields{
					"school_name": req.SchoolName,
				})
				return nil, errors.New("failed to create school")
			}
			logger.LogInfo("School created successfully", logrus.Fields{
				"school_id":   fmt.Sprintf("%d", school.ID),
				"school_name": school.Name,
			})
		} else {
			logger.LogError(err, "Failed to find school", logrus.Fields{
				"school_name": req.SchoolName,
			})
			return nil, errors.New("failed to find school")
		}
	}

	// Create new student
	student := models.Student{
		StudentNo: studentNo,
		FirstName: req.Firstname, // Note: field name difference
		LastName:  req.Lastname,  // Note: field name difference
		SchoolID:  &school.ID,
		GenderID:  req.GenderID,
		PrefixID:  req.PrefixID,
	}

	if err := configs.DB.Create(&student).Error; err != nil {
		logger.LogError(err, "Failed to create student", logrus.Fields{
			"student_no": req.StudentNo,
			"school_id":  fmt.Sprintf("%d", school.ID),
		})
		return nil, errors.New("failed to create student")
	}

	logger.LogInfo("Student created successfully", logrus.Fields{
		"student_id": fmt.Sprintf("%d", student.ID),
		"student_no": student.StudentNo,
		"school_id":  fmt.Sprintf("%d", school.ID),
	})

	// Log activity automatically - ใช้ teacherID จาก context หรือ system user
	// TODO: Get actual teacher ID from context/session
	var systemTeacherID uint = 1 // Temporary - should get from auth context
	logger.LogActivity(systemTeacherID, models.LogActionCreateStudent, fmt.Sprintf("สร้างนักเรียนใหม่: %s %s (รหัส: %s)", req.Firstname, req.Lastname, studentNo), &school.ID)

	return &student, nil
}

func (s *StudentService) UpdateStudent(id uint, req *requests.StudentUpdateRequest) (*models.Student, error) {
	var student models.Student
	if err := configs.DB.Where("id = ?", id).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("failed to find student")
	}

	// Check if student number is being changed and if it already exists
	if req.StudentNo != student.StudentNo {
		var existingStudent models.Student
		if err := configs.DB.Where("student_no = ? AND id != ?", req.StudentNo, id).First(&existingStudent).Error; err == nil {
			return nil, errors.New("student with this student number already exists")
		}
	}

	// Find or create school by name
	var school models.School
	if err := configs.DB.Where("name = ?", req.SchoolName).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new school
			school = models.School{
				Name: req.SchoolName,
			}
			if err := configs.DB.Create(&school).Error; err != nil {
				return nil, errors.New("failed to create school")
			}
		} else {
			return nil, errors.New("failed to find school")
		}
	}

	// Update fields
	student.StudentNo = req.StudentNo
	student.FirstName = req.Firstname // Note: field name difference
	student.LastName = req.Lastname   // Note: field name difference
	student.SchoolID = &school.ID

	if err := configs.DB.Save(&student).Error; err != nil {
		return nil, errors.New("failed to update student")
	}

	// Log activity automatically
	var systemTeacherID uint = 1 // TODO: Get actual teacher ID from context/session
	logger.LogActivity(systemTeacherID, models.LogActionUpdateStudent,
		fmt.Sprintf("อัพเดทข้อมูลนักเรียน: %s %s (รหัส: %s)", req.Firstname, req.Lastname, req.StudentNo),
		&school.ID)

	return &student, nil
}

func (s *StudentService) DeleteStudent(id uint) error {
	var student models.Student
	if err := configs.DB.Where("id = ?", id).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("student not found")
		}
		return errors.New("failed to find student")
	}

	// Log activity automatically before deletion
	var systemTeacherID uint = 1 // TODO: Get actual teacher ID from context/session
	logger.LogActivity(systemTeacherID, models.LogActionDeleteStudent,
		fmt.Sprintf("ลบข้อมูลนักเรียน: %s %s (รหัส: %s)", student.FirstName, student.LastName, student.StudentNo),
		student.SchoolID)

	if err := configs.DB.Delete(&student).Error; err != nil {
		return errors.New("failed to delete student")
	}

	return nil
}
