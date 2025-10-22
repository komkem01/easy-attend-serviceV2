package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables using GORM
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Paginate struct {
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// Legacy Unix timestamp models (keeping for compatibility)
type CreateUpdateUnixTimestamp struct {
	CreateUnixTimestamp
	UpdateUnixTimestamp
}

type CreateUnixTimestamp struct {
	CreatedAt int64 `json:"created_at" bun:",notnull,default:EXTRACT(EPOCH FROM NOW())"`
}

type UpdateUnixTimestamp struct {
	UpdatedAt int64 `json:"updated_at" bun:",notnull,default:EXTRACT(EPOCH FROM NOW())"`
}

// type SoftDelete struct {
// 	DeletedAt *time.Time `json:"deleted_at" bun:",soft_delete,nullzero"`
// }

type SoftDelete struct {
	DeletedAt int64 `json:"deleted_at" bun:",soft_delete,nullzero"`
}

func (t *CreateUnixTimestamp) SetCreated(ts int64) {
	t.CreatedAt = ts
}

func (t *CreateUnixTimestamp) SetCreatedNow() {
	t.SetCreated(time.Now().Unix())
}

func (t *UpdateUnixTimestamp) SetUpdate(ts int64) {
	t.UpdatedAt = ts
}

func (t *UpdateUnixTimestamp) SetUpdateNow() {
	t.SetUpdate(time.Now().Unix())
}
