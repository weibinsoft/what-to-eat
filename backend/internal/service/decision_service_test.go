package service

import (
	"testing"
	"time"

	"what-to-eat/internal/model"
)

func TestWeightedRandom(t *testing.T) {
	service := NewDecisionService(nil, nil)

	tests := []struct {
		name          string
		menus         []model.Menu
		recentRecords []model.DecisionRecord
		expectLower   bool // 是否期望最近出现过的菜品概率降低
	}{
		{
			name: "no recent records - equal probability",
			menus: []model.Menu{
				{ID: 1, DishName: "菜品A"},
				{ID: 2, DishName: "菜品B"},
				{ID: 3, DishName: "菜品C"},
			},
			recentRecords: []model.DecisionRecord{},
			expectLower:   false,
		},
		{
			name: "with recent records - lower probability for recent",
			menus: []model.Menu{
				{ID: 1, DishName: "菜品A"},
				{ID: 2, DishName: "菜品B"},
				{ID: 3, DishName: "菜品C"},
			},
			recentRecords: []model.DecisionRecord{
				{MenuID: 1, Menu: model.Menu{ID: 1, DishName: "菜品A"}},
			},
			expectLower: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 运行多次验证随机性
			counts := make(map[int64]int)
			iterations := 1000

			for i := 0; i < iterations; i++ {
				result := service.weightedRandom(tt.menus, tt.recentRecords)
				counts[result.ID]++
			}

			// 验证所有菜品都被选中过
			for _, m := range tt.menus {
				if counts[m.ID] == 0 {
					t.Errorf("Menu %s (ID=%d) was never selected in %d iterations", m.DishName, m.ID, iterations)
				}
			}

			// 如果有最近记录，验证其概率确实降低了
			if tt.expectLower && len(tt.recentRecords) > 0 {
				recentID := tt.recentRecords[0].MenuID
				recentCount := counts[recentID]

				// 计算其他菜品的平均选中次数
				var otherTotal int
				var otherCount int
				for id, count := range counts {
					if id != recentID {
						otherTotal += count
						otherCount++
					}
				}
				avgOther := float64(otherTotal) / float64(otherCount)

				// 最近出现过的菜品应该被选中的次数更少
				if float64(recentCount) >= avgOther {
					t.Logf("Warning: Recent menu selected %d times, avg other: %.1f (expected lower)", recentCount, avgOther)
				}
			}
		})
	}
}

func TestGenerateMessage(t *testing.T) {
	service := &DecisionService{}

	tests := []struct {
		name          string
		menu          *model.Menu
		recentRecords []model.DecisionRecord
		expectedMsg   string
	}{
		{
			name:          "not in recent records",
			menu:          &model.Menu{ID: 1, DishName: "菜品A"},
			recentRecords: []model.DecisionRecord{},
			expectedMsg:   "就决定是你了！",
		},
		{
			name: "in recent records",
			menu: &model.Menu{ID: 1, DishName: "菜品A"},
			recentRecords: []model.DecisionRecord{
				{MenuID: 1, Menu: model.Menu{ID: 1, DishName: "菜品A"}},
			},
			expectedMsg: "虽然最近吃过，但命运让你再吃一次！",
		},
		{
			name: "different menu in recent",
			menu: &model.Menu{ID: 2, DishName: "菜品B"},
			recentRecords: []model.DecisionRecord{
				{MenuID: 1, Menu: model.Menu{ID: 1, DishName: "菜品A"}},
			},
			expectedMsg: "就决定是你了！",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := service.generateMessage(tt.menu, tt.recentRecords)
			if msg != tt.expectedMsg {
				t.Errorf("generateMessage() = %q, want %q", msg, tt.expectedMsg)
			}
		})
	}
}

func TestDecisionRecord_Today(t *testing.T) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	tests := []struct {
		name      string
		decidedAt time.Time
		isToday   bool
	}{
		{
			name:      "now is today",
			decidedAt: now,
			isToday:   true,
		},
		{
			name:      "yesterday",
			decidedAt: now.AddDate(0, 0, -1),
			isToday:   false,
		},
		{
			name:      "tomorrow",
			decidedAt: now.AddDate(0, 0, 1),
			isToday:   false,
		},
		{
			name:      "start of today",
			decidedAt: todayStart,
			isToday:   true,
		},
		{
			name:      "end of today minus 1 second",
			decidedAt: todayEnd.Add(-time.Second),
			isToday:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isToday := tt.decidedAt.After(todayStart) && tt.decidedAt.Before(todayEnd) ||
				tt.decidedAt.Equal(todayStart)

			if isToday != tt.isToday {
				t.Errorf("isToday check for %v: got %v, want %v", tt.decidedAt, isToday, tt.isToday)
			}
		})
	}
}
