<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="isOpen = !isOpen"
      class="flex items-center gap-2 px-4 py-2 rounded-xl bg-dark-800/80 border border-dark-700 hover:border-dark-600 transition-colors"
    >
      <span class="text-dark-400 text-sm">当前群:</span>
      <span class="text-white font-medium">
        {{ currentGroupName }}
      </span>
      <svg 
        class="w-4 h-4 text-dark-400 transition-transform"
        :class="{ 'rotate-180': isOpen }"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- 下拉菜单 -->
    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-1"
    >
      <div
        v-if="isOpen"
        class="absolute top-full left-0 mt-2 w-64 bg-dark-800 border border-dark-700 rounded-xl shadow-xl z-50 overflow-hidden"
      >
        <div class="p-2 border-b border-dark-700">
          <div class="text-xs text-dark-400 uppercase tracking-wider px-2">选择群</div>
        </div>
        
        <div class="max-h-64 overflow-y-auto">
          <button
            v-for="group in groups"
            :key="group.id"
            @click="selectGroup(group)"
            class="w-full flex items-center gap-3 px-4 py-3 hover:bg-dark-700/50 transition-colors text-left"
            :class="{ 'bg-blue-600/20': group.id === currentGroupId }"
          >
            <!-- 群图标 -->
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-bold">
              {{ getGroupInitial(group) }}
            </div>
            
            <div class="flex-1 min-w-0">
              <div class="text-white font-medium truncate">
                {{ group.name || `群 ${group.id}` }}
              </div>
              <div class="text-xs text-dark-400">
                {{ group.member_count || 0 }} 成员
              </div>
            </div>

            <!-- 选中标记 -->
            <svg
              v-if="group.id === currentGroupId"
              class="w-5 h-5 text-blue-400"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>

          <!-- 空状态 -->
          <div
            v-if="groups.length === 0"
            class="px-4 py-8 text-center text-dark-400"
          >
            <svg class="w-12 h-12 mx-auto mb-3 text-dark-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <div class="text-sm">暂无可用群</div>
            <div class="text-xs mt-1">请配置 ALLOWED_GROUPS 环境变量</div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  groups: {
    type: Array,
    default: () => []
  },
  currentGroupId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['select'])

const isOpen = ref(false)
const dropdownRef = ref(null)

const currentGroupName = computed(() => {
  if (!props.currentGroupId) return '请选择'
  const group = props.groups.find(g => g.id === props.currentGroupId)
  return group?.name || `群 ${props.currentGroupId}`
})

function getGroupInitial(group) {
  if (group.name) {
    return group.name.charAt(0).toUpperCase()
  }
  return String(group.id).charAt(0)
}

function selectGroup(group) {
  emit('select', group.id)
  isOpen.value = false
}

// 点击外部关闭
function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
