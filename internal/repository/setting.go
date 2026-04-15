package repository

import (
	"fmt"

	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// SettingRepository - 设置数据访问层
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository - 创建设置仓储
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetSetting - 获取设置值
func (r *SettingRepository) GetSetting(key string) (string, error) {
	var setting model.Setting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

// SetSetting - 设置值（upsert）
func (r *SettingRepository) SetSetting(key, value string) error {
	setting := model.Setting{Key: key, Value: value}
	return r.db.Save(&setting).Error
}

// DeleteSetting - 删除设置
func (r *SettingRepository) DeleteSetting(key string) error {
	return r.db.Where("key = ?", key).Delete(&model.Setting{}).Error
}

// GetAllSettings - 获取所有设置
func (r *SettingRepository) GetAllSettings() (map[string]string, error) {
	var settings []model.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// HasSettings - 检查是否存在任何设置（用于判断是否已配置）
func (r *SettingRepository) HasSettings() (bool, error) {
	var count int64
	if err := r.db.Model(&model.Setting{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetMusicDir - 获取音乐目录
func (r *SettingRepository) GetMusicDir() (string, error) {
	return r.GetSetting("music_dir")
}

// SetMusicDir - 设置音乐目录
func (r *SettingRepository) SetMusicDir(path string) error {
	return r.SetSetting("music_dir", path)
}

// GetDBPath - 获取数据库路径
func (r *SettingRepository) GetDBPath() (string, error) {
	return r.GetSetting("db_path")
}

// SetDBPath - 设置数据库路径
func (r *SettingRepository) SetDBPath(path string) error {
	return r.SetSetting("db_path", path)
}

// GetLastScanTime - 获取上次扫描时间
func (r *SettingRepository) GetLastScanTime() (int64, error) {
	val, err := r.GetSetting("last_scan_time")
	if err != nil {
		return 0, nil // 返回 0 表示从未扫描
	}
	var t int64
	if _, err := fmt.Sscanf(val, "%d", &t); err != nil {
		return 0, nil
	}
	return t, nil
}

// SetLastScanTime - 设置上次扫描时间
func (r *SettingRepository) SetLastScanTime(t int64) error {
	return r.SetSetting("last_scan_time", fmt.Sprintf("%d", t))
}
