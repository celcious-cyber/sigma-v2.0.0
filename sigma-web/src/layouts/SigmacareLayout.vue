<script setup lang="ts">
import { onMounted } from 'vue'
import SigmacareSidebar from '../components/SigmacareSidebar.vue'
import GlobalNavbar from '../components/GlobalNavbar.vue'
import { useRoute } from 'vue-router'

const route = useRoute()

onMounted(() => {
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
})
</script>

<template>
  <div class="flex h-screen bg-[#FFF1F2] text-slate-800 overflow-hidden font-sans transition-colors duration-300">
    <SigmacareSidebar />
    
    <main class="flex-1 overflow-y-auto custom-scrollbar flex flex-col">
      <GlobalNavbar moduleName="Sigma Care" />
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
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(225, 29, 72, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(225, 29, 72, 0.2); }
</style>
