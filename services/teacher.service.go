package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TeacherService struct{}

func NewTeacherService() *TeacherService {
	return &TeacherService{}
}

func (s *TeacherService) GetAllTeachers(page, limit int) ([]models.Teacher, int64, error) {
	var teachers []models.Teacher
	var total int64

	// Count total records
	if err := configs.DB.Model(&models.Teacher{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count teachers")
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get teachers with pagination
	if err := configs.DB.Offset(offset).Limit(limit).Find(&teachers).Error; err != nil {
		return nil, 0, errors.New("failed to get teachers")
	}

	return teachers, total, nil
}

func (s *TeacherService) GetTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := configs.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to get teacher")
	}
	return &teacher, nil
}

func (s *TeacherService) CreateTeacher(req *requests.TeacherCreateRequest) (*models.Teacher, error) {
	// Check if teacher already exists
	var existingTeacher models.Teacher
	if err := configs.DB.Where("email = ?", req.Email).First(&existingTeacher).Error; err == nil {
		return nil, errors.New("teacher with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new teacher
	teacher := models.Teacher{
		SchoolID:  &req.SchoolID,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	if err := configs.DB.Create(&teacher).Error; err != nil {
		return nil, errors.New("failed to create teacher")
	}

	// Log activity automatically
	logger.LogActivity(teacher.ID, models.LogActionCreateTeacher, fmt.Sprintf("สร้างครูใหม่: %s %s (%s)", req.FirstName, req.LastName, req.Email), teacher.SchoolID)

	return &teacher, nil
}

func (s *TeacherService) UpdateTeacher(id string, req *requests.TeacherUpdateRequest) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := configs.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to find teacher")
	}

	// Check if email is being changed and if it already exists
	if req.Email != teacher.Email {
		var existingTeacher models.Teacher
		if err := configs.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingTeacher).Error; err == nil {
			return nil, errors.New("teacher with this email already exists")
		}
	}

	// Find or create school if school name is provided
	if req.SchoolName != "" {
		var school models.School
		err := configs.DB.Where("name = ?", req.SchoolName).First(&school).Error
		if err != nil {
			// School doesn't exist, create new one
			school = models.School{
				Name: req.SchoolName,
			}
			if err := configs.DB.Create(&school).Error; err != nil {
				return nil, errors.New("failed to create school")
			}
		}
		teacher.SchoolID = &school.ID
	}
	teacher.Email = req.Email
	teacher.FirstName = req.FirstName
	teacher.LastName = req.LastName
	teacher.Phone = req.Phone

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		teacher.Password = hashedPassword
	}

	if err := configs.DB.Save(&teacher).Error; err != nil {
		return nil, errors.New("failed to update teacher")
	}

	// Log activity automatically
	logger.LogActivity(teacher.ID, models.LogActionUpdateTeacher, fmt.Sprintf("อัพเดทข้อมูลครู: %s %s (%s)", teacher.FirstName, teacher.LastName, teacher.Email), teacher.SchoolID)

	return &teacher, nil
}

func (s *TeacherService) DeleteTeacher(id string) error {
	var teacher models.Teacher
	if err := configs.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("teacher not found")
		}
		return errors.New("failed to find teacher")
	}

	// Log activity automatically before deletion
	logger.LogActivity(teacher.ID, models.LogActionDeleteTeacher, fmt.Sprintf("ลบข้อมูลครู: %s %s (%s)", teacher.FirstName, teacher.LastName, teacher.Email), teacher.SchoolID)

	if err := configs.DB.Delete(&teacher).Error; err != nil {
		return errors.New("failed to delete teacher")
	}

	return nil
}

// TeacherInfo represents comprehensive teacher information
type TeacherInfo struct {
	Teacher    models.Teacher  `json:"teacher"`
	School     models.School   `json:"school"`
	Classrooms []ClassroomInfo `json:"classrooms"`
	TotalStats TeacherStats    `json:"total_stats"`
}

type ClassroomInfo struct {
	Classroom    models.Classroom `json:"classroom"`
	StudentCount int64            `json:"student_count"`
	Students     []models.Student `json:"students"`
}

type TeacherStats struct {
	TotalClassrooms int64 `json:"total_classrooms"`
	TotalStudents   int64 `json:"total_students"`
	TotalAttendance int64 `json:"total_attendance"`
}

// GetTeacherInfo gets comprehensive information for a specific teacher
func (s *TeacherService) GetTeacherInfo(teacherID uint) (*TeacherInfo, error) {
	// Get teacher with relationships
	var teacher models.Teacher
	if err := configs.DB.Preload("School").Preload("Gender").Preload("Prefix").Where("id = ?", teacherID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to get teacher")
	}

	// Get teacher's classrooms with preloaded students (optimized to prevent N+1 query)
	var classrooms []models.Classroom
	if err := configs.DB.
		Preload("Students", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Gender").Preload("Prefix")
		}).
		Where("teacher_id = ?", teacherID).
		Find(&classrooms).Error; err != nil {
		return nil, errors.New("failed to get classrooms")
	}

	// Get attendance counts for all classrooms in a single query (optimized)
	type AttendanceCountResult struct {
		ClassroomID uint
		Count       int64
	}
	var attendanceCounts []AttendanceCountResult
	configs.DB.Model(&models.Attendance{}).
		Select("classroom_id, COUNT(*) as count").
		Where("teacher_id = ?", teacherID).
		Group("classroom_id").
		Scan(&attendanceCounts)

	// Create a map for quick lookup
	attendanceCountMap := make(map[uint]int64)
	for _, ac := range attendanceCounts {
		attendanceCountMap[ac.ClassroomID] = ac.Count
	}

	// Get detailed classroom info with students
	var classroomInfos []ClassroomInfo
	var totalStudents int64
	var totalAttendance int64

	for _, classroom := range classrooms {
		// Students are already loaded via Preload
		students := classroom.Students
		if students == nil {
			students = []models.Student{}
		}

		// Get attendance count from map
		attendanceCount := attendanceCountMap[classroom.ID]

		classroomInfo := ClassroomInfo{
			Classroom:    classroom,
			StudentCount: int64(len(students)),
			Students:     students,
		}

		classroomInfos = append(classroomInfos, classroomInfo)
		totalStudents += int64(len(students))
		totalAttendance += attendanceCount
	}

	// Prepare teacher info
	teacherInfo := &TeacherInfo{
		Teacher:    teacher,
		School:     *teacher.School,
		Classrooms: classroomInfos,
		TotalStats: TeacherStats{
			TotalClassrooms: int64(len(classrooms)),
			TotalStudents:   totalStudents,
			TotalAttendance: totalAttendance,
		},
	}

	return teacherInfo, nil
}
