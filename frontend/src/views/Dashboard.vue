<template>
  <div class="dashboard-container">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h3>DdOaListDownload</h3>
      </div>
      <nav class="sidebar-nav">
        <ul>
          <!-- 权限管理 -->
          <li class="nav-item">
            <div class="nav-title">权限管理</div>
            <ul class="nav-children">
              <li>
                <router-link to="/permission/user" active-class="active">用户管理</router-link>
              </li>
              <li>
                <router-link to="/permission/role" active-class="active">角色管理</router-link>
              </li>
              <li>
                <router-link to="/permission/field" active-class="active">字段权限设置</router-link>
              </li>
              <li>
                <router-link to="/permission/dictionary" active-class="active">数据字典管理</router-link>
              </li>
            </ul>
          </li>
          <!-- API测试 -->
          <li class="nav-item">
            <div class="nav-title">API测试</div>
            <ul class="nav-children">
              <li>
                <router-link to="/api/test" active-class="active">API测试</router-link>
              </li>
              <li>
                <router-link to="/api/test-case" active-class="active">测试用例管理</router-link>
              </li>
              <li>
                <router-link to="/api/test-history" active-class="active">测试历史记录</router-link>
              </li>
            </ul>
          </li>
          <!-- 下载任务 -->
          <li class="nav-item">
            <div class="nav-title">下载任务</div>
            <ul class="nav-children">
              <li>
                <router-link to="/download/task" active-class="active">任务创建和管理</router-link>
              </li>
              <li>
                <router-link to="/download/progress" active-class="active">任务进度展示</router-link>
              </li>
              <li>
                <router-link to="/download/result" active-class="active">结果下载</router-link>
              </li>
            </ul>
          </li>
        </ul>
      </nav>
    </aside>
    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 顶部导航栏 -->
      <header class="top-header">
        <div class="header-left">
          <router-link to="/dashboard" class="logo">
            DdOaListDownload
          </router-link>
        </div>
        <div class="header-right">
          <div class="user-info">
            <span class="username">{{ userStore.userInfo.username || '未知用户' }}</span>
            <button class="logout-btn" @click="handleLogout">退出登录</button>
          </div>
        </div>
      </header>
      <!-- 内容区域 -->
      <div class="content-wrapper">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { useUserStore } from '../stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dashboard-container {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

/* 侧边栏样式 */
.sidebar {
  width: 250px;
  background-color: #304156;
  color: white;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.sidebar-header {
  padding: 1rem;
  background-color: #263445;
  text-align: center;
  border-bottom: 1px solid #1f2d3d;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
}

.sidebar-nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  margin-bottom: 0.5rem;
}

.nav-title {
  padding: 0.8rem 1rem;
  font-size: 0.9rem;
  color: #a0a6ad;
  font-weight: bold;
  text-transform: uppercase;
}

.nav-children {
  margin-left: 0;
}

.nav-children li {
  padding: 0;
}

.nav-children a {
  display: block;
  padding: 0.8rem 1rem 0.8rem 2rem;
  color: white;
  text-decoration: none;
  transition: all 0.3s;
  border-left: 3px solid transparent;
}

.nav-children a:hover {
  background-color: #409eff;
  border-left-color: #409eff;
}

.nav-children a.active {
  background-color: #409eff;
  border-left-color: white;
}

/* 主内容区样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
  overflow: hidden;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2rem;
  height: 60px;
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.header-left .logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #409eff;
  text-decoration: none;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.username {
  font-size: 1rem;
  color: #333;
}

.logout-btn {
  padding: 0.5rem 1rem;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background-color: #f78989;
}

.content-wrapper {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}
</style>