package repository

import (
	"nas-manager/internal/model"

	"gorm.io/gorm"
)

// BatchRepository - 批量操作数据访问层
type BatchRepository struct {
	db *gorm.DB
}

// NewBatchRepository - 创建批量操作仓储
func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

// Create - 创建批量操作记录
func (r *BatchRepository) Create(batch *model.BatchOperation) error {
	return r.db.Create(batch).Error
}

// GetByID - 根据ID获取批量操作
func (r *BatchRepository) GetByID(id uint) (*model.BatchOperation, error) {
	var batch model.BatchOperation
	if err := r.db.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetLatest - 获取最新的批量操作
func (r *BatchRepository) GetLatest() (*model.BatchOperation, error) {
	var batch model.BatchOperation
	if err := r.db.Order("created_at DESC").First(&batch).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetAll - 获取所有批量操作（按时间倒序）
func (r *BatchRepository) GetAll() ([]model.BatchOperation, error) {
	var batches []model.BatchOperation
	if err := r.db.Order("created_at DESC").Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// Delete - 删除批量操作
func (r *BatchRepository) Delete(id uint) error {
	return r.db.Delete(&model.BatchOperation{}, id).Error
}

// DeleteOlderThan - 删除指定时间之前的批量操作
func (r *BatchRepository) DeleteOlderThan(limit int) error {
	subquery := r.db.Model(&model.BatchOperation{}).
		Order("created_at DESC").
		Limit(limit).
		Select("id")
	return r.db.Where("id NOT IN (?)", subquery).Delete(&model.BatchOperation{}).Error
}