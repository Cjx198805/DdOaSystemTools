import service from './index'

const userApi = {
  // 获取用户列表
  getUsers(params) {
    return service.get('/api/v1/user', { params })
  },
  // 获取单个用户
  getUser(userId) {
    return service.get(`/api/v1/user/${userId}`)
  },
  // 创建用户
  createUser(data) {
    return service.post('/api/v1/user', data)
  },
  // 更新用户
  updateUser(userId, data) {
    return service.put(`/api/v1/user/${userId}`, data)
  },
  // 删除用户
  deleteUser(userId) {
    return service.delete(`/api/v1/user/${userId}`)
  },
  // 获取用户角色
  getUserRoles(userId) {
    return service.get(`/api/v1/user/${userId}/roles`)
  },
  // 分配角色
  assignRoles(userId, data) {
    return service.put(`/api/v1/user/${userId}/assign-roles`, data)
  }
}

export default userApi
