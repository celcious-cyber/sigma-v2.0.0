<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  LogOut, Settings, ArrowRight, DatabaseZap, School, Ambulance, BookOpen, Wallet, ShieldCheck, NotepadText
} from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const user = ref({
  name: 'User',
  role: 'Guest'
})

// Module access mapping based on roles
const moduleMapping = {
  'Admin': ['Sigmabase', 'Sigmaedu', 'Sigmaflow', 'Sigmadesk', 'Sigmaguard', 'Sigmacare', 'Sigmalit'],
  'Teacher': ['Sigmabase', 'Sigmaedu'],
  'Finance': ['Sigmaflow'],
  'Musyrif': ['Sigmaguard'],
  'Staff': ['Sigmadesk', 'Sigmalit'],
  'Health': ['Sigmacare']
}

const panels = ref([
  { name: 'Sigmabase', icon: DatabaseZap, color: 'emerald', desc: 'Data Induk Santri & Wali', url: '/sigmabase' },
  { name: 'Sigmaedu', icon: School, color: 'blue', desc: 'Kurikulum & Nilai Akademik', url: '/sigmaedu' },
  { name: 'Sigmaflow', icon: Wallet, color: 'indigo', desc: 'Manajemen Keuangan & SPP', url: '/sigmaflow' },
  { name: 'Sigmadesk', icon: NotepadText, color: 'purple', desc: 'Tamu & Surat Menyurat', url: '#' },
  { name: 'Sigmaguard', icon: ShieldCheck, color: 'rose', desc: 'Perizinan & Keamanan', url: '/sigmaguard/violations' },
  { name: 'Sigmacare', icon: Ambulance, color: 'rose', desc: 'Kesehatan & Kebugaran', url: '#' },
  { name: 'Sigmalit', icon: BookOpen, color: 'rose', desc: 'Manajemen Perpustakaan', url: '#' },
])

const filteredPanels = ref([...panels.value])

onMounted(() => {
  const savedUser = localStorage.getItem('sigma_user')
  if (savedUser) {
    user.value = JSON.parse(savedUser)
    
    // Filter panels based on role
    if (user.value.role !== 'Admin') {
      const allowedModules = moduleMapping[user.value.role as keyof typeof moduleMapping] || []
      filteredPanels.value = panels.value.filter(panel => allowedModules.includes(panel.name))
    }
  }
})

const handleLogout = () => {
  localStorage.removeItem('sigma_token')
  localStorage.removeItem('sigma_user')
  router.push('/login')
}
</script>

<template>
  <div class="bg-slate-950 text-white min-h-screen relative overflow-x-hidden selection:bg-emerald-500/30">
    
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden -z-10">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-emerald-500/10 rounded-full blur-[120px]"></div>
      <div class="absolute top-[20%] -right-[10%] w-[30%] h-[30%] bg-blue-500/10 rounded-full blur-[100px]"></div>
      <div class="absolute -bottom-[10%] left-[20%] w-[35%] h-[35%] bg-indigo-500/10 rounded-full blur-[110px]"></div>
    </div>

    <header class="p-6 md:px-12 flex justify-between items-center bg-slate-950/50 backdrop-blur-md sticky top-0 z-50">
      <div class="flex items-center gap-3">
        <img src="/logo/SIGMA.svg" class="w-10 h-10 drop-shadow-[0_0_8px_rgba(16,185,129,0.3)]" alt="SIGMA Logo" />
        <h1 class="font-bold text-2xl tracking-tighter">SIGMA <span class="text-slate-500">PORTAL</span></h1>
      </div>

      <div class="flex items-center gap-4">
        <div class="hidden md:block text-right">
          <p class="text-sm font-semibold">{{ user.name }}</p>
          <p class="text-xs text-slate-400 uppercase tracking-widest">{{ user.role }}</p>
        </div>
        <button @click="handleLogout" 
                class="p-2.5 rounded-xl bg-red-500/10 text-red-400 border border-red-500/20 hover:bg-red-500 hover:text-white transition-all duration-300 flex items-center gap-2 group">
          <LogOut class="w-5 h-5" />
          <span class="hidden md:inline font-medium">Keluar</span>
        </button>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-6 pt-12 pb-8">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-6xl font-extrabold mb-4 bg-clip-text text-transparent bg-gradient-to-r from-white to-slate-500">
          Pusat Kendali Sistem
        </h2>
        <p class="text-slate-400 text-lg">Silakan pilih modul akses sesuai hak akses Anda.</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <div v-for="panel in filteredPanels" :key="panel.name" 
             @click="panel.url !== '#' && router.push(panel.url)"
             class="glass-card group relative p-6 rounded-[2.5rem] flex flex-col items-start transition-all duration-500 hover:-translate-y-2 cursor-pointer">
          
          <div :class="`w-16 h-16 rounded-2xl flex items-center justify-center mb-6 border transition-all duration-300 group-hover:scale-110
                        bg-${panel.color}-500/10 border-${panel.color}-500/20 text-${panel.color}-400 group-hover:bg-${panel.color}-500 group-hover:text-white`">
            <component :is="panel.icon" class="w-8 h-8" />
          </div>

          <h3 class="text-2xl font-bold mb-2">{{ panel.name }}</h3>
          <p class="text-slate-400 text-sm leading-relaxed mb-8">{{ panel.desc }}</p>

          <div :class="`mt-auto flex items-center gap-2 font-semibold transition-all duration-300 group-hover:gap-4 text-${panel.color}-400`">
            Akses Modul <ArrowRight class="w-5 h-5" />
          </div>
        </div>

        <div v-if="user.role === 'Administrator' || user.role === 'Admin'"
             class="lg:col-span-2 bg-gradient-to-br from-slate-800/50 to-slate-900/50 border border-slate-700/50 p-6 rounded-[2.5rem] flex flex-col group hover:border-amber-500/50 transition-all duration-500">
          <div class="w-16 h-16 bg-amber-500/10 text-amber-500 rounded-2xl flex items-center justify-center mb-6 border border-amber-500/20 group-hover:bg-amber-500 group-hover:text-white transition-all">
            <Settings class="w-8 h-8" />
          </div>
          <h3 class="text-2xl font-bold mb-2">System Setting</h3>
          <p class="text-slate-400 text-sm mb-8">Konfigurasi hak akses dan manajemen user sistem.</p>
          <div class="mt-auto flex items-center gap-2 text-amber-500 font-semibold group-hover:gap-4 transition-all">
            Administrator Only <ArrowRight class="w-5 h-5" />
          </div>
        </div>
      </div>
    </main>

    <footer class="p-12 text-center text-slate-600 text-xs tracking-widest uppercase">
      &copy; 2026 SIGMA - Sistem Integrasi Global Manajemen Al-Hikmah
    </footer>
  </div>
</template>

<style scoped>
.glass-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
.glass-card:hover {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>