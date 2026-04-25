<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  TrendingUp, BookOpen, Calendar, UserCheck, GraduationCap, Clock, 
  CalendarRange, School, ArrowRight, ArrowLeft
} from 'lucide-vue-next'
import axios from 'axios'

// State
const router = useRouter()
const stats = ref({
  total_subjects: 0,
  total_schedules: 0,
  total_attendance_today: 0
})

const isLoading = ref(true)
const user = ref({ name: 'Admin', role: 'Administrator' })

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/v1/admin/edu/stats')
    stats.value = response.data
  } catch (err) {
    console.error('Gagal mengambil statistik akademik:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  const savedUser = localStorage.getItem('sigma_user')
  if (savedUser) user.value = JSON.parse(savedUser)
  
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
  
  fetchStats()
})

</script>

<template>
  <div class="p-8 space-y-10">
    <!-- Header Actions -->
    <div class="flex items-center justify-between mb-2">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-blue-500 rounded-full"></div>
        <h2 class="text-sm font-black uppercase tracking-[0.3em] text-sigma-muted">Academic Module</h2>
      </div>
            <div class="flex items-center gap-3">
        <!-- Back to Portal -->
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-blue-500/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Welcome Banner -->
    <div class="relative overflow-hidden rounded-[2rem] bg-gradient-to-r from-blue-600 to-indigo-900/80 p-8 shadow-2xl shadow-blue-500/10 border border-white/5">
      <div class="absolute top-0 right-0 -translate-y-1/2 translate-x-1/4 w-80 h-80 bg-blue-400/10 rounded-full blur-[100px]"></div>
      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <span class="inline-block px-3 py-1 rounded-full bg-white/10 backdrop-blur-md text-[10px] font-bold uppercase tracking-[0.2em] text-blue-100/90 mb-2 border border-white/10">
            Sigmaedu Intelligence
          </span>
          <h1 class="text-2xl md:text-4xl font-black text-white leading-tight">
            Pusat Akademik, <span class="text-blue-300 font-black italic">{{ user.name }}!</span>
          </h1>
          <p class="text-blue-100/60 text-sm font-medium uppercase tracking-widest">Manajemen Kurikulum, Jadwal & Penilaian.</p>
        </div>
        
        <div class="flex items-center gap-5 bg-black/20 backdrop-blur-2xl p-5 rounded-3xl border border-white/5">
          <div class="text-right">
            <p class="text-[10px] uppercase tracking-widest text-blue-400 font-black">Academic Year</p>
            <p class="text-sm font-bold text-white">2024 / 2025</p>
          </div>
          <div class="w-10 h-10 bg-blue-500/20 rounded-2xl flex items-center justify-center text-blue-400">
            <Calendar class="w-5 h-5" />
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <!-- Total Subjects -->
      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] flex flex-col gap-6 group hover:border-blue-500/30 transition-all shadow-sm">
        <div class="flex items-center justify-between">
          <div class="w-14 h-14 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-500 group-hover:bg-blue-500 group-hover:text-white transition-all duration-500 shadow-lg shadow-blue-500/5">
            <BookOpen class="w-7 h-7" />
          </div>
          <div class="text-right">
            <p class="text-[10px] text-sigma-muted font-black uppercase tracking-[0.2em]">Mata Pelajaran</p>
            <h3 class="text-4xl font-black italic">{{ stats.total_subjects }}</h3>
          </div>
        </div>
        <div class="flex items-center justify-between pt-4 border-t border-sigma-border">
          <span class="text-[10px] font-black text-sigma-muted uppercase tracking-widest">Aktif di Kurikulum</span>
          <TrendingUp class="w-4 h-4 text-blue-400" />
        </div>
      </div>

      <!-- Total Schedules -->
      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] flex flex-col gap-6 group hover:border-indigo-500/30 transition-all shadow-sm">
        <div class="flex items-center justify-between">
          <div class="w-14 h-14 bg-indigo-500/10 rounded-2xl flex items-center justify-center text-indigo-500 group-hover:bg-indigo-500 group-hover:text-white transition-all duration-500 shadow-lg shadow-indigo-500/5">
            <Clock class="w-7 h-7" />
          </div>
          <div class="text-right">
            <p class="text-[10px] text-sigma-muted font-black uppercase tracking-[0.2em]">Jadwal KBM</p>
            <h3 class="text-4xl font-black italic">{{ stats.total_schedules }}</h3>
          </div>
        </div>
        <div class="flex items-center justify-between pt-4 border-t border-sigma-border">
          <span class="text-[10px] font-black text-sigma-muted uppercase tracking-widest">Sesi Terdaftar / Pekan</span>
          <TrendingUp class="w-4 h-4 text-indigo-400" />
        </div>
      </div>

      <!-- Attendance Today -->
      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] flex flex-col gap-6 group hover:border-emerald-500/30 transition-all shadow-sm">
        <div class="flex items-center justify-between">
          <div class="w-14 h-14 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-500 group-hover:bg-emerald-500 group-hover:text-white transition-all duration-500 shadow-lg shadow-emerald-500/5">
            <UserCheck class="w-7 h-7" />
          </div>
          <div class="text-right">
            <p class="text-[10px] text-sigma-muted font-black uppercase tracking-[0.2em]">Presensi Hari Ini</p>
            <h3 class="text-4xl font-black italic">{{ stats.total_attendance_today }}</h3>
          </div>
        </div>
        <div class="flex items-center justify-between pt-4 border-t border-sigma-border">
          <span class="text-[10px] font-black text-sigma-muted uppercase tracking-widest">Santri Masuk Kelas</span>
          <div class="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- Quick Actions / Academic Preview -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="bg-sigma-surface border border-sigma-border p-10 rounded-[3rem] space-y-8 shadow-sm">
        <div class="flex items-center justify-between">
          <h3 class="text-2xl font-black flex items-center gap-4 italic uppercase">
            <div class="w-2 h-8 bg-blue-500 rounded-full"></div>
            Informasi <span class="text-blue-500 not-italic">Terbaru</span>
          </h3>
        </div>
        
        <div class="space-y-4">
          <div class="p-6 rounded-3xl bg-sigma-surface-alt border border-sigma-border flex items-center gap-6 hover:border-blue-500/20 transition-all cursor-pointer group">
             <div class="w-12 h-12 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-500 group-hover:bg-blue-500 group-hover:text-white transition-all">
                <CalendarRange class="w-6 h-6" />
             </div>
             <div class="flex-1">
                <h4 class="font-bold text-sm">Persiapan Ujian Tengah Semester</h4>
                <p class="text-xs text-sigma-muted mt-1">Ganjil 2024/2025 • Mulai 15 Okt</p>
             </div>
             <ArrowRight class="w-5 h-5 text-sigma-muted group-hover:text-blue-500 transition-all" />
          </div>

          <div class="p-6 rounded-3xl bg-sigma-surface-alt border border-sigma-border flex items-center gap-6 hover:border-blue-500/20 transition-all cursor-pointer group">
             <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center text-indigo-500 group-hover:bg-indigo-500 group-hover:text-white transition-all">
                <GraduationCap class="w-6 h-6" />
             </div>
             <div class="flex-1">
                <h4 class="font-bold text-sm">Rapat Kerja Kurikulum KMI</h4>
                <p class="text-xs text-sigma-muted mt-1">Auditorium Pusat • Besok 09:00</p>
             </div>
             <ArrowRight class="w-5 h-5 text-sigma-muted group-hover:text-indigo-500 transition-all" />
          </div>
        </div>
      </div>

      <div class="bg-sigma-surface border border-sigma-border p-10 rounded-[3rem] flex flex-col items-center justify-center text-center space-y-6 shadow-sm overflow-hidden relative">
        <div class="absolute -top-10 -right-10 w-40 h-40 bg-blue-500/5 rounded-full blur-[40px]"></div>
        <div class="w-24 h-24 bg-blue-500/10 rounded-[2rem] flex items-center justify-center text-blue-500 mb-2">
          <School class="w-12 h-12" />
        </div>
        <div class="space-y-2">
          <h3 class="text-2xl font-black uppercase italic tracking-tighter">Academic System</h3>
          <p class="text-sigma-muted text-sm max-w-xs mx-auto">Kelola data kurikulum, penilaian raport, dan progres tahfidz santri secara terpusat.</p>
        </div>
        <button @click="router.push('/sigmaedu/subjects')" class="px-8 py-4 bg-blue-600 text-white rounded-2xl font-bold uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20">
          Mulai Pengaturan Akademik
        </button>
      </div>
    </div>
  </div>
</template>


<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }
</style>
