<script setup lang="ts">
import { ref } from 'vue'
import type { Restaurant } from '../api'

interface Props {
  restaurants: Restaurant[]
  loading: boolean
}

defineProps<Props>()
const emit = defineEmits<{
  (e: 'add', name: string): void
  (e: 'delete', id: number): void
}>()

const newRestaurantName = ref('')

function handleAdd() {
  const name = newRestaurantName.value.trim()
  if (name) {
    emit('add', name)
    newRestaurantName.value = ''
  }
}

function handleDelete(id: number) {
  emit('delete', id)
}
</script>

<template>
  <div class="restaurant-list">
    <h3 class="section-title">🍜 餐厅列表</h3>
    
    <div class="add-form">
      <input
        v-model="newRestaurantName"
        type="text"
        placeholder="输入餐厅名称"
        class="add-input"
        @keyup.enter="handleAdd"
      />
      <button class="add-btn" @click="handleAdd" :disabled="!newRestaurantName.trim()">
        添加
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="restaurants.length === 0" class="empty">
      暂无餐厅，请先添加
    </div>
    
    <ul v-else class="list">
      <li v-for="restaurant in restaurants" :key="restaurant.id" class="list-item">
        <span class="item-name">{{ restaurant.name }}</span>
        <button class="delete-btn" @click="handleDelete(restaurant.id)">
          ×
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.restaurant-list {
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

.add-form {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.add-input {
  flex: 1;
  padding: 10px 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.add-input:focus {
  outline: none;
  border-color: #667eea;
}

.add-btn {
  padding: 10px 20px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.add-btn:hover:not(:disabled) {
  background: #5a6fd6;
}

.add-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.loading, .empty {
  text-align: center;
  color: #888;
  padding: 20px;
}

.list {
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 300px;
  overflow-y: auto;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 8px;
  transition: background 0.2s;
}

.list-item:hover {
  background: #e9ecef;
}

.item-name {
  font-size: 15px;
  color: #333;
}

.delete-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: #ff6b6b;
  color: white;
  border-radius: 50%;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.3s, transform 0.2s;
}

.delete-btn:hover {
  background: #ee5a5a;
  transform: scale(1.1);
}
</style>
