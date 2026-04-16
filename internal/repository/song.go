package repository

import (
	"context"
	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// SongRepository - 歌曲数据访问层
type SongRepository struct {
	db *gorm.DB
}

// NewSongRepository - 创建歌曲仓储
func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{db: db}
}

// Create - 创建歌曲记录
func (r *SongRepository) Create(song *model.Song) error {
	return r.db.Create(song).Error
}

// GetByFilePath - 根据文件路径获取歌曲
func (r *SongRepository) GetByFilePath(filePath string) (*model.Song, error) {
	var song model.Song
	if err := r.db.Where("file_path = ?", filePath).First(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// ExistsByFilePath - 检查歌曲是否存在
func (r *SongRepository) ExistsByFilePath(filePath string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Song{}).Where("file_path = ?", filePath).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAll - 获取所有歌曲
func (r *SongRepository) GetAll() ([]model.Song, error) {
	var songs []model.Song
	if err := r.db.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

// GetByID - 根据ID获取歌曲
func (r *SongRepository) GetByID(id uint) (*model.Song, error) {
	var song model.Song
	if err := r.db.First(&song, id).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// GetByIDWithContext - 根据ID获取歌曲（带上下文超时）
func (r *SongRepository) GetByIDWithContext(ctx context.Context, id uint64) (*model.Song, error) {
	var song model.Song
	if err := r.db.WithContext(ctx).First(&song, id).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// Update - 更新歌曲
func (r *SongRepository) Update(song *model.Song) error {
	return r.db.Save(song).Error
}

// Delete - 删除歌曲
func (r *SongRepository) Delete(id uint) error {
	return r.db.Delete(&model.Song{}, id).Error
}

// GetByArtist - 根据艺术家获取歌曲
func (r *SongRepository) GetByArtist(artist string) ([]model.Song, error) {
	var songs []model.Song
	if err := r.db.Where("artist = ?", artist).Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

// GetByAlbum - 根据专辑获取歌曲
func (r *SongRepository) GetByAlbum(album string) ([]model.Song, error) {
	var songs []model.Song
	if err := r.db.Where("album = ?", album).Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

// SearchByFileName - 根据文件名搜索歌曲（模糊匹配，带分页）
func (r *SongRepository) SearchByFileName(keyword string, limit, offset int) ([]model.Song, error) {
	var songs []model.Song
	if err := r.db.Where("file_path LIKE ?", "%"+keyword+"%").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}
