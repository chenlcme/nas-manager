package model

import (
	"time"
)

// BatchOperation - 批量操作记录
type BatchOperation struct {
	ID        uint      `gorm:"primaryKey"`
	Type      string    // "update", "delete"
	TargetIDs string    // JSON array of song IDs
	OldValues string    // JSON of previous values
	NewValues string    // JSON of new values
	CreatedAt time.Time
}
