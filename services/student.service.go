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

// generateStudentNo creates a new student number per classroom in the format STD001, STD002, etc.
func (s *StudentService) generateStudentNo(classroomID uint) (string, error) {
	var lastStudent models.Student

	// Find the student with the highest student_no in the specific classroom
	err := configs.DB.Where("classroom_id = ? AND student_no LIKE 'STD%'", classroomID).
		Order("CAST(SUBSTRING(student_no, 4) AS INTEGER) DESC").
		First(&lastStudent).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// If no students exist in this classroom, start with STD001
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

	// Increment and format as STD001, STD002, etc. for this classroom
	nextNum := lastNum + 1
	return fmt.Sprintf("STD%03d", nextNum), nil
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
	// Verify classroom exists
	var classroom models.Classroom
	if err := configs.DB.Where("id = ?", req.ClassroomID).First(&classroom).Error; err != nil {
		return nil, errors.New("classroom not found")
	}

	// Generate student number if not provided (per classroom)
	studentNo := req.StudentNo
	if studentNo == "" {
		generatedNo, err := s.generateStudentNo(req.ClassroomID)
		if err != nil {
			logger.LogError(err, "Failed to generate student number", logrus.Fields{"classroom_id": req.ClassroomID})
			return nil, errors.New("failed to generate student number")
		}
		studentNo = generatedNo
	}

	logger.LogInfo("Creating new student", logrus.Fields{
		"student_no":   studentNo,
		"classroom_id": req.ClassroomID,
		"school_name":  req.SchoolName,
		"first_name":   req.Firstname,
		"last_name":    req.Lastname,
	})

	// Check if student already exists by student number in the same classroom
	var existingStudent models.Student
	if err := configs.DB.Where("student_no = ? AND classroom_id = ?", studentNo, req.ClassroomID).First(&existingStudent).Error; err == nil {
		logger.LogWarning("Student creation failed - student number already exists in classroom", logrus.Fields{
			"student_no":   studentNo,
			"classroom_id": req.ClassroomID,
		})
		return nil, errors.New("student with this student number already exists in this classroom")
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
		StudentNo:   studentNo,
		FirstName:   req.Firstname, // Note: field name difference
		LastName:    req.Lastname,  // Note: field name difference
		SchoolID:    &school.ID,
		ClassroomID: &req.ClassroomID,
		GenderID:    req.GenderID,
		PrefixID:    req.PrefixID,
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

	// Log activity automatically - use teacherID from service parameter if available
	// Note: For now using a default system user ID. Should be passed from controller context.
	var systemTeacherID uint = 1 // Default system user
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
	var systemTeacherID uint = 1 // Default system user
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
	var systemTeacherID uint = 1 // Default system user
	logger.LogActivity(systemTeacherID, models.LogActionDeleteStudent,
		fmt.Sprintf("ลบข้อมูลนักเรียน: %s %s (รหัส: %s)", student.FirstName, student.LastName, student.StudentNo),
		student.SchoolID)

	if err := configs.DB.Delete(&student).Error; err != nil {
		return errors.New("failed to delete student")
	}

	return nil
}

// TestCreateStudentWithAutoClassroom creates a student with auto classroom creation (for testing)
func (s *StudentService) TestCreateStudentWithAutoClassroom(schoolName, firstname, lastname string, genderID, prefixID *uint) (*models.Student, error) {
	// Find or create school
	var school models.School
	if err := configs.DB.Where("name = ?", schoolName).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new school
			school = models.School{
				Name: schoolName,
			}
			if err := configs.DB.Create(&school).Error; err != nil {
				return nil, errors.New("failed to create school")
			}
		} else {
			return nil, errors.New("failed to find school")
		}
	}

	// Find or create a default classroom for this school
	var classroom models.Classroom
	classroomName := "ห้องเรียนทดสอบ " + schoolName
	if err := configs.DB.Where("school_id = ? AND name = ?", school.ID, classroomName).First(&classroom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default classroom (need a teacher first)
			var teacher models.Teacher
			if err := configs.DB.Where("school_id = ?", school.ID).First(&teacher).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create default teacher
					teacher = models.Teacher{
						Email:     "test@" + schoolName + ".com",
						FirstName: "ครูทดสอบ",
						LastName:  "ระบบ",
						Phone:     "081-000-0000",
						SchoolID:  &school.ID,
					}
					if err := configs.DB.Create(&teacher).Error; err != nil {
						return nil, errors.New("failed to create default teacher")
					}
				} else {
					return nil, errors.New("failed to find teacher")
				}
			}

			// Create classroom
			classroom = models.Classroom{
				SchoolID:  &school.ID,
				TeacherID: &teacher.ID,
				Name:      classroomName,
				Grade:     "ม.1",
			}
			if err := configs.DB.Create(&classroom).Error; err != nil {
				return nil, errors.New("failed to create classroom")
			}
		} else {
			return nil, errors.New("failed to find classroom")
		}
	}

	// Generate student number for this classroom
	studentNo, err := s.generateStudentNo(classroom.ID)
	if err != nil {
		return nil, errors.New("failed to generate student number")
	}

	// Create student
	student := models.Student{
		StudentNo:   studentNo,
		FirstName:   firstname,
		LastName:    lastname,
		SchoolID:    &school.ID,
		ClassroomID: &classroom.ID,
		GenderID:    genderID,
		PrefixID:    prefixID,
	}

	if err := configs.DB.Create(&student).Error; err != nil {
		return nil, errors.New("failed to create student")
	}

	logger.LogInfo("Test student created successfully", logrus.Fields{
		"student_id":   student.ID,
		"student_no":   student.StudentNo,
		"classroom_id": classroom.ID,
		"school_id":    school.ID,
	})

	return &student, nil
}

