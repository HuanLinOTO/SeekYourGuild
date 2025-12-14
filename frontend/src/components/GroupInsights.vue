<template>
  <div class="card p-6 min-h-[380px]">
    <h3 class="text-dark-400 text-sm uppercase tracking-wider mb-6 flex items-center gap-2">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      群活跃度洞察
    </h3>

    <div class="space-y-6">
      <!-- 群健康指数 -->
      <div class="relative">
        <div class="flex justify-between items-center mb-2">
          <span class="text-dark-400 text-sm">群健康指数</span>
          <span class="text-2xl font-bold" :class="healthColor">{{ healthScore }}</span>
        </div>
        <div class="h-2 bg-dark-700 rounded-full overflow-hidden">
          <div 
            class="h-full rounded-full transition-all duration-1000 ease-out"
            :class="healthGradient"
            :style="{ width: `${healthScore}%` }"
          ></div>
        </div>
        <div class="flex justify-between text-xs text-dark-500 mt-1">
          <span>低活跃</span>
          <span>高活跃</span>
        </div>
      </div>

      <!-- 关键指标 -->
      <div class="grid grid-cols-2 gap-4">
        <!-- 消息密度 -->
        <div class="bg-dark-800/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-lg bg-purple-500/20 flex items-center justify-center">
              <svg class="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
              </svg>
            </div>
            <span class="text-xs text-dark-400">消息密度</span>
          </div>
          <div class="text-xl font-bold text-purple-400">
            {{ messageDensity }}
          </div>
          <div class="text-xs text-dark-500">条/人/日</div>
        </div>

        <!-- 互动率 -->
        <div class="bg-dark-800/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-lg bg-cyan-500/20 flex items-center justify-center">
              <svg class="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z" />
              </svg>
            </div>
            <span class="text-xs text-dark-400">互动率</span>
          </div>
          <div class="text-xl font-bold text-cyan-400">
            {{ interactionRate }}%
          </div>
          <div class="text-xs text-dark-500">活跃/总人数</div>
        </div>

        <!-- 净增长 -->
        <div class="bg-dark-800/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-lg bg-green-500/20 flex items-center justify-center">
              <svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
            <span class="text-xs text-dark-400">今日净增</span>
          </div>
          <div class="text-xl font-bold" :class="netGrowth >= 0 ? 'text-green-400' : 'text-red-400'">
            {{ netGrowth >= 0 ? '+' : '' }}{{ netGrowth }}
          </div>
          <div class="text-xs text-dark-500">成员变化</div>
        </div>

        <!-- 留存率 -->
        <div class="bg-dark-800/50 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-lg bg-amber-500/20 flex items-center justify-center">
              <svg class="w-4 h-4 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <span class="text-xs text-dark-400">留存率</span>
          </div>
          <div class="text-xl font-bold text-amber-400">
            {{ retentionRate }}%
          </div>
          <div class="text-xs text-dark-500">成员稳定性</div>
        </div>
      </div>

      <!-- 时段分析 -->
      <div class="pt-4 border-t border-dark-700">
        <div class="text-xs text-dark-400 mb-3">活跃时段分析</div>
        <div class="flex items-center gap-3">
          <div class="flex-1">
            <div class="text-sm text-dark-300 mb-1">🌙 深夜档 (0-6点)</div>
            <div class="h-1.5 bg-dark-700 rounded-full overflow-hidden">
              <div class="h-full bg-indigo-500 rounded-full" :style="{ width: `${nightActivity}%` }"></div>
            </div>
          </div>
        </div>
        <div class="flex items-center gap-3 mt-2">
          <div class="flex-1">
            <div class="text-sm text-dark-300 mb-1">☀️ 日间档 (6-18点)</div>
            <div class="h-1.5 bg-dark-700 rounded-full overflow-hidden">
              <div class="h-full bg-amber-500 rounded-full" :style="{ width: `${dayActivity}%` }"></div>
            </div>
          </div>
        </div>
        <div class="flex items-center gap-3 mt-2">
          <div class="flex-1">
            <div class="text-sm text-dark-300 mb-1">🌆 黄金档 (18-24点)</div>
            <div class="h-1.5 bg-dark-700 rounded-full overflow-hidden">
              <div class="h-full bg-orange-500 rounded-full" :style="{ width: `${eveningActivity}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  members: { type: Number, default: 0 },
  activeUsers: { type: Number, default: 0 },
  todayMessages: { type: Number, default: 0 },
  todayJoined: { type: Number, default: 0 },
  todayLeft: { type: Number, default: 0 },
  hourlyData: { type: Array, default: () => [] }
})

// 群健康指数 (0-100)
const healthScore = computed(() => {
  if (props.members === 0) return 0
  
  // 活跃率贡献 (40分)
  const activeRate = (props.activeUsers / props.members) * 100
  const activeScore = Math.min(activeRate * 4, 40)
  
  // 消息密度贡献 (30分)
  const density = props.todayMessages / props.members
  const densityScore = Math.min(density * 10, 30)
  
  // 留存率贡献 (30分)
  const retention = props.todayLeft === 0 ? 100 : 
    Math.max(0, 100 - (props.todayLeft / props.members) * 100)
  const retentionScore = (retention / 100) * 30
  
  return Math.round(activeScore + densityScore + retentionScore)
})

const healthColor = computed(() => {
  if (healthScore.value >= 70) return 'text-green-400'
  if (healthScore.value >= 40) return 'text-amber-400'
  return 'text-red-400'
})

const healthGradient = computed(() => {
  if (healthScore.value >= 70) return 'bg-gradient-to-r from-green-500 to-emerald-400'
  if (healthScore.value >= 40) return 'bg-gradient-to-r from-amber-500 to-yellow-400'
  return 'bg-gradient-to-r from-red-500 to-orange-400'
})

// 消息密度 (条/人/日)
const messageDensity = computed(() => {
  if (props.members === 0) return '0.0'
  return (props.todayMessages / props.members).toFixed(1)
})

// 互动率
const interactionRate = computed(() => {
  if (props.members === 0) return '0.0'
  return ((props.activeUsers / props.members) * 100).toFixed(1)
})

// 净增长
const netGrowth = computed(() => {
  return props.todayJoined - props.todayLeft
})

// 留存率
const retentionRate = computed(() => {
  if (props.members === 0) return '100.0'
  const lostRate = (props.todayLeft / props.members) * 100
  return Math.max(0, 100 - lostRate).toFixed(1)
})

// 时段活跃度分析
const totalHourlyMessages = computed(() => {
  return props.hourlyData.reduce((sum, h) => sum + (h.count || 0), 0) || 1
})

const nightActivity = computed(() => {
  const nightMessages = props.hourlyData
    .filter(h => h.hour >= 0 && h.hour < 6)
    .reduce((sum, h) => sum + (h.count || 0), 0)
  return Math.round((nightMessages / totalHourlyMessages.value) * 100)
})

const dayActivity = computed(() => {
  const dayMessages = props.hourlyData
    .filter(h => h.hour >= 6 && h.hour < 18)
    .reduce((sum, h) => sum + (h.count || 0), 0)
  return Math.round((dayMessages / totalHourlyMessages.value) * 100)
})

const eveningActivity = computed(() => {
  const eveningMessages = props.hourlyData
    .filter(h => h.hour >= 18 && h.hour < 24)
    .reduce((sum, h) => sum + (h.count || 0), 0)
  return Math.round((eveningMessages / totalHourlyMessages.value) * 100)
})
</script>
