import service from './index'

const roleApi = {
  // 获取角色列表
  getRoles(params) {
    return service.get('/role', { params })
  },
  // 获取单个角色
  getRole(roleId) {
    return service.get(`/role/${roleId}`)
  },
  // 创建角色
  createRole(data) {
    return service.post('/role', data)
  },
  // 更新角色
  updateRole(roleId, data) {
    return service.put(`/role/${roleId}`, data)
  },
  // 删除角色
  deleteRole(roleId) {
    return service.delete(`/role/${roleId}`)
  },
  // 获取角色菜单
  getRoleMenus(roleId) {
    return service.get(`/role/${roleId}/menus`)
  },
  // 分配角色菜单
  assignRolesMenus(roleId, data) {
    return service.put(`/role/${roleId}/assign-menus`, data)
  }
}

export default roleApi
