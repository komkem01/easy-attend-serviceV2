package models

type Student struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID  *uint  `gorm:"not null" json:"school_id"`
	StudentNo string `gorm:"type:varchar(20);not null;unique" json:"student_no"`
	FirstName string `gorm:"type:varchar(100);not null" json:"firstname"`
	LastName  string `gorm:"type:varchar(100);not null" json:"lastname"`
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
	ClassroomMembers []ClassroomMember `gorm:"foreignKey:StudentID" json:"classroom_members,omitempty"`
	Attendances      []Attendance      `gorm:"foreignKey:StudentID" json:"attendances,omitempty"`
}

func (s *Student) TableName() string {
	return "students"
}
