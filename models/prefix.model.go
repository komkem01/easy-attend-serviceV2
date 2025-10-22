package models

type Prefix struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(20);not null;unique" json:"name"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *int64 `gorm:"index" json:"deleted_at,omitempty"`
}

func (Prefix) TableName() string {
	return "prefixes"
}
