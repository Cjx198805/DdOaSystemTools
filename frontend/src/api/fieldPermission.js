import service from './index'

const fieldPermissionApi = {
  // 获取字段权限列表
  getFieldPermissions(params) {
    return service.get('/field-permissions', { params })
  },
  // 获取单个字段权限
  getFieldPermission(permissionId) {
    return service.get(`/field-permissions/${permissionId}`)
  },
  // 创建字段权限
  createFieldPermission(data) {
    return service.post('/field-permissions', data)
  },
  // 更新字段权限
  updateFieldPermission(permissionId, data) {
    return service.put(`/field-permissions/${permissionId}`, data)
  },
  // 删除字段权限
  deleteFieldPermission(permissionId) {
    return service.delete(`/field-permissions/${permissionId}`)
  },
  // 获取角色的字段权限
  getRoleFieldPermissions(roleId) {
    return service.get(`/roles/${roleId}/field-permissions`)
  },
  // 更新角色的字段权限
  updateRoleFieldPermissions(roleId, data) {
    return service.put(`/roles/${roleId}/field-permissions`, data)
  }
}

export default fieldPermissionApi
