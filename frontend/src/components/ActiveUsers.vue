<template>
  <div class="card p-6 h-full min-h-[380px]">
    <h3 class="text-dark-400 text-sm uppercase tracking-wider mb-4">
      活跃用户排行
    </h3>

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="space-y-3">
        <div v-for="i in 5" :key="i" class="flex items-center gap-3 p-3">
          <div class="skeleton w-10 h-10 rounded-full"></div>
          <div class="flex-1">
            <div class="skeleton h-4 w-24 rounded mb-2"></div>
            <div class="skeleton h-3 w-16 rounded"></div>
          </div>
          <div class="skeleton w-12 h-6 rounded"></div>
        </div>
      </div>
    </template>

    <!-- 用户列表 -->
    <template v-else>
      <div class="space-y-2">
        <div 
          v-for="(user, index) in users"
          :key="user.user_id"
          class="relative"
        >
          <!-- 进度条背景 -->
          <div 
            class="absolute inset-0 rounded-lg bg-gradient-to-r opacity-20"
            :class="getRankGradient(index)"
            :style="{ width: getProgressWidth(user.message_count) }"
          ></div>
          
          <!-- 内容 -->
          <div class="relative flex items-center gap-3 p-3 rounded-lg hover:bg-dark-800/30 transition-colors">
            <!-- 排名徽章 -->
            <div 
              class="w-6 h-6 rounded-full flex items-center justify-center font-bold text-xs absolute -left-1 -top-1 z-10 shadow-lg"
              :class="getRankClass(index)"
            >
              {{ index + 1 }}
            </div>

            <!-- 用户头像 -->
            <img 
              :src="`https://q1.qlogo.cn/g?b=qq&nk=${user.user_id}&s=640`"
              :alt="user.nickname || user.user_id"
              class="w-10 h-10 rounded-full object-cover ring-2 ring-dark-600"
              @error="(e) => e.target.src = 'https://q1.qlogo.cn/g?b=qq&nk=0&s=640'"
            />

            <!-- 用户信息 -->
            <div class="flex-1 min-w-0">
              <div class="text-white font-medium truncate">
                {{ user.nickname || `用户 ${user.user_id}` }}
              </div>
              <div class="text-xs text-dark-500">QQ: {{ user.user_id }}</div>
            </div>

            <!-- 消息数 -->
            <div class="text-right">
              <div class="text-lg font-bold number-display" :class="getRankTextColor(index)">
                {{ user.message_count }}
              </div>
              <div class="text-xs text-dark-500">条消息</div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div 
          v-if="users.length === 0"
          class="text-center text-dark-500 py-8"
        >
          暂无数据
        </div>
      </div>

      <!-- 时间段切换 -->
      <div class="mt-4 pt-4 border-t border-dark-700 flex justify-center gap-2">
        <button
          v-for="period in periods"
          :key="period.value"
          class="px-3 py-1 rounded-full text-sm transition-colors"
          :class="activePeriod === period.value 
            ? 'bg-blue-600 text-white' 
            : 'text-dark-400 hover:text-white hover:bg-dark-700'"
          @click="$emit('changePeriod', period.value)"
        >
          {{ period.label }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  users: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['changePeriod'])

const activePeriod = ref('day')

const periods = [
  { value: 'day', label: '今日' },
  { value: 'week', label: '本周' },
  { value: 'month', label: '本月' }
]

const maxCount = computed(() => {
  return Math.max(...props.users.map(u => u.message_count), 1)
})

function getProgressWidth(count) {
  return `${(count / maxCount.value) * 100}%`
}

function getRankClass(index) {
  if (index === 0) return 'bg-gradient-to-br from-yellow-400 to-yellow-600 text-white'
  if (index === 1) return 'bg-gradient-to-br from-gray-300 to-gray-500 text-white'
  if (index === 2) return 'bg-gradient-to-br from-amber-600 to-amber-800 text-white'
  return 'bg-dark-700 text-dark-300'
}

function getRankTextColor(index) {
  if (index === 0) return 'text-yellow-400'
  if (index === 1) return 'text-gray-300'
  if (index === 2) return 'text-amber-500'
  return 'text-blue-400'
}

function getRankGradient(index) {
  if (index === 0) return 'from-yellow-500 to-transparent'
  if (index === 1) return 'from-gray-400 to-transparent'
  if (index === 2) return 'from-amber-600 to-transparent'
  return 'from-blue-500 to-transparent'
}
</script>
