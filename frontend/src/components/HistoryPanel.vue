<script setup lang="ts">
import type { DecisionRecord } from '../api'

interface Props {
  records: DecisionRecord[]
  loading: boolean
}

defineProps<Props>()

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hours}:${minutes}`
}

function getDaysAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '昨天'
  if (diffDays === 2) return '前天'
  return `${diffDays}天前`
}
</script>

<template>
  <div class="history-panel">
    <h3 class="section-title">📋 最近5天吃过啥</h3>
    
    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="records.length === 0" class="empty">
      还没有记录，快去摇一摇吧！
    </div>
    
    <ul v-else class="history-list">
      <li v-for="record in records" :key="record.id" class="history-item">
        <div class="item-content">
          <div class="menu-info">
            <span class="dish-name">{{ record.menu?.dish_name }}</span>
            <span class="restaurant-name">@ {{ record.menu?.restaurant?.name }}</span>
          </div>
          <span class="days-ago">{{ getDaysAgo(record.decided_at) }}</span>
        </div>
        <span class="item-time">{{ formatDate(record.decided_at) }}</span>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.history-panel {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 1.2rem;
  color: #333;
}

.loading, .empty {
  text-align: center;
  color: #888;
  padding: 20px;
}

.history-list {
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 8px;
  margin-bottom: 8px;
  border-left: 4px solid #667eea;
}

.item-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.menu-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.dish-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.restaurant-name {
  font-size: 12px;
  color: #666;
}

.days-ago {
  font-size: 12px;
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
}

.item-time {
  font-size: 12px;
  color: #888;
}
</style>
