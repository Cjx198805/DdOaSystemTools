import service from './index'

const fieldPermissionApi = {
  // 获取字段权限列表
  getFieldPermissions(params) {
    return service.get('/field-permission', { params })
  },
  // 获取单个字段权限
  getFieldPermission(permissionId) {
    return service.get(`/field-permission/${permissionId}`)
  },
  // 创建字段权限
  createFieldPermission(data) {
    return service.post('/field-permission', data)
  },
  // 更新字段权限
  updateFieldPermission(permissionId, data) {
    return service.put(`/field-permission/${permissionId}`, data)
  },
  // 删除字段权限
  deleteFieldPermission(permissionId) {
    return service.delete(`/field-permission/${permissionId}`)
  },
  // 获取角色的字段权限
  getRoleFieldPermissions(roleId, module) {
    return service.get(`/field-permission/role/${roleId}/module/${module}`)
  },
  // 更新角色的字段权限
  updateRoleFieldPermissions(roleId, data) {
    return service.put(`/roles/${roleId}/field-permissions`, data)
  }
}

export default fieldPermissionApi
