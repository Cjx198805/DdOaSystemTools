<template>
  <div class="field-permission">
    <div class="page-header">
      <h2>字段权限设置</h2>
    </div>
    
    <!-- 筛选条件 -->
    <el-card class="filter-section">
      <el-form :inline="true" :model="filter" class="demo-form-inline">
        <el-form-item label="模块">
          <el-select v-model="filter.module" placeholder="全部" clearable style="width: 150px">
            <el-option label="用户管理" value="user" />
            <el-option label="角色管理" value="role" />
            <el-option label="下载任务" value="task" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="filter.roleId" placeholder="选择角色" clearable style="width: 150px">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">筛选</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 字段权限列表 -->
    <el-card class="permission-list">
      <el-table :data="fieldPermissions" v-loading="loading" style="width: 100%">
        <el-table-column prop="module" label="模块" width="120">
          <template #default="{ row }">
            {{ getModuleName(row.module) }}
          </template>
        </el-table-column>
        <el-table-column prop="field" label="字段标识" min-width="150" />
        <el-table-column label="显示权限" width="100">
          <template #default="{ row }">
            <el-switch 
              v-model="row.viewable" 
              :active-value="1" 
              :inactive-value="0"
              @change="handlePermissionChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="编辑权限" width="100">
          <template #default="{ row }">
            <el-switch 
              v-model="row.editable" 
              :active-value="1" 
              :inactive-value="0"
              @change="handlePermissionChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="报表展示" width="100">
          <template #default="{ row }">
            <el-switch 
              v-model="row.report_visible" 
              :active-value="1" 
              :inactive-value="0"
              @change="handlePermissionChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="特殊编辑" width="120">
          <template #header>
            <el-tooltip content="优先级高于数据字典的编辑限制" placement="top">
              <span class="tooltip-label">
                特殊编辑 <el-icon><InfoFilled /></el-icon>
              </span>
            </el-tooltip>
          </template>
          <template #default="{ row }">
            <el-switch 
              v-model="row.special_edit" 
              :active-value="1" 
              :inactive-value="0"
              active-color="#13ce66"
              @change="handlePermissionChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="getFieldPermissions"
        />
      </div>
    </el-card>
    
    <!-- 编辑字段权限弹窗 -->
    <el-dialog
      v-model="showEditModal"
      title="编辑字段权限"
      width="500px"
    >
      <el-form :model="formData" label-width="100px">
        <el-form-item label="模块">
          <el-input v-model="formData.module" disabled />
        </el-form-item>
        <el-form-item label="字段标识">
          <el-input v-model="formData.field" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="formData.role_id" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-divider>权限设置</el-divider>
        <el-form-item label="显示权限">
          <el-switch v-model="formData.viewable" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="编辑权限">
          <el-switch v-model="formData.editable" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="报表展示">
          <el-switch v-model="formData.report_visible" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="特殊编辑">
          <el-switch v-model="formData.special_edit" :active-value="1" :inactive-value="0" active-color="#13ce66" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditModal = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import fieldPermissionApi from '../../api/fieldPermission'
import roleApi from '../../api/role'

// 数据状态
const fieldPermissions = ref([])
const roles = ref([])
const loading = ref(false)
const total = ref(0)
const pageSize = ref(10)

// 筛选条件
const filter = ref({
  module: '',
  roleId: ''
})

// 分页
const currentPage = ref(1)

// 弹窗状态
const showEditModal = ref(false)

// 表单数据
const formData = ref({
  id: '',
  role_id: '',
  module: '',
  field: '',
  viewable: 1,
  editable: 1,
  report_visible: 1,
  special_edit: 0
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
      page_size: pageSize.value,
      module: filter.value.module,
      role_id: filter.value.roleId
    })
    // 适配后端返回的 list 和 total 格式
    fieldPermissions.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    // 拦截器已处理错误提示
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const getRoles = async () => {
  try {
    const res = await roleApi.getRoles()
    roles.value = res.data
  } catch (error) {}
}

// 筛选
const handleFilter = () => {
  currentPage.value = 1
  getFieldPermissions()
}

// 处理权限变更
const handlePermissionChange = async (row) => {
  try {
    await fieldPermissionApi.updateFieldPermission(row.id, {
      viewable: row.viewable,
      editable: row.editable,
      report_visible: row.report_visible,
      special_edit: row.special_edit,
      role_id: row.role_id,
      module: row.module,
      field: row.field
    })
    ElMessage.success('更新成功')
  } catch (error) {
    getFieldPermissions() // 恢复原始状态
  }
}

// 编辑弹窗
const handleEdit = (row) => {
  formData.value = { ...row }
  showEditModal.value = true
}

// 提交表单
const handleSubmit = async () => {
  try {
    await fieldPermissionApi.updateFieldPermission(formData.value.id, formData.value)
    ElMessage.success('保存成功')
    showEditModal.value = false
    getFieldPermissions()
  } catch (error) {}
}

// 初始化
onMounted(() => {
  getFieldPermissions()
  getRoles()
})
</script>

<style scoped>
.field-permission {
  padding: 10px;
}
.page-header {
  margin-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 10px;
}
.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}
.filter-section {
  margin-bottom: 20px;
}
.permission-list {
  min-height: 400px;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
.tooltip-label {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
