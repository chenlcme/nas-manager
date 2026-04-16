package repository

import (
	"fmt"

	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// ArtistWithCount - 艺术家及其歌曲数量
type ArtistWithCount struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	SongCount int    `json:"songCount"`
}

// ArtistRepository - 艺术家数据访问层
type ArtistRepository struct {
	db *gorm.DB
}

// NewArtistRepository - 创建艺术家仓储
func NewArtistRepository(db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

// GetAllArtistsWithSongCount - 获取所有艺术家及其歌曲数量
func (r *ArtistRepository) GetAllArtistsWithSongCount(orderAsc bool) ([]ArtistWithCount, error) {
	var results []ArtistWithCount

	query := r.db.Model(&model.Song{}).
		Select("artist as name, COUNT(*) as song_count, 0 as id").
		Where("artist IS NOT NULL AND artist != '' AND artist != '   ' AND artist != '　'").
		Group("artist")

	if orderAsc {
		query = query.Order("artist ASC")
	} else {
		query = query.Order("artist DESC")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Assign IDs based on row number (for consistent ordering)
	for i := range results {
		results[i].ID = uint(i + 1)
	}

	return results, nil
}

// GetSongsByArtist - 根据艺术家名获取歌曲列表，支持排序
// sortBy: title, duration, created_at
// order: asc, desc
// Returns error if sortBy or order is invalid
func (r *ArtistRepository) GetSongsByArtist(artistName string, sortBy string, order string) ([]model.Song, error) {
	var songs []model.Song

	// 校验排序参数，使用白名单
	validSortFields := map[string]bool{
		"title":     true,
		"duration":   true,
		"created_at": true,
	}
	validOrders := map[string]bool{"asc": true, "desc": true}

	// 验证 sortBy 参数，无效则返回错误（不静默默认值）
	if !validSortFields[sortBy] {
		return nil, fmt.Errorf("invalid sort_by parameter: %s", sortBy)
	}
	// 验证 order 参数，无效则返回错误
	if !validOrders[order] {
		return nil, fmt.Errorf("invalid order parameter: %s", order)
	}

	query := r.db.Where("artist = ?", artistName)
	// 使用 gorm.Expr 安全地构建 ORDER 子句，避免 SQL 注入风险
	query = query.Order(gorm.Expr("? ?", sortBy, order))

	err := query.Find(&songs).Error
	return songs, err
}

// GetArtistByName - 根据艺术家名获取艺术家信息
func (r *ArtistRepository) GetArtistByName(name string) (*model.Artist, error) {
	var artist model.Artist
	err := r.db.Where("name = ?", name).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

// Create - 创建艺术家记录
func (r *ArtistRepository) Create(artist *model.Artist) error {
	return r.db.Create(artist).Error
}

// GetArtistNameByID - 根据动态分配的 ID 获取艺术家名
// ID 是基于 GetAllArtistsWithSongCount 返回结果的顺序
func (r *ArtistRepository) GetArtistNameByID(id uint) (string, error) {
	artists, err := r.GetAllArtistsWithSongCount(false)
	if err != nil {
		return "", err
	}
	if int(id) > len(artists) || id < 1 {
		return "", gorm.ErrRecordNotFound
	}
	return artists[id-1].Name, nil
}
