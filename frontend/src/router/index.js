import { createRouter, createWebHistory } from 'vue-router'

// 导入页面组件
const Login = () => import('../views/auth/Login.vue')
const Layout = () => import('../components/Layout.vue')

// 权限管理页面
const UserManagement = () => import('../views/permission/UserManagement.vue')
const RoleManagement = () => import('../views/permission/RoleManagement.vue')
const FieldPermission = () => import('../views/permission/FieldPermission.vue')
const DataDictionary = () => import('../views/permission/DataDictionary.vue')

// API测试页面
const ApiTest = () => import('../views/api/ApiTest.vue')
const TestCaseManagement = () => import('../views/api/TestCaseManagement.vue')
const TestHistory = () => import('../views/api/TestHistory.vue')

// 下载任务页面
const TaskManagement = () => import('../views/download/TaskManagement.vue')
const TaskProgress = () => import('../views/download/TaskProgress.vue')
const ResultDownload = () => import('../views/download/ResultDownload.vue')

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录' }
  },
  {
    path: '/',
    name: 'Layout',
    component: Layout,
    redirect: '/permission/user',
    meta: { title: '首页', requireAuth: true },
    children: [
      // 权限管理路由
      {
        path: '/permission/user',
        name: 'UserManagement',
        component: UserManagement,
        meta: { title: '用户管理', requireAuth: true }
      },
      {
        path: '/permission/role',
        name: 'RoleManagement',
        component: RoleManagement,
        meta: { title: '角色管理', requireAuth: true }
      },
      {
        path: '/permission/field',
        name: 'FieldPermission',
        component: FieldPermission,
        meta: { title: '字段权限', requireAuth: true }
      },
      {
        path: '/permission/dictionary',
        name: 'DataDictionary',
        component: DataDictionary,
        meta: { title: '数据字典', requireAuth: true }
      },
      // API测试路由
      {
        path: '/api/test',
        name: 'ApiTest',
        component: ApiTest,
        meta: { title: 'API测试', requireAuth: true }
      },
      {
        path: '/api/testcase',
        name: 'TestCaseManagement',
        component: TestCaseManagement,
        meta: { title: '测试用例管理', requireAuth: true }
      },
      {
        path: '/api/history',
        name: 'TestHistory',
        component: TestHistory,
        meta: { title: '测试历史记录', requireAuth: true }
      },
      // 下载任务路由
      {
        path: '/download/task',
        name: 'TaskManagement',
        component: TaskManagement,
        meta: { title: '任务管理', requireAuth: true }
      },
      {
        path: '/download/progress',
        name: 'TaskProgress',
        component: TaskProgress,
        meta: { title: '任务进度', requireAuth: true }
      },
      {
        path: '/download/result',
        name: 'ResultDownload',
        component: ResultDownload,
        meta: { title: '结果下载', requireAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - DdOaListDownload` : 'DdOaListDownload'
  
  // 检查是否需要登录
  if (to.meta.requireAuth) {
    const token = localStorage.getItem('token')
    if (token) {
      next()
    } else {
      next({ path: '/login' })
    }
  } else {
    next()
  }
})

export default router