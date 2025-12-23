import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: (() => {
      try {
        const stored = localStorage.getItem('userInfo')
        if (stored && stored !== 'undefined') {
          return JSON.parse(stored)
        }
        return {}
      } catch (e) {
        console.error('解析用户信息失败:', e)
        localStorage.removeItem('userInfo')
        return {}
      }
    })(),
    menus: []
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    hasPermission: (state) => (permission) => {
      // 简单的权限检查逻辑，实际项目中需要根据后端返回的权限信息进行调整
      return true
    }
  },
  actions: {
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    },
    setUserInfo(userInfo) {
      this.userInfo = userInfo
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    setMenus(menus) {
      this.menus = menus
    },
    logout() {
      this.token = ''
      this.userInfo = {}
      this.menus = []
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }
})
