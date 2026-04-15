package service

import (
	"errors"
	"os"
	"path/filepath"

	"nas-manager/internal/repository"
)

// SetupStatus - 设置状态
type SetupStatus struct {
	NeedsSetup bool   `json:"needs_setup"`
	MusicDir   string `json:"music_dir"`
	DBPath     string `json:"db_path"`
}

// SetupConfig - 配置请求
type SetupConfig struct {
	MusicDir string `json:"music_dir"`
	DBPath   string `json:"db_path"`
}

// SettingService - 设置服务
type SettingService struct {
	repo *repository.SettingRepository
}

// NewSettingService - 创建设置服务
func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// CheckSetupRequired - 检查是否需要首次配置
func (s *SettingService) CheckSetupRequired() (*SetupStatus, error) {
	musicDir, err := s.repo.GetMusicDir()
	if err != nil {
		// 设置不存在，需要配置
		return &SetupStatus{
			NeedsSetup: true,
			MusicDir:   "",
			DBPath:     "",
		}, nil
	}

	return &SetupStatus{
		NeedsSetup: musicDir == "",
		MusicDir:   musicDir,
		DBPath:     "", // db_path 可以为空，使用默认
	}, nil
}

// SaveSetupConfig - 保存配置
func (s *SettingService) SaveSetupConfig(cfg *SetupConfig) error {
	// 验证音乐目录
	if cfg.MusicDir == "" {
		return errors.New("music_dir cannot be empty")
	}

	// 检查目录是否存在
	info, err := os.Stat(cfg.MusicDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("music directory does not exist")
		}
		return errors.New("cannot access music directory: " + err.Error())
	}
	if !info.IsDir() {
		return errors.New("music path is not a directory")
	}

	// 检查目录是否可读
	if !isDirReadable(cfg.MusicDir) {
		return errors.New("music directory is not readable")
	}

	// 如果指定了 db_path，检查是否可写
	if cfg.DBPath != "" {
		dir := filepath.Dir(cfg.DBPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New("cannot create database directory: " + err.Error())
		}
		if !isDirWritable(dir) {
			return errors.New("database directory is not writable")
		}
	}

	// 保存配置
	if err := s.repo.SetMusicDir(cfg.MusicDir); err != nil {
		return errors.New("failed to save music_dir: " + err.Error())
	}

	if cfg.DBPath != "" {
		if err := s.repo.SetDBPath(cfg.DBPath); err != nil {
			return errors.New("failed to save db_path: " + err.Error())
		}
	}

	return nil
}

// isDirReadable - 检查目录是否可读
func isDirReadable(path string) bool {
	// 尝试读取目录
	_, err := os.ReadDir(path)
	return err == nil
}

// isDirWritable - 检查目录是否可写
func isDirWritable(path string) bool {
	// 尝试创建临时文件测试
	testFile := filepath.Join(path, ".write_test")
	defer os.Remove(testFile)
	_, err := os.Create(testFile)
	return err == nil
}
