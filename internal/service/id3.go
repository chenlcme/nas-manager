package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"nas-manager/pkg/id3"
)

// ID3Service - ID3 解析服务
type ID3Service struct {
	parser   *id3.Parser
	songRepo *repository.SongRepository
}

// NewID3Service - 创建 ID3 服务
func NewID3Service(songRepo *repository.SongRepository) *ID3Service {
	return &ID3Service{
		parser:   id3.NewParser(),
		songRepo: songRepo,
	}
}

// ParseSongMetadata - 解析并更新歌曲元数据
func (s *ID3Service) ParseSongMetadata(song *model.Song) error {
	// 检查文件是否存在
	if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", song.FilePath)
	}

	// 解析 ID3 标签
	metadata, err := s.parser.ParseFile(song.FilePath)
	if err != nil {
		log.Printf("Warning: failed to parse ID3 for %s: %v", song.FilePath, err)
		// 即使解析失败也继续，不阻塞其他文件
		return nil
	}

	// 更新歌曲元数据
	song.Title = metadata.Title
	song.Artist = metadata.Artist
	song.Album = metadata.Album
	song.Year = metadata.Year
	song.TrackNum = metadata.TrackNum
	song.Genre = metadata.Genre
	song.Duration = metadata.Duration
	song.Lyrics = metadata.Lyrics
	song.FileHash = metadata.FileHash
	song.FileSize = metadata.FileSize

	// 保存封面（如果有）
	if len(metadata.Cover) > 0 {
		coverPath, err := s.saveCover(song.ID, metadata.Cover)
		if err != nil {
			log.Printf("Warning: failed to save cover for %s: %v", song.FilePath, err)
		} else {
			song.CoverPath = coverPath
		}
	}

	// 更新数据库
	if err := s.songRepo.Update(song); err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	return nil
}

// saveCover - 保存封面图片到本地
func (s *ID3Service) saveCover(songID uint, coverData []byte) (string, error) {
	// 创建封面目录
	coverDir := "./covers"
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%d.jpg", songID, time.Now().Unix())
	coverPath := coverDir + "/" + filename

	// 写入文件
	if err := os.WriteFile(coverPath, coverData, 0644); err != nil {
		return "", err
	}

	return coverPath, nil
}

// ParseAllPendingSongs - 解析所有待处理的歌曲
func (s *ID3Service) ParseAllPendingSongs() (int, error) {
	songs, err := s.songRepo.GetAll()
	if err != nil {
		return 0, err
	}

	successCount := 0
	errorCount := 0

	for _, song := range songs {
		// 只解析还没有标题的文件
		if song.Title == "" {
			if err := s.ParseSongMetadata(&song); err != nil {
				log.Printf("Error parsing %s: %v", song.FilePath, err)
				errorCount++
			} else {
				successCount++
			}
		}
	}

	return successCount, nil
}
