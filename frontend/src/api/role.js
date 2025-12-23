import service from './index'

const roleApi = {
  // 获取角色列表
  getRoles(params) {
    return service.get('/api/v1/role', { params })
  },
  // 获取单个角色
  getRole(roleId) {
    return service.get(`/api/v1/role/${roleId}`)
  },
  // 创建角色
  createRole(data) {
    return service.post('/api/v1/role', data)
  },
  // 更新角色
  updateRole(roleId, data) {
    return service.put(`/api/v1/role/${roleId}`, data)
  },
  // 删除角色
  deleteRole(roleId) {
    return service.delete(`/api/v1/role/${roleId}`)
  },
  // 获取角色菜单
  getRoleMenus(roleId) {
    return service.get(`/api/v1/role/${roleId}/menus`)
  },
  // 分配角色菜单
  assignRolesMenus(roleId, data) {
    return service.put(`/api/v1/role/${roleId}/assign-menus`, data)
  }
}

export default roleApi
