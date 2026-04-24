<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  TrendingUp, Users2, GraduationCap, MapPin, Briefcase, UserCheck,
  Sun, Moon, ArrowLeft
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'
import VueApexCharts from 'vue3-apexcharts'
import { ApexOptions } from 'apexcharts'

const router = useRouter()
const isDark = ref(true)
const user = ref({ name: 'Admin', role: 'Staff' })

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


// Stats state
const stats = ref<any>({
  total_santri: 0,
  total_guru: 0,
  total_alumni: 0,
  chart_gender: { L: 0, P: 0 },
  chart_trend: { labels: [], data: [] },
  chart_teacher_status: [],
  chart_alumni_service: { Mengabdi: 0, Lainnya: 0 },
  alumni_regions: []
})

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/v1/admin/base/stats')
    stats.value = response.data
  } catch (err) {
    console.error('Gagal mengambil statistik:', err)
  }
}

// Chart Options: Registration Trend
const trendChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'area', toolbar: { show: false }, background: 'transparent' },
  colors: ['#10b981'],
  stroke: { curve: 'smooth', width: 3 },
  fill: {
    type: 'gradient',
    gradient: { shadeIntensity: 1, opacityFrom: 0.45, opacityTo: 0.05, stops: [20, 100] }
  },
  xaxis: {
    categories: stats.value.chart_trend.labels,
    labels: { style: { colors: '#64748b', fontFamily: 'Inter' } }
  },
  yaxis: { show: false },
  grid: { show: false },
  dataLabels: { enabled: false },
  tooltip: { theme: 'dark' }
}))

const trendChartSeries = computed(() => [{
  name: 'Registrasi',
  data: stats.value.chart_trend.data
}])

// Chart Options: Gender Distribution
const genderChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut', background: 'transparent' },
  labels: ['Laki-laki', 'Perempuan'],
  colors: ['#10b981', '#3b82f6'],
  stroke: { show: false },
  legend: { show: false },
  plotOptions: {
    pie: {
      donut: {
        size: '75%',
        labels: {
          show: true,
          name: { show: false },
          total: { show: true, label: 'Gender', formatter: () => 'SANTRI', color: '#64748b', fontSize: '12px' }
        }
      }
    }
  },
  dataLabels: { enabled: false },
  tooltip: { theme: 'dark' }
}))

const genderChartSeries = computed(() => [
  stats.value.chart_gender.L || 0, 
  stats.value.chart_gender.P || 0
])

// Chart Options: Teacher Status
const teacherChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut', background: 'transparent' },
  labels: stats.value.chart_teacher_status.map((s: any) => s.Status),
  colors: ['#10b981', '#f59e0b', '#3b82f6', '#8b5cf6'],
  stroke: { show: false },
  legend: { show: false },
  plotOptions: {
    pie: {
      donut: {
        size: '75%',
        labels: {
          show: true,
          name: { show: false },
          total: { show: true, label: 'Status', formatter: () => 'GURU', color: '#64748b', fontSize: '12px' }
        }
      }
    }
  },
  dataLabels: { enabled: false },
  tooltip: { theme: 'dark' }
}))

const teacherChartSeries = computed(() => stats.value.chart_teacher_status.map((s: any) => s.Count))

// Database Data
const recentStudents = ref<any[]>([])
const fetchRecentStudents = async () => {
  try {
    const response = await axios.get('/api/v1/admin/base/students')
    recentStudents.value = response.data.slice(0, 5)
  } catch (err) {
    console.error('Gagal mengambil data santri:', err)
  }
}

onMounted(() => {
  const savedUser = localStorage.getItem('sigma_user')
  if (savedUser) user.value = JSON.parse(savedUser)
  
  const savedTheme = localStorage.getItem('sigma_theme')
  isDark.value = savedTheme !== 'light'
  
  fetchStats()
  fetchRecentStudents()
})
</script>


