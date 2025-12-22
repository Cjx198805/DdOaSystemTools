<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">
        <h2>DdOaListDownload</h2>
      </div>
      <nav class="nav">
        <div class="nav-group">
          <h3>权限管理</h3>
          <ul>
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
        </div>
        <div class="nav-group">
          <h3>API测试</h3>
          <ul>
            <li>
              <router-link to="/api/test-case" active-class="active">测试用例管理</router-link>
            </li>
            <li>
              <router-link to="/api/test" active-class="active">测试执行</router-link>
            </li>
            <li>
              <router-link to="/api/test-history" active-class="active">测试历史记录</router-link>
            </li>
          </ul>
        </div>
        <div class="nav-group">
          <h3>下载任务</h3>
          <ul>
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
        </div>
      </nav>
    </aside>
    <main class="main">
      <header class="header">
        <div class="header-left">
          <button class="menu-btn" @click="toggleMenu">
            <span class="menu-icon"></span>
          </button>
          <h1>{{ currentTitle }}</h1>
        </div>
        <div class="header-right">
          <span class="user-info">欢迎, {{ username }}</span>
          <button class="logout-btn" @click="logout">退出登录</button>
        </div>
      </header>
      <div class="content">
        <slot></slot>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()
const isMenuOpen = ref(true)

const username = computed(() => userStore.username || '管理员')

const currentTitle = computed(() => {
  const path = router.currentRoute.value.path
  const titleMap = {
    '/permission/user': '用户管理',
    '/permission/role': '角色管理',
    '/permission/field': '字段权限设置',
    '/permission/dictionary': '数据字典管理',
    '/api/test-case': '测试用例管理',
    '/api/test': '测试执行',
    '/api/test-history': '测试历史记录',
    '/download/task': '任务创建和管理',
    '/download/progress': '任务进度展示',
    '/download/result': '结果下载'
  }
  return titleMap[path] || '首页'
})

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: 240px;
  background-color: #2c3e50;
  color: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
}

.sidebar.collapsed {
  width: 60px;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #34495e;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.nav {
  flex: 1;
  padding: 20px 0;
  overflow-y: auto;
}

.nav-group {
  margin-bottom: 20px;
}

.nav-group h3 {
  padding: 0 20px;
  font-size: 14px;
  color: #bdc3c7;
  margin-bottom: 10px;
}

.nav-group ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-group li {
  margin: 0;
}

.nav-group a {
  display: block;
  padding: 10px 20px;
  color: white;
  text-decoration: none;
  transition: background-color 0.2s ease;
}

.nav-group a:hover {
  background-color: #34495e;
}

.nav-group a.active {
  background-color: #3498db;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #ecf0f1;
  overflow: hidden;
}

.header {
  background-color: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
}

.menu-btn {
  background: none;
  border: none;
  cursor: pointer;
  margin-right: 10px;
  padding: 5px;
}

.menu-icon {
  display: block;
  width: 20px;
  height: 2px;
  background-color: #333;
  position: relative;
}

.menu-icon::before,
.menu-icon::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  background-color: #333;
  left: 0;
}

.menu-icon::before {
  top: -6px;
}

.menu-icon::after {
  top: 6px;
}

.header h1 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  margin-right: 20px;
  color: #666;
}

.logout-btn {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.logout-btn:hover {
  background-color: #c0392b;
}

.content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}
</style>