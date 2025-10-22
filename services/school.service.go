package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"errors"

	"gorm.io/gorm"
)

type SchoolService struct{}

func NewSchoolService() *SchoolService {
	return &SchoolService{}
}

func (s *SchoolService) GetAllSchools(page, limit int) ([]models.School, int64, error) {
	var schools []models.School
	var total int64

	// Count total records
	if err := configs.DB.Model(&models.School{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count schools")
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get schools with pagination
	if err := configs.DB.Offset(offset).Limit(limit).Find(&schools).Error; err != nil {
		return nil, 0, errors.New("failed to get schools")
	}

	return schools, total, nil
}

func (s *SchoolService) GetSchoolByID(id uint) (*models.School, error) {
	var school models.School
	if err := configs.DB.Where("id = ?", id).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, errors.New("failed to get school")
	}
	return &school, nil
}

func (s *SchoolService) CreateSchool(name string) (*models.School, error) {
	// Check if school already exists
	var existingSchool models.School
	if err := configs.DB.Where("name = ?", name).First(&existingSchool).Error; err == nil {
		return nil, errors.New("school with this name already exists")
	}

	// Create new school
	school := models.School{
		Name: name,
	}

	if err := configs.DB.Create(&school).Error; err != nil {
		return nil, errors.New("failed to create school")
	}

	return &school, nil
}

func (s *SchoolService) UpdateSchool(id uint, name string) (*models.School, error) {
	var school models.School
	if err := configs.DB.Where("id = ?", id).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, errors.New("failed to find school")
	}

	// Check if name is being changed and if it already exists
	if name != school.Name {
		var existingSchool models.School
		if err := configs.DB.Where("name = ? AND id != ?", name, id).First(&existingSchool).Error; err == nil {
			return nil, errors.New("school with this name already exists")
		}
	}

	// Update fields
	school.Name = name

	if err := configs.DB.Save(&school).Error; err != nil {
		return nil, errors.New("failed to update school")
	}

	return &school, nil
}

func (s *SchoolService) DeleteSchool(id uint) error {
	var school models.School
	if err := configs.DB.Where("id = ?", id).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("school not found")
		}
		return errors.New("failed to find school")
	}

	if err := configs.DB.Delete(&school).Error; err != nil {
		return errors.New("failed to delete school")
	}

	return nil
}
