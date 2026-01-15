<template>
  <div class="role-management">
    <div class="page-header">
      <h2>角色管理</h2>
      <button class="add-btn" @click="showAddModal = true">添加角色</button>
    </div>
    
    <!-- 搜索和筛选 -->
    <div class="search-filter">
      <input 
        type="text" 
        v-model="searchQuery" 
        placeholder="搜索角色名称" 
        class="search-input"
      />
      <button class="search-btn" @click="handleSearch">搜索</button>
    </div>
    
    <!-- 角色列表 -->
    <div class="role-list">
      <table class="role-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>角色名称</th>
            <th>角色标识</th>
            <th>描述</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="role in roles" :key="role.id">
            <td>{{ role.id }}</td>
            <td>{{ role.name }}</td>
            <td>{{ role.code }}</td>
            <td>{{ role.description || '-' }}</td>
            <td>{{ formatDate(role.created_at) }}</td>
            <td class="action-buttons">
              <button class="edit-btn" @click="handleEdit(role)">编辑</button>
              <button class="permission-btn" @click="handlePermission(role)">权限设置</button>
              <button class="delete-btn" @click="handleDelete(role.id)">删除</button>
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
    
    <!-- 添加/编辑角色弹窗 -->
    <div class="modal" v-if="showAddModal || showEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ showAddModal ? '添加角色' : '编辑角色' }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="name">角色名称</label>
              <input 
                type="text" 
                id="name" 
                v-model="formData.name" 
                placeholder="请输入角色名称" 
                required
              />
            </div>
            <div class="form-group">
              <label for="code">角色标识</label>
              <input 
                type="text" 
                id="code" 
                v-model="formData.code" 
                placeholder="请输入角色标识" 
                required
              />
            </div>
            <div class="form-group">
              <label for="description">描述</label>
              <textarea 
                id="description" 
                v-model="formData.description" 
                placeholder="请输入角色描述" 
                rows="3"
              ></textarea>
            </div>
            <div class="modal-footer">
              <button type="submit" class="submit-btn">
                {{ showAddModal ? '添加' : '保存' }}
              </button>
              <button type="button" class="cancel-btn" @click="closeModal">取消</button>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- 权限设置弹窗 -->
    <div class="modal permission-modal" v-if="showPermissionModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>权限设置 - {{ selectedRole.name }}</h3>
          <button class="close-btn" @click="closePermissionModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="permission-tree">
            <div v-for="menu in menus" :key="menu.id" class="menu-item">
              <div class="menu-header">
                <label class="checkbox-label">
                  <input 
                    type="checkbox" 
                    v-model="menu.checked" 
                    @change="handleMenuCheck(menu)"
                  />
                  {{ menu.name }}
                </label>
              </div>
              <div v-if="menu.children && menu.children.length > 0" class="sub-menu">
                <div v-for="subMenu in menu.children" :key="subMenu.id" class="sub-menu-item">
                  <label class="checkbox-label">
                    <input 
                      type="checkbox" 
                      v-model="subMenu.checked" 
                      @change="handleSubMenuCheck(subMenu, menu)"
                    />
                    {{ subMenu.name }}
                  </label>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="submit-btn" @click="handlePermissionSubmit">保存权限</button>
            <button type="button" class="cancel-btn" @click="closePermissionModal">取消</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import roleApi from '../../api/role'
import menuApi from '../../api/menu'

// 数据状态
const roles = ref([])
const menus = ref([])
const searchQuery = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const loading = ref(false)

// 弹窗状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showPermissionModal = ref(false)

// 表单数据
const formData = ref({
  name: '',
  code: '',
  description: ''
})

