<template>
  <div class="download-tasks">
    <div class="page-header">
      <h2>下载任务管理</h2>
      <button class="create-btn" @click="createTask">创建新任务</button>
    </div>
    
    <!-- 任务筛选 -->
    <div class="task-filters">
      <div class="filter-row">
        <div class="form-group">
          <label for="task-status">任务状态</label>
          <select id="task-status" v-model="filters.status">
            <option value="">全部</option>
            <option value="pending">待执行</option>
            <option value="running">执行中</option>
            <option value="completed">已完成</option>
            <option value="failed">失败</option>
          </select>
        </div>
        <div class="form-group">
          <label for="task-type">任务类型</label>
          <select id="task-type" v-model="filters.type">
            <option value="">全部</option>
            <option value="excel">Excel</option>
            <option value="csv">CSV</option>
            <option value="json">JSON</option>
          </select>
        </div>
        <div class="form-group">
          <label for="task-search">搜索</label>
          <input 
            type="text" 
            id="task-search" 
            v-model="filters.search" 
            placeholder="任务名称或描述"
            @input="handleSearch"
          />
        </div>
        <button class="reset-btn" @click="resetFilters">重置</button>
      </div>
    </div>
    
    <!-- 任务列表 -->
    <div class="task-list">
      <table class="tasks-table">
        <thead>
          <tr>
            <th>任务名称</th>
            <th>任务类型</th>
            <th>创建时间</th>
            <th>执行状态</th>
            <th>进度</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in tasks" :key="task.id">
            <td class="task-name">{{ task.name }}</td>
            <td class="task-type">{{ formatType(task.type) }}</td>
            <td class="task-time">{{ formatDate(task.createdAt) }}</td>
            <td class="task-status">
              <span class="status-badge" :class="task.status">
                {{ formatStatus(task.status) }}
              </span>
            </td>
            <td class="task-progress">
              <div class="progress-container">
                <div 
                  class="progress-bar" 
                  :style="{ width: `${task.progress}%` }"
                  :class="getProgressClass(task.progress)"
                ></div>
                <span class="progress-text">{{ task.progress }}%</span>
              </div>
            </td>
            <td class="task-actions">
              <button class="detail-btn" @click="viewTaskDetail(task.id)">详情</button>
              <button 
                class="download-btn" 
                @click="downloadResult(task.id)"
                :disabled="task.status !== 'completed'"
              >
                下载
              </button>
              <button 
                class="delete-btn" 
                @click="deleteTask(task.id)"
                :disabled="task.status === 'running'"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- 空状态 -->
      <div class="empty-state" v-if="tasks.length === 0">
        <p>暂无下载任务</p>
      </div>
    </div>
    
    <!-- 分页 -->
    <div class="pagination" v-if="total > pageSize">
      <button 
        class="page-btn" 
        @click="changePage(currentPage - 1)"
        :disabled="currentPage === 1"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页
      </span>
      <button 
        class="page-btn" 
        @click="changePage(currentPage + 1)"
        :disabled="currentPage >= Math.ceil(total / pageSize)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import downloadApi from '../../api/download'

const router = useRouter()

// 任务列表
const tasks = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 筛选条件
const filters = ref({
  status: '',
  type: '',
  search: ''
})

// 获取任务列表
const getTasks = async () => {
  try {
    const response = await downloadApi.getTasks({
      page: currentPage.value,
      pageSize: pageSize.value,
      ...filters.value
    })
    tasks.value = response.data.tasks
    total.value = response.data.total
  } catch (error) {
    console.error('获取任务列表失败:', error)
  }
}

// 格式化任务类型
const formatType = (type) => {
  const typeMap = {
    excel: 'Excel',
    csv: 'CSV',
    json: 'JSON'
  }
  return typeMap[type] || type
}

// 格式化任务状态
const formatStatus = (status) => {
  const statusMap = {
    pending: '待执行',
    running: '执行中',
    completed: '已完成',
    failed: '失败'
  }
  return statusMap[status] || status
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取进度条样式类
const getProgressClass = (progress) => {
  if (progress === 100) return 'completed'
  if (progress > 0) return 'running'
  return ''
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getTasks()
}

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    status: '',
    type: '',
    search: ''
  }
  currentPage.value = 1
  getTasks()
}

