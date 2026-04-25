<script setup lang="ts">
import { onMounted, ref } from 'vue'
import SigmaflowSidebar from '../components/SigmaflowSidebar.vue'
import GlobalNavbar from '../components/GlobalNavbar.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isDark = ref(true)

onMounted(() => {
  isDark.value = false
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
})
</script>

<template>
  <div class="flex h-screen bg-[#E0E7FF] text-sigma-text overflow-hidden font-sans transition-colors duration-300">
    <SigmaflowSidebar />
    
    <main class="flex-1 overflow-y-auto custom-scrollbar flex flex-col">
      <GlobalNavbar moduleName="SigmaFlow" />
      <div class="flex-1">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" :key="route.path" />
          </transition>
        </router-view>
      </div>
    </main>
  </div>
</template>

<style scoped>
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(79, 70, 229, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(79, 70, 229, 0.2); }
</style>
