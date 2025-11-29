<template>
  <div 
    class="card p-6 relative overflow-hidden group cursor-pointer"
    :class="colorClasses"
  >
    <!-- 背景装饰 -->
    <div 
      class="absolute -right-4 -top-4 text-6xl opacity-20 group-hover:opacity-30 transition-opacity"
    >
      {{ icon }}
    </div>

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="skeleton h-4 w-16 rounded mb-2"></div>
      <div class="skeleton h-10 w-24 rounded"></div>
    </template>

    <!-- 内容 -->
    <template v-else>
      <div class="text-dark-400 text-sm uppercase tracking-wider mb-2">
        {{ title }}
      </div>
      <div class="text-4xl font-bold number-display" :class="valueColorClass">
        {{ formattedValue }}
      </div>
    </template>

    <!-- 悬停动画边框 -->
    <div 
      class="absolute inset-0 border-2 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"
      :class="borderColorClass"
    ></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: Number,
    default: 0
  },
  icon: {
    type: String,
    default: '📊'
  },
  color: {
    type: String,
    default: 'blue' // blue, green, red, purple
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const colorClasses = computed(() => {
  const colors = {
    blue: 'hover:border-blue-500/50',
    green: 'hover:border-green-500/50',
    red: 'hover:border-red-500/50',
    purple: 'hover:border-purple-500/50'
  }
  return colors[props.color] || colors.blue
})

const valueColorClass = computed(() => {
  const colors = {
    blue: 'text-blue-400',
    green: 'text-green-400',
    red: 'text-red-400',
    purple: 'text-purple-400'
  }
  return colors[props.color] || colors.blue
})

const borderColorClass = computed(() => {
  const colors = {
    blue: 'border-blue-500/30',
    green: 'border-green-500/30',
    red: 'border-red-500/30',
    purple: 'border-purple-500/30'
  }
  return colors[props.color] || colors.blue
})

const formattedValue = computed(() => {
  if (props.color === 'green') return '+' + props.value
  if (props.color === 'red') return '-' + props.value
  return props.value
})
</script>
