<template>
  <div class="dictionary-management">
    <div class="page-header">
      <h2>数据字典管理</h2>
      <button class="add-btn" @click="showAddModal = true">添加字典</button>
    </div>
    
    <!-- 搜索和筛选 -->
    <div class="search-filter">
      <input 
        type="text" 
        v-model="searchQuery" 
        placeholder="搜索字典名称或标识" 
        class="search-input"
      />
      <button class="search-btn" @click="handleSearch">搜索</button>
    </div>
    
    <!-- 数据字典列表 -->
    <div class="dictionary-list">
      <table class="dictionary-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>字典名称</th>
            <th>字典标识</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="dict in dictionaries" :key="dict.id">
            <td>{{ dict.id }}</td>
            <td>{{ dict.name }}</td>
            <td>{{ dict.code }}</td>
            <td>
              <span :class="['status-tag', dict.status === 1 ? 'active' : 'inactive']">
                {{ dict.status === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td>{{ formatDate(dict.created_at) }}</td>
            <td class="action-buttons">
              <button class="edit-btn" @click="handleEdit(dict)">编辑</button>
              <button class="detail-btn" @click="handleDetail(dict)">字典项</button>
              <button class="delete-btn" @click="handleDelete(dict.id)">删除</button>
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
    
    <!-- 添加/编辑字典弹窗 -->
    <div class="modal" v-if="showAddModal || showEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ showAddModal ? '添加字典' : '编辑字典' }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="name">字典名称</label>
              <input 
                type="text" 
                id="name" 
                v-model="formData.name" 
                placeholder="请输入字典名称" 
                required
              />
            </div>
            <div class="form-group">
              <label for="code">字典标识</label>
              <input 
                type="text" 
                id="code" 
                v-model="formData.code" 
                placeholder="请输入字典标识" 
                required
              />
            </div>
            <div class="form-group">
              <label for="description">描述</label>
              <textarea 
                id="description" 
                v-model="formData.description" 
                placeholder="请输入字典描述" 
                rows="3"
              ></textarea>
            </div>
            <div class="form-group">
              <label for="status">状态</label>
              <select id="status" v-model="formData.status">
                <option value="1">启用</option>
                <option value="0">禁用</option>
              </select>
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
    
    <!-- 字典项管理弹窗 -->
    <div class="modal dictionary-item-modal" v-if="showItemModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>字典项管理 - {{ selectedDict.name }}</h3>
          <button class="close-btn" @click="closeItemModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="item-header">
            <button class="add-item-btn" @click="showAddItemModal = true">添加字典项</button>
          </div>
          
          <!-- 字典项列表 -->
          <table class="item-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>字典项名称</th>
                <th>字典项值</th>
                <th>排序</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in dictionaryItems" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.name }}</td>
                <td>{{ item.value }}</td>
                <td>{{ item.sort }}</td>
                <td>
                  <span :class="['status-tag', item.status === 1 ? 'active' : 'inactive']">
                    {{ item.status === 1 ? '启用' : '禁用' }}
                  </span>
                </td>
                <td class="action-buttons">
                  <button class="edit-item-btn" @click="handleEditItem(item)">编辑</button>
                  <button class="delete-item-btn" @click="handleDeleteItem(item.id)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- 添加/编辑字典项弹窗 -->
    <div class="modal" v-if="showAddItemModal || showEditItemModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ showAddItemModal ? '添加字典项' : '编辑字典项' }}</h3>
          <button class="close-btn" @click="closeItemFormModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleItemSubmit">
            <div class="form-group">
              <label for="item-name">字典项名称</label>
              <input 
                type="text" 
                id="item-name" 
                v-model="itemFormData.name" 
                placeholder="请输入字典项名称" 
                required
              />
            </div>
            <div class="form-group">
              <label for="item-value">字典项值</label>
              <input 
                type="text" 
                id="item-value" 
                v-model="itemFormData.value" 
                placeholder="请输入字典项值" 
                required
              />
            </div>
            <div class="form-group">
              <label for="item-sort">排序</label>
              <input 
                type="number" 
                id="item-sort" 
                v-model.number="itemFormData.sort" 
                placeholder="请输入排序值" 
                required
              />
            </div>
            <div class="form-group">
              <label for="item-status">状态</label>
              <select id="item-status" v-model="itemFormData.status">
                <option value="1">启用</option>
                <option value="0">禁用</option>
              </select>
            </div>
            <div class="modal-footer">
              <button type="submit" class="submit-btn">
                {{ showAddItemModal ? '添加' : '保存' }}
              </button>
              <button type="button" class="cancel-btn" @click="closeItemFormModal">取消</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import dictionaryApi from '../../api/dataDictionary'

// 数据状态
const dictionaries = ref([])
const dictionaryItems = ref([])
const searchQuery = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const loading = ref(false)

// 弹窗状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showItemModal = ref(false)
const showAddItemModal = ref(false)
const showEditItemModal = ref(false)

// 选中的字典
const selectedDict = ref({})

// 表单数据
const formData = ref({
  id: '',
  name: '',
  code: '',
  description: '',
  status: '1'
})

// 字典项表单数据
const itemFormData = ref({
  id: '',
  name: '',
  value: '',
  sort: 0,
  status: '1'
})

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取数据字典列表
const getDictionaries = async () => {
  loading.value = true
  try {
    const res = await dictionaryApi.getDictionaries({
      page: currentPage.value,
      keyword: searchQuery.value
    })
    dictionaries.value = res.data.dictionaries
    totalPages.value = res.data.total_pages
  } catch (error) {
    console.error('获取数据字典失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取字典项列表
const getDictionaryItems = async (dictId) => {
  try {
    const res = await dictionaryApi.getDictionaryItems(dictId)
    dictionaryItems.value = res.data.items
  } catch (error) {
    console.error('获取字典项失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getDictionaries()
}

// 编辑字典
const handleEdit = (dict) => {
  showEditModal.value = true
  formData.value = {
    id: dict.id,
    name: dict.name,
    code: dict.code,
    description: dict.description,
    status: dict.status.toString()
  }
}

// 删除字典
const handleDelete = async (dictId) => {
  if (confirm('确定要删除这个字典吗？')) {
    try {
      await dictionaryApi.deleteDictionary(dictId)
      getDictionaries()
    } catch (error) {
      console.error('删除字典失败:', error)
    }
  }
}

// 查看字典项
const handleDetail = (dict) => {
  selectedDict.value = dict
  showItemModal.value = true
  getDictionaryItems(dict.id)
}

// 提交表单
const handleSubmit = async () => {
  try {
    if (showAddModal.value) {
      await dictionaryApi.createDictionary(formData.value)
    } else {
      await dictionaryApi.updateDictionary(formData.value.id, formData.value)
    }
    closeModal()
    getDictionaries()
  } catch (error) {
    console.error('保存字典失败:', error)
  }
}

// 关闭弹窗
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  formData.value = {
    id: '',
    name: '',
    code: '',
    description: '',
    status: '1'
  }
}

// 关闭字典项弹窗
const closeItemModal = () => {
  showItemModal.value = false
  selectedDict.value = {}
  dictionaryItems.value = []
}

// 编辑字典项
const handleEditItem = (item) => {
  showEditItemModal.value = true
  itemFormData.value = {
    id: item.id,
    name: item.name,
    value: item.value,
    sort: item.sort,
    status: item.status.toString()
  }
}

// 删除字典项
const handleDeleteItem = async (itemId) => {
  if (confirm('确定要删除这个字典项吗？')) {
    try {
      await dictionaryApi.deleteDictionaryItem(itemId)
      getDictionaryItems(selectedDict.value.id)
    } catch (error) {
      console.error('删除字典项失败:', error)
    }
  }
}

// 提交字典项表单
const handleItemSubmit = async () => {
  try {
    const data = {
      ...itemFormData.value,
      dictionary_id: selectedDict.value.id
    }
    
    if (showAddItemModal.value) {
      await dictionaryApi.createDictionaryItem(data)
    } else {
      await dictionaryApi.updateDictionaryItem(itemFormData.value.id, data)
    }
    closeItemFormModal()
    getDictionaryItems(selectedDict.value.id)
  } catch (error) {
    console.error('保存字典项失败:', error)
  }
}

// 关闭字典项表单弹窗
const closeItemFormModal = () => {
  showAddItemModal.value = false
  showEditItemModal.value = false
  itemFormData.value = {
    id: '',
    name: '',
    value: '',
    sort: 0,
    status: '1'
  }
}

// 初始化
onMounted(() => {
  getDictionaries()
})
</script>

<style scoped>
.dictionary-management {
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

.dictionary-list {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.dictionary-table,
.item-table {
  width: 100%;
  border-collapse: collapse;
}

.dictionary-table th,
.dictionary-table td,
.item-table th,
.item-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.dictionary-table th,
.item-table th {
  background-color: #f5f7fa;
  font-weight: bold;
  color: #333;
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

.edit-btn,
.edit-item-btn {
  padding: 0.4rem 0.8rem;
  background-color: #e6a23c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.edit-btn:hover,
.edit-item-btn:hover {
  background-color: #ebb563;
}

.detail-btn {
  padding: 0.4rem 0.8rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.detail-btn:hover {
  background-color: #66b1ff;
}

.delete-btn,
.delete-item-btn {
  padding: 0.4rem 0.8rem;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.delete-btn:hover,
.delete-item-btn:hover {
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

.dictionary-item-modal .modal-content {
  max-width: 800px;
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

/* 字典项管理弹窗样式 */
.item-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}

.add-item-btn {
  padding: 0.6rem 1.2rem;
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.add-item-btn:hover {
  background-color: #85ce61;
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
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.6rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
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
</style>
