package model

// Artist - 艺术家
type Artist struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}
