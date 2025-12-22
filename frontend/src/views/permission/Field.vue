<template>
  <div class="field-permission">
    <div class="page-header">
      <h2>字段权限设置</h2>
    </div>
    
    <!-- 筛选条件 -->
    <div class="filter-section">
      <div class="filter-item">
        <label for="module">模块</label>
        <select id="module" v-model="filter.module">
          <option value="">全部</option>
          <option value="user">用户管理</option>
          <option value="role">角色管理</option>
          <option value="task">下载任务</option>
        </select>
      </div>
      <div class="filter-item">
        <label for="role">角色</label>
        <select id="role" v-model="filter.roleId">
          <option value="">选择角色</option>
          <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
        </select>
      </div>
      <button class="filter-btn" @click="handleFilter">筛选</button>
    </div>
    
    <!-- 字段权限列表 -->
    <div class="permission-list">
      <table class="permission-table">
        <thead>
          <tr>
            <th>模块</th>
            <th>字段名称</th>
            <th>字段标识</th>
            <th>显示权限</th>
            <th>编辑权限</th>
            <th>报表权限</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="permission in fieldPermissions" :key="permission.id">
            <td>{{ getModuleName(permission.module) }}</td>
            <td>{{ permission.field_name }}</td>
            <td>{{ permission.field_code }}</td>
            <td>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="permission.view_permission"
                  @change="handlePermissionChange(permission)"
                />
                <span class="slider"></span>
              </label>
            </td>
            <td>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="permission.edit_permission"
                  @change="handlePermissionChange(permission)"
                />
                <span class="slider"></span>
              </label>
            </td>
            <td>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="permission.report_permission"
                  @change="handlePermissionChange(permission)"
                />
                <span class="slider"></span>
              </label>
            </td>
            <td>
              <button class="edit-btn" @click="handleEdit(permission)">编辑</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- 分页 -->
    <div class="pagination">
      <button class="prev-btn" @click="currentPage > 1 && currentPage--">上一页</button>
      <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
      <button class="next-btn" @click="currentPage < totalPages && currentPage++">下一页</button>
    </div>
    
    <!-- 编辑字段权限弹窗 -->
    <div class="modal" v-if="showEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>编辑字段权限</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label>模块</label>
              <input type="text" v-model="formData.module" readonly class="readonly-input" />
            </div>
            <div class="form-group">
              <label>字段名称</label>
              <input type="text" v-model="formData.field_name" readonly class="readonly-input" />
            </div>
            <div class="form-group">
              <label>字段标识</label>
              <input type="text" v-model="formData.field_code" readonly class="readonly-input" />
            </div>
            <div class="form-group">
              <label>角色</label>
              <select v-model="formData.role_id">
                <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
              </select>
            </div>
            <div class="permission-section">
              <h4>权限设置</h4>
              <div class="permission-item">
                <label class="permission-label">
                  <input 
                    type="checkbox" 
                    v-model="formData.view_permission" 
                  />
                  显示权限
                </label>
              </div>
              <div class="permission-item">
                <label class="permission-label">
                  <input 
                    type="checkbox" 
                    v-model="formData.edit_permission" 
                  />
                  编辑权限
                </label>
              </div>
              <div class="permission-item">
                <label class="permission-label">
                  <input 
                    type="checkbox" 
                    v-model="formData.report_permission" 
                  />
                  报表权限
                </label>
              </div>
            </div>
            <div class="modal-footer">
              <button type="submit" class="submit-btn">保存</button>
              <button type="button" class="cancel-btn" @click="closeModal">取消</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import fieldPermissionApi from '../../api/fieldPermission'
import roleApi from '../../api/role'

// 数据状态
const fieldPermissions = ref([])
const roles = ref([])
const loading = ref(false)

// 筛选条件
const filter = ref({
  module: '',
  roleId: ''
})

// 分页
const currentPage = ref(1)
const totalPages = ref(1)

// 弹窗状态
const showEditModal = ref(false)

// 表单数据
const formData = ref({
  id: '',
  module: '',
  field_name: '',
  field_code: '',
  role_id: '',
  view_permission: false,
  edit_permission: false,
  report_permission: false
})

// 模块名称映射
const moduleNames = {
  'user': '用户管理',
  'role': '角色管理',
  'task': '下载任务'
}

// 获取模块名称
const getModuleName = (moduleCode) => {
  return moduleNames[moduleCode] || moduleCode
}

