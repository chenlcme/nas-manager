package repository

import (
	"path"
	"sort"
	"strings"

	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// FolderWithCount - 文件夹及其歌曲数量
type FolderWithCount struct {
	ID        uint   `json:"id"`
	Path      string `json:"path"`
	SongCount int    `json:"songCount"`
}

// FolderRepository - 文件夹数据访问层
type FolderRepository struct {
	db *gorm.DB
}

// NewFolderRepository - 创建文件夹仓储
func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// getParentDir - 从文件路径提取父目录
// 使用 Go 的 path.Dir 函数，兼容 / 和 \ 路径分隔符
func getParentDir(filePath string) string {
	// 规范化路径：去除尾部斜杠，将反斜杠转为正斜杠
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	filePath = strings.TrimRight(filePath, "/")

	// 处理空路径
	if filePath == "" {
		return "/"
	}

	dir := path.Dir(filePath)

	// path.Dir(".") 返回 "."，需要转为 "/"
	if dir == "." {
		return "/"
	}

	// 去除双斜杠 (e.g., "/music//rock" -> "/music/rock")
	dir = strings.ReplaceAll(dir, "//", "/")

	return dir
}

// GetAllFoldersWithSongCount - 获取所有文件夹及其歌曲数量
// 由于 SQLite 不支持 REVERSE 函数，采用在 Go 层面提取目录的方式
func (r *FolderRepository) GetAllFoldersWithSongCount(orderAsc bool) ([]FolderWithCount, error) {
	var songs []model.Song

	err := r.db.Model(&model.Song{}).
		Where("file_path IS NOT NULL AND file_path != ''").
		Find(&songs).Error
	if err != nil {
		return nil, err
	}

	// 按父目录分组统计
	folderMap := make(map[string]int)
	for _, song := range songs {
		parentDir := getParentDir(song.FilePath)
		folderMap[parentDir]++
	}

	// 转换为结果切片
	results := []FolderWithCount{}
	for folderPath, count := range folderMap {
		results = append(results, FolderWithCount{
			Path:      folderPath,
			SongCount: count,
		})
	}

	// 排序
	if orderAsc {
		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].Path) < strings.ToLower(results[j].Path)
		})
	} else {
		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].Path) > strings.ToLower(results[j].Path)
		})
	}

	// 分配 ID
	for i := range results {
		results[i].ID = uint(i + 1)
	}

	return results, nil
}

// GetSongsByFolder - 根据文件夹路径获取歌曲列表
func (r *FolderRepository) GetSongsByFolder(folderPath string) ([]model.Song, error) {
	var songs []model.Song

	// 规范化路径分隔符（统一使用 /）
	folderPath = strings.ReplaceAll(folderPath, "\\", "/")
	folderPath = strings.TrimRight(folderPath, "/")

	// 转义 LIKE 特殊字符 (_ 和 %)
	escapedPath := strings.ReplaceAll(folderPath, "_", "_")
	escapedPath = strings.ReplaceAll(escapedPath, "%", "\\%")

	var err error
	if folderPath == "/" || folderPath == "" {
		// 根目录：匹配直接在根目录下的文件（只有一个 / 的路径）
		// 即 file_path 格式为 /filename.ext（根目录下直接文件）
		err = r.db.Where("file_path LIKE '/%' AND file_path NOT LIKE '/%/%'").Find(&songs).Error
	} else {
		// 子目录：匹配所有以 folderPath/ 开头的文件
		err = r.db.Where("file_path LIKE ? ESCAPE '\\'", escapedPath+"/%").Find(&songs).Error
	}
	if err != nil {
		return nil, err
	}

	return songs, nil
}

// GetFolderPathByID - 根据动态分配的 ID 获取文件夹路径
// ID 是基于 GetAllFoldersWithSongCount 返回结果的顺序
func (r *FolderRepository) GetFolderPathByID(id uint) (string, error) {
	folders, err := r.GetAllFoldersWithSongCount(false)
	if err != nil {
		return "", err
	}
	if int(id) > len(folders) || id < 1 {
		return "", gorm.ErrRecordNotFound
	}
	return folders[id-1].Path, nil
}
