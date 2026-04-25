<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { 
  LayoutDashboard, BookOpen, Calendar, 
  UserCheck, FileText, Heart, Briefcase,
  ChevronRight, Users,
  PanelLeftClose, PanelLeftOpen, Building2, Clock, GraduationCap, Printer
} from 'lucide-vue-next'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const isCollapsed = ref(false)
const openMenus = ref<string[]>([])

const toggleMenu = (name: string) => {
  if (openMenus.value.includes(name)) {
    openMenus.value = openMenus.value.filter(m => m !== name)
  } else {
    openMenus.value.push(name)
  }
}

const menuItems = ref([
  { name: 'Dashboard Edu', icon: LayoutDashboard, url: '/sigmaedu' },
  { 
    name: 'Sistem Kurikulum', 
    icon: BookOpen, 
    children: [
      { name: 'Mata Pelajaran', icon: FileText, url: '/sigmaedu/subjects' },
      { name: 'Jam Pelajaran', icon: Clock, url: '/sigmaedu/hours' },
      { name: 'Kalender Akademik', icon: Calendar, url: '/sigmaedu/calendar' },
    ]
  },
  {
    name: 'Ruang Kelas',
    icon: Building2,
    children: [
      { name: 'Data Kelas', icon: Building2, url: '/sigmaedu/classes' },
      { name: 'Data Santri', icon: Users, url: '/sigmaedu/students' },
      { name: 'Jadwal Pelajaran', icon: Calendar, url: '/sigmaedu/schedules' },
      { name: 'Presensi Santri', icon: UserCheck, url: '/sigmaedu/attendance' },
    ]
  },
  {
    name: 'Laporan & Penilaian',
    icon: GraduationCap,
    children: [
      { name: 'Nilai Pelajaran', icon: GraduationCap, url: '/sigmaedu/grades' },
      { name: 'Nilai Karakter', icon: Heart, url: '/sigmaedu/attitude' },
      { name: 'Hafalan Quran', icon: BookOpen, url: '/sigmaedu/tahfidz' },
      { name: 'Hafalan Pelajaran', icon: FileText, url: '/sigmaedu/memorization' },
    ]
  },
  {
    name: 'Manajemen Guru',
    icon: Briefcase,
    children: [
      { name: 'Presensi Guru', icon: UserCheck, url: '/sigmaedu/teacher-attendance' },
      { name: 'Jurnal Mengajar', icon: BookOpen, url: '/sigmaedu/journal' },
    ]
  },
  {
    name: 'Cetak Rapot',
    icon: Printer,
    children: [
      { name: 'Rapot KMI', icon: FileText, url: '/sigmaedu/reports/kmi' },
      { name: 'Rapot Tahfidz', icon: BookOpen, url: '/sigmaedu/reports/tahfidz' },
      { name: 'Rapot Ekskul', icon: Heart, url: '/sigmaedu/reports/ekskul' },
    ]
  }
])

const isDark = ref(true)

// Computed active item based on route
const currentActiveItem = computed(() => {
  const mainItem = menuItems.value.find(item => item.url === route.path)
  if (mainItem) return mainItem.name

  for (const item of menuItems.value) {
    if (item.children) {
      const subItem = item.children.find(sub => sub.url === route.path)
      if (subItem) return subItem.name
    }
  }
  return ''
})

// Computed active parent for opening menus
const currentActiveParent = computed(() => {
  for (const item of menuItems.value) {
    if (item.children) {
      const isActive = item.children.some(sub => sub.url === route.path)
      if (isActive) return item.name
    }
  }
  return ''
})

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

onMounted(() => {
  if (currentActiveParent.value && !openMenus.value.includes(currentActiveParent.value)) {
    openMenus.value.push(currentActiveParent.value)
  }

  isDark.value = false
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
})
</script>

<template>
  <aside 
    class="border-r border-sigma-border bg-sigma-surface/50 backdrop-blur-xl flex flex-col hidden lg:flex h-screen sticky top-0 transition-all duration-500 ease-in-out z-[60] group/sidebar"
    :class="isCollapsed ? 'w-24' : 'w-72'"
  >
    <!-- Toggle Button -->
    <button 
      @click="toggleSidebar"
      class="absolute -right-3 top-10 w-6 h-6 bg-blue-500 text-white rounded-full flex items-center justify-center shadow-lg shadow-blue-500/40 hover:scale-110 transition-all z-50 border border-white/10"
    >
      <PanelLeftClose v-if="!isCollapsed" class="w-3.5 h-3.5" />
      <PanelLeftOpen v-else class="w-3.5 h-3.5" />
    </button>

    <div class="p-6 flex-1 flex flex-col overflow-hidden">
      <!-- Logo Section -->
      <div class="flex items-center gap-3 mb-10 overflow-hidden" :class="isCollapsed ? 'justify-center' : ''">
        <img src="/logo/SIGMA.svg" 
             class="w-10 h-10 min-w-[2.5rem] drop-shadow-[0_0_8px_rgba(59,130,246,0.3)] cursor-pointer hover:scale-110 transition-transform filter hue-rotate-[180deg]" 
             alt="SIGMA Logo" 
             @click="router.push('/portal')" />
        <h1 v-if="!isCollapsed" class="font-black text-xl tracking-tighter uppercase italic animate-in fade-in slide-in-from-left-4 duration-500">Sigmaedu</h1>
      </div>

      <!-- Navigation -->
      <nav class="space-y-2 flex-1 overflow-y-auto custom-scrollbar pr-1">
        <div v-for="item in menuItems" :key="item.name">
          <button 
            @click="item.children ? toggleMenu(item.name) : router.push(item.url)"
            class="w-full flex items-center p-3 rounded-xl transition-all group"
            :class="[
              currentActiveItem === item.name || currentActiveParent === item.name ? 'bg-blue-500 text-white shadow-lg shadow-blue-500/20' : 'hover:bg-sigma-surface-alt text-sigma-muted hover:text-sigma-text',
              isCollapsed ? 'justify-center' : 'justify-between'
            ]"
          >
            <div class="flex items-center gap-3">
              <component :is="item.icon" class="w-5 h-5 min-w-[1.25rem]" />
              <span v-if="!isCollapsed" class="font-bold text-sm tracking-tight animate-in fade-in slide-in-from-left-2 duration-300">{{ item.name }}</span>
            </div>
            <ChevronRight 
              v-if="item.children && !isCollapsed" 
              class="w-4 h-4 text-slate-600 group-hover:text-blue-400 transition-all duration-300" 
              :class="openMenus.includes(item.name) ? 'rotate-90' : ''"
            />
          </button>

          <!-- Submenu -->
          <div v-if="item.children && !isCollapsed && openMenus.includes(item.name)" class="mt-2 ml-4 pl-4 border-l border-sigma-border space-y-1 animate-in fade-in slide-in-from-top-2 duration-300">
            <button 
              v-for="sub in item.children" 
              :key="sub.name"
              @click="sub.url !== '#' && router.push(sub.url)"
              class="w-full flex items-center gap-3 p-2.5 rounded-lg text-xs font-bold transition-all text-left uppercase tracking-widest"
              :class="currentActiveItem === sub.name ? 'text-blue-400 bg-blue-500/5' : 'text-sigma-muted hover:text-blue-400 hover:bg-blue-500/5'"
            >
              <component :is="sub.icon" class="w-3.5 h-3.5" />
              {{ sub.name }}
            </button>
          </div>
        </div>
      </nav>

    </div>
  </aside>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 3px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }
</style>