<template>
  <div class="p-8 space-y-10">
    <!-- Header Actions -->
    <div class="flex items-center justify-between mb-2">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-sigma-emerald rounded-full"></div>
        <h2 class="text-sm font-black uppercase tracking-[0.3em] text-sigma-muted">Data Center Module</h2>
      </div>
      
      <div class="flex items-center gap-3">
        <!-- Theme Toggle -->
        <button @click="toggleTheme" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-text hover:border-sigma-emerald/30 transition-all shadow-sm group">
          <div class="w-6 h-6 rounded-lg bg-sigma-emerald/10 flex items-center justify-center text-sigma-emerald">
            <Sun v-if="isDark" class="w-3.5 h-3.5" />
            <Moon v-else class="w-3.5 h-3.5" />
          </div>
          <span class="text-[10px] font-black uppercase tracking-widest">{{ isDark ? 'Light' : 'Dark' }}</span>
        </button>

        <!-- Back to Portal -->
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-sigma-emerald/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Welcome Banner -->
    <div class="relative overflow-hidden rounded-[2rem] bg-gradient-to-r from-emerald-600 to-emerald-900/80 p-8 shadow-2xl shadow-emerald-500/10 border border-sigma-border">
      <div class="absolute top-0 right-0 -translate-y-1/2 translate-x-1/4 w-80 h-80 bg-emerald-400/10 rounded-full blur-[100px]"></div>
      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="space-y-1">
          <span class="inline-block px-3 py-1 rounded-full bg-white/10 backdrop-blur-md text-[10px] font-bold uppercase tracking-[0.2em] text-emerald-100/90 mb-2 border border-sigma-border">
            Sigmabase Intelligence
          </span>
          <h1 class="text-2xl md:text-4xl font-black text-white leading-tight">
            Halo, <span class="text-emerald-300 font-black italic">{{ user.name }}!</span>
          </h1>
          <p class="text-emerald-100/60 text-sm font-medium">Pusat Data Center santri, guru, dan alumni.</p>
        </div>
        
        <div class="flex items-center gap-5 bg-black/20 backdrop-blur-2xl p-5 rounded-3xl border border-sigma-border">
          <div class="text-right">
            <p class="text-[10px] uppercase tracking-widest text-emerald-400 font-black">System Status</p>
            <p class="text-sm font-bold text-white">Operational</p>
          </div>
          <div class="relative flex h-3 w-3">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Stats Row -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div @click="router.push('/sigmabase/students')" class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group cursor-pointer shadow-sm">
        <div class="w-14 h-14 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-400 group-hover:bg-emerald-500 group-hover:text-white transition-all duration-500">
          <Users2 class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-bold text-sigma-muted uppercase tracking-widest">Mahasantri Aktif</p>
          <h3 class="text-3xl font-black mt-1">{{ stats.total_santri }}</h3>
        </div>
      </div>

      <div @click="router.push('/sigmabase/teachers')" class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group cursor-pointer shadow-sm">
        <div class="w-14 h-14 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-400 group-hover:bg-blue-500 group-hover:text-white transition-all duration-500">
          <Briefcase class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-bold text-sigma-muted uppercase tracking-widest">Total Pendidik</p>
          <h3 class="text-3xl font-black mt-1">{{ stats.total_guru }}</h3>
        </div>
      </div>

      <div @click="router.push('/sigmabase/alumni')" class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group cursor-pointer shadow-sm">
        <div class="w-14 h-14 bg-orange-500/10 rounded-2xl flex items-center justify-center text-orange-400 group-hover:bg-orange-500 group-hover:text-white transition-all duration-500">
          <GraduationCap class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-bold text-sigma-muted uppercase tracking-widest">Database Alumni</p>
          <h3 class="text-3xl font-black mt-1">{{ stats.total_alumni }}</h3>
        </div>
      </div>
    </div>

    <!-- Secondary Stats & Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Alumni Regions Summary -->
      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-6 shadow-sm">
        <div class="flex items-center justify-between">
          <h3 class="text-xl font-bold flex items-center gap-3">
            <MapPin class="w-5 h-5 text-sigma-emerald" /> Sebaran Wilayah
          </h3>
          <span class="text-[10px] text-sigma-muted font-black uppercase tracking-widest">Top 5 Alumni</span>
        </div>
        <div class="space-y-4">
          <div v-for="(region, idx) in stats.alumni_regions" :key="idx" class="flex items-center gap-4 group">
            <div class="w-8 h-8 rounded-lg bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-[10px] font-black group-hover:bg-sigma-emerald group-hover:text-white transition-all">
              {{ Number(idx) + 1 }}
            </div>
            <div class="flex-1">
              <div class="flex justify-between text-xs font-bold mb-1">
                <span>{{ region.name }}</span>
                <span class="text-sigma-emerald">{{ region.value }}</span>
              </div>
              <div class="h-1 w-full bg-sigma-surface-alt rounded-full overflow-hidden">
                <div class="h-full bg-sigma-emerald transition-all duration-1000" :style="{ width: (region.value / (stats.total_alumni || 1) * 100) + '%' }"></div>
              </div>
            </div>
          </div>
          <div v-if="stats.alumni_regions.length === 0" class="py-10 text-center text-sigma-muted text-xs italic">Belum ada data wilayah alumni.</div>
        </div>
      </div>

      <!-- Trend Chart -->
      <div class="lg:col-span-2 bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-6 shadow-sm">
        <div class="flex items-center justify-between">
          <h3 class="text-xl font-bold flex items-center gap-3">
            <TrendingUp class="w-5 h-5 text-sigma-emerald" /> Trend Registrasi
          </h3>
        </div>
        <div class="h-[250px] w-full">
          <VueApexCharts width="100%" height="250" :options="trendChartOptions" :series="trendChartSeries" />
        </div>
      </div>

      <!-- Donut Charts Row -->
      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-8 shadow-sm">
        <h3 class="text-xl font-bold flex items-center gap-3">
          <UserCheck class="w-5 h-5 text-sigma-emerald" /> Komposisi Santri
        </h3>
        <div class="relative w-full flex justify-center">
          <VueApexCharts type="donut" width="220" :options="genderChartOptions" :series="genderChartSeries" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-sigma-surface-alt p-4 rounded-2xl text-center border border-sigma-border/50">
            <p class="text-[9px] text-sigma-muted font-bold uppercase mb-1">Laki-laki</p>
            <p class="text-xl font-black text-emerald-400">{{ stats.chart_gender.L }}</p>
          </div>
          <div class="bg-sigma-surface-alt p-4 rounded-2xl text-center border border-sigma-border/50">
            <p class="text-[9px] text-sigma-muted font-bold uppercase mb-1">Perempuan</p>
            <p class="text-xl font-black text-blue-400">{{ stats.chart_gender.P }}</p>
          </div>
        </div>
      </div>

      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-8 shadow-sm">
        <h3 class="text-xl font-bold flex items-center gap-3">
          <Briefcase class="w-5 h-5 text-sigma-emerald" /> Status Pendidik
        </h3>
        <div class="relative w-full flex justify-center">
          <VueApexCharts type="donut" width="220" :options="teacherChartOptions" :series="teacherChartSeries" />
        </div>
        <div class="flex flex-wrap gap-2 justify-center">
          <span v-for="(s, idx) in stats.chart_teacher_status" :key="idx" class="px-3 py-1 rounded-full bg-sigma-surface-alt border border-sigma-border text-[9px] font-black uppercase tracking-widest text-sigma-muted">
            {{ s.Status }}: {{ s.Count }}
          </span>
        </div>
      </div>

      <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-8 shadow-sm flex flex-col justify-between">
        <h3 class="text-xl font-bold flex items-center gap-3">
          <GraduationCap class="w-5 h-5 text-sigma-emerald" /> Status Alumni
        </h3>
        <div class="space-y-6">
           <div class="p-6 bg-gradient-to-br from-emerald-500/20 to-transparent rounded-3xl border border-emerald-500/20 text-center">
              <p class="text-[10px] font-black uppercase tracking-[0.3em] text-emerald-400 mb-2">Total Mengabdi</p>
              <p class="text-5xl font-black text-white italic">{{ stats.chart_alumni_service.Mengabdi }}</p>
           </div>
           <div class="p-6 bg-sigma-surface-alt rounded-3xl border border-sigma-border text-center">
              <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted mb-2">Lainnya / Kuliah</p>
              <p class="text-3xl font-black text-sigma-text">{{ stats.chart_alumni_service.Lainnya }}</p>
           </div>
        </div>
      </div>

      <!-- Recent Records -->
      <div class="lg:col-span-3 bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
        <div class="p-8 pb-4 flex justify-between items-center bg-gradient-to-r from-sigma-surface-alt/50 to-transparent">
          <h3 class="text-xl font-bold italic tracking-tight uppercase">REGISTRASI <span class="text-sigma-emerald not-italic">TERBARU</span></h3>
          <button @click="router.push('/sigmabase/students')" class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-emerald hover:text-emerald-300 transition-all flex items-center gap-2">
            View All Database <TrendingUp class="w-4 h-4" />
          </button>
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.3em] border-b border-sigma-border bg-sigma-surface-alt/20">
                <th class="p-8">Profil</th>
                <th class="p-8">Identitas</th>
                <th class="p-8">Klasifikasi</th>
                <th class="p-8 text-right">Timestamp</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-sigma-border">
              <tr v-for="student in recentStudents" :key="student.ID" class="hover:bg-sigma-surface-alt transition-all group">
                <td class="p-8">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center font-black text-sigma-emerald">
                      {{ student.name.charAt(0) }}
                    </div>
                    <div class="font-black text-sigma-text text-sm uppercase tracking-tight group-hover:text-sigma-emerald transition-colors">{{ student.name }}</div>
                  </div>
                </td>
                <td class="p-8">
                  <div class="text-[10px] font-mono text-sigma-muted">NIS: {{ student.nis }}</div>
                  <div class="text-[10px] font-bold text-sigma-text/60 mt-1 uppercase">{{ student.gender === 'L' ? 'Laki-laki' : 'Perempuan' }}</div>
                </td>
                <td class="p-8">
                  <span class="px-3 py-1 rounded-lg bg-sigma-surface-alt border border-sigma-border text-[9px] font-black uppercase tracking-widest text-sigma-muted group-hover:text-sigma-text transition-colors">
                    {{ student.class }} / {{ student.entry_year }}
                  </span>
                </td>
                <td class="p-8 text-right text-[10px] font-bold text-sigma-muted uppercase tracking-widest">
                  Active
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.05); border-radius: 10px; }
</style>
