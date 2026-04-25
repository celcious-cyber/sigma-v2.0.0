<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { 
  Users, Activity, Target, UserCog,
  Building2, ChevronRight, Calendar,
  UserCheck, Clock, FileText, Layers, LayoutDashboard
} from 'lucide-vue-next'
import axios from 'axios'
import apexchart from 'vue3-apexcharts'

// State
const stats = ref({
  total_subjects: 0,
  total_schedules: 0,
  total_students: 0,
  total_classrooms: 0,
  total_teachers: 0,
  attendance_today: [] as any[],
  tahfidz_today: [] as any[],
  lesson_mem_today: [] as any[],
  top_absent_students: [] as any[],
  character_stats: [] as any[],
  journals_today_count: 0,
  assessments_summary: [] as any[],
  gender_stats: [] as any[],
  classroom_stats: [] as any[],
  teacher_gender_stats: [] as any[],
  teacher_attendance_today: [] as any[],
  upcoming_events: [] as any[],
  quran_chart: { labels: [] as string[], datasets: [] as any[] },
  lesson_chart: { labels: [] as string[], datasets: [] as any[] },
  attendance_chart: { labels: [] as string[], datasets: [] as any[] }
})

const formatDate = (dateStr: string) => {
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

const getRatio = (val1: number, val2: number) => {
  if (val1 === 0 && val2 === 0) return '0:0'
  if (val1 === 0) return `0:${val2}`
  if (val2 === 0) return `${val1}:0`
  
  const gcd = (a: number, b: number): number => b === 0 ? a : gcd(b, a % b)
  const common = gcd(val1, val2)
  return `${val1 / common}:${val2 / common}`
}

const isLoading = ref(true)
const lastUpdated = ref(new Date().toLocaleTimeString())
const user = ref({ name: 'Michael', role: 'Administrator', avatar: '' })

// Chart Periods
const quranPeriod = ref('Week')
const lessonPeriod = ref('Week')
const attendancePeriod = ref('Week')

// Carousel Logic for Agenda
const activeEventIndex = ref(0)
let carouselInterval: any = null

const startCarousel = () => {
  if (carouselInterval) clearInterval(carouselInterval)
  carouselInterval = setInterval(() => {
    if (stats.value.upcoming_events.length > 0) {
      activeEventIndex.value = (activeEventIndex.value + 1) % stats.value.upcoming_events.length
    }
  }, 5000)
}

const fetchStats = async () => {
  isLoading.value = true
  try {
    const params = new URLSearchParams()
    params.append('quran_filter', quranPeriod.value.toLowerCase())
    params.append('lesson_filter', lessonPeriod.value.toLowerCase())
    params.append('attendance_filter', attendancePeriod.value.toLowerCase())
    
    const response = await axios.get(`/api/v1/admin/edu/stats?${params.toString()}`)
    const data = response.data
    const arrayKeys = [
      'attendance_today', 'tahfidz_today', 'lesson_mem_today', 
      'top_absent_students', 'character_stats', 'gender_stats', 
      'classroom_stats', 'teacher_gender_stats', 'teacher_attendance_today',
      'assessments_summary', 'upcoming_events'
    ]
    arrayKeys.forEach(key => {
      if (!data[key]) data[key] = []
    })
    
    // Fallback labels if empty
    if (!data.quran_chart) data.quran_chart = { labels: [], datasets: [] }
    if (!data.lesson_chart) data.lesson_chart = { labels: [], datasets: [] }
    if (!data.attendance_chart) data.attendance_chart = { labels: [], datasets: [] }
    
    stats.value = data
    lastUpdated.value = new Date().toLocaleTimeString()
  } catch (err) {
    console.error('Gagal mengambil statistik akademik:', err)
  } finally {
    isLoading.value = false
  }
}

watch([quranPeriod, lessonPeriod, attendancePeriod], () => {
  fetchStats()
})

const getColor = (code: string) => {
  const map: any = {
    emerald: '#10b981',
    pink: '#f472b6',
    blue: '#3b82f6',
    indigo: '#6366f1',
    amber: '#f59e0b',
    rose: '#f43f5e',
    cyan: '#06b6d4',
    teal: '#14b8a6'
  }
  return map[code] || '#cbd5e1'
}

onMounted(() => {
  const savedUser = localStorage.getItem('sigma_user')
  if (savedUser) user.value = JSON.parse(savedUser)
  
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
  
  fetchStats().then(() => {
    startCarousel()
  })
})

</script>

<template>
  <div class="min-h-screen font-sans pb-20">
    <!-- MAIN CONTENT AREA (Glassmorphism) -->
    <main class="max-w-[1600px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- ACADEMIC AGENDA BANNER (Glass Style) -->
      <section class="relative overflow-hidden rounded-[3rem] bg-white/40 backdrop-blur-2xl p-8 lg:p-12 shadow-2xl border border-white/60 animate-reveal">
        <div class="absolute -top-24 -right-24 w-80 h-80 bg-blue-400/20 rounded-full blur-[100px]"></div>
        <div class="absolute -bottom-24 -left-24 w-60 h-60 bg-pink-400/10 rounded-full blur-[80px]"></div>
        
        <div class="relative z-10 flex flex-col lg:flex-row lg:items-center gap-10">
          <!-- Featured Carousel Item -->
          <div class="flex-1 space-y-6">
            <div class="flex items-center gap-2">
               <span class="px-3 py-1 rounded-full bg-blue-500/10 text-blue-600 border border-blue-500/20 text-[9px] font-black uppercase tracking-widest">
                 Agenda Mendatang
               </span>
            </div>
            
            <div class="relative h-40">
               <TransitionGroup name="slide-fade">
                 <div v-for="(event, idx) in stats.upcoming_events" :key="event.id" 
                      v-show="activeEventIndex === idx"
                      class="absolute inset-0 flex flex-col gap-4">
                   <h1 class="text-4xl lg:text-5xl font-black text-slate-900 tracking-tighter italic">{{ event.title }}</h1>
                   <p class="text-slate-500 text-sm font-medium max-w-xl leading-relaxed">
                     {{ event.description || 'Laksanakan persiapan terbaik untuk mensukseskan kegiatan akademik ini demi kemajuan santri.' }}
                   </p>
                   <div class="flex items-center gap-6 mt-4">
                      <div class="flex items-center gap-3">
                         <div class="w-10 h-10 rounded-xl bg-blue-500 text-white flex items-center justify-center shadow-lg shadow-blue-500/20">
                           <Calendar class="w-5 h-5" />
                         </div>
                         <span class="text-xs font-black text-slate-700">{{ formatDate(event.start_date) }}</span>
                      </div>
                      <button class="px-6 py-2.5 rounded-xl bg-[#0F2942] text-white text-[10px] font-black uppercase tracking-widest hover:bg-blue-600 transition-all shadow-xl">
                        View Details
                      </button>
                   </div>
                 </div>
               </TransitionGroup>
               <div v-if="!stats.upcoming_events.length" class="flex items-center h-full text-slate-400 italic font-medium">
                  Belum ada agenda terdaftar untuk periode ini.
               </div>
            </div>
          </div>

          <!-- Mini List of Next Agendas -->
          <div class="lg:w-80 flex flex-col gap-3">
             <div v-for="(event, idx) in stats.upcoming_events.slice(0, 2)" :key="event.id" 
                  @click="activeEventIndex = idx"
                  :class="activeEventIndex === idx ? 'bg-white border-blue-500/30 shadow-xl scale-105' : 'bg-white/40 border-transparent hover:bg-white/60'"
                  class="p-4 rounded-3xl border transition-all duration-500 cursor-pointer flex items-center gap-4 group">
                <div :class="activeEventIndex === idx ? 'bg-blue-500 text-white' : 'bg-slate-100 text-slate-400'" 
                     class="w-12 h-12 rounded-2xl flex flex-col items-center justify-center transition-colors">
                   <span class="text-xs font-black leading-none">{{ formatDate(event.start_date).split(' ')[0] }}</span>
                   <span class="text-[8px] font-bold uppercase">{{ formatDate(event.start_date).split(' ')[1] }}</span>
                </div>
                <div class="flex-1 overflow-hidden">
                   <p class="text-[8px] font-bold text-slate-400 uppercase tracking-widest mb-0.5">{{ event.type || 'Kegiatan' }}</p>
                   <h4 class="text-xs font-black text-slate-800 truncate">{{ event.title }}</h4>
                </div>
             </div>
          </div>
        </div>
      </section>

      <!-- ANALYTICS CHARTS (ApexCharts 3-Grid) -->
      <section class="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-reveal reveal-delay-1">
        <!-- Chart 1: Hafalan Quran -->
        <div class="bg-emerald-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-emerald-100 shadow-xl overflow-hidden group">
           <div class="flex items-center justify-between mb-8">
              <div>
                <h3 class="text-sm font-black text-slate-800 italic">Hafalan Quran</h3>
                <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Progres Hafalan</p>
              </div>
              <div class="flex bg-slate-100/50 p-1 rounded-xl">
                 <button v-for="t in ['Minggu', 'Bulan']" :key="t" 
                         @click="quranPeriod = t === 'Minggu' ? 'Week' : 'Month'"
                         :class="(quranPeriod === 'Week' && t === 'Minggu') || (quranPeriod === 'Month' && t === 'Bulan') ? 'bg-white shadow-sm text-emerald-600' : 'text-slate-400'"
                         class="px-3 py-1 rounded-lg text-[9px] font-black uppercase transition-all">
                   {{ t }}
                 </button>
              </div>
           </div>

           <div class="h-48">
              <apexchart 
                type="area" 
                height="100%" 
                :options="{
                   chart: { toolbar: { show: false }, sparkline: { enabled: false }, animations: { enabled: true, speed: 800 } },
                   dataLabels: { enabled: false },
                   stroke: { curve: 'smooth', width: 2 },
                   fill: { 
                    type: 'gradient', 
                    gradient: { shadeIntensity: 1, opacityFrom: 0.15, opacityTo: 0.05 } 
                   },
                   colors: stats.quran_chart.datasets.map(d => getColor(d.color_code)),
                   xaxis: { categories: stats.quran_chart.labels, axisBorder: { show: false }, axisTicks: { show: false }, labels: { style: { colors: '#94a3b8', fontSize: '9px', fontWeight: 600 } } },
                   yaxis: { show: false },
                   grid: { show: false },
                   tooltip: { theme: 'light', x: { show: true } },
                   legend: { 
                     show: true, position: 'top', horizontalAlign: 'right', fontSize: '9px', fontWeight: 700, 
                     markers: { strokeWidth: 0 }, itemMargin: { horizontal: 5 } 
                   }
                 }" 
                 :series="stats.quran_chart.datasets.map(d => ({ name: d.label, data: d.data }))"
              />
           </div>
        </div>

        <!-- Chart 2: Hafalan Pelajaran -->
        <div class="bg-blue-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-blue-100 shadow-xl overflow-hidden group">
           <div class="flex items-center justify-between mb-8">
              <div>
                <h3 class="text-sm font-black text-slate-800 italic">Hafalan Pelajaran</h3>
                <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Target Kurikulum</p>
              </div>
              <div class="flex bg-slate-100/50 p-1 rounded-xl">
                 <button v-for="t in ['Minggu', 'Bulan']" :key="t" 
                         @click="lessonPeriod = t === 'Minggu' ? 'Week' : 'Month'"
                         :class="(lessonPeriod === 'Week' && t === 'Minggu') || (lessonPeriod === 'Month' && t === 'Bulan') ? 'bg-white shadow-sm text-blue-600' : 'text-slate-400'"
                         class="px-3 py-1 rounded-lg text-[9px] font-black uppercase transition-all">
                   {{ t }}
                 </button>
              </div>
           </div>

           <div class="h-48">
              <apexchart 
                type="area" 
                height="100%" 
                :options="{
                   chart: { toolbar: { show: false }, animations: { enabled: true, speed: 800 } },
                   dataLabels: { enabled: false },
                   stroke: { curve: 'smooth', width: 2 },
                   fill: { 
                    type: 'gradient', 
                    gradient: { shadeIntensity: 1, opacityFrom: 0.15, opacityTo: 0.05 } 
                   },
                   colors: stats.lesson_chart.datasets.map(d => getColor(d.color_code)),
                   xaxis: { categories: stats.lesson_chart.labels, labels: { style: { colors: '#94a3b8', fontSize: '9px', fontWeight: 600 } } },
                   yaxis: { show: false },
                   grid: { show: false },
                   tooltip: { theme: 'light' },
                   legend: { 
                     show: true, position: 'top', horizontalAlign: 'right', fontSize: '9px', fontWeight: 700, 
                     markers: { strokeWidth: 0 }, itemMargin: { horizontal: 5 } 
                   }
                 }" 
                 :series="stats.lesson_chart.datasets.map(d => ({ name: d.label, data: d.data }))"
              />
           </div>
        </div>

        <!-- Chart 3: Presensi Santri -->
        <div class="bg-indigo-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-indigo-100 shadow-xl overflow-hidden group">
           <div class="flex items-center justify-between mb-8">
              <div>
                <h3 class="text-sm font-black text-slate-800 italic">Presensi Santri</h3>
                <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Tren Kehadiran</p>
              </div>
              <div class="flex bg-slate-100/50 p-1 rounded-xl">
                 <button v-for="t in ['Minggu', 'Bulan']" :key="t" 
                         @click="attendancePeriod = t === 'Minggu' ? 'Week' : 'Month'"
                         :class="(attendancePeriod === 'Week' && t === 'Minggu') || (attendancePeriod === 'Month' && t === 'Bulan') ? 'bg-white shadow-sm text-indigo-600' : 'text-slate-400'"
                         class="px-3 py-1 rounded-lg text-[9px] font-black uppercase transition-all">
                   {{ t }}
                 </button>
              </div>
           </div>

           <div class="h-48">
              <apexchart 
                type="area" 
                height="100%" 
                :options="{
                   chart: { toolbar: { show: false }, animations: { enabled: true, speed: 800 } },
                   dataLabels: { enabled: false },
                   stroke: { curve: 'smooth', width: 2 },
                   fill: { 
                    type: 'gradient', 
                    gradient: { shadeIntensity: 1, opacityFrom: 0.15, opacityTo: 0.05 } 
                   },
                   colors: stats.attendance_chart.datasets.map(d => getColor(d.color_code)),
                   xaxis: { categories: stats.attendance_chart.labels, labels: { style: { colors: '#94a3b8', fontSize: '9px', fontWeight: 600 } } },
                   yaxis: { show: false },
                   grid: { show: false },
                   tooltip: { theme: 'light' },
                   legend: { 
                     show: true, position: 'top', horizontalAlign: 'right', fontSize: '9px', fontWeight: 700, 
                     markers: { strokeWidth: 0 }, itemMargin: { horizontal: 5 } 
                   }
                 }" 
                 :series="stats.attendance_chart.datasets.map(d => ({ name: d.label, data: d.data }))"
              />
           </div>
        </div>
      </section>

      <!-- RADIAL ANALYTICS (ApexCharts 3-Grid) -->
      <section class="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-reveal reveal-delay-2">
        <!-- Chart 1: Gender Santri -->
        <div class="bg-cyan-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-cyan-100 shadow-xl overflow-hidden group flex flex-col items-center">
           <div class="w-full mb-4 text-center">
              <h3 class="text-sm font-black text-slate-800 italic">Gender Santri</h3>
              <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Rasio Santri Laki-laki & Perempuan</p>
           </div>
           <div class="h-60 w-full">
              <apexchart 
                type="donut" 
                height="100%" 
                :options="{
                  chart: { sparkline: { enabled: true } },
                  stroke: { show: false },
                  colors: ['#3b82f6', '#f472b6'],
                  labels: ['LAKI-LAKI', 'PEREMPUAN'],
                  plotOptions: {
                    pie: {
                      donut: {
                        size: '80%',
                        labels: {
                          show: true,
                          name: { show: true, fontSize: '9px', fontWeight: 800, color: '#94a3b8', offsetY: 25 },
                          value: { 
                            show: true, fontSize: '24px', fontWeight: 900, color: '#1e293b', offsetY: -10,
                            formatter: () => getRatio(stats.gender_stats.find(g => g.gender === 'L')?.count || 0, stats.gender_stats.find(g => g.gender === 'P')?.count || 0)
                          },
                          total: {
                            show: true, label: 'RATIO', fontSize: '9px', fontWeight: 800, color: '#94a3b8',
                            formatter: () => getRatio(stats.gender_stats.find(g => g.gender === 'L')?.count || 0, stats.gender_stats.find(g => g.gender === 'P')?.count || 0)
                          }
                        }
                      }
                    }
                  },
                  legend: { show: false },
                  dataLabels: { enabled: false },
                  tooltip: { enabled: true, theme: 'light' }
                }" 
                :series="[
                  stats.gender_stats.find(g => g.gender === 'L')?.count || 0,
                  stats.gender_stats.find(g => g.gender === 'P')?.count || 0
                ]"
              />
           </div>
           <div class="flex gap-4 mt-2">
              <div class="flex items-center gap-2">
                 <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                 <span class="text-[9px] font-black text-slate-400 uppercase">L: {{ stats.gender_stats.find(g => g.gender === 'L')?.count || 0 }}</span>
              </div>
              <div class="flex items-center gap-2">
                 <div class="w-2 h-2 rounded-full bg-pink-400"></div>
                 <span class="text-[9px] font-black text-slate-400 uppercase">P: {{ stats.gender_stats.find(g => g.gender === 'P')?.count || 0 }}</span>
              </div>
           </div>
        </div>

        <!-- Chart 2: Gender Guru -->
        <div class="bg-teal-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-teal-100 shadow-xl overflow-hidden group flex flex-col items-center">
           <div class="w-full mb-4 text-center">
              <h3 class="text-sm font-black text-slate-800 italic">Gender Guru</h3>
              <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Distribusi Tenaga Pengajar</p>
           </div>
           <div class="h-60 w-full">
              <apexchart 
                type="donut" 
                height="100%" 
                :options="{
                  chart: { sparkline: { enabled: true } },
                  stroke: { show: false },
                  colors: ['#059669', '#fb923c'],
                  labels: ['LAKI-LAKI', 'PEREMPUAN'],
                  plotOptions: {
                    pie: {
                      donut: {
                        size: '80%',
                        labels: {
                          show: true,
                          name: { show: true, fontSize: '9px', fontWeight: 800, color: '#94a3b8', offsetY: 25 },
                          value: { 
                            show: true, fontSize: '24px', fontWeight: 900, color: '#1e293b', offsetY: -10,
                            formatter: () => getRatio(stats.teacher_gender_stats.find(g => g.gender === 'L')?.count || 0, stats.teacher_gender_stats.find(g => g.gender === 'P')?.count || 0)
                          },
                          total: {
                            show: true, label: 'RATIO', fontSize: '9px', fontWeight: 800, color: '#94a3b8',
                            formatter: () => getRatio(stats.teacher_gender_stats.find(g => g.gender === 'L')?.count || 0, stats.teacher_gender_stats.find(g => g.gender === 'P')?.count || 0)
                          }
                        }
                      }
                    }
                  },
                  legend: { show: false },
                  dataLabels: { enabled: false }
                }" 
                :series="[
                  stats.teacher_gender_stats.find(g => g.gender === 'L')?.count || 0,
                  stats.teacher_gender_stats.find(g => g.gender === 'P')?.count || 0
                ]"
              />
           </div>
           <div class="flex gap-4 mt-2">
              <div class="flex items-center gap-2">
                 <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                 <span class="text-[9px] font-black text-slate-400 uppercase">L: {{ stats.teacher_gender_stats.find(g => g.gender === 'L')?.count || 0 }}</span>
              </div>
              <div class="flex items-center gap-2">
                 <div class="w-2 h-2 rounded-full bg-orange-400"></div>
                 <span class="text-[9px] font-black text-slate-400 uppercase">P: {{ stats.teacher_gender_stats.find(g => g.gender === 'P')?.count || 0 }}</span>
              </div>
           </div>
        </div>

        <!-- Chart 3: Karakter Santri -->
        <div class="bg-amber-50/80 backdrop-blur-md rounded-[3rem] p-8 border border-amber-100 shadow-xl relative overflow-hidden group flex flex-col items-center">
           <div class="absolute -right-10 -top-10 w-40 h-40 bg-amber-500/10 rounded-full blur-3xl"></div>
           <div class="w-full mb-4 text-center relative z-10">
              <h3 class="text-sm font-black text-slate-800 italic">Karakter Santri</h3>
              <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Rata-rata Skor Karakter</p>
           </div>
           <div class="h-64 w-full relative z-10">
              <apexchart 
                type="radialBar" 
                height="100%" 
                :options="{
                  chart: { sparkline: { enabled: true } },
                  plotOptions: {
                    radialBar: {
                      hollow: { size: '60%', background: 'transparent' },
                      track: { background: '#00000005', strokeWidth: '100%' },
                      dataLabels: {
                        name: { show: true, fontSize: '9px', fontWeight: 800, color: '#94a3b8', offsetY: 30 },
                        value: { show: true, fontSize: '28px', fontWeight: 900, color: '#1e293b', offsetY: -10, formatter: (val: number) => val + ' pts' }
                      }
                    }
                  },
                  colors: ['#f59e0b'],
                  labels: ['AVG SCORE'],
                  stroke: { lineCap: 'round' }
                }" 
                :series="[82]"
              />
           </div>
           <div class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-2 relative z-10">Indeks Kebiasaan Positif</div>
        </div>
      </section>

      <section class="reveal-delay-2">
         <div class="flex items-center justify-between">
           <h3 class="text-lg font-black text-slate-800 tracking-tight italic">Statistik Cepat</h3>
           <div class="flex gap-2">
              <button class="w-8 h-8 rounded-lg bg-white/40 flex items-center justify-center text-slate-400 hover:text-blue-600 transition-all shadow-sm">
                <Layers class="w-4 h-4" />
              </button>
              <button class="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center text-white shadow-md">
                <LayoutDashboard class="w-4 h-4" />
              </button>
           </div>
         </div>

         <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
             <div v-for="card in [
               { label: 'Populasi Santri', value: stats.total_students, sub: 'Santri Aktif', color: 'pink', icon: Users, gradient: 'from-pink-500 to-rose-400' },
               { label: 'SDM Pengajar', value: stats.total_teachers, sub: 'Tenaga Guru', color: 'blue', icon: UserCog, gradient: 'from-blue-500 to-indigo-400' },
               { label: 'Total Kelas', value: stats.total_classrooms, sub: 'Ruang Belajar', color: 'emerald', icon: Building2, gradient: 'from-emerald-500 to-teal-400' },
               { label: 'Mata Pelajaran', value: stats.total_subjects, sub: 'Kurikulum Terpadu', color: 'orange', icon: Target, gradient: 'from-orange-500 to-amber-400' },
               { label: 'Hadir Hari Ini', value: stats.attendance_today.find(a => a.status === 'Hadir')?.count || 0, sub: 'Presensi Realtime', color: 'cyan', icon: UserCheck, gradient: 'from-cyan-500 to-blue-400' },
               { label: 'Jadwal Aktif', value: stats.total_schedules, sub: 'Sesi Belajar', color: 'indigo', icon: Clock, gradient: 'from-indigo-500 to-purple-400' },
               { label: 'Jurnal Terisi', value: stats.journals_today_count, sub: 'Laporan Guru', color: 'violet', icon: FileText, gradient: 'from-violet-500 to-fuchsia-400' },
               { label: 'Agenda Sekolah', value: stats.upcoming_events.length, sub: 'Kegiatan Terdekat', color: 'amber', icon: Calendar, gradient: 'from-amber-500 to-orange-400' }
             ]" :key="card.label" class="bg-white/60 backdrop-blur-md rounded-[2.5rem] border border-white/60 p-6 shadow-xl hover:shadow-2xl transition-all duration-500 group flex flex-col gap-6 overflow-hidden relative">
               <div :class="`bg-gradient-to-br ${card.gradient}`" class="absolute top-0 right-0 w-24 h-24 -mr-12 -mt-12 rounded-full opacity-10 group-hover:scale-150 transition-transform duration-700"></div>
               
               <div class="flex items-center gap-4 relative z-10">
                  <div :class="`bg-gradient-to-br ${card.gradient}`" class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shadow-lg">
                     <component :is="card.icon" class="w-7 h-7" />
                  </div>
                  <div>
                    <h4 class="text-sm font-black text-slate-800 leading-none mb-1 group-hover:text-blue-600 transition-colors">{{ card.label }}</h4>
                    <p class="text-[10px] font-medium text-slate-400 italic">{{ card.sub }}</p>
                  </div>
               </div>

               <div class="space-y-3 relative z-10">
                  <div class="flex justify-between items-end">
                     <div class="text-3xl font-black italic tracking-tighter text-slate-900">{{ card.value }}</div>
                     <div class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Update Hari Ini</div>
                  </div>
                  <div class="h-2 bg-slate-100 rounded-full overflow-hidden">
                     <div :class="`bg-gradient-to-r ${card.gradient}`" class="h-full w-3/4 rounded-full"></div>
                  </div>
               </div>
            </div>
         </div>
      </section>

      <!-- DETAILED ANALYSIS GRID -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 reveal-delay-3 pb-10">
        
        <!-- Presence Bar -->
        <div class="bg-white/60 backdrop-blur-md rounded-[3rem] border border-white/60 p-10 shadow-xl space-y-8">
           <div class="flex items-center justify-between">
              <h3 class="text-lg font-black text-slate-800 italic flex items-center gap-3">
                 <div class="w-8 h-8 rounded-xl bg-blue-500 text-white flex items-center justify-center shadow-lg"><Activity class="w-4 h-4" /></div>
                 Monitoring Kehadiran
              </h3>
              <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">{{ lastUpdated }}</span>
           </div>
           
           <div class="grid grid-cols-2 gap-6">
              <div v-for="item in [
                { key: 'Sakit', label: 'Sakit', val: stats.attendance_today.find(a => a.status === 'Sakit')?.count || 0, color: 'orange' },
                { key: 'Izin', label: 'Izin', val: stats.attendance_today.find(a => a.status === 'Izin')?.count || 0, color: 'blue' },
                { key: 'Alpa', label: 'Alpha', val: stats.attendance_today.find(a => a.status === 'Alpa')?.count || 0, color: 'rose' },
                { key: 'Piket', label: 'Piket', val: stats.attendance_today.find(a => a.status === 'Piket')?.count || 0, color: 'purple' }
              ]" :key="item.key" class="p-6 rounded-[2rem] bg-white shadow-sm border border-slate-100 flex flex-col gap-2 hover:shadow-xl transition-all duration-500 group">
                 <div class="flex justify-between items-center">
                    <span class="text-[9px] font-black uppercase text-slate-400 tracking-widest">{{ item.label }}</span>
                    <ChevronRight class="w-3 h-3 text-slate-200 group-hover:text-blue-500 transition-colors" />
                 </div>
                 <div class="text-4xl font-black italic tracking-tighter text-slate-900">{{ item.val }}</div>
                 <div class="h-1 w-full bg-slate-50 rounded-full overflow-hidden mt-2">
                    <div :class="`bg-${item.color}-500`" class="h-full" :style="{ width: `${(item.val / (stats.total_students || 1)) * 100}%` }"></div>
                 </div>
              </div>
           </div>
        </div>

        <!-- Teacher Monitoring -->
        <div class="bg-[#0F2942] rounded-[3rem] p-10 shadow-2xl text-white space-y-8 relative overflow-hidden">
           <div class="absolute right-0 top-0 w-1/2 h-full bg-blue-600/10 blur-[80px]"></div>
           <div class="relative z-10 flex flex-col h-full justify-between">
              <div class="flex items-center justify-between mb-8">
                 <h3 class="text-lg font-black italic">Presensi Guru</h3>
                 <div class="flex gap-2">
                    <div v-for="i in 3" :key="i" class="w-1.5 h-1.5 rounded-full bg-white/20"></div>
                 </div>
              </div>
              
              <div class="space-y-6">
                 <div v-for="t in [
                   { label: 'Hadir', key: 'Hadir', color: 'bg-blue-400' },
                   { label: 'Izin / Sakit', key: 'Izin', color: 'bg-white/20' }
                 ]" :key="t.key" class="space-y-2">
                    <div class="flex justify-between items-end">
                       <span class="text-xs font-bold text-white/60">{{ t.label }}</span>
                       <span class="text-xl font-black italic">{{ stats.teacher_attendance_today.find(s => s.status === t.key)?.count || 0 }} Guru</span>
                    </div>
                    <div class="h-2 bg-white/5 rounded-full overflow-hidden">
                       <div :class="t.color" class="h-full" :style="{ width: `${((stats.teacher_attendance_today.find(s => s.status === t.key)?.count || 0) / (stats.total_teachers || 1)) * 100}%` }"></div>
                    </div>
                 </div>
              </div>

              <div class="mt-10 p-6 rounded-[2rem] bg-white/5 border border-white/10 backdrop-blur-xl flex items-center justify-between">
                 <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-blue-500/20 text-blue-400 flex items-center justify-center">
                       <UserCog class="w-6 h-6" />
                    </div>
                    <div>
                       <p class="text-xs font-black italic leading-none mb-1">Progres Mengajar</p>
                       <p class="text-[9px] font-bold text-white/30 uppercase tracking-widest">{{ stats.journals_today_count }} Jurnal Terinput Hari Ini</p>
                    </div>
                 </div>
                 <ChevronRight class="w-5 h-5 text-white/20" />
              </div>
           </div>
        </div>

      </div>
    </main>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700;900&family=Outfit:wght@400;600;900&display=swap');

.font-sans {
  font-family: 'Outfit', sans-serif;
}

/* Custom Scrollbar for Glass Effect */
.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 20px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

/* Transition Animations */
.slide-fade-enter-active { transition: all 0.6s cubic-bezier(0.2, 0.8, 0.2, 1); }
.slide-fade-leave-active { transition: all 0.5s cubic-bezier(1, 0.5, 0.8, 1); }
.slide-fade-enter-from { transform: translateY(20px); opacity: 0; }
.slide-fade-leave-to { transform: translateY(-20px); opacity: 0; }

@keyframes fade-up {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
}

.animate-reveal {
  animation: fade-up 1s cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
}

.reveal-delay-1 { animation-delay: 0.1s; }
.reveal-delay-2 { animation-delay: 0.3s; }
.reveal-delay-3 { animation-delay: 0.5s; }

/* Utilities */
.tracking-tighter { letter-spacing: -0.05em; }
.tracking-tight { letter-spacing: -0.025em; }

/* Safelist colors */
.bg-orange-500 { background-color: #F97316; }
.bg-blue-500 { background-color: #3B82F6; }
.bg-rose-500 { background-color: #F43F5E; }
.bg-purple-500 { background-color: #A855F7; }
</style>
