package models

type Classroom struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID  *uint  `gorm:"not null" json:"school_id"`
	TeacherID *uint  `gorm:"not null" json:"teacher_id"`
	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	Grade     string `gorm:"type:varchar(10);not null" json:"grade"` // ชั้นเรียน เช่น "ม.1", "ป.6"
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *int64 `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key Relationships
	School  *School  `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"school,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"teacher,omitempty"`

	// Has Many Relationships
	Members     []ClassroomMember `gorm:"foreignKey:ClassroomID" json:"members,omitempty"`
	Attendances []Attendance      `gorm:"foreignKey:ClassroomID" json:"attendances,omitempty"`
}

func (c *Classroom) TableName() string {
	return "classrooms"
}