// 获取字段权限列表
const getFieldPermissions = async () => {
  loading.value = true
  try {
    const res = await fieldPermissionApi.getFieldPermissions({
      page: currentPage.value,
      module: filter.value.module,
      role_id: filter.value.roleId
    })
    fieldPermissions.value = res.data.permissions
    totalPages.value = res.data.total_pages
  } catch (error) {
    console.error('获取字段权限失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const getRoles = async () => {
  try {
    const res = await roleApi.getRoles()
    roles.value = res.data
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

// 筛选
const handleFilter = () => {
  currentPage.value = 1
  getFieldPermissions()
}

// 处理权限变更
const handlePermissionChange = async (permission) => {
  try {
    await fieldPermissionApi.updateFieldPermission(permission.id, {
      view_permission: permission.view_permission,
      edit_permission: permission.edit_permission,
      report_permission: permission.report_permission
    })
  } catch (error) {
    console.error('更新权限失败:', error)
    // 恢复原始状态
    getFieldPermissions()
  }
}

// 编辑字段权限
const handleEdit = (permission) => {
  showEditModal.value = true
  formData.value = {
    id: permission.id,
    module: permission.module,
    field_name: permission.field_name,
    field_code: permission.field_code,
    role_id: permission.role_id,
    view_permission: permission.view_permission,
    edit_permission: permission.edit_permission,
    report_permission: permission.report_permission
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await fieldPermissionApi.updateFieldPermission(formData.value.id, {
      role_id: formData.value.role_id,
      view_permission: formData.value.view_permission,
      edit_permission: formData.value.edit_permission,
      report_permission: formData.value.report_permission
    })
    closeModal()
    getFieldPermissions()
  } catch (error) {
    console.error('保存权限失败:', error)
  }
}

// 关闭弹窗
const closeModal = () => {
  showEditModal.value = false
  formData.value = {
    id: '',
    module: '',
    field_name: '',
    field_code: '',
    role_id: '',
    view_permission: false,
    edit_permission: false,
    report_permission: false
  }
}

// 初始化
onMounted(() => {
  getFieldPermissions()
  getRoles()
})
</script>

<style scoped>
.field-permission {
  background-color: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.page-header {
  margin-bottom: 1.5rem;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

/* 筛选区域 */
.filter-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-item label {
  font-size: 0.9rem;
  color: #606266;
}

.filter-item select {
  padding: 0.6rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  min-width: 150px;
}

.filter-btn {
  padding: 0.6rem 1.2rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  align-self: flex-end;
}

.filter-btn:hover {
  background-color: #66b1ff;
}

/* 权限列表 */
.permission-list {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.permission-table {
  width: 100%;
  border-collapse: collapse;
}

.permission-table th,
.permission-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.permission-table th {
  background-color: #f5f7fa;
  font-weight: bold;
  color: #333;
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 20px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 20px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #409eff;
}

input:focus + .slider {
  box-shadow: 0 0 1px #409eff;
}

input:checked + .slider:before {
  transform: translateX(20px);
}

/* 操作按钮 */
.edit-btn {
  padding: 0.4rem 0.8rem;
  background-color: #e6a23c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.edit-btn:hover {
  background-color: #ebb563;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
}

.prev-btn,
.next-btn {
  padding: 0.5rem 1rem;
  background-color: #f5f7fa;
  color: #606266;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.prev-btn:hover,
.next-btn:hover {
  background-color: #ecf5ff;
  color: #409eff;
  border-color: #d9ecff;
}

.page-info {
  color: #606266;
}

/* 弹窗样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #909399;
}

.close-btn:hover {
  color: #606266;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #606266;
  font-size: 0.9rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.6rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #409eff;
}

.readonly-input {
  background-color: #f5f7fa;
  cursor: not-allowed;
}

/* 权限设置区域 */
.permission-section {
  margin: 1.5rem 0;
}

.permission-section h4 {
  margin: 0 0 1rem 0;
  color: #333;
  font-size: 1rem;
}

.permission-item {
  margin-bottom: 0.8rem;
}

.permission-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.permission-label:hover {
  color: #409eff;
}

/* 弹窗底部按钮 */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.submit-btn {
  padding: 0.6rem 1.2rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-btn:hover {
  background-color: #66b1ff;
}

.cancel-btn {
  padding: 0.6rem 1.2rem;
  background-color: #f5f7fa;
  color: #606266;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.cancel-btn:hover {
  background-color: #ecf5ff;
  color: #409eff;
  border-color: #d9ecff;
}
</style>
