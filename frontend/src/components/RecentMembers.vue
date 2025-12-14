<template>
  <div class="card p-6 min-h-[280px]">
    <h3 class="text-dark-400 text-sm uppercase tracking-wider mb-4">
      最新入群成员
    </h3>

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="space-y-3">
        <div v-for="i in 5" :key="i" class="flex items-center gap-3 p-2">
          <div class="skeleton w-10 h-10 rounded-full"></div>
          <div class="flex-1">
            <div class="skeleton h-4 w-24 rounded mb-2"></div>
            <div class="skeleton h-3 w-16 rounded"></div>
          </div>
        </div>
      </div>
    </template>

    <!-- 成员列表 -->
    <template v-else>
      <div class="space-y-3">
        <div 
          v-for="(member, index) in members"
          :key="member.user_id"
          class="flex items-center gap-3 p-2 rounded-lg hover:bg-dark-800/50 transition-colors group"
          :style="{ animationDelay: `${index * 0.1}s` }"
        >
          <!-- 头像 -->
          <div class="relative">
            <div 
              class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-bold"
            >
              {{ getInitial(member) }}
            </div>
            <!-- 新成员标记 -->
            <div 
              v-if="isNew(member)"
              class="absolute -top-1 -right-1 w-4 h-4 bg-green-500 rounded-full flex items-center justify-center"
            >
              <span class="text-xs">✓</span>
            </div>
          </div>

          <!-- 信息 -->
          <div class="flex-1 min-w-0">
            <div class="text-white font-medium truncate group-hover:text-blue-400 transition-colors">
              {{ member.card || member.nickname || 'Unknown' }}
            </div>
            <div class="text-xs text-dark-400">
              {{ formatTime(member.join_time) }}
            </div>
          </div>

          <!-- 序号 -->
          <div class="text-3xl font-bold text-dark-700 number-display">
            {{ String(index + 1).padStart(2, '0') }}
          </div>
        </div>

        <!-- 空状态 -->
        <div 
          v-if="members.length === 0"
          class="text-center text-dark-500 py-8"
        >
          暂无新成员
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
const props = defineProps({
  members: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

function getInitial(member) {
  const name = member.card || member.nickname || 'U'
  return name.charAt(0).toUpperCase()
}

function formatTime(timestamp) {
  if (!timestamp) return '未知'
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = now - date

  // 小于1小时
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}分钟前`
  }
  
  // 小于24小时
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}小时前`
  }

  // 显示日期
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  })
}

function isNew(member) {
  if (!member.join_time) return false
  const joinDate = new Date(member.join_time * 1000)
  const now = new Date()
  return (now - joinDate) < 86400000 // 24小时内
}
</script>
