package models

type ClassroomMember struct {
	TeacherID   *uint `gorm:"default:null" json:"teacher_id"`
	StudentID   *uint `gorm:"default:null" json:"student_id"`
	ClassroomID *uint `gorm:"not null" json:"classroom_id"`

	// Foreign Key Relationships
	Teacher   *Teacher   `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teacher,omitempty"`
	Student   *Student   `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"student,omitempty"`
	Classroom *Classroom `gorm:"foreignKey:ClassroomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"classroom,omitempty"`
}

func (cm *ClassroomMember) TableName() string {
	return "classroom_members"
}
