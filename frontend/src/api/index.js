import axios from 'axios'

const api = axios.create({
    baseURL: '/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// 响应拦截器
api.interceptors.response.use(
    (response) => {
        const { data } = response
        if (data.code !== 0) {
            return Promise.reject(new Error(data.message || '请求失败'))
        }
        return data.data
    },
    (error) => {
        console.error('API Error:', error)
        return Promise.reject(error)
    }
)

export default api

// 统计相关 API
export const statsApi = {
    // 获取可用群列表
    getGroups() {
        return api.get('/groups')
    },

    // 获取群概览
    getOverview(groupId) {
        return api.get('/stats/overview', { params: { group_id: groupId } })
    },

    // 获取成员变动
    getMemberChanges(groupId, params = {}) {
        return api.get('/stats/member-changes', {
            params: { group_id: groupId, ...params }
        })
    },

    // 获取消息统计
    getMessageStats(groupId, params = {}) {
        return api.get('/stats/messages', {
            params: { group_id: groupId, ...params }
        })
    },

    // 获取活跃用户
    getActiveUsers(groupId, params = {}) {
        return api.get('/stats/active-users', {
            params: { group_id: groupId, ...params }
        })
    }
}
