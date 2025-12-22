<template>
  <div class="api-test">
    <div class="page-header">
      <h2>API测试</h2>
    </div>
    
    <!-- API选择和参数配置 -->
    <div class="test-config">
      <div class="form-row">
        <div class="form-group">
          <label for="api-method">请求方法</label>
          <select id="api-method" v-model="testConfig.method">
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
          </select>
        </div>
        <div class="form-group">
          <label for="api-url">API URL</label>
          <input 
            type="text" 
            id="api-url" 
            v-model="testConfig.url" 
            placeholder="请输入API URL"
          />
        </div>
      </div>
      
      <!-- 请求参数 -->
      <div class="params-section">
        <h3>请求参数</h3>
        <div class="params-list">
          <div v-for="(param, index) in testConfig.params" :key="index" class="param-item">
            <input 
              type="text" 
              v-model="param.key" 
              placeholder="参数名"
              class="param-key"
            />
            <input 
              type="text" 
              v-model="param.value" 
              placeholder="参数值"
              class="param-value"
            />
            <button class="remove-param" @click="removeParam(index)">删除</button>
          </div>
          <button class="add-param" @click="addParam">添加参数</button>
        </div>
      </div>
      
      <!-- 请求头 -->
      <div class="headers-section">
        <h3>请求头</h3>
        <div class="headers-list">
          <div v-for="(header, index) in testConfig.headers" :key="index" class="header-item">
            <input 
              type="text" 
              v-model="header.key" 
              placeholder="Header名"
              class="header-key"
            />
            <input 
              type="text" 
              v-model="header.value" 
              placeholder="Header值"
              class="header-value"
            />
            <button class="remove-header" @click="removeHeader(index)">删除</button>
          </div>
          <button class="add-header" @click="addHeader">添加Header</button>
        </div>
      </div>
      
      <!-- 请求体 -->
      <div class="body-section">
        <h3>请求体</h3>
        <textarea 
          v-model="testConfig.body" 
          placeholder="请输入请求体JSON"
          rows="6"
          class="body-textarea"
        ></textarea>
      </div>
      
      <!-- 测试按钮 -->
      <div class="test-actions">
        <button class="test-btn" @click="runTest" :disabled="isRunning">
          {{ isRunning ? '测试中...' : '运行测试' }}
        </button>
        <button class="save-btn" @click="saveTestCase">保存为测试用例</button>
      </div>
    </div>
    
    <!-- 测试结果 -->
    <div class="test-result" v-if="testResult">
      <h3>测试结果</h3>
      <div class="result-header">
        <span class="status" :class="testResult.success ? 'success' : 'error'">
          {{ testResult.success ? '成功' : '失败' }}
        </span>
        <span class="time">响应时间: {{ testResult.responseTime }}ms</span>
      </div>
      <div class="result-content">
        <pre>{{ JSON.stringify(testResult.data, null, 2) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import apiTestApi from '../../api/apiTest'

// 测试配置
const testConfig = ref({
  method: 'GET',
  url: '',
  params: [],
  headers: [
    { key: 'Content-Type', value: 'application/json' }
  ],
  body: ''
})

// 测试结果
const testResult = ref(null)
const isRunning = ref(false)

// 添加参数
const addParam = () => {
  testConfig.value.params.push({ key: '', value: '' })
}

// 删除参数
const removeParam = (index) => {
  testConfig.value.params.splice(index, 1)
}

// 添加请求头
const addHeader = () => {
  testConfig.value.headers.push({ key: '', value: '' })
}

// 删除请求头
const removeHeader = (index) => {
  testConfig.value.headers.splice(index, 1)
}

// 运行测试
const runTest = async () => {
  isRunning.value = true
  testResult.value = null
  
  try {
    const startTime = Date.now()
    
    // 构建请求配置
    const requestConfig = {
      method: testConfig.value.method,
      url: testConfig.value.url,
      headers: testConfig.value.reduce((acc, header) => {
        if (header.key) {
          acc[header.key] = header.value
        }
        return acc
      }, {})
    }
    
    // 添加查询参数
    if (testConfig.value.method === 'GET' && testConfig.value.params.length > 0) {
      requestConfig.params = testConfig.value.params.reduce((acc, param) => {
        if (param.key) {
          acc[param.key] = param.value
        }
        return acc
      }, {})
    }
    
    // 添加请求体
    if (['POST', 'PUT', 'DELETE'].includes(testConfig.value.method) && testConfig.value.body) {
      requestConfig.data = JSON.parse(testConfig.value.body)
    }
    
    const response = await apiTestApi.runTest(requestConfig)
    const endTime = Date.now()
    
    testResult.value = {
      success: true,
      data: response.data,
      responseTime: endTime - startTime
    }
  } catch (error) {
    const endTime = Date.now()
    
    testResult.value = {
      success: false,
      data: error.response?.data || error.message,
      responseTime: endTime - startTime
    }
  } finally {
    isRunning.value = false
  }
}

// 保存测试用例
const saveTestCase = () => {
  // 这里可以添加保存测试用例的逻辑
  console.log('保存测试用例:', testConfig.value)
}
</script>

<style scoped>
.api-test {
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

/* 测试配置区域 */
.test-config {
  background-color: #f5f7fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
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

/* 参数和请求头区域 */
.params-section,
.headers-section,
.body-section {
  margin: 1.5rem 0;
}

.params-section h3,
.headers-section h3,
.body-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
  font-size: 1rem;
}

.params-list,
.headers-list {
  background-color: white;
  padding: 1rem;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.param-item,
.header-item {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.8rem;
  align-items: center;
}

.param-key,
.param-value,
.header-key,
.header-value {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 0.9rem;
}

.param-key:focus,
.param-value:focus,
.header-key:focus,
.header-value:focus {
  outline: none;
  border-color: #409eff;
}

.remove-param,
.remove-header {
  padding: 0.5rem;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
}

.remove-param:hover,
.remove-header:hover {
  background-color: #f78989;
}

.add-param,
.add-header {
  padding: 0.6rem 1.2rem;
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.add-param:hover,
.add-header:hover {
  background-color: #85ce61;
}

/* 请求体区域 */
.body-textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 0.9rem;
  font-family: monospace;
  resize: vertical;
  background-color: white;
}

.body-textarea:focus {
  outline: none;
  border-color: #409eff;
}

/* 测试按钮 */
.test-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.test-btn {
  padding: 0.8rem 1.5rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 1rem;
}

.test-btn:hover {
  background-color: #66b1ff;
}

.test-btn:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}

.save-btn {
  padding: 0.8rem 1.5rem;
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 1rem;
}

.save-btn:hover {
  background-color: #85ce61;
}

/* 测试结果区域 */
.test-result {
  background-color: #f5f7fa;
  padding: 1.5rem;
  border-radius: 8px;
}

.test-result h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e4e7ed;
}

.status {
  padding: 0.3rem 0.8rem;
  border-radius: 12px;
  font-size: 0.9rem;
  font-weight: bold;
}

.status.success {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status.error {
  background-color: #fef0f0;
  color: #f56c6c;
}

.result-content {
  background-color: white;
  padding: 1rem;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
  overflow-x: auto;
}

.result-content pre {
  margin: 0;
  font-family: monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #333;
}
</style>
