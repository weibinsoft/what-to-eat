<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

async function handleSubmit() {
  error.value = ''
  
  if (!username.value || !password.value) {
    error.value = '请填写用户名和密码'
    return
  }

  if (isRegister.value && password.value !== confirmPassword.value) {
    error.value = '两次密码输入不一致'
    return
  }

  if (password.value.length < 6) {
    error.value = '密码长度至少6位'
    return
  }

  loading.value = true

  try {
    if (isRegister.value) {
      await userStore.register(username.value, password.value)
    } else {
      await userStore.login(username.value, password.value)
    }
    router.push('/home')
  } catch (err: any) {
    error.value = err.response?.data?.message || '操作失败，请重试'
  } finally {
    loading.value = false
  }
}

function toggleMode() {
  isRegister.value = !isRegister.value
  error.value = ''
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="login-title">🍽️ 今天吃什么</h1>
      <p class="login-subtitle">让选择困难症不再困难</p>

      <form @submit.prevent="handleSubmit" class="login-form">
        <div class="form-group">
          <input
            v-model="username"
            type="text"
            placeholder="用户名"
            class="form-input"
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <input
            v-model="password"
            type="password"
            placeholder="密码"
            class="form-input"
            :disabled="loading"
          />
        </div>

        <div v-if="isRegister" class="form-group">
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="确认密码"
            class="form-input"
            :disabled="loading"
          />
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '处理中...' : (isRegister ? '注册' : '登录') }}
        </button>
      </form>

      <p class="toggle-text">
        {{ isRegister ? '已有账号？' : '没有账号？' }}
        <span class="toggle-link" @click="toggleMode">
          {{ isRegister ? '去登录' : '去注册' }}
        </span>
      </p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-title {
  font-size: 2rem;
  text-align: center;
  margin-bottom: 8px;
  color: #333;
}

.login-subtitle {
  text-align: center;
  color: #666;
  margin-bottom: 30px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-input {
  padding: 14px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
}

.form-input:disabled {
  background: #f5f5f5;
}

.error-message {
  color: #e74c3c;
  font-size: 14px;
  text-align: center;
  margin: 0;
}

.submit-btn {
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(102, 126, 234, 0.4);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.toggle-text {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.toggle-link {
  color: #667eea;
  cursor: pointer;
  font-weight: bold;
}

.toggle-link:hover {
  text-decoration: underline;
}
</style>