// 选中的角色
const selectedRole = ref({})

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取角色列表
const getRoles = async () => {
  loading.value = true
  try {
    const res = await roleApi.getRoles({
      page: currentPage.value,
      keyword: searchQuery.value
    })
    roles.value = res.data.roles
    totalPages.value = res.data.total_pages
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取菜单列表
const getMenus = async () => {
  try {
    const res = await menuApi.getMenus()
    menus.value = res.data
    // 初始化菜单的checked状态
    menus.value.forEach(menu => {
      menu.checked = false
      if (menu.children) {
        menu.children.forEach(subMenu => {
          subMenu.checked = false
        })
      }
    })
  } catch (error) {
    console.error('获取菜单列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getRoles()
}

// 编辑角色
const handleEdit = (role) => {
  showEditModal.value = true
  formData.value = {
    id: role.id,
    name: role.name,
    code: role.code,
    description: role.description
  }
}

// 删除角色
const handleDelete = async (roleId) => {
  if (confirm('确定要删除这个角色吗？')) {
    try {
      await roleApi.deleteRole(roleId)
      getRoles()
    } catch (error) {
      console.error('删除角色失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    if (showAddModal.value) {
      await roleApi.createRole(formData.value)
    } else {
      await roleApi.updateRole(formData.value.id, formData.value)
    }
    closeModal()
    getRoles()
  } catch (error) {
    console.error('保存角色失败:', error)
  }
}

// 关闭弹窗
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  formData.value = {
    name: '',
    code: '',
    description: ''
  }
}

// 打开权限设置弹窗
const handlePermission = async (role) => {
  selectedRole.value = role
  showPermissionModal.value = true
  
  // 获取角色的现有权限
  try {
    const res = await roleApi.getRolePermissions(role.id)
    const rolePermissions = res.data.permissions
    
    // 更新菜单的checked状态
    menus.value.forEach(menu => {
      const menuHasPermission = rolePermissions.some(p => p.menu_id === menu.id)
      menu.checked = menuHasPermission
      
      if (menu.children) {
        menu.children.forEach(subMenu => {
          const subMenuHasPermission = rolePermissions.some(p => p.menu_id === subMenu.id)
          subMenu.checked = subMenuHasPermission
        })
      }
    })
  } catch (error) {
    console.error('获取角色权限失败:', error)
  }
}

// 关闭权限设置弹窗
const closePermissionModal = () => {
  showPermissionModal.value = false
  selectedRole.value = {}
}

// 处理菜单勾选
const handleMenuCheck = (menu) => {
  if (menu.children) {
    menu.children.forEach(subMenu => {
      subMenu.checked = menu.checked
    })
  }
}

// 处理子菜单勾选
const handleSubMenuCheck = (subMenu, menu) => {
  // 检查所有子菜单是否都被勾选
  const allChecked = menu.children.every(child => child.checked)
  menu.checked = allChecked
}

// 提交权限设置
const handlePermissionSubmit = async () => {
  // 收集选中的权限
  const permissions = []
  
  menus.value.forEach(menu => {
    // 检查父菜单是否被勾选
    if (menu.checked) {
      permissions.push(menu.id)
    }
    
    // 检查子菜单是否被勾选
    if (menu.children) {
      menu.children.forEach(subMenu => {
        if (subMenu.checked) {
          permissions.push(subMenu.id)
        }
      })
    }
  })
  
  try {
    await roleApi.updateRolePermissions(selectedRole.value.id, { menu_ids: permissions })
    closePermissionModal()
  } catch (error) {
    console.error('保存权限失败:', error)
  }
}

// 初始化
onMounted(() => {
  getRoles()
  getMenus()
})
</script>

<style scoped>
.role-management {
  background-color: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

.add-btn {
  padding: 0.6rem 1.2rem;
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.add-btn:hover {
  background-color: #85ce61;
}

.search-filter {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.search-input {
  flex: 1;
  padding: 0.6rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.search-btn {
  padding: 0.6rem 1.2rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.search-btn:hover {
  background-color: #66b1ff;
}

.role-list {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.role-table {
  width: 100%;
  border-collapse: collapse;
}

.role-table th,
.role-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.role-table th {
  background-color: #f5f7fa;
  font-weight: bold;
  color: #333;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

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

.permission-btn {
  padding: 0.4rem 0.8rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.permission-btn:hover {
  background-color: #66b1ff;
}

.delete-btn {
  padding: 0.4rem 0.8rem;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.delete-btn:hover {
  background-color: #f78989;
}

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

.permission-modal .modal-content {
  max-width: 700px;
  max-height: 80vh;
  overflow-y: auto;
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
.form-group textarea {
  width: 100%;
  padding: 0.6rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #409eff;
}

.form-group textarea {
  resize: vertical;
}

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

/* 权限树样式 */
.permission-tree {
  margin: 1rem 0;
}

.menu-item {
  margin-bottom: 1rem;
}

.menu-header {
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.sub-menu {
  margin-left: 1.5rem;
}

.sub-menu-item {
  margin-bottom: 0.5rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.checkbox-label:hover {
  color: #409eff;
}
</style>
