import service from './index'

const dataDictionaryApi = {
  // 获取数据字典列表
  getDataDictionaries(params) {
    return service.get('/data-dictionary', { params })
  },
  // 获取单个数据字典
  getDataDictionary(dictionaryId) {
    return service.get(`/data-dictionary/${dictionaryId}`)
  },
  // 创建数据字典
  createDataDictionary(data) {
    return service.post('/data-dictionary', data)
  },
  // 更新数据字典
  updateDataDictionary(dictionaryId, data) {
    return service.put(`/data-dictionary/${dictionaryId}`, data)
  },
  // 删除数据字典
  deleteDataDictionary(dictionaryId) {
    return service.delete(`/data-dictionary/${dictionaryId}`)
  }
}

export default dataDictionaryApi
