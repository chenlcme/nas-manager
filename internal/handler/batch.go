package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// BatchHandler - 批量操作处理
type BatchHandler struct {
	songRepo  *repository.SongRepository
	batchRepo *repository.BatchRepository
}

// NewBatchHandler - 创建批量操作处理器
func NewBatchHandler(songRepo *repository.SongRepository, batchRepo *repository.BatchRepository) *BatchHandler {
	return &BatchHandler{
		songRepo:  songRepo,
		batchRepo: batchRepo,
	}
}

// BatchUpdateResult - 批量更新结果
type BatchUpdateResult struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

// BatchUpdateRequest - 批量更新请求
type BatchUpdateRequest struct {
	IDs            []uint `json:"ids" binding:"required,min=1"`
	Title          *string `json:"title"`
	Artist         *string `json:"artist"`
	Album          *string `json:"album"`
	Year           *int    `json:"year"`
	Genre          *string `json:"genre"`
	TrackNum       *int    `json:"trackNum"`
	CoverPath      *string `json:"coverPath"`
	Lyrics         *string `json:"lyrics"`
}

// BatchUpdate - 批量更新歌曲
// POST /api/songs/batch-update
func (h *BatchHandler) BatchUpdate(c *gin.Context) {
	start := time.Now()

	var req BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[BatchHandler] BatchUpdate invalid request: %v", err)
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "无效的请求参数")
		return
	}

	// Deduplicate IDs
	seen := make(map[uint]bool)
	uniqueIDs := make([]uint, 0, len(req.IDs))
	for _, id := range req.IDs {
		if !seen[id] {
			seen[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	log.Printf("[BatchHandler] BatchUpdate count=%d", len(uniqueIDs))

	result := &BatchUpdateResult{
		Total:     len(uniqueIDs),
		Succeeded: 0,
		Failed:    0,
	}

	// Collect old values for undo
	oldValues := make(map[uint]map[string]interface{})

	// First pass: collect old values
	for _, id := range uniqueIDs {
		song, err := h.songRepo.GetByID(id)
		if err != nil {
			log.Printf("[BatchHandler] BatchUpdate song not found: id=%d", id)
			result.Failed++
			continue
		}
		oldValues[id] = map[string]interface{}{
			"title":     song.Title,
			"artist":    song.Artist,
			"album":     song.Album,
			"year":      song.Year,
			"genre":     song.Genre,
			"trackNum":  song.TrackNum,
			"coverPath": song.CoverPath,
			"lyrics":    song.Lyrics,
		}
	}

	// Second pass: update songs
	for _, id := range uniqueIDs {
		song, err := h.songRepo.GetByID(id)
		if err != nil {
			result.Failed++
			continue
		}

		// Update fields if provided
		if req.Title != nil {
			song.Title = *req.Title
		}
		if req.Artist != nil {
			song.Artist = *req.Artist
		}
		if req.Album != nil {
			song.Album = *req.Album
		}
		if req.Year != nil {
			song.Year = *req.Year
		}
		if req.Genre != nil {
			song.Genre = *req.Genre
		}
		if req.TrackNum != nil {
			song.TrackNum = *req.TrackNum
		}
		if req.CoverPath != nil {
			song.CoverPath = *req.CoverPath
		}
		if req.Lyrics != nil {
			song.Lyrics = *req.Lyrics
		}

		if err := h.songRepo.Update(song); err != nil {
			log.Printf("[BatchHandler] BatchUpdate failed to update id=%d: %v", id, err)
			result.Failed++
			continue
		}
		result.Succeeded++
	}

	// Save batch operation for undo if any updates succeeded
	if result.Succeeded > 0 {
		// Serialize target IDs
		idsJSON, err := json.Marshal(uniqueIDs)
		if err != nil {
			log.Printf("[BatchHandler] BatchUpdate failed to marshal ids: %v", err)
		} else {
			// Serialize old values
			oldValuesJSON, err := json.Marshal(oldValues)
			if err != nil {
				log.Printf("[BatchHandler] BatchUpdate failed to marshal old values: %v", err)
			} else {
				// Serialize new values (only the fields that were updated)
				newValues := make(map[uint]map[string]interface{})
				for _, id := range uniqueIDs {
					newValues[id] = make(map[string]interface{})
					if req.Title != nil {
						newValues[id]["title"] = *req.Title
					}
					if req.Artist != nil {
						newValues[id]["artist"] = *req.Artist
					}
					if req.Album != nil {
						newValues[id]["album"] = *req.Album
					}
					if req.Year != nil {
						newValues[id]["year"] = *req.Year
					}
					if req.Genre != nil {
						newValues[id]["genre"] = *req.Genre
					}
					if req.TrackNum != nil {
						newValues[id]["trackNum"] = *req.TrackNum
					}
					if req.CoverPath != nil {
						newValues[id]["coverPath"] = *req.CoverPath
					}
					if req.Lyrics != nil {
						newValues[id]["lyrics"] = *req.Lyrics
					}
				}
				newValuesJSON, err := json.Marshal(newValues)
				if err != nil {
					log.Printf("[BatchHandler] BatchUpdate failed to marshal new values: %v", err)
				} else {
					batchOp := &model.BatchOperation{
						Type:      "update",
						TargetIDs: string(idsJSON),
						OldValues: string(oldValuesJSON),
						NewValues: string(newValuesJSON),
					}

					if err := h.batchRepo.Create(batchOp); err != nil {
						log.Printf("[BatchHandler] BatchUpdate failed to save batch operation: %v", err)
						// Don't fail the whole operation, just log the error
					}
				}
			}
		}
	}

	log.Printf("[BatchHandler] BatchUpdate completed: total=%d, succeeded=%d, failed=%d, duration=%v",
		result.Total, result.Succeeded, result.Failed, time.Since(start))
	response.Success(c, result)
}

// UndoBatch - 撤销批量操作
// POST /api/songs/undo/:batchId
func (h *BatchHandler) UndoBatch(c *gin.Context) {
	start := time.Now()

	batchIdStr := c.Param("batchId")
	batchId, err := strconv.ParseUint(batchIdStr, 10, 64)
	if err != nil {
		log.Printf("[BatchHandler] UndoBatch invalid batch ID: %s, error: %v", batchIdStr, err)
		response.Error(c, http.StatusBadRequest, "INVALID_BATCH_ID", "无效的批量操作ID")
		return
	}

	log.Printf("[BatchHandler] UndoBatch batchId=%d", batchId)

	// Get batch operation
	batch, err := h.batchRepo.GetByID(uint(batchId))
	if err != nil {
		log.Printf("[BatchHandler] UndoBatch batch not found: batchId=%d", batchId)
		response.Error(c, http.StatusNotFound, "BATCH_NOT_FOUND", "批量操作不存在")
		return
	}

	// Parse old values
	var oldValues map[uint]map[string]interface{}
	if err := json.Unmarshal([]byte(batch.OldValues), &oldValues); err != nil {
		log.Printf("[BatchHandler] UndoBatch failed to parse old values: %v", err)
		response.Error(c, http.StatusInternalServerError, "PARSE_ERROR", "解析旧数据失败")
		return
	}

	// Restore old values
	succeeded := 0
	failed := 0

	for id, values := range oldValues {
		song, err := h.songRepo.GetByID(id)
		if err != nil {
			log.Printf("[BatchHandler] UndoBatch song not found: id=%d", id)
			failed++
			continue
		}

		// Restore all fields from old values
		if v, ok := values["title"].(string); ok {
			song.Title = v
		}
		if v, ok := values["artist"].(string); ok {
			song.Artist = v
		}
		if v, ok := values["album"].(string); ok {
			song.Album = v
		}
		if v, ok := values["year"]; ok && v != nil {
			switch vv := v.(type) {
			case float64:
				song.Year = int(vv)
			case int:
				song.Year = vv
			case int64:
				song.Year = int(vv)
			}
		}
		if v, ok := values["genre"].(string); ok {
			song.Genre = v
		}
		if v, ok := values["trackNum"]; ok && v != nil {
			switch vv := v.(type) {
			case float64:
				song.TrackNum = int(vv)
			case int:
				song.TrackNum = vv
			case int64:
				song.TrackNum = int(vv)
			}
		}
		if v, ok := values["coverPath"].(string); ok {
			song.CoverPath = v
		}
		if v, ok := values["lyrics"].(string); ok {
			song.Lyrics = v
		}

		if err := h.songRepo.Update(song); err != nil {
			log.Printf("[BatchHandler] UndoBatch failed to restore song id=%d: %v", id, err)
			failed++
			continue
		}
		succeeded++
	}

	// Delete the batch operation after undo
	if err := h.batchRepo.Delete(uint(batchId)); err != nil {
		log.Printf("[BatchHandler] UndoBatch failed to delete batch: %v", err)
		// Don't fail, just log
	}

	result := map[string]int{
		"succeeded": succeeded,
		"failed":    failed,
	}

	log.Printf("[BatchHandler] UndoBatch completed: batchId=%d, succeeded=%d, failed=%d, duration=%v",
		batchId, succeeded, failed, time.Since(start))
	response.Success(c, result)
}

// GetLatestBatch - 获取最新的批量操作（用于撤销）
// GET /api/batches/latest
func (h *BatchHandler) GetLatestBatch(c *gin.Context) {
	batch, err := h.batchRepo.GetLatest()
	if err != nil {
		log.Printf("[BatchHandler] GetLatestBatch error: %v", err)
		response.Error(c, http.StatusNotFound, "NO_BATCH", "没有可撤销的操作")
		return
	}
	response.Success(c, batch)
}