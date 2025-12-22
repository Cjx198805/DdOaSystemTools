import service from './index'

const roleApi = {
  // 获取角色列表
  getRoles(params) {
    return service.get('/roles', { params })
  },
  // 获取单个角色
  getRole(roleId) {
    return service.get(`/roles/${roleId}`)
  },
  // 创建角色
  createRole(data) {
    return service.post('/roles', data)
  },
  // 更新角色
  updateRole(roleId, data) {
    return service.put(`/roles/${roleId}`, data)
  },
  // 删除角色
  deleteRole(roleId) {
    return service.delete(`/roles/${roleId}`)
  },
  // 获取角色权限
  getRolePermissions(roleId) {
    return service.get(`/roles/${roleId}/permissions`)
  },
  // 更新角色权限
  updateRolePermissions(roleId, data) {
    return service.put(`/roles/${roleId}/permissions`, data)
  }
}

export default roleApi
