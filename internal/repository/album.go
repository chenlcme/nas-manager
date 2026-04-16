package repository

import (
	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// AlbumWithCount - 专辑及其歌曲数量
type AlbumWithCount struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	SongCount int    `json:"songCount"`
}

// AlbumRepository - 专辑数据访问层
type AlbumRepository struct {
	db *gorm.DB
}

// NewAlbumRepository - 创建专辑仓储
func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

// GetAllAlbumsWithSongCount - 获取所有专辑及其歌曲数量
func (r *AlbumRepository) GetAllAlbumsWithSongCount(orderAsc bool) ([]AlbumWithCount, error) {
	var results []AlbumWithCount

	query := r.db.Model(&model.Song{}).
		Select("TRIM(album) as name, TRIM(artist) as artist, COUNT(*) as song_count, 0 as id").
		Where("TRIM(album) IS NOT NULL AND TRIM(album) != ''").
		Group("TRIM(album), TRIM(artist)")

	if orderAsc {
		query = query.Order("album ASC")
	} else {
		query = query.Order("album DESC")
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

// GetSongsByAlbum - 根据专辑名和艺术家获取歌曲列表
func (r *AlbumRepository) GetSongsByAlbum(albumName string, artistName string) ([]model.Song, error) {
	var songs []model.Song
	err := r.db.Where("album = ? AND artist = ?", albumName, artistName).Find(&songs).Error
	return songs, err
}

// GetAlbumByNameAndArtist - 根据专辑名和艺术家获取专辑信息
func (r *AlbumRepository) GetAlbumByNameAndArtist(name string, artist string) (*model.Album, error) {
	var album model.Album
	err := r.db.Where("name = ? AND artist = ?", name, artist).First(&album).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

// Create - 创建专辑记录
func (r *AlbumRepository) Create(album *model.Album) error {
	return r.db.Create(album).Error
}

// GetAlbumNameAndArtistByID - 根据动态分配的 ID 获取专辑名和艺术家名
// ID 是基于 GetAllAlbumsWithSongCount 返回结果的顺序
func (r *AlbumRepository) GetAlbumNameAndArtistByID(id uint) (string, string, error) {
	albums, err := r.GetAllAlbumsWithSongCount(false)
	if err != nil {
		return "", "", err
	}
	if int(id) > len(albums) || id < 1 {
		return "", "", gorm.ErrRecordNotFound
	}
	return albums[id-1].Name, albums[id-1].Artist, nil
}
