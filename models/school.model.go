package models

type School struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(255);not null;unique" json:"name"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *int64 `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Students []Student `gorm:"foreignKey:SchoolID" json:"students,omitempty"`
	Teachers []Teacher `gorm:"foreignKey:SchoolID" json:"teachers,omitempty"`
}

func (s *School) TableName() string {
	return "schools"
}
