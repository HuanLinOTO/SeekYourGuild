<template>
  <div class="card p-6 min-h-[380px]">
    <h3 class="text-dark-400 text-sm uppercase tracking-wider mb-6">
      入退群趋势（近7天）
    </h3>

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="h-40 flex items-center justify-center mb-4">
        <div class="flex items-end gap-6 w-full px-4">
          <div 
            v-for="i in 7" 
            :key="i" 
            class="skeleton flex-1 rounded-t"
            :style="{ height: `${30 + (i % 3) * 25}%` }"
          ></div>
        </div>
      </div>
      <div class="flex justify-between px-1 mb-4">
        <div class="skeleton h-4 w-10 rounded"></div>
        <div class="skeleton h-4 w-10 rounded"></div>
        <div class="skeleton h-4 w-10 rounded"></div>
      </div>
      <div class="border-t border-dark-700 pt-4 flex justify-center gap-6 mb-4">
        <div class="flex items-center gap-2">
          <div class="skeleton w-3 h-3 rounded-full"></div>
          <div class="skeleton h-3 w-8 rounded"></div>
        </div>
        <div class="flex items-center gap-2">
          <div class="skeleton w-3 h-3 rounded-full"></div>
          <div class="skeleton h-3 w-8 rounded"></div>
        </div>
      </div>
      <div class="grid grid-cols-3 gap-4 text-center">
        <div>
          <div class="skeleton h-8 w-12 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-10 mx-auto rounded"></div>
        </div>
        <div>
          <div class="skeleton h-8 w-12 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-10 mx-auto rounded"></div>
        </div>
        <div>
          <div class="skeleton h-8 w-12 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-10 mx-auto rounded"></div>
        </div>
      </div>
    </template>

    <!-- 图表 -->
    <template v-else>
      <div class="relative">
        <!-- 折线图区域 -->
        <div class="h-40 relative mb-4">
          <svg class="w-full h-full" :viewBox="`0 0 ${svgWidth} ${svgHeight}`" preserveAspectRatio="none">
            <!-- 网格线 -->
            <line 
              v-for="i in 4" 
              :key="`grid-${i}`"
              :x1="0" 
              :y1="svgHeight * i / 4" 
              :x2="svgWidth" 
              :y2="svgHeight * i / 4"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            
            <!-- 入群折线 - 渐变填充 -->
            <defs>
              <linearGradient id="joinGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" stop-color="rgba(34, 197, 94, 0.3)" />
                <stop offset="100%" stop-color="rgba(34, 197, 94, 0)" />
              </linearGradient>
              <linearGradient id="leftGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" stop-color="rgba(239, 68, 68, 0.3)" />
                <stop offset="100%" stop-color="rgba(239, 68, 68, 0)" />
              </linearGradient>
            </defs>
            
            <!-- 入群区域填充 -->
            <path 
              :d="joinAreaPath" 
              fill="url(#joinGradient)"
            />
            
            <!-- 退群区域填充 -->
            <path 
              :d="leftAreaPath" 
              fill="url(#leftGradient)"
            />
            
            <!-- 入群折线 -->
            <path 
              :d="joinLinePath" 
              fill="none" 
              stroke="#22c55e" 
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            
            <!-- 退群折线 -->
            <path 
              :d="leftLinePath" 
              fill="none" 
              stroke="#ef4444" 
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            
            <!-- 入群数据点 -->
            <circle 
              v-for="(point, index) in joinPoints" 
              :key="`join-${index}`"
              :cx="point.x" 
              :cy="point.y" 
              r="4"
              fill="#22c55e"
              class="transition-all duration-200"
              :class="{ 'opacity-100': hoveredIndex === index, 'opacity-0': hoveredIndex !== index && hoveredIndex !== null }"
            />
            
            <!-- 退群数据点 -->
            <circle 
              v-for="(point, index) in leftPoints" 
              :key="`left-${index}`"
              :cx="point.x" 
              :cy="point.y" 
              r="4"
              fill="#ef4444"
              class="transition-all duration-200"
              :class="{ 'opacity-100': hoveredIndex === index, 'opacity-0': hoveredIndex !== index && hoveredIndex !== null }"
            />
          </svg>
          
          <!-- 悬停区域 -->
          <div class="absolute inset-0 flex">
            <div 
              v-for="(item, index) in chartData" 
              :key="`hover-${index}`"
              class="flex-1 cursor-pointer"
              @mouseenter="hoveredIndex = index"
              @mouseleave="hoveredIndex = null"
            >
              <!-- 提示框 -->
              <div 
                v-if="hoveredIndex === index"
                class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-dark-800 rounded-lg text-xs whitespace-nowrap z-10 shadow-lg border border-dark-700"
                :style="{ left: `${(index + 0.5) / chartData.length * 100}%` }"
              >
                <div class="font-medium text-white mb-1">{{ item.date }}</div>
                <div class="flex gap-3">
                  <span class="text-green-400">+{{ item.joined }}</span>
                  <span class="text-red-400">-{{ item.left }}</span>
                  <span :class="item.net >= 0 ? 'text-green-400' : 'text-red-400'">
                    净{{ item.net >= 0 ? '+' : '' }}{{ item.net }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- X轴标签 -->
        <div class="flex justify-between text-xs text-dark-500 px-1">
          <span v-for="(item, index) in chartData" :key="`label-${index}`" class="text-center" style="width: 14.28%;">
            {{ formatDateLabel(item.date, index) }}
          </span>
        </div>
      </div>

      <!-- 图例 -->
      <div class="mt-4 pt-4 border-t border-dark-700 flex justify-center gap-6">
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 rounded-full bg-green-500"></div>
          <span class="text-xs text-dark-400">入群</span>
        </div>
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 rounded-full bg-red-500"></div>
          <span class="text-xs text-dark-400">退群</span>
        </div>
      </div>

      <!-- 统计摘要 -->
      <div class="mt-4 grid grid-cols-3 gap-4 text-center">
        <div>
          <div class="text-2xl font-bold text-green-400 number-display">
            +{{ totalJoined }}
          </div>
          <div class="text-xs text-dark-400 mt-1">总入群</div>
        </div>
        <div>
          <div class="text-2xl font-bold text-red-400 number-display">
            -{{ totalLeft }}
          </div>
          <div class="text-xs text-dark-400 mt-1">总退群</div>
        </div>
        <div>
          <div class="text-2xl font-bold number-display" :class="netChange >= 0 ? 'text-green-400' : 'text-red-400'">
            {{ netChange >= 0 ? '+' : '' }}{{ netChange }}
          </div>
          <div class="text-xs text-dark-400 mt-1">净增长</div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const hoveredIndex = ref(null)
const svgWidth = 280
const svgHeight = 160
const padding = 10

// 处理数据，确保有7天的数据，按日期正序排列
const chartData = computed(() => {
  const result = []
  const today = new Date()
  
  // 生成最近7天的日期
  for (let i = 6; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    const dateStr = date.toISOString().split('T')[0]
    
    const found = props.data.find(d => d.date === dateStr)
    result.push({
      date: dateStr,
      joined: found?.joined || 0,
      left: found?.left || 0,
      net: (found?.joined || 0) - (found?.left || 0)
    })
  }
  
  return result
})

// 计算最大值用于缩放
const maxValue = computed(() => {
  let max = 1
  chartData.value.forEach(item => {
    max = Math.max(max, item.joined, item.left)
  })
  return max
})

// 计算点的位置
function getPoints(values) {
  const points = []
  const dataLen = values.length
  const xStep = (svgWidth - padding * 2) / (dataLen - 1)
  
  values.forEach((val, index) => {
    const x = padding + index * xStep
    const y = svgHeight - padding - (val / maxValue.value) * (svgHeight - padding * 2)
    points.push({ x, y })
  })
  
  return points
}

const joinPoints = computed(() => getPoints(chartData.value.map(d => d.joined)))
const leftPoints = computed(() => getPoints(chartData.value.map(d => d.left)))

// 生成折线路径
function getLinePath(points) {
  if (points.length === 0) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ')
}

// 生成填充区域路径
function getAreaPath(points) {
  if (points.length === 0) return ''
  const linePath = getLinePath(points)
  const lastPoint = points[points.length - 1]
  const firstPoint = points[0]
  return `${linePath} L ${lastPoint.x} ${svgHeight - padding} L ${firstPoint.x} ${svgHeight - padding} Z`
}

const joinLinePath = computed(() => getLinePath(joinPoints.value))
const leftLinePath = computed(() => getLinePath(leftPoints.value))
const joinAreaPath = computed(() => getAreaPath(joinPoints.value))
const leftAreaPath = computed(() => getAreaPath(leftPoints.value))

// 格式化日期标签
function formatDateLabel(dateStr, index) {
  const date = new Date(dateStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  
  // 只显示第一个、中间和最后一个日期，其他用点表示
  if (index === 0 || index === 3 || index === 6) {
    return `${month}/${day}`
  }
  return ''
}

// 统计数据
const totalJoined = computed(() => chartData.value.reduce((sum, d) => sum + d.joined, 0))
const totalLeft = computed(() => chartData.value.reduce((sum, d) => sum + d.left, 0))
const netChange = computed(() => totalJoined.value - totalLeft.value)
</script>
