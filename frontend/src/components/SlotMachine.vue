<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { soundManager } from '../utils/sound'

interface Props {
  items: string[]
  result: string
  isSpinning: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'spinEnd'): void
}>()

const displayText = ref('')
const showResult = ref(false)
const isShaking = ref(false)

// 监听 isSpinning 变化
watch(() => props.isSpinning, (spinning) => {
  if (spinning && props.items.length > 0) {
    startSpinning()
  }
})

function startSpinning() {
  showResult.value = false
  isShaking.value = false
  
  // 播放开始音效
  soundManager.playStart()
  
  const totalDuration = 3000 // 总时长3秒
  const startTime = Date.now()
  
  function spin() {
    const elapsed = Date.now() - startTime
    const progress = elapsed / totalDuration
    
    if (progress >= 1) {
      // 停止，显示最终结果
      displayText.value = props.result
      showResult.value = true
      
      // 播放结果揭晓音效
      soundManager.playResult()
      
      // 触发震动效果
      setTimeout(() => {
        isShaking.value = true
        // 震动结束后
        setTimeout(() => {
          emit('spinEnd')
        }, 900)
      }, 100)
      return
    }
    
    // 播放滴答音效
    soundManager.playTick()
    
    // 随机显示一个餐厅名
    const randomIndex = Math.floor(Math.random() * props.items.length)
    displayText.value = props.items[randomIndex] || ''
    
    // 计算下一次间隔：开始快，结束慢
    // 从50ms逐渐增加到300ms
    const minInterval = 50
    const maxInterval = 300
    const interval = minInterval + (maxInterval - minInterval) * Math.pow(progress, 2)
    
    setTimeout(spin, interval)
  }
  
  spin()
}

const displayClass = computed(() => ({
  'result-display': true,
  'show-result': showResult.value,
  'shake': isShaking.value,
}))
</script>

<template>
  <div class="slot-machine">
    <div class="slot-window">
      <div :class="displayClass">
        {{ displayText || '？？？' }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.slot-machine {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.slot-window {
  background: linear-gradient(145deg, #1a1a2e, #16213e);
  border-radius: 20px;
  padding: 30px 50px;
  box-shadow: 
    0 10px 30px rgba(0, 0, 0, 0.3),
    inset 0 2px 10px rgba(255, 255, 255, 0.1);
  border: 3px solid #e94560;
  min-width: 300px;
  text-align: center;
}

.result-display {
  font-size: 2.5rem;
  font-weight: bold;
  color: #ffd700;
  text-shadow: 0 0 20px rgba(255, 215, 0, 0.5);
  transition: transform 0.1s;
}

.result-display.show-result {
  color: #00ff88;
  text-shadow: 0 0 30px rgba(0, 255, 136, 0.8);
}

.result-display.shake {
  animation: shake 0.3s ease-in-out 3, scaleUp 0.3s ease-out forwards;
}

@keyframes shake {
  0%, 100% { 
    transform: translateX(0) scale(1.2); 
  }
  10%, 30%, 50%, 70%, 90% { 
    transform: translateX(-8px) scale(1.2); 
  }
  20%, 40%, 60%, 80% { 
    transform: translateX(8px) scale(1.2); 
  }
}

@keyframes scaleUp {
  from {
    transform: scale(1);
  }
  to {
    transform: scale(1.2);
  }
}
</style>
