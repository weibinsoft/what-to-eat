package repository

import (
	"time"

	"what-to-eat/internal/model"

	"gorm.io/gorm"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{db: db}
}

// Create 创建决策记录
func (r *DecisionRepository) Create(record *model.DecisionRecord) error {
	return r.db.Create(record).Error
}

// CreateOrUpdateToday 创建或更新当天的决策记录（每天只保留一条）
func (r *DecisionRepository) CreateOrUpdateToday(userID int64, menuID int64) (*model.DecisionRecord, error) {
	// 获取今天的开始和结束时间
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var existingRecord model.DecisionRecord
	err := r.db.Where("user_id = ? AND decided_at >= ? AND decided_at < ?", userID, todayStart, todayEnd).
		First(&existingRecord).Error

	if err == gorm.ErrRecordNotFound {
		// 今天没有记录，创建新记录
		newRecord := &model.DecisionRecord{
			UserID:    userID,
			MenuID:    menuID,
			DecidedAt: now,
		}
		if err := r.db.Create(newRecord).Error; err != nil {
			return nil, err
		}
		return newRecord, nil
	} else if err != nil {
		return nil, err
	}

	// 今天已有记录，更新它
	existingRecord.MenuID = menuID
	existingRecord.DecidedAt = now
	if err := r.db.Save(&existingRecord).Error; err != nil {
		return nil, err
	}

	return &existingRecord, nil
}

// GetRecentByUserID 获取用户最近N条决策记录
func (r *DecisionRepository) GetRecentByUserID(userID int64, limit int) ([]model.DecisionRecord, error) {
	var records []model.DecisionRecord
	err := r.db.Where("user_id = ?", userID).
		Order("decided_at DESC").
		Limit(limit).
		Preload("Menu").
		Preload("Menu.Restaurant").
		Find(&records).Error
	return records, err
}

// GetByUserIDAndDays 获取用户最近N天的决策记录（每天一条）
func (r *DecisionRepository) GetByUserIDAndDays(userID int64, days int) ([]model.DecisionRecord, error) {
	var records []model.DecisionRecord
	startTime := time.Now().AddDate(0, 0, -days)
	err := r.db.Where("user_id = ? AND decided_at >= ?", userID, startTime).
		Order("decided_at DESC").
		Preload("Menu").
		Preload("Menu.Restaurant").
		Find(&records).Error
	return records, err
}

// CountByUserIDAndDays 统计用户最近N天的决策记录数量
func (r *DecisionRepository) CountByUserIDAndDays(userID int64, days int) (int64, error) {
	var count int64
	startTime := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&model.DecisionRecord{}).
		Where("user_id = ? AND decided_at >= ?", userID, startTime).
		Count(&count).Error
	return count, err
}

// GetTodayRecord 获取用户今天的决策记录
func (r *DecisionRepository) GetTodayRecord(userID int64) (*model.DecisionRecord, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var record model.DecisionRecord
	err := r.db.Where("user_id = ? AND decided_at >= ? AND decided_at < ?", userID, todayStart, todayEnd).
		Preload("Menu").
		Preload("Menu.Restaurant").
		First(&record).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}
