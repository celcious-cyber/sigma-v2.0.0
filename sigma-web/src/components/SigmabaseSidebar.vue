<script setup lang="ts">
import { ref } from 'vue'
import { 
  Users, LayoutDashboard, GraduationCap, 
  Briefcase, BookOpen, ChevronRight, ArrowLeft,
  Sun, Moon
} from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

defineProps<{
  activeItem?: string
}>()

const router = useRouter()

const menuItems = ref([
  { name: 'Dashboard', icon: LayoutDashboard, url: '/sigmabase' },
  { 
    name: 'Data Induk', 
    icon: BookOpen, 
    children: [
      { name: 'Data Santri', icon: GraduationCap, url: '/sigmabase/students' },
      { name: 'Data Guru', icon: Briefcase, url: '/sigmabase/teachers' },
      { name: 'Data Alumni', icon: Users, url: '/sigmabase/alumni' },
    ]
  }
])

const isDark = ref(true)

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('sigma_theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('sigma_theme', 'light')
  }
}

onMounted(() => {
  const savedTheme = localStorage.getItem('sigma_theme')
  if (savedTheme === 'light') {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  } else {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
})
</script>

<template>
  <aside class="w-72 border-r border-sigma-border bg-sigma-surface/50 backdrop-blur-xl flex flex-col hidden lg:flex h-screen sticky top-0 transition-colors duration-300">
    <div class="p-8">
      <div class="flex items-center gap-3 mb-10">
        <img src="/logo/SIGMA.svg" 
             class="w-10 h-10 drop-shadow-[0_0_8px_rgba(16,185,129,0.3)] cursor-pointer hover:scale-110 transition-transform" 
             alt="SIGMA Logo" 
             @click="router.push('/portal')" />
        <h1 class="font-bold text-xl tracking-tighter uppercase">Sigmabase</h1>
      </div>

      <nav class="space-y-2">
        <div v-for="item in menuItems" :key="item.name">
          <button 
            @click="!item.children && router.push(item.url)"
            class="w-full flex items-center justify-between p-3 rounded-xl transition-all group"
            :class="activeItem === item.name ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'hover:bg-sigma-surface-alt text-sigma-muted hover:text-sigma-text'"
          >
            <div class="flex items-center gap-3">
              <component :is="item.icon" class="w-5 h-5" />
              <span class="font-semibold">{{ item.name }}</span>
            </div>
            <ChevronRight v-if="item.children" class="w-4 h-4 text-slate-600 group-hover:text-slate-400 transition-all" />
          </button>

          <!-- Submenu -->
          <div v-if="item.children" class="mt-2 ml-4 pl-4 border-l border-slate-800 space-y-1">
            <button 
              v-for="sub in item.children" 
              :key="sub.name"
              @click="sub.url !== '#' && router.push(sub.url)"
              class="w-full flex items-center gap-3 p-2.5 rounded-lg text-sm transition-all text-left"
              :class="activeItem === sub.name ? 'text-emerald-400 font-bold' : 'text-sigma-muted hover:text-emerald-400 hover:bg-emerald-500/5'"
            >
              <component :is="sub.icon" class="w-4 h-4" />
              {{ sub.name }}
            </button>
          </div>
        </div>
      </nav>
    </div>

    <div class="mt-auto p-6 space-y-4">
      <button @click="toggleTheme" 
              class="w-full flex items-center justify-between p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border text-sigma-text hover:border-sigma-emerald/30 transition-all group shadow-sm">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-sigma-emerald/10 flex items-center justify-center text-sigma-emerald">
            <Sun v-if="isDark" class="w-4 h-4" />
            <Moon v-else class="w-4 h-4" />
          </div>
          <span class="text-sm font-bold uppercase tracking-wider">{{ isDark ? 'Light Mode' : 'Dark Mode' }}</span>
        </div>
        <div class="w-10 h-6 bg-sigma-emerald/10 rounded-full relative p-1 flex items-center transition-all" :class="{'justify-end bg-sigma-emerald/20': !isDark}">
          <div class="w-4 h-4 bg-sigma-emerald rounded-full shadow-lg shadow-sigma-emerald/40 animate-in fade-in zoom-in duration-300"></div>
        </div>
      </button>

      <button @click="router.push('/portal')" 
              class="w-full flex items-center gap-3 p-4 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-sigma-border transition-all group">
        <ArrowLeft class="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
        <span class="font-medium">Kembali ke Portal</span>
      </button>
    </div>
  </aside>
</template>
