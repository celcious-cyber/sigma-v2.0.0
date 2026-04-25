<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  Activity, ArrowLeft, Search, Bell, Settings,
  User, LogOut, Info, Clock
} from 'lucide-vue-next'
import axios from 'axios'

const router = useRouter()
const route = useRoute()

const props = defineProps({
  moduleName: {
    type: String,
    default: 'SigmaEdu'
  }
})

const pageTitle = computed(() => {
  return (route.meta.title as string) || (route.name as string) || 'Dashboard'
})

const user = ref({ name: 'User', role: 'Loading...', avatar: '' })

const fetchUser = async () => {
  try {
    const res = await axios.get('/api/v1/auth/me')
    const data = res.data.data
    user.value = {
      name: data.name,
      role: res.data.type === 'admin' ? data.role : 'Santri',
      avatar: data.avatar_url || ''
    }
  } catch (err) {
    console.error('Failed to fetch user', err)
  }
}

// Dropdown States
const showNotifications = ref(false)
const showProfileMenu = ref(false)

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value
  showProfileMenu.value = false
}

const toggleProfileMenu = () => {
  showProfileMenu.value = !showProfileMenu.value
  showNotifications.value = false
}

const handleLogout = () => {
  localStorage.removeItem('sigma_token')
  router.push('/login')
}

// Close on click outside
const closeDropdowns = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.dropdown-container')) {
    showNotifications.value = false
    showProfileMenu.value = false
  }
}

onMounted(() => {
  window.addEventListener('click', closeDropdowns)
  fetchUser()
})
onUnmounted(() => window.removeEventListener('click', closeDropdowns))

const notifications = [
  { id: 1, title: 'Presensi Selesai', desc: 'Jurnal kelas X-A telah diverifikasi', time: '5m ago', icon: Activity, color: 'text-blue-500' },
  { id: 2, title: 'Agenda Baru', desc: 'Ujian Tengah Semester dijadwalkan', time: '1h ago', icon: Info, color: 'text-amber-500' },
  { id: 3, title: 'Sistem Update', desc: 'Versi 2.0.1 telah tersedia', time: '2h ago', icon: Settings, color: 'text-emerald-500' },
]
</script>

