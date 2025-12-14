<template>
  <div class="card p-6 h-full min-h-[420px]">
    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="space-y-6">
        <div class="skeleton h-5 w-32 rounded"></div>
        <div class="flex justify-center">
          <div class="skeleton h-24 w-48 rounded-lg"></div>
        </div>
        <div class="skeleton h-4 w-24 mx-auto rounded"></div>
        <div class="flex justify-center">
          <div class="skeleton h-48 w-48 rounded-full"></div>
        </div>
        <div class="flex justify-between">
          <div class="skeleton h-12 w-20 rounded"></div>
          <div class="skeleton h-12 w-20 rounded"></div>
        </div>
      </div>
    </template>

    <!-- 实际内容 -->
    <template v-else>
      <h2 class="text-dark-400 text-sm uppercase tracking-wider mb-6">群成员概览</h2>
      
      <!-- 大数字展示 -->
      <div class="text-center mb-8">
        <div class="text-7xl lg:text-8xl font-bold number-display gradient-text mb-2">
          {{ data?.current_members || 0 }}
        </div>
        <div class="text-dark-400">当前成员数</div>
      </div>

      <!-- 圆形装饰 -->
      <div class="relative mx-auto w-48 h-48 mb-6">
        <svg class="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
          <!-- 背景圆 -->
          <circle
            cx="50"
            cy="50"
            r="45"
            fill="none"
            stroke="rgba(30, 41, 59, 0.5)"
            stroke-width="8"
          />
          <!-- 进度圆 -->
          <circle
            cx="50"
            cy="50"
            r="45"
            fill="none"
            stroke="url(#gradient)"
            stroke-width="8"
            stroke-linecap="round"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="progressOffset"
            class="transition-all duration-1000"
          />
          <defs>
            <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stop-color="#3b82f6" />
              <stop offset="50%" stop-color="#8b5cf6" />
              <stop offset="100%" stop-color="#ec4899" />
            </linearGradient>
          </defs>
        </svg>
        
        <!-- 中心数据 -->
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <div class="text-3xl font-bold text-white">
            {{ netChange >= 0 ? '+' : '' }}{{ netChange }}
          </div>
          <div class="text-sm text-dark-400">今日净变化</div>
        </div>
      </div>

      <!-- 底部统计 -->
      <div class="flex justify-between text-center">
        <div>
          <div class="text-2xl font-bold text-green-400">
            +{{ data?.today_joined || 0 }}
          </div>
          <div class="text-xs text-dark-400 mt-1">新增</div>
        </div>
        <div class="w-px bg-dark-700"></div>
        <div>
          <div class="text-2xl font-bold text-red-400">
            -{{ data?.today_left || 0 }}
          </div>
          <div class="text-xs text-dark-400 mt-1">离开</div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  data: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const circumference = 2 * Math.PI * 45

// 假设群人数上限为1000，计算进度
const progressOffset = computed(() => {
  const members = props.data?.current_members || 0
  const max = 1000
  const progress = Math.min(members / max, 1)
  return circumference * (1 - progress)
})

const netChange = computed(() => {
  const joined = props.data?.today_joined || 0
  const left = props.data?.today_left || 0
  return joined - left
})
</script>
