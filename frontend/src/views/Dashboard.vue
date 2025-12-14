<template>
  <div class="min-h-screen p-6 lg:p-8">
    <!-- 头部 -->
    <header class="mb-8 text-center">
      <h1 class="text-4xl lg:text-6xl font-bold mb-3">
        <span class="gradient-text">SeekYourGuild</span>
      </h1>
      <p class="text-dark-400 text-base">群数据实时监控面板</p>
      
      <!-- 状态指示器和群选择器 -->
      <div class="flex items-center justify-center gap-4 mt-4 flex-wrap">
        <!-- 群选择器 -->
        <GroupSelector 
          :groups="availableGroups"
          :currentGroupId="groupId"
          @select="store.setGroupId"
        />
        
        <div class="text-dark-600 hidden sm:block">|</div>
        
        <div class="flex items-center gap-2">
          <span 
            class="w-2.5 h-2.5 rounded-full"
            :class="botOnline ? 'bg-green-500 animate-pulse' : 'bg-red-500'"
          ></span>
          <span class="text-sm text-dark-400">
            机器人 {{ botOnline ? '在线' : '离线' }}
          </span>
        </div>
        <div class="text-dark-600 hidden sm:block">|</div>
        <div class="text-sm text-dark-400">
          最后更新: {{ lastUpdateText }}
        </div>
      </div>
    </header>

    <!-- 未选择群提示 -->
    <div v-if="!groupId" class="text-center py-20">
      <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-dark-800 flex items-center justify-center">
        <svg class="w-10 h-10 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
      <h2 class="text-2xl font-bold text-white mb-2">请选择一个群</h2>
      <p class="text-dark-400">点击上方的群选择器选择要查看的群</p>
    </div>

    <!-- 主内容区域 -->
    <template v-else>
      <!-- 核心数据统计 -->
      <div class="mb-8">
        <StatsGrid 
          :members="overview?.current_members || 0"
          :joined="overview?.today_joined || 0"
          :left="overview?.today_left || 0"
          :messages="messageStats?.stats?.[0]?.total_messages || 0"
          :activeCount="activeUsers?.length || 0"
        />
      </div>

      <!-- 主要数据网格 -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-8">
        <!-- 左侧大卡片 - 群成员概览 -->
        <div class="lg:col-span-5">
          <MemberOverview 
            :data="overview"
            :loading="loading.overview"
          />
        </div>

        <!-- 右侧区域 -->
        <div class="lg:col-span-7 space-y-6">
        <!-- 今日数据 -->
        <div class="grid grid-cols-2 gap-4">
          <StatCard
            title="今日新增"
            :value="overview?.today_joined || 0"
            icon="👋"
            color="green"
            :loading="loading.overview"
          />
          <StatCard
            title="今日离开"
            :value="overview?.today_left || 0"
            icon="👋"
            color="red"
            :loading="loading.overview"
          />
        </div>

        <!-- 最新入群成员 -->
        <RecentMembers 
          :members="overview?.recent_members || []"
          :loading="loading.overview"
        />
      </div>
    </div>

    <!-- 发言统计区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-8">
      <!-- 小时分布图 -->
      <div class="lg:col-span-8">
        <HourlyChart 
          :data="messageStats?.hourly_distribution || []"
          :loading="loading.messageStats"
        />
      </div>

      <!-- 活跃用户排行 -->
      <div class="lg:col-span-4">
        <ActiveUsers 
          :users="activeUsers"
          :loading="loading.activeUsers"
        />
      </div>
    </div>

    <!-- 入退群趋势 -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-8">
      <!-- 入退群趋势图 -->
      <div class="lg:col-span-8">
        <MemberTrendChart 
          :data="memberChanges"
          :loading="loading.memberChanges"
        />
      </div>

      <!-- 群活跃度洞察 -->
      <div class="lg:col-span-4">
        <GroupInsights 
          :members="overview?.current_members || 0"
          :activeUsers="activeUsers?.length || 0"
          :todayMessages="messageStats?.stats?.[0]?.total_messages || 0"
          :todayJoined="overview?.today_joined || 0"
          :todayLeft="overview?.today_left || 0"
          :hourlyData="messageStats?.hourly_distribution || []"
        />
      </div>
    </div>
    </template>

    <!-- 底部信息 -->
    <footer class="text-center text-dark-500 text-sm py-8">
      <p>SeekYourGuild © 2024 - 群数据监控系统</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useStatsStore } from '@/stores/stats'
import { storeToRefs } from 'pinia'

import GroupSelector from '@/components/GroupSelector.vue'
import StatsGrid from '@/components/StatsGrid.vue'
import MemberOverview from '@/components/MemberOverview.vue'
import StatCard from '@/components/StatCard.vue'
import RecentMembers from '@/components/RecentMembers.vue'
import HourlyChart from '@/components/HourlyChart.vue'
import MemberTrendChart from '@/components/MemberTrendChart.vue'
import ActiveUsers from '@/components/ActiveUsers.vue'
import GroupInsights from '@/components/GroupInsights.vue'

const store = useStatsStore()
const { groupId, availableGroups, overview, memberChanges, messageStats, activeUsers, loading } = storeToRefs(store)

// 机器人状态
const botOnline = computed(() => overview.value?.bot_status === 'online')

// 最后更新时间
const lastUpdateText = computed(() => {
  if (!overview.value?.last_update) return '未知'
  const date = new Date(overview.value.last_update * 1000)
  return date.toLocaleTimeString('zh-CN')
})

// 定时刷新
let refreshInterval = null

onMounted(() => {
  store.init()
  refreshInterval = setInterval(() => {
    store.refreshAll()
  }, 30000) // 30秒刷新一次
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
