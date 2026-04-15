package model

import (
	"time"
)

// Song - 音乐文件元数据
type Song struct {
	ID        uint      `gorm:"primaryKey"`
	FilePath  string    `gorm:"uniqueIndex;not null"`
	Title     string
	Artist    string
	Album     string
	Year      int
	Genre     string
	TrackNum  int
	Duration  int // 秒
	CoverPath string
	Lyrics    string
	FileHash  string    `gorm:"index"`
	FileSize  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
