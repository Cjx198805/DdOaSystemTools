import axios from 'axios'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
  baseURL: '/api/v1', // 后端 API 版本前缀
  timeout: 30000 // 增加超时时间到 30s
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    // 如果 code 不为 200，则判定为错误
    if (res.code && res.code !== 200) {
      ElMessage.error(res.message || '服务异常')

      // 401: Token 失效
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        location.reload() // 触发路由守卫跳转登录
      }

      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res
  },
  error => {
    console.error('响应拦截错误:', error)
    let message = error.message
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message
    } else if (message.includes('timeout')) {
      message = '请求超时，请检查网络或重试'
    }
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default service
