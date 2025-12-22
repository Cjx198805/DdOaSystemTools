import service from './index'

const dataDictionaryApi = {
  // 获取数据字典列表
  getDataDictionaries(params) {
    return service.get('/data-dictionaries', { params })
  },
  // 获取单个数据字典
  getDataDictionary(dictionaryId) {
    return service.get(`/data-dictionaries/${dictionaryId}`)
  },
  // 创建数据字典
  createDataDictionary(data) {
    return service.post('/data-dictionaries', data)
  },
  // 更新数据字典
  updateDataDictionary(dictionaryId, data) {
    return service.put(`/data-dictionaries/${dictionaryId}`, data)
  },
  // 删除数据字典
  deleteDataDictionary(dictionaryId) {
    return service.delete(`/data-dictionaries/${dictionaryId}`)
  }
}

export default dataDictionaryApi
