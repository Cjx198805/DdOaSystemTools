import { defineStore } from 'pinia'
import userApi from '../api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    currentUser: null,
    loading: false,
    error: null
  }),
  
  getters: {
    getUserById: (state) => (id) => {
      return state.users.find(user => user.id === id)
    }
  },
  
  actions: {
    async fetchUsers(params) {
      this.loading = true
      this.error = null
      try {
        const response = await userApi.getUsers(params)
        this.users = response.data
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
    
    async fetchUser(userId) {
      this.loading = true
      this.error = null
      try {
        const response = await userApi.getUser(userId)
        const user = response.data
        // 更新users数组中的用户信息
        const index = this.users.findIndex(u => u.id === userId)
        if (index !== -1) {
          this.users[index] = user
        }
        this.currentUser = user
        return user
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
    
    async createUser(userData) {
      this.loading = true
      this.error = null
      try {
        const response = await userApi.createUser(userData)
        this.users.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
    
    async updateUser(userId, userData) {
      this.loading = true
      this.error = null
      try {
        const response = await userApi.updateUser(userId, userData)
        // 更新users数组中的用户信息
        const index = this.users.findIndex(u => u.id === userId)
        if (index !== -1) {
          this.users[index] = response.data
        }
        // 如果是当前用户，更新currentUser
        if (this.currentUser && this.currentUser.id === userId) {
          this.currentUser = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
    
    async deleteUser(userId) {
      this.loading = true
      this.error = null
      try {
        await userApi.deleteUser(userId)
        // 从users数组中移除用户
        this.users = this.users.filter(user => user.id !== userId)
        // 如果是当前用户，清空currentUser
        if (this.currentUser && this.currentUser.id === userId) {
          this.currentUser = null
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
    
    async assignRoles(userId, roleIds) {
      this.loading = true
      this.error = null
      try {
        const response = await userApi.assignRoles(userId, roleIds)
        // 更新用户的角色信息
        const index = this.users.findIndex(u => u.id === userId)
        if (index !== -1) {
          this.users[index].roles = response.data.roles
        }
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
