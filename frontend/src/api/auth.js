import service from './index'

const authApi = {
  // 登录
  login(data) {
    return service.post('/user/login', data)
  },
  // 获取用户信息
  getUserInfo() {
    return service.get('/auth/user-info')
  },
  // 退出登录
  logout() {
    return service.post('/auth/logout')
  }
}

export default authApi
