package models

type Teacher struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID  *uint  `gorm:"not null" json:"school_id"`
	Email     string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string `gorm:"type:varchar(255);not null" json:"-"` // Hide password from JSON
	FirstName string `gorm:"type:varchar(100);not null" json:"firstname"`
	LastName  string `gorm:"type:varchar(100);not null" json:"lastname"`
	Phone     string `gorm:"type:varchar(20)" json:"phone"`
	GenderID  *uint  `json:"gender_id"`
	PrefixID  *uint  `json:"prefix_id"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *int64 `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key Relationships
	School *School `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"school,omitempty"`
	Gender *Gender `gorm:"foreignKey:GenderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"gender,omitempty"`
	Prefix *Prefix `gorm:"foreignKey:PrefixID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"prefix,omitempty"`

	// Has Many Relationships
	Classrooms       []Classroom       `gorm:"foreignKey:TeacherID" json:"classrooms,omitempty"`
	ClassroomMembers []ClassroomMember `gorm:"foreignKey:TeacherID" json:"classroom_members,omitempty"`
	Attendances      []Attendance      `gorm:"foreignKey:TeacherID" json:"attendances,omitempty"`
}

func (t *Teacher) TableName() string {
	return "teachers"
}