// TestCreateStudent creates a student with auto-generated classroom for testing purposes
func (s *StudentService) TestCreateStudent(schoolName, firstname, lastname, studentNo *string, genderID, prefixID *uint) (*models.Student, error) {
	// Find or create school
	var school models.School
	if err := configs.DB.Where("name = ?", *schoolName).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new school
			school = models.School{Name: *schoolName}
			if err := configs.DB.Create(&school).Error; err != nil {
				return nil, errors.New("failed to create school")
			}
		} else {
			return nil, errors.New("failed to find school")
		}
	}

	// Find or create a default classroom for this school
	var classroom models.Classroom
	classroomName := "ห้องทดสอบ - " + *schoolName
	if err := configs.DB.Where("name = ? AND school_id = ?", classroomName, school.ID).First(&classroom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new classroom (need a teacher first)
			// Find first teacher in this school or create system teacher
			var teacher models.Teacher
			if err := configs.DB.Where("school_id = ?", school.ID).First(&teacher).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create system teacher
					teacher = models.Teacher{
						Email:     fmt.Sprintf("system@%s.com", school.Name),
						Password:  "system123", // This should be hashed in real implementation
						FirstName: "ระบบ",
						LastName:  "ทดสอบ",
						Phone:     "000-000-0000",
						SchoolID:  &school.ID,
					}
					if err := configs.DB.Create(&teacher).Error; err != nil {
						return nil, errors.New("failed to create system teacher")
					}
				} else {
					return nil, errors.New("failed to find teacher")
				}
			}

			// Create classroom
			classroom = models.Classroom{
				SchoolID:  &school.ID,
				TeacherID: &teacher.ID,
				Name:      classroomName,
				Grade:     "ทดสอบ",
			}
			if err := configs.DB.Create(&classroom).Error; err != nil {
				return nil, errors.New("failed to create classroom")
			}
		} else {
			return nil, errors.New("failed to find classroom")
		}
	}

	// Generate student number if not provided
	finalStudentNo := ""
	if studentNo != nil && *studentNo != "" {
		finalStudentNo = *studentNo
	} else {
		generatedNo, err := s.generateStudentNo(classroom.ID)
		if err != nil {
			return nil, errors.New("failed to generate student number")
		}
		finalStudentNo = generatedNo
	}

	// Create student
	student := models.Student{
		StudentNo:   finalStudentNo,
		FirstName:   *firstname,
		LastName:    *lastname,
		SchoolID:    &school.ID,
		ClassroomID: &classroom.ID,
		GenderID:    genderID,
		PrefixID:    prefixID,
	}

	if err := configs.DB.Create(&student).Error; err != nil {
		return nil, errors.New("failed to create student")
	}

	// Load relationships for response
	if err := configs.DB.Preload("School").Preload("Classroom").Preload("Gender").Preload("Prefix").First(&student, student.ID).Error; err != nil {
		return &student, nil // Return even if preload fails
	}

	return &student, nil
}

// GetStudentsByTeacherPaginated gets all students taught by a specific teacher with pagination
func (s *StudentService) GetStudentsByTeacherPaginated(teacherID uint, page, limit int) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	// Count total students for this teacher
	if err := configs.DB.
		Model(&models.Student{}).
		Joins("JOIN classrooms ON students.classroom_id = classrooms.id").
		Where("classrooms.teacher_id = ?", teacherID).
		Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count students")
	}

	// Get paginated students
	offset := (page - 1) * limit
	if err := configs.DB.
		Preload("School").
		Preload("Classroom").
		Preload("Gender").
		Preload("Prefix").
		Joins("JOIN classrooms ON students.classroom_id = classrooms.id").
		Where("classrooms.teacher_id = ?", teacherID).
		Offset(offset).
		Limit(limit).
		Find(&students).Error; err != nil {
		return nil, 0, errors.New("failed to get students by teacher")
	}

	return students, total, nil
}