// 切换页码
const changePage = (page) => {
  if (page < 1 || page > Math.ceil(total.value / pageSize.value)) {
    return
  }
  currentPage.value = page
  getTasks()
}

// 查看任务详情
const viewTaskDetail = (taskId) => {
  router.push(`/download/${taskId}`)
}

// 下载结果
const downloadResult = async (taskId) => {
  try {
    await downloadApi.downloadResult(taskId)
  } catch (error) {
    console.error('下载结果失败:', error)
  }
}

// 删除任务
const deleteTask = async (taskId) => {
  if (confirm('确定要删除该任务吗？')) {
    try {
      await downloadApi.deleteTask(taskId)
      getTasks()
    } catch (error) {
      console.error('删除任务失败:', error)
    }
  }
}

// 创建任务
const createTask = () => {
  router.push('/download/create')
}

// 初始化数据
onMounted(() => {
  getTasks()
  // 每秒更新一次任务进度
  setInterval(() => {
    tasks.value.forEach(task => {
      if (task.status === 'running') {
        getTasks()
      }
    })
  }, 1000)
})
</script>

<style scoped>
.download-tasks {
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

.create-btn {
  padding: 0.7rem 1.5rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 1rem;
}

.create-btn:hover {
  background-color: #66b1ff;
}

/* 筛选区域 */
.task-filters {
  background-color: #f5f7fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.filter-row {
  display: flex;
  gap: 1rem;
  align-items: end;
}

.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-width: 200px;
}

.form-group label {
  font-size: 0.9rem;
  color: #606266;
}

.form-group input,
.form-group select {
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

.reset-btn {
  padding: 0.6rem 1.2rem;
  background-color: #909399;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 1rem;
  align-self: end;
}

.reset-btn:hover {
  background-color: #a6a9ad;
}

/* 任务列表 */
.task-list {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.tasks-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.tasks-table th,
.tasks-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.tasks-table th {
  background-color: #f5f7fa;
  color: #606266;
  font-weight: 600;
  font-size: 0.9rem;
  white-space: nowrap;
}

.tasks-table td {
  color: #303133;
}

/* 任务状态标签 */
.status-badge {
  padding: 0.3rem 0.8rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
  white-space: nowrap;
}

.status-badge.pending {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-badge.running {
  background-color: #ecf5ff;
  color: #409eff;
}

.status-badge.completed {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-badge.failed {
  background-color: #fef0f0;
  color: #f56c6c;
}

/* 进度条 */
.progress-container {
  position: relative;
  width: 100%;
  height: 20px;
  background-color: #f0f2f5;
  border-radius: 10px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: 10px;
}

.progress-bar.running {
  background-color: #409eff;
}

.progress-bar.completed {
  background-color: #67c23a;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 0.8rem;
  color: #606266;
  font-weight: 500;
}

/* 操作按钮 */
.task-actions {
  display: flex;
  gap: 0.5rem;
}

.detail-btn,
.download-btn,
.delete-btn {
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.detail-btn {
  background-color: #ecf5ff;
  color: #409eff;
}

.detail-btn:hover {
  background-color: #d9ecff;
}

.download-btn {
  background-color: #f0f9eb;
  color: #67c23a;
}

.download-btn:hover:not(:disabled) {
  background-color: #e1f3d8;
}

.download-btn:disabled {
  background-color: #f5f7fa;
  color: #c0c4cc;
  cursor: not-allowed;
}

.delete-btn {
  background-color: #fef0f0;
  color: #f56c6c;
}

.delete-btn:hover:not(:disabled) {
  background-color: #fbc4c4;
}

.delete-btn:disabled {
  background-color: #f5f7fa;
  color: #c0c4cc;
  cursor: not-allowed;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #909399;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.page-btn {
  padding: 0.6rem 1.2rem;
  background-color: #f5f7fa;
  color: #606266;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 0.9rem;
}

.page-btn:hover:not(:disabled) {
  background-color: #ecf5ff;
  color: #409eff;
  border-color: #c6e2ff;
}

.page-btn:disabled {
  background-color: #f5f7fa;
  color: #c0c4cc;
  cursor: not-allowed;
  border-color: #e4e7ed;
}

.page-info {
  color: #606266;
  font-size: 0.9rem;
}
</style>