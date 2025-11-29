<template>
  <div class="flex gap-12 px-6">
    <div
      v-for="(stat, index) in stats"
      :key="index"
      class="flex items-center gap-4 py-4"
    >
      <span class="text-4xl lg:text-5xl font-bold number-display gradient-text">
        {{ formatNumber(stat.value) }}
      </span>
      <span class="text-dark-400 text-sm uppercase tracking-wider">
        {{ stat.label }}
      </span>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  stats: {
    type: Array,
    default: () => []
  }
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
</script>
