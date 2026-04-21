package service

import (
	"io/fs"
	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

// ScanMode - 扫描模式
type ScanMode string

const (
	ScanModeFull       ScanMode = "full"
	ScanModeIncremental ScanMode = "incremental"
)

// 支持的音频格式
var supportedExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".ape":  true,
	".ogg":  true,
	".m4a":  true,
	".wav":  true,
	".aiff": true,
	".wma":  true,
	".aac":  true,
}

// ScanResult - 扫描结果
type ScanResult struct {
	Found  int      `json:"found"`
	New    int      `json:"new"`
	Updated int     `json:"updated"`
	Errors []string `json:"errors,omitempty"`
}

// ScannerService - 扫描服务
type ScannerService struct {
	id3Service *ID3Service
	songRepo   *repository.SongRepository
	lastScanTime int64
}

// NewScannerService - 创建扫描服务
func NewScannerService(id3Service *ID3Service, songRepo *repository.SongRepository) *ScannerService {
	return &ScannerService{
		id3Service: id3Service,
		songRepo:   songRepo,
	}
}

// IsMusicFile - 判断是否为支持的音频文件
func (s *ScannerService) IsMusicFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return supportedExtensions[ext]
}

// GetMusicFiles - 获取目录下所有音乐文件
func (s *ScannerService) GetMusicFiles(rootDir string) ([]string, error) {
	var musicFiles []string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // 跳过权限错误
		}

		if d.IsDir() {
			return nil
		}

		if s.IsMusicFile(path) {
			musicFiles = append(musicFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return musicFiles, nil
}

// GetMusicFilesWithErrors - 获取音乐文件并记录错误
func (s *ScannerService) GetMusicFilesWithErrors(rootDir string) ([]string, []string, error) {
	var musicFiles []string
	var errors []string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errors = append(errors, "权限不足: "+path)
			return nil // 继续扫描
		}

		if d.IsDir() {
			return nil
		}

		if s.IsMusicFile(path) {
			musicFiles = append(musicFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, errors, err
	}

	return musicFiles, errors, nil
}

// GetFileModTime - 获取文件修改时间
func (s *ScannerService) GetFileModTime(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}

// SetLastScanTime - 设置上次扫描时间
func (s *ScannerService) SetLastScanTime(t int64) {
	s.lastScanTime = t
}

// ScanFiles - 扫描文件并处理（支持全量/增量模式）
func (s *ScannerService) ScanFiles(rootDir string, mode ScanMode) (*ScanResult, error) {
	result := &ScanResult{}

	// 获取所有音乐文件
	files, scanErrors, err := s.GetMusicFilesWithErrors(rootDir)
	if err != nil {
		return nil, err
	}
	result.Errors = scanErrors
	result.Found = len(files)

	// 设置增量扫描的起始时间
	if mode == ScanModeIncremental && s.lastScanTime == 0 {
		// 无法获取上次扫描时间，退化为全量扫描
		mode = ScanModeFull
	}

	for _, filePath := range files {
		fileModTime, err := s.GetFileModTime(filePath)
		if err != nil {
			result.Errors = append(result.Errors, "无法获取文件时间: "+filePath)
			continue
		}

		// 增量扫描：跳过未修改的文件
		if mode == ScanModeIncremental && fileModTime <= s.lastScanTime {
			continue
		}

		// 检查文件是否已存在
		exists, err := s.songRepo.ExistsByFilePath(filePath)
		if err != nil {
			result.Errors = append(result.Errors, "数据库错误: "+filePath)
			continue
		}

		if !exists {
			// 新文件：创建记录并解析 ID3
			song := &model.Song{
				FilePath: filePath,
				Folder:   s.extractFolder(filePath, rootDir),
			}
			if err := s.songRepo.Create(song); err != nil {
				result.Errors = append(result.Errors, "创建记录失败: "+filePath)
				continue
			}
			// 解析 ID3 标签
			if err := s.id3Service.ParseSongMetadata(song); err != nil {
				result.Errors = append(result.Errors, "解析ID3失败: "+filePath)
			}
			result.New++
		} else if mode == ScanModeFull {
			// 全量扫描：更新已有记录
			song, err := s.songRepo.GetByFilePath(filePath)
			if err != nil {
				result.Errors = append(result.Errors, "获取记录失败: "+filePath)
				continue
			}
			// 强制重新解析 ID3 标签
			if err := s.id3Service.ParseSongMetadata(song); err != nil {
				result.Errors = append(result.Errors, "解析ID3失败: "+filePath)
				continue
			}
			result.Updated++
		}
	}

	return result, nil
}

// CleanupResult - 清理结果
type CleanupResult struct {
	Cleaned int      `json:"cleaned"`
	Errors  []string `json:"errors,omitempty"`
}

// CleanupOrphanRecords - 清理孤立的数据库记录（文件已不存在）
func (s *ScannerService) CleanupOrphanRecords() (*CleanupResult, error) {
	result := &CleanupResult{}

	// 获取所有数据库记录
	songs, err := s.songRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, song := range songs {
		// 检查文件是否存在
		if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
			// 文件不存在，删除记录
			if err := s.songRepo.Delete(song.ID); err != nil {
				result.Errors = append(result.Errors, "删除记录失败: "+song.FilePath)
				continue
			}
			result.Cleaned++
		}
	}

	return result, nil
}

// extractFolder - 从文件路径提取文件夹路径（相对于根目录）
func (s *ScannerService) extractFolder(filePath, rootDir string) string {
	// 获取文件的目录
	dir := filepath.Dir(filePath)

	// 如果文件就在根目录下，返回空字符串或 "."
	if dir == rootDir {
		return ""
	}

	// 返回相对于根目录的路径
	relPath, err := filepath.Rel(rootDir, dir)
	if err != nil {
		return dir // 如果计算失败，返回完整目录
	}

	return relPath
}
