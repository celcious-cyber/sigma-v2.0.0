<template>
  <div class="flex h-screen bg-sigma-dark text-slate-100 overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-72 bg-white/5 border-r border-white/10 flex flex-col z-20 overflow-y-auto">
      <div class="p-8">
        <h1 class="text-2xl font-black premium-gradient-text tracking-tighter">SIGMA v2.0</h1>
        <p class="text-[10px] text-slate-500 uppercase tracking-[0.2em] mt-1 font-bold">Modular Ecosystem</p>
      </div>

      <nav class="flex-1 px-4 space-y-2">
        <div class="pb-4">
          <p class="px-4 text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-3">Core Modules</p>
          <router-link to="/dashboard" class="nav-item">
            <LayoutDashboard :size="20" />
            <span>Beranda</span>
          </router-link>
          <router-link to="/base/students" class="nav-item">
            <Users :size="20" />
            <span>Sigmabase</span>
          </router-link>
          <router-link to="/flow/invoices" class="nav-item">
            <CreditCard :size="20" />
            <span>Sigmaflow</span>
          </router-link>
        </div>

        <div class="pb-4">
          <p class="px-4 text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-3">Guard & Care</p>
          <router-link to="/guard/violations" class="nav-item">
            <ShieldAlert :size="20" />
            <span>Sigmaguard</span>
          </router-link>
          <router-link to="/edu/assessments" class="nav-item">
            <GraduationCap :size="20" />
            <span>Sigmaedu</span>
          </router-link>
        </div>

        <div class="pb-4">
          <p class="px-4 text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-3">System</p>
          <router-link to="/settings/global" class="nav-item">
            <Settings :size="20" />
            <span>Pengaturan</span>
          </router-link>
        </div>
      </nav>

      <div class="p-6 border-t border-white/10 bg-black/20">
        <div class="flex items-center space-x-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-gradient-to-tr from-indigo-500 to-sky-500 flex items-center justify-center font-bold">
            A
          </div>
          <div class="overflow-hidden">
            <p class="text-sm font-bold truncate">Administrator</p>
            <p class="text-[10px] text-slate-500 truncate">admin@sigma.com</p>
          </div>
        </div>
        <button @click="logout" class="w-full text-xs font-bold py-2 rounded-lg bg-white/5 hover:bg-red-500/10 hover:text-red-400 transition-all border border-white/5">
          LOGOUT
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col relative overflow-hidden">
      <!-- Topbar -->
      <header class="h-16 border-b border-white/10 flex items-center justify-between px-8 bg-sigma-dark/50 backdrop-blur-md z-10">
        <h2 class="text-sm font-bold text-slate-400 uppercase tracking-widest">{{ routeName }}</h2>
        <div class="flex items-center space-x-4">
          <Bell :size="18" class="text-slate-500 cursor-pointer" />
          <div class="h-4 w-[1px] bg-white/10"></div>
          <span class="text-xs font-bold text-slate-500">v2.0.0 Stable</span>
        </div>
      </header>

      <!-- Content Area -->
      <section class="flex-1 overflow-y-auto p-8 relative">
        <router-view />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  LayoutDashboard, Users, CreditCard, ShieldAlert, GraduationCap, Settings, Bell 
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const routeName = computed(() => route.name?.toString() || 'Dashboard')

const logout = () => {
  localStorage.removeItem('sigma_token')
  router.push('/login')
}
</script>

<style scoped>
@reference "tailwindcss";

.nav-item {
  @apply flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 text-slate-400 hover:text-white hover:bg-white/5 font-medium mb-1;
}

.router-link-active.nav-item {
  @apply bg-indigo-600 text-white shadow-lg shadow-indigo-600/20;
}
</style>
