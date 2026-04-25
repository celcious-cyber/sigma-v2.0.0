<script setup lang="ts">
import { onMounted, ref } from 'vue'
import SigmaeduSidebar from '../components/SigmaeduSidebar.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isDark = ref(true)

// Centralized theme logic to prevent flashing
onMounted(() => {
  isDark.value = false
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
})
</script>

<template>
  <div class="flex h-screen bg-sigma-app text-sigma-text overflow-hidden font-sans transition-colors duration-300">
    <SigmaeduSidebar />
    
    <!-- Using a key on router-view can help with transitions, but for the layout itself, 
         not including the sidebar inside the router-view is what prevents the 'bounce' -->
    <main class="flex-1 overflow-y-auto custom-scrollbar">
      <router-view v-slot="{ Component }">
        <transition name="fade-slide" mode="out-in">
          <component :is="Component" :key="route.path" />
        </transition>
      </router-view>
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
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }
</style>
