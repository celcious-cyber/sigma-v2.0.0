<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { 
  LayoutDashboard, Activity, Stethoscope, ClipboardList,
  UserCheck, Thermometer, FlaskConical, History,
  Database, Calendar, Bed, Pill, PackagePlus,
  ChevronRight, PanelLeftClose, PanelLeftOpen,
  FileText, Microscope
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
  { name: 'Dasbor Klinis', icon: LayoutDashboard, url: '/sigmacare' },
  { 
    name: 'Poliklinik & Observasi', 
    icon: Stethoscope, 
    children: [
      { name: 'Anamnesis & Diagnosis', icon: ClipboardList, url: '/sigmacare/visits' },
      { name: 'Observasi Klinis', icon: Bed, url: '/sigmacare/observations' },
      { name: 'Sertifikasi Medis', icon: FileText, url: '/sigmacare/certificates' },
    ]
  },
  {
    name: 'Medical Check-Up',
    icon: UserCheck,
    children: [
      { name: 'Prosedur Antropometri', icon: Activity, url: '/sigmacare/mcu/physical' },
      { name: 'Profil Serologi', icon: Microscope, url: '/sigmacare/mcu/biometric' },
      { name: 'Evaluasi Pertumbuhan', icon: Database, url: '/sigmacare/mcu/reports' },
    ]
  },
  {
    name: 'Rekam Medis Elektronik',
    icon: History,
    children: [
      { name: 'Riwayat Terintegrasi', icon: History, url: '/sigmacare/history' },
      { name: 'Katalog Patologi', icon: Database, url: '/sigmacare/icd10' },
    ]
  },
  {
    name: 'Instalasi Farmasi',
    icon: Pill,
    children: [
      { name: 'Inventaris Farmakologi', icon: Pill, url: '/sigmacare/pharmacy/inventory' },
      { name: 'Mutasi Perbekalan', icon: PackagePlus, url: '/sigmacare/pharmacy/mutations' },
    ]
  },
  {
    name: 'Administrasi Faskes',
    icon: Calendar,
    children: [
      { name: 'Kalender MCU', icon: Calendar, url: '/sigmacare/admin/calendar' },
      { name: 'Manajemen Fasilitas', icon: Bed, url: '/sigmacare/admin/facilities' },
    ]
  }
])

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
      class="absolute -right-3 top-10 w-6 h-6 bg-rose-600 text-white rounded-full flex items-center justify-center shadow-lg shadow-rose-500/40 hover:scale-110 transition-all z-50 border border-white/10"
    >
      <PanelLeftClose v-if="!isCollapsed" class="w-3.5 h-3.5" />
      <PanelLeftOpen v-else class="w-3.5 h-3.5" />
    </button>

    <div class="p-6 flex-1 flex flex-col overflow-hidden">
      <!-- Logo Section -->
      <div class="flex items-center gap-3 mb-10 overflow-hidden" :class="isCollapsed ? 'justify-center' : ''">
        <img src="/logo/SIGMA.svg" 
             class="w-10 h-10 min-w-[2.5rem] drop-shadow-[0_0_8px_rgba(225,29,72,0.3)] cursor-pointer hover:scale-110 transition-transform filter hue-rotate-[320deg]" 
             alt="SIGMA Logo" 
             @click="router.push('/portal')" />
        <h1 v-if="!isCollapsed" class="font-black text-xl tracking-tighter uppercase italic animate-in fade-in slide-in-from-left-4 duration-500 text-slate-800">Sigma<span class="text-rose-600">care</span></h1>
      </div>

      <!-- Navigation -->
      <nav class="space-y-2 flex-1 overflow-y-auto custom-scrollbar pr-1">
        <div v-for="item in menuItems" :key="item.name">
          <button 
            @click="item.children ? toggleMenu(item.name) : (item.url !== '#' && router.push(item.url))"
            class="w-full flex items-center p-3 rounded-xl transition-all group"
            :class="[
              currentActiveItem === item.name || currentActiveParent === item.name ? 'bg-rose-600 text-white shadow-lg shadow-rose-500/20' : 'hover:bg-sigma-surface-alt text-sigma-muted hover:text-sigma-text',
              isCollapsed ? 'justify-center' : 'justify-between'
            ]"
          >
            <div class="flex items-center gap-3">
              <component :is="item.icon" class="w-5 h-5 min-w-[1.25rem]" />
              <span v-if="!isCollapsed" class="font-bold text-sm tracking-tight animate-in fade-in slide-in-from-left-2 duration-300">{{ item.name }}</span>
            </div>
            <ChevronRight 
              v-if="item.children && !isCollapsed" 
              class="w-4 h-4 text-slate-600 group-hover:text-white transition-all duration-300" 
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
              :class="currentActiveItem === sub.name ? 'text-rose-600 bg-rose-600/5' : 'text-sigma-muted hover:text-rose-600 hover:bg-rose-600/5'"
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
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(225, 29, 72, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(225, 29, 72, 0.2); }
</style>
