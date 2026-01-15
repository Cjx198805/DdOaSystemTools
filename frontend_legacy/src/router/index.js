import { createRouter, createWebHistory } from 'vue-router'

// 导入页面组件
const Login = () => import('../views/auth/Login.vue')
const Layout = () => import('../components/Layout.vue')

// 权限管理页面
const UserManagement = () => import('../views/permission/User.vue')
const RoleManagement = () => import('../views/permission/Role.vue')
const FieldPermission = () => import('../views/permission/FieldPermission.vue')
const DataDictionary = () => import('../views/permission/Dictionary.vue')

// API测试页面
const ApiTest = () => import('../views/api/Test.vue')
const TestCaseManagement = () => import('../views/api/TestCaseManagement.vue')
const TestHistory = () => import('../views/api/TestHistory.vue')

// 下载任务页面
const TaskManagement = () => import('../views/download/Index.vue')
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
    redirect: '/system/user',
    meta: { title: '首页', requireAuth: true },
    children: [
      // 系统管理
      {
        path: '/system/user',
        name: 'UserManagement',
        component: UserManagement,
        meta: { title: '用户管理', requireAuth: true }
      },
      {
        path: '/system/role',
        name: 'RoleManagement',
        component: RoleManagement,
        meta: { title: '角色管理', requireAuth: true }
      },
      {
        path: '/system/field-permission',
        name: 'FieldPermission',
        component: FieldPermission,
        meta: { title: '字段权限', requireAuth: true }
      },
      {
        path: '/system/dict',
        name: 'DataDictionary',
        component: DataDictionary,
        meta: { title: '数据字典', requireAuth: true }
      },
      // 业务功能
      {
        path: '/business/download-task',
        name: 'TaskManagement',
        component: TaskManagement,
        meta: { title: '下载任务', requireAuth: true }
      },
      {
        path: '/business/progress',
        name: 'TaskProgress',
        component: TaskProgress,
        meta: { title: '任务进度', requireAuth: true }
      },
      {
        path: '/business/result',
        name: 'ResultDownload',
        component: ResultDownload,
        meta: { title: '结果下载', requireAuth: true }
      },
      // API测试
      {
        path: '/api-test/case',
        name: 'TestCaseManagement',
        component: TestCaseManagement,
        meta: { title: '测试用例', requireAuth: true }
      },
      {
        path: '/api-test/history',
        name: 'TestHistory',
        component: TestHistory,
        meta: { title: '测试历史', requireAuth: true }
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