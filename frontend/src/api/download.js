import api from './index'

export default {
    // 获取任务列表
    getTasks(params) {
        return api.get('/download-task', { params })
    },

    // 获取任务详情
    getTask(id) {
        return api.get(`/download-task/${id}`)
    },

    // 创建任务
    createTask(data) {
        return api.post('/download-task', data)
    },

    // 删除任务
    deleteTask(id) {
        return api.delete(`/download-task/${id}`)
    },

    // 获取用户任务
    getUserTasks(userId) {
        return api.get(`/download-task/user/${userId}`)
    },

    // 下载结果
    downloadResult(taskId) {
        return api.get(`/download-task/result/${taskId}`, {
            responseType: 'blob'
        })
    }
}
