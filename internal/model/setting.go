package model

// Setting - 应用配置
type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string
}
