<template>
  <div class="min-h-screen relative">
    <!-- 动态光晕背景 -->
    <div class="glow-bg">
      <div class="glow-orb glow-orb-1"></div>
      <div class="glow-orb glow-orb-2"></div>
      <div class="glow-orb glow-orb-3"></div>
    </div>

    <!-- 浮动装饰元素 -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden z-0">
      <div class="circle-decoration w-96 h-96 -top-20 -right-20 floating"></div>
      <div class="circle-decoration w-64 h-64 top-1/3 -left-10 floating-delay-1"></div>
      <div class="circle-decoration w-48 h-48 bottom-20 right-1/4 floating-delay-2"></div>
    </div>

    <!-- 主内容 -->
    <div class="relative z-10">
      <router-view />
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
</style>
