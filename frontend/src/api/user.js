import service from './index'

const userApi = {
  // 获取用户列表
  getUsers(params) {
    return service.get('/user', { params })
  },
  // 获取单个用户
  getUser(userId) {
    return service.get(`/user/${userId}`)
  },
  // 登录
  login(data) {
    return service.post('/user/login', data)
  },
  // 创建用户
  createUser(data) {
    return service.post('/user', data)
  },
  // 更新用户
  updateUser(userId, data) {
    return service.put(`/user/${userId}`, data)
  },
  // 删除用户
  deleteUser(userId) {
    return service.delete(`/user/${userId}`)
  },
  // 获取用户角色
  getUserRoles(userId) {
    return service.get(`/user/${userId}/roles`)
  },
  // 分配角色
  assignRoles(userId, data) {
    return service.put(`/api/v1/user/${userId}/assign-roles`, data)
  }
}

export default userApi
