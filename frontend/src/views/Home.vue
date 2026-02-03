<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { menuApi, restaurantApi, decisionApi, type Menu, type Restaurant, type DecisionRecord } from '../api'
import SlotMachine from '../components/SlotMachine.vue'
import MenuList from '../components/MenuList.vue'
import HistoryPanel from '../components/HistoryPanel.vue'

const router = useRouter()
const userStore = useUserStore()

// 状态
const menus = ref<Menu[]>([])
const restaurants = ref<Restaurant[]>([])
const historyRecords = ref<DecisionRecord[]>([])
const loadingMenus = ref(false)
const loadingHistory = ref(false)
const isSpinning = ref(false)
const decisionResult = ref('')
const resultMessage = ref('')
const showResultMessage = ref(false)

// 计算属性：菜品显示名称（餐厅 - 菜品）
const menuDisplayNames = computed(() => 
  menus.value.map(m => `${m.restaurant?.name} - ${m.dish_name}`)
)
const canDecide = computed(() => menus.value.length > 0 && !isSpinning.value)

// 加载数据
async function loadMenus() {
  loadingMenus.value = true
  try {
    const response = await menuApi.list()
    menus.value = response.data.data || []
  } catch (err) {
    console.error('Failed to load menus:', err)
  } finally {
    loadingMenus.value = false
  }
}

async function loadRestaurants() {
  try {
    const response = await restaurantApi.list()
    restaurants.value = response.data.data || []
  } catch (err) {
    console.error('Failed to load restaurants:', err)
  }
}

async function loadHistory() {
  loadingHistory.value = true
  try {
    const response = await decisionApi.history()
    historyRecords.value = response.data.data?.records || []
  } catch (err) {
    console.error('Failed to load history:', err)
  } finally {
    loadingHistory.value = false
  }
}

// 添加菜单
async function handleAddMenu(restaurantName: string, dishName: string) {
  try {
    const response = await menuApi.create(restaurantName, dishName)
    const newMenu = response.data.data.menu
    menus.value.push(newMenu)
    
    // 如果创建了新餐厅，刷新餐厅列表
    if (response.data.data.is_new_restaurant) {
      loadRestaurants()
    }
  } catch (err: any) {
    console.error('Failed to add menu:', err)
    const message = err.response?.data?.message || '添加失败，请重试'
    alert(message)
  }
}

// 删除菜单
async function handleDeleteMenu(id: number) {
  try {
    await menuApi.delete(id)
    menus.value = menus.value.filter(m => m.id !== id)
  } catch (err) {
    console.error('Failed to delete menu:', err)
    alert('删除失败，请重试')
  }
}

// 执行决策
async function handleDecide() {
  if (!canDecide.value) return
  
  showResultMessage.value = false
  
  try {
    const response = await decisionApi.decide()
    const result = response.data.data
    // 显示格式：餐厅 - 菜品
    decisionResult.value = `${result.menu.restaurant?.name} - ${result.menu.dish_name}`
    resultMessage.value = result.message
    isSpinning.value = true
  } catch (err) {
    console.error('Failed to decide:', err)
    alert('决策失败，请重试')
  }
}

// 滚动结束
function handleSpinEnd() {
  isSpinning.value = false
  showResultMessage.value = true
  loadHistory() // 刷新历史记录
}

// 登出
function handleLogout() {
  userStore.logout()
  router.push('/login')
}

// 初始化
onMounted(() => {
  loadMenus()
  loadRestaurants()
  loadHistory()
})
</script>

<template>
  <div class="home-container">
    <!-- 头部 -->
    <header class="header">
      <h1 class="app-title">🍽️ 今天吃什么</h1>
      <div class="user-info">
        <span class="username">{{ userStore.user?.username }}</span>
        <button class="logout-btn" @click="handleLogout">退出</button>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="main-content">
      <!-- 决策区域 -->
      <section class="decision-section">
        <SlotMachine
          :items="menuDisplayNames"
          :result="decisionResult"
          :isSpinning="isSpinning"
          @spinEnd="handleSpinEnd"
        />
        
        <p v-if="showResultMessage" class="result-message">
          {{ resultMessage }}
        </p>

        <button
          class="decide-btn"
          :disabled="!canDecide"
          @click="handleDecide"
        >
          {{ isSpinning ? '决策中...' : '🎰 开始决策' }}
        </button>
      </section>

      <!-- 侧边栏 -->
      <aside class="sidebar">
        <MenuList
          :menus="menus"
          :restaurants="restaurants"
          :loading="loadingMenus"
          @add="handleAddMenu"
          @delete="handleDeleteMenu"
        />

        <HistoryPanel
          :records="historyRecords"
          :loading="loadingHistory"
        />
      </aside>
    </main>
  </div>
</template>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto 30px;
  padding: 0 20px;
}

.app-title {
  color: white;
  font-size: 1.8rem;
  margin: 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.logout-btn {
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 420px;
  gap: 30px;
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 900px) {
  .main-content {
    grid-template-columns: 1fr;
  }
}

.decision-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.result-message {
  color: #00ff88;
  font-size: 1.2rem;
  margin: 20px 0;
  text-align: center;
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.decide-btn {
  margin-top: 30px;
  padding: 20px 60px;
  font-size: 1.5rem;
  font-weight: bold;
  color: white;
  background: linear-gradient(135deg, #e94560 0%, #ff6b6b 100%);
  border: none;
  border-radius: 50px;
  cursor: pointer;
  box-shadow: 0 10px 30px rgba(233, 69, 96, 0.4);
  transition: transform 0.3s, box-shadow 0.3s;
}

.decide-btn:hover:not(:disabled) {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 15px 40px rgba(233, 69, 96, 0.5);
}

.decide-btn:active:not(:disabled) {
  transform: translateY(0) scale(0.98);
}

.decide-btn:disabled {
  background: linear-gradient(135deg, #666 0%, #888 100%);
  box-shadow: none;
  cursor: not-allowed;
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>
