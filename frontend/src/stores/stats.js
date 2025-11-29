import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { statsApi } from '@/api'

export const useStatsStore = defineStore('stats', () => {
    // 当前群号 - 从 URL 参数获取
    const urlParams = new URLSearchParams(window.location.search)
    const groupId = ref(parseInt(urlParams.get('group_id')) || null)

    // 可用群列表
    const availableGroups = ref([])

    // 数据状态
    const overview = ref(null)
    const memberChanges = ref([])
    const messageStats = ref(null)
    const activeUsers = ref([])

    // 加载状态
    const loading = ref({
        overview: false,
        memberChanges: false,
        messageStats: false,
        activeUsers: false
    })

    // 错误状态
    const error = ref(null)

    // 获取概览数据
    async function fetchOverview() {
        loading.value.overview = true
        error.value = null
        try {
            overview.value = await statsApi.getOverview(groupId.value)
        } catch (e) {
            error.value = e.message
            console.error('Failed to fetch overview:', e)
        } finally {
            loading.value.overview = false
        }
    }

    // 获取成员变动
    async function fetchMemberChanges(params = {}) {
        loading.value.memberChanges = true
        try {
            const data = await statsApi.getMemberChanges(groupId.value, params)
            memberChanges.value = data.items || []
        } catch (e) {
            console.error('Failed to fetch member changes:', e)
        } finally {
            loading.value.memberChanges = false
        }
    }

    // 获取消息统计
    async function fetchMessageStats(params = {}) {
        loading.value.messageStats = true
        try {
            messageStats.value = await statsApi.getMessageStats(groupId.value, params)
        } catch (e) {
            console.error('Failed to fetch message stats:', e)
        } finally {
            loading.value.messageStats = false
        }
    }

    // 获取活跃用户
    async function fetchActiveUsers(params = {}) {
        loading.value.activeUsers = true
        try {
            const data = await statsApi.getActiveUsers(groupId.value, params)
            activeUsers.value = data.users || []
        } catch (e) {
            console.error('Failed to fetch active users:', e)
        } finally {
            loading.value.activeUsers = false
        }
    }

    // 获取可用群列表
    async function fetchAvailableGroups() {
        try {
            const data = await statsApi.getGroups()
            availableGroups.value = data.groups || []
            // 如果当前没有选中群，选择第一个
            if (!groupId.value && availableGroups.value.length > 0) {
                groupId.value = availableGroups.value[0].id
            }
        } catch (e) {
            console.error('Failed to fetch groups:', e)
        }
    }

    // 刷新所有数据
    async function refreshAll() {
        if (!groupId.value) return
        await Promise.all([
            fetchOverview(),
            fetchMemberChanges(),
            fetchMessageStats(),
            fetchActiveUsers()
        ])
    }

    // 初始化
    async function init() {
        await fetchAvailableGroups()
        if (groupId.value) {
            await refreshAll()
        }
    }

    // 设置群号
    function setGroupId(id) {
        groupId.value = id
        // 更新 URL
        const url = new URL(window.location)
        url.searchParams.set('group_id', id)
        window.history.replaceState({}, '', url)
        refreshAll()
    }

    return {
        groupId,
        availableGroups,
        overview,
        memberChanges,
        messageStats,
        activeUsers,
        loading,
        error,
        fetchOverview,
        fetchMemberChanges,
        fetchMessageStats,
        fetchActiveUsers,
        fetchAvailableGroups,
        refreshAll,
        init,
        setGroupId
    }
})