<template>
  <header class="sticky top-0 z-50 flex items-center justify-between py-4 px-8 bg-[#E0E7FF]/80 backdrop-blur-xl border-b border-white/20">
    <div class="flex items-center gap-4">
      <div class="w-10 h-10 bg-white rounded-xl shadow-sm flex items-center justify-center text-blue-600">
        <Activity class="w-5 h-5" />
      </div>
      <div>
        <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest leading-none mb-1">{{ moduleName }}</p>
        <h2 class="text-xl font-black text-slate-800 tracking-tight flex items-center gap-2">
           <ArrowLeft @click="router.back()" class="w-4 h-4 cursor-pointer hover:text-blue-600" />
           {{ pageTitle }}
        </h2>
      </div>
    </div>

    <div class="flex items-center gap-6">
      <div class="hidden md:flex items-center gap-4 bg-white/40 backdrop-blur-md px-4 py-2 rounded-2xl border border-white/50 shadow-sm">
        <Search class="w-4 h-4 text-slate-400" />
        <input type="text" placeholder="Search..." class="bg-transparent border-none outline-none text-sm font-medium w-40">
      </div>
      
      <div class="flex items-center gap-3">
         <!-- Notification Button -->
         <div class="relative dropdown-container">
           <button 
             @click.stop="toggleNotifications"
             class="w-10 h-10 rounded-xl bg-white/40 backdrop-blur-md border border-white/50 shadow-sm flex items-center justify-center text-slate-500 hover:text-blue-600 transition-all relative"
             :class="{ 'bg-white shadow-md text-blue-600': showNotifications }"
           >
             <Bell class="w-5 h-5" />
             <span class="absolute top-2 right-2 w-2 h-2 bg-rose-500 rounded-full border-2 border-white"></span>
           </button>

           <!-- Notification Dropdown -->
           <div v-if="showNotifications" class="absolute right-0 mt-3 w-80 bg-white/90 backdrop-blur-2xl rounded-3xl shadow-2xl border border-white/60 p-4 animate-in fade-in zoom-in duration-200">
             <div class="flex items-center justify-between mb-4 px-2">
                <h3 class="text-sm font-black italic text-slate-800">Notifications</h3>
                <span class="text-[10px] font-bold text-blue-600 uppercase">Mark all read</span>
             </div>
             <div class="space-y-2">
                <div v-for="n in notifications" :key="n.id" class="p-3 rounded-2xl hover:bg-white transition-all cursor-pointer border border-transparent hover:border-slate-100 flex gap-4">
                   <div :class="n.color" class="w-10 h-10 rounded-xl bg-slate-50 flex items-center justify-center shrink-0">
                      <component :is="n.icon" class="w-5 h-5" />
                   </div>
                   <div class="overflow-hidden">
                      <p class="text-xs font-black text-slate-800 leading-none mb-1">{{ n.title }}</p>
                      <p class="text-[10px] text-slate-500 truncate">{{ n.desc }}</p>
                      <div class="flex items-center gap-1 mt-1 text-[9px] font-bold text-slate-400 uppercase">
                         <Clock class="w-2.5 h-2.5" /> {{ n.time }}
                      </div>
                   </div>
                </div>
             </div>
             <button class="w-full mt-4 py-3 rounded-xl bg-slate-50 text-[10px] font-black text-slate-400 uppercase tracking-widest hover:bg-blue-50 hover:text-blue-600 transition-all">
                View All Activity
             </button>
           </div>
         </div>



         <!-- Profile Button & Menu -->
         <div class="relative dropdown-container">
            <div 
              @click.stop="toggleProfileMenu"
              class="flex items-center gap-3 pl-4 border-l border-white/20 cursor-pointer group"
            >
               <div class="text-right hidden sm:block">
                  <p class="text-xs font-black text-slate-800 leading-none group-hover:text-blue-600 transition-colors">{{ user.name }}</p>
                  <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">{{ user.role }}</p>
               </div>
               <div 
                 class="w-10 h-10 rounded-xl bg-blue-600 flex items-center justify-center text-white font-black overflow-hidden shadow-lg shadow-blue-500/20 group-hover:scale-105 transition-transform"
                 :class="{ 'ring-2 ring-blue-500 ring-offset-2': showProfileMenu }"
               >
                  <img v-if="user.avatar" :src="user.avatar" class="w-full h-full object-cover">
                  <span v-else>M</span>
               </div>
            </div>

            <!-- Profile Dropdown -->
            <div v-if="showProfileMenu" class="absolute right-0 mt-3 w-56 bg-white/90 backdrop-blur-2xl rounded-3xl shadow-2xl border border-white/60 p-2 animate-in fade-in zoom-in duration-200">
               <div class="p-4 mb-2 border-b border-slate-100 sm:hidden">
                  <p class="text-xs font-black text-slate-800 leading-none">{{ user.name }}</p>
                  <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ user.role }}</p>
               </div>
               <div class="space-y-1">
                  <button @click="router.push('/profile')" class="w-full flex items-center gap-3 p-3 rounded-2xl hover:bg-blue-50 text-slate-600 hover:text-blue-600 transition-all">
                     <User class="w-4 h-4" />
                     <span class="text-xs font-bold">Profil Saya</span>
                  </button>

                  <div class="h-px bg-slate-100 my-2 mx-2"></div>
                  <button @click="handleLogout" class="w-full flex items-center gap-3 p-3 rounded-2xl hover:bg-rose-50 text-rose-500 transition-all">
                     <LogOut class="w-4 h-4" />
                     <span class="text-xs font-bold">Logout System</span>
                  </button>
               </div>
            </div>
         </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
header {
  font-family: 'Outfit', sans-serif;
}

/* Glassmorphism Animation */
.animate-in {
  animation: reveal 0.3s cubic-bezier(0.2, 0.8, 0.2, 1);
}

@keyframes reveal {
  from { opacity: 0; transform: translateY(10px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
</style>

