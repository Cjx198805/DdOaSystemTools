<template>
  <div class="user-management">
    <div class="page-header">
      <h2>用户管理</h2>
      <button class="add-btn" @click="showAddModal = true">添加用户</button>
    </div>
    
    <!-- 搜索和筛选 -->
    <div class="search-filter">
      <input 
        type="text" 
        v-model="searchQuery" 
        placeholder="搜索用户名或邮箱" 
        class="search-input"
      />
      <button class="search-btn" @click="handleSearch">搜索</button>
    </div>
    
    <!-- 用户列表 -->
    <div class="user-list">
      <table class="user-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span v-for="role in user.roles" :key="role.id" class="role-tag">{{ role.name }}</span>
            </td>
            <td>
              <span :class="['status-tag', user.status === 1 ? 'active' : 'inactive']">
                {{ user.status === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td class="action-buttons">
              <button class="edit-btn" @click="handleEdit(user)">编辑</button>
              <button class="delete-btn" @click="handleDelete(user.id)">删除</button>
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
    
    <!-- 添加/编辑用户弹窗 -->
    <div class="modal" v-if="showAddModal || showEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ showAddModal ? '添加用户' : '编辑用户' }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="username">用户名</label>
              <input 
                type="text" 
                id="username" 
                v-model="formData.username" 
                placeholder="请输入用户名" 
                required
              />
            </div>
            <div class="form-group">
              <label for="email">邮箱</label>
              <input 
                type="email" 
                id="email" 
                v-model="formData.email" 
                placeholder="请输入邮箱" 
                required
              />
            </div>
            <div class="form-group">
              <label for="password">密码</label>
              <input 
                type="password" 
                id="password" 
                v-model="formData.password" 
                placeholder="请输入密码" 
                :required="showAddModal"
              />
            </div>
            <div class="form-group">
              <label for="status">状态</label>
              <select id="status" v-model="formData.status">
                <option value="1">启用</option>
                <option value="0">禁用</option>
              </select>
            </div>
            <div class="form-group">
              <label>角色分配</label>
              <div class="role-select">
                <label v-for="role in roles" :key="role.id" class="role-option">
                  <input 
                    type="checkbox" 
                    v-model="formData.role_ids" 
                    :value="role.id"
                  />
                  {{ role.name }}
                </label>
              </div>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import userApi from '../../api/user'
import roleApi from '../../api/role'

// 数据状态
const users = ref([])
const roles = ref([])
const searchQuery = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const loading = ref(false)

// 弹窗状态
const showAddModal = ref(false)
const showEditModal = ref(false)

// 表单数据
const formData = ref({
  username: '',
  email: '',
  password: '',
  status: '1',
  role_ids: []
})

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取用户列表
const getUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.getUsers({
      page: currentPage.value,
      keyword: searchQuery.value
    })
    users.value = res.data.users
    totalPages.value = res.data.total_pages
  } catch (error) {
    console.error('获取用户列表失败:', error)
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

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getUsers()
}

// 编辑用户
const handleEdit = (user) => {
  showEditModal.value = true
  formData.value = {
    id: user.id,
    username: user.username,
    email: user.email,
    password: '',
    status: user.status.toString(),
    role_ids: user.roles.map(role => role.id)
  }
}

// 删除用户
const handleDelete = async (userId) => {
  if (confirm('确定要删除这个用户吗？')) {
    try {
      await userApi.deleteUser(userId)
      getUsers()
    } catch (error) {
      console.error('删除用户失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    if (showAddModal.value) {
      await userApi.createUser(formData.value)
    } else {
      await userApi.updateUser(formData.value.id, formData.value)
    }
    closeModal()
    getUsers()
  } catch (error) {
    console.error('保存用户失败:', error)
  }
}

// 关闭弹窗
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  formData.value = {
    username: '',
    email: '',
    password: '',
    status: '1',
    role_ids: []
  }
}

// 初始化
onMounted(() => {
  getUsers()
  getRoles()
})
</script>

<style scoped>
.user-management {
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

.user-list {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.user-table {
  width: 100%;
  border-collapse: collapse;
}

.user-table th,
.user-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.user-table th {
  background-color: #f5f7fa;
  font-weight: bold;
  color: #333;
}

.role-tag {
  display: inline-block;
  padding: 0.2rem 0.5rem;
  background-color: #ecf5ff;
  color: #409eff;
  border-radius: 10px;
  font-size: 0.8rem;
  margin-right: 0.5rem;
}

.status-tag {
  display: inline-block;
  padding: 0.2rem 0.5rem;
  border-radius: 10px;
  font-size: 0.8rem;
}

.status-tag.active {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-tag.inactive {
  background-color: #fef0f0;
  color: #f56c6c;
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

.role-select {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.role-option {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  color: #606266;
  font-size: 0.9rem;
  cursor: pointer;
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
</style>