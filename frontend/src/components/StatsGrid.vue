<template>
  <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
    <!-- 总成员数 -->
    <div class="stats-card group">
      <div class="stats-icon bg-gradient-to-br from-blue-500 to-cyan-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
      <div class="stats-content">
        <div class="stats-value text-blue-400">
          {{ formatNumber(members) }}
        </div>
        <div class="stats-label">总成员</div>
      </div>
      <div class="stats-trend text-dark-500">
        <span class="text-xs">群活跃度</span>
      </div>
    </div>

    <!-- 今日新增 -->
    <div class="stats-card group">
      <div class="stats-icon bg-gradient-to-br from-green-500 to-emerald-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
        </svg>
      </div>
      <div class="stats-content">
        <div class="stats-value text-green-400">
          +{{ joined }}
        </div>
        <div class="stats-label">今日新增</div>
      </div>
      <div class="stats-trend" :class="joined > 0 ? 'text-green-500' : 'text-dark-500'">
        <span v-if="joined > 0">↑ 新人加入</span>
        <span v-else>暂无变动</span>
      </div>
    </div>

    <!-- 今日离开 -->
    <div class="stats-card group">
      <div class="stats-icon bg-gradient-to-br from-red-500 to-orange-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" />
        </svg>
      </div>
      <div class="stats-content">
        <div class="stats-value text-red-400">
          -{{ left }}
        </div>
        <div class="stats-label">今日离开</div>
      </div>
      <div class="stats-trend" :class="left > 0 ? 'text-red-500' : 'text-dark-500'">
        <span v-if="left > 0">↓ 成员流失</span>
        <span v-else>零流失 ✓</span>
      </div>
    </div>

    <!-- 今日消息 -->
    <div class="stats-card group">
      <div class="stats-icon bg-gradient-to-br from-purple-500 to-pink-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <div class="stats-content">
        <div class="stats-value text-purple-400">
          {{ formatNumber(messages) }}
        </div>
        <div class="stats-label">今日消息</div>
      </div>
      <div class="stats-trend text-purple-500">
        <span>{{ avgPerHour }} 条/时</span>
      </div>
    </div>

    <!-- 活跃用户 -->
    <div class="stats-card group col-span-2 md:col-span-1">
      <div class="stats-icon bg-gradient-to-br from-amber-500 to-yellow-400">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      <div class="stats-content">
        <div class="stats-value text-amber-400">
          {{ activeCount }}
        </div>
        <div class="stats-label">活跃用户</div>
      </div>
      <div class="stats-trend text-amber-500">
        <span>{{ activeRate }}% 活跃率</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  members: { type: Number, default: 0 },
  joined: { type: Number, default: 0 },
  left: { type: Number, default: 0 },
  messages: { type: Number, default: 0 },
  activeCount: { type: Number, default: 0 }
})

function formatNumber(num) {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const avgPerHour = computed(() => {
  return Math.round(props.messages / 24)
})

const activeRate = computed(() => {
  if (props.members === 0) return 0
  return ((props.activeCount / props.members) * 100).toFixed(1)
})
</script>

<style scoped>
.stats-card {
  @apply relative overflow-hidden rounded-xl p-4 transition-all duration-300;
  @apply bg-gradient-to-br from-dark-800/80 to-dark-900/80;
  @apply border border-dark-700/50 hover:border-dark-600;
  @apply hover:shadow-lg hover:shadow-blue-500/5;
  @apply hover:-translate-y-1;
}

.stats-icon {
  @apply w-12 h-12 rounded-xl flex items-center justify-center text-white;
  @apply shadow-lg mb-3;
  @apply group-hover:scale-110 transition-transform duration-300;
}

.stats-content {
  @apply relative z-10;
}

.stats-value {
  @apply text-3xl font-bold tracking-tight;
  font-variant-numeric: tabular-nums;
}

.stats-label {
  @apply text-dark-400 text-sm mt-1;
}

.stats-trend {
  @apply text-xs mt-2;
}
</style>
