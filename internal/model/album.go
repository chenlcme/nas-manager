package model

// Album - 专辑
type Album struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"index"`
	Artist string
}
