<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <!-- <img src="/favicon.ico" alt="Logo" class="logo" v-if="hasLogo" /> -->
        <h2>DdOaListDownload</h2>
        <p>集团数字化办公列表下载系统</p>
      </div>
      
      <el-form :model="loginForm" @submit.prevent="handleLogin" label-position="top">
        <el-form-item label="用户名">
          <el-input 
            v-model="loginForm.username" 
            placeholder="请输入用户名" 
            prefix-icon="User"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="请输入密码" 
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-button 
          type="primary" 
          class="login-btn" 
          :loading="loading" 
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form>
      
      <div class="login-footer">
        <p>&copy; 2025 cjx 项目开发</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../../stores/user'
import authApi from '../../api/auth'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const hasLogo = ref(false)

const loginForm = ref({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const res = await authApi.login(loginForm.value)
    userStore.setToken(res.data.token)
    userStore.setUserInfo(res.data.userInfo)
    userStore.setMenus(res.data.menus || [])
    
    ElMessage.success('登录成功，欢迎回来')
    // 之前是 /dashboard，但在 router/index.js 中根路径是 / 且重定向到 /permission/user
    router.push('/')
  } catch (error) {
    // 拦截器已处理大部分错误，这里可以处理特定 UI 逻辑
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  background-color: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 400px;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  margin: 10px 0;
  color: #2c3e50;
  font-size: 24px;
}

.login-header p {
  color: #7f8c8d;
  font-size: 14px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  margin-top: 20px;
  background: linear-gradient(to right, #4facfe 0%, #00f2fe 100%);
  border: none;
  transition: transform 0.2s;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
}

.login-footer {
  margin-top: 30px;
  text-align: center;
  color: #95a5a6;
  font-size: 12px;
}
</style>