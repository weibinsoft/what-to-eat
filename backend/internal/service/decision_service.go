package service

import (
	"errors"
	"math/rand"
	"time"

	"what-to-eat/internal/model"
	"what-to-eat/internal/repository"
)

var (
	ErrNoMenus = errors.New("没有可选择的菜品")
)

type DecisionService struct {
	decisionRepo *repository.DecisionRepository
	menuRepo     *repository.MenuRepository
	rng          *rand.Rand
}

func NewDecisionService(decisionRepo *repository.DecisionRepository, menuRepo *repository.MenuRepository) *DecisionService {
	return &DecisionService{
		decisionRepo: decisionRepo,
		menuRepo:     menuRepo,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Decide 执行决策（加权随机算法）
func (s *DecisionService) Decide(userID int64, menuIDs []int64) (*model.DecideResponse, error) {
	// 获取候选菜单列表
	var menus []model.Menu
	var err error

	if len(menuIDs) > 0 {
		menus, err = s.menuRepo.GetByIDs(menuIDs)
	} else {
		menus, err = s.menuRepo.GetAll()
	}

	if err != nil {
		return nil, err
	}

	if len(menus) == 0 {
		return nil, ErrNoMenus
	}

	// 获取用户最近3次决策记录
	recentRecords, err := s.decisionRepo.GetRecentByUserID(userID, 3)
	if err != nil {
		return nil, err
	}

	// 执行加权随机选择
	selected := s.weightedRandom(menus, recentRecords)

	// 保存决策记录
	record := &model.DecisionRecord{
		UserID:    userID,
		MenuID:    selected.ID,
		DecidedAt: time.Now(),
	}
	if err := s.decisionRepo.Create(record); err != nil {
		return nil, err
	}

	// 生成响应消息
	message := s.generateMessage(selected, recentRecords)

	return &model.DecideResponse{
		Menu:    *selected,
		Message: message,
	}, nil
}

// weightedRandom 加权随机算法
// 如果菜品在最近3次记录中出现过，其被选中的概率降低50%
func (s *DecisionService) weightedRandom(menus []model.Menu, recentRecords []model.DecisionRecord) *model.Menu {
	// 统计最近3次中各菜品出现次数
	recentCount := make(map[int64]int)
	for _, record := range recentRecords {
		recentCount[record.MenuID]++
	}

	// 计算权重
	weights := make([]float64, len(menus))
	totalWeight := 0.0

	for i, m := range menus {
		if recentCount[m.ID] > 0 {
			// 在最近3次中出现过，权重降为0.5
			weights[i] = 0.5
		} else {
			weights[i] = 1.0
		}
		totalWeight += weights[i]
	}

	// 加权随机选择
	randomValue := s.rng.Float64() * totalWeight
	cumulativeWeight := 0.0

	for i, w := range weights {
		cumulativeWeight += w
		if randomValue <= cumulativeWeight {
			return &menus[i]
		}
	}

	// 兜底：返回最后一个
	return &menus[len(menus)-1]
}

// generateMessage 生成响应消息
func (s *DecisionService) generateMessage(menu *model.Menu, recentRecords []model.DecisionRecord) string {
	// 检查是否在最近记录中
	for _, record := range recentRecords {
		if record.MenuID == menu.ID {
			return "虽然最近吃过，但命运让你再吃一次！"
		}
	}
	return "就决定是你了！"
}

// GetHistory 获取用户最近5天的决策历史
func (s *DecisionService) GetHistory(userID int64) (*model.HistoryResponse, error) {
	records, err := s.decisionRepo.GetByUserIDAndDays(userID, 5)
	if err != nil {
		return nil, err
	}

	return &model.HistoryResponse{
		Records: records,
		Total:   int64(len(records)),
	}, nil
}

// GetRecentRecords 获取用户最近N条决策记录
func (s *DecisionService) GetRecentRecords(userID int64, limit int) ([]model.DecisionRecord, error) {
	return s.decisionRepo.GetRecentByUserID(userID, limit)
}
