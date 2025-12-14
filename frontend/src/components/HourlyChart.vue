<template>
  <div class="card p-6 min-h-[340px]">
    <h3 class="text-dark-400 text-sm uppercase tracking-wider mb-6">
      24小时消息分布
    </h3>

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="flex items-end gap-1 h-40 mb-4">
        <div 
          v-for="i in 24" 
          :key="i" 
          class="skeleton flex-1 rounded-t"
          :style="{ height: `${20 + (i % 5) * 15}%` }"
        ></div>
      </div>
      <div class="flex justify-between mb-6">
        <div class="skeleton h-4 w-12 rounded"></div>
        <div class="skeleton h-4 w-12 rounded"></div>
        <div class="skeleton h-4 w-12 rounded"></div>
        <div class="skeleton h-4 w-12 rounded"></div>
        <div class="skeleton h-4 w-16 rounded"></div>
      </div>
      <div class="border-t border-dark-700 pt-6 grid grid-cols-3 gap-4">
        <div class="text-center">
          <div class="skeleton h-8 w-16 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-12 mx-auto rounded"></div>
        </div>
        <div class="text-center">
          <div class="skeleton h-8 w-16 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-12 mx-auto rounded"></div>
        </div>
        <div class="text-center">
          <div class="skeleton h-8 w-16 mx-auto rounded mb-1"></div>
          <div class="skeleton h-3 w-12 mx-auto rounded"></div>
        </div>
      </div>
    </template>

    <!-- 图表 -->
    <template v-else>
      <div class="relative">
        <!-- 柱状图 -->
        <div class="flex items-end gap-1 h-40 mb-4">
          <div
            v-for="(item, index) in chartData"
            :key="index"
            class="flex-1 group cursor-pointer"
            @mouseenter="hoveredHour = index"
            @mouseleave="hoveredHour = null"
          >
            <div class="relative">
              <!-- 提示框 -->
              <div 
                v-if="hoveredHour === index"
                class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-dark-800 rounded text-xs whitespace-nowrap z-10"
              >
                {{ item.label }} - {{ item.count }}条
              </div>
              
              <!-- 柱子 -->
              <div 
                class="w-full rounded-t transition-all duration-300"
                :class="[
                  item.isNow 
                    ? 'bg-gradient-to-t from-green-500 to-emerald-400 ring-2 ring-green-400/50' 
                    : hoveredHour === index 
                      ? 'bg-gradient-to-t from-blue-600 to-purple-500' 
                      : 'bg-gradient-to-t from-blue-600/60 to-purple-500/60'
                ]"
                :style="{ 
                  height: `${getBarHeight(item.count)}px`,
                  minHeight: '4px'
                }"
              ></div>
            </div>
          </div>
        </div>

        <!-- X轴标签 -->
        <div class="flex justify-between text-xs text-dark-500">
          <span>{{ xLabels[0] }}</span>
          <span>{{ xLabels[1] }}</span>
          <span>{{ xLabels[2] }}</span>
          <span>{{ xLabels[3] }}</span>
          <span class="text-green-400 font-medium">{{ xLabels[4] }} ◀ 现在</span>
        </div>
      </div>

      <!-- 统计摘要 -->
      <div class="mt-6 pt-6 border-t border-dark-700 grid grid-cols-3 gap-4 text-center">
        <div>
          <div class="text-2xl font-bold text-blue-400 number-display">
            {{ totalMessages }}
          </div>
          <div class="text-xs text-dark-400 mt-1">总消息</div>
        </div>
        <div>
          <div class="text-2xl font-bold text-purple-400 number-display">
            {{ peakHour }}
          </div>
          <div class="text-xs text-dark-400 mt-1">高峰时段</div>
        </div>
        <div>
          <div class="text-2xl font-bold text-green-400 number-display">
            {{ avgPerHour }}
          </div>
          <div class="text-xs text-dark-400 mt-1">小时均值</div>
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

const hoveredHour = ref(null)

// 获取当前小时
const currentHour = computed(() => new Date().getHours())

// 生成过去24小时的数据（从左到右：24小时前 -> 现在）
const chartData = computed(() => {
  const now = currentHour.value
  const hours = []
  
  // 从24小时前到现在
  for (let i = 23; i >= 0; i--) {
    const hour = (now - i + 24) % 24
    const found = props.data.find(d => d.hour === hour)
    hours.push({
      hour,
      count: found?.count || 0,
      label: `${String(hour).padStart(2, '0')}:00`,
      isNow: i === 0
    })
  }
  
  return hours
})

// 生成X轴标签（显示5个时间点）
const xLabels = computed(() => {
  const now = currentHour.value
  return [
    `${String((now - 23 + 24) % 24).padStart(2, '0')}:00`,  // 24小时前
    `${String((now - 17 + 24) % 24).padStart(2, '0')}:00`,  // 18小时前
    `${String((now - 11 + 24) % 24).padStart(2, '0')}:00`,  // 12小时前
    `${String((now - 5 + 24) % 24).padStart(2, '0')}:00`,   // 6小时前
    `${String(now).padStart(2, '0')}:00`                     // 现在
  ]
})

const maxCount = computed(() => {
  return Math.max(...chartData.value.map(d => d.count), 1)
})

function getBarHeight(count) {
  const maxHeight = 160
  return (count / maxCount.value) * maxHeight
}

const totalMessages = computed(() => {
  return chartData.value.reduce((sum, item) => sum + item.count, 0)
})

const peakHour = computed(() => {
  let maxHour = 0
  let maxVal = 0
  chartData.value.forEach((item) => {
    if (item.count > maxVal) {
      maxVal = item.count
      maxHour = item.hour
    }
  })
  return `${String(maxHour).padStart(2, '0')}:00`
})

const avgPerHour = computed(() => {
  const total = totalMessages.value
  return Math.round(total / 24)
})
</script>
