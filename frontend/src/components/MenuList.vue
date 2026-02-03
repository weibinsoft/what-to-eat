<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Menu, Restaurant } from '../api'

interface Props {
  menus: Menu[]
  restaurants: Restaurant[]
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'add', restaurantName: string, dishName: string): void
  (e: 'delete', id: number): void
}>()

const restaurantName = ref('')
const dishName = ref('')
const showSuggestions = ref(false)

// 过滤餐厅建议
const filteredRestaurants = computed(() => {
  if (!restaurantName.value) return props.restaurants
  const query = restaurantName.value.toLowerCase()
  return props.restaurants.filter(r => r.name.toLowerCase().includes(query))
})

function handleAdd() {
  const rName = restaurantName.value.trim()
  const dName = dishName.value.trim()
  
  if (rName && dName) {
    emit('add', rName, dName)
    restaurantName.value = ''
    dishName.value = ''
  }
}

function selectRestaurant(name: string) {
  restaurantName.value = name
  showSuggestions.value = false
}

function handleDelete(id: number) {
  emit('delete', id)
}

function handleInputFocus() {
  showSuggestions.value = true
}

function handleInputBlur() {
  // 延迟关闭，以便点击建议项
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}

// 按餐厅分组菜单
const groupedMenus = computed(() => {
  const groups: { [key: string]: Menu[] } = {}
  for (const menu of props.menus) {
    const restaurantName = menu.restaurant?.name || '未知餐厅'
    if (!groups[restaurantName]) {
      groups[restaurantName] = []
    }
    groups[restaurantName].push(menu)
  }
  return groups
})

const canAdd = computed(() => {
  return restaurantName.value.trim() && dishName.value.trim()
})
</script>

<template>
  <div class="menu-list">
    <h3 class="section-title">🍜 菜单列表</h3>
    
    <div class="add-form">
      <div class="input-group">
        <div class="input-wrapper">
          <input
            v-model="restaurantName"
            type="text"
            placeholder="餐厅名称"
            class="add-input"
            @focus="handleInputFocus"
            @blur="handleInputBlur"
          />
          <ul v-if="showSuggestions && filteredRestaurants.length > 0" class="suggestions">
            <li 
              v-for="r in filteredRestaurants" 
              :key="r.id"
              @mousedown="selectRestaurant(r.name)"
              class="suggestion-item"
            >
              {{ r.name }}
            </li>
          </ul>
        </div>
        <input
          v-model="dishName"
          type="text"
          placeholder="菜品名称"
          class="add-input"
          @keyup.enter="handleAdd"
        />
      </div>
      <button class="add-btn" @click="handleAdd" :disabled="!canAdd">
        添加
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="menus.length === 0" class="empty">
      暂无菜单，请先添加餐厅和菜品
    </div>
    
    <div v-else class="menu-groups">
      <div v-for="(items, restaurant) in groupedMenus" :key="restaurant" class="menu-group">
        <div class="group-header">🏪 {{ restaurant }}</div>
        <ul class="list">
          <li v-for="menu in items" :key="menu.id" class="list-item">
            <span class="item-name">{{ menu.dish_name }}</span>
            <button class="delete-btn" @click="handleDelete(menu.id)">
              ×
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.menu-list {
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

.input-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-wrapper {
  position: relative;
}

.add-input {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.add-input:focus {
  outline: none;
  border-color: #667eea;
}

.suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  max-height: 150px;
  overflow-y: auto;
  z-index: 100;
  list-style: none;
  padding: 0;
  margin: 4px 0 0 0;
}

.suggestion-item {
  padding: 10px 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.suggestion-item:hover {
  background: #f0f0f0;
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
  align-self: flex-end;
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

.menu-groups {
  max-height: 350px;
  overflow-y: auto;
}

.menu-group {
  margin-bottom: 16px;
}

.group-header {
  font-size: 14px;
  font-weight: 600;
  color: #667eea;
  padding: 8px 0;
  border-bottom: 1px solid #e0e0e0;
  margin-bottom: 8px;
}

.list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 6px;
  transition: background 0.2s;
}

.list-item:hover {
  background: #e9ecef;
}

.item-name {
  font-size: 14px;
  color: #333;
}

.delete-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: #ff6b6b;
  color: white;
  border-radius: 50%;
  font-size: 16px;
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
