<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Filter, Edit, TrendingUp, Users2, ShieldAlert
} from 'lucide-vue-next'
import SigmabaseSidebar from '../../components/SigmabaseSidebar.vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import VueApexCharts from 'vue3-apexcharts'
import { ApexOptions } from 'apexcharts'

const router = useRouter()
const user = ref({ name: 'Admin', role: 'Staff' })

// Stats state
const stats = ref<any>({
  total_santri: 0,
  santri_aktif: 0,
  perizinan_hari_ini: 0,
  chart_gender: { L: 0, P: 0 },
  chart_trend: { labels: [], data: [] }
})

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/v1/admin/base/stats')
    stats.value = response.data
  } catch (err) {
    console.error('Gagal mengambil statistik:', err)
  }
}

// Chart Options
const trendChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'area' as const,
    toolbar: { show: false },
    sparkline: { enabled: false },
    background: 'transparent'
  },
  colors: ['#10b981'],
  stroke: { curve: 'smooth' as const, width: 3 },
  fill: {
    type: 'gradient' as const,
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.45,
      opacityTo: 0.05,
      stops: [20, 100, 100, 100]
    }
  },
  xaxis: {
    categories: stats.value.chart_trend.labels.length > 0 ? stats.value.chart_trend.labels : ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
    axisBorder: { show: false },
    axisTicks: { show: false },
    labels: { style: { colors: '#64748b', fontFamily: 'Inter' } }
  },
  yaxis: { show: false },
  grid: { show: false },
  dataLabels: { enabled: false },
  tooltip: { theme: 'dark' }
}))

const trendChartSeries = computed(() => [{
  name: 'Registrasi',
  data: stats.value.chart_trend.data.length > 0 ? stats.value.chart_trend.data : [10, 25, 45, 30, 60, 85]
}])

const genderChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut' as const, background: 'transparent' },
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
          total: {
            show: true,
            label: 'Gender',
            formatter: () => 'SANTRI',
            color: '#64748b',
            fontSize: '12px'
          }
        }
      }
    }
  },
  dataLabels: { enabled: false },
  tooltip: { theme: 'dark' }
}))

const genderChartSeries = computed(() => [
  stats.value.chart_gender.L || 50, 
  stats.value.chart_gender.P || 50
])

// Database Data
const students = ref<any[]>([])
const fetchStudents = async () => {
  try {
    const response = await axios.get('/api/v1/admin/base/students')
    students.value = response.data.slice(0, 5)
  } catch (err) {
    console.error('Gagal mengambil data santri:', err)
  }
}

onMounted(() => {
  const savedUser = localStorage.getItem('sigma_user')
  if (savedUser) user.value = JSON.parse(savedUser)
  fetchStats()
  fetchStudents()
})


</script>

<template>
  <div class="flex h-screen bg-sigma-app text-sigma-text overflow-hidden font-sans transition-colors duration-300">
    
    <!-- Sidebar -->
    <SigmabaseSidebar activeItem="Dashboard" />

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto p-8 space-y-10 custom-scrollbar">
      
      <!-- Welcome Banner (Slim Version) -->
      <div class="relative overflow-hidden rounded-[2rem] bg-gradient-to-r from-emerald-600 to-emerald-900/80 p-8 shadow-2xl shadow-emerald-500/10 border border-sigma-border">
        <div class="absolute top-0 right-0 -translate-y-1/2 translate-x-1/4 w-80 h-80 bg-emerald-400/10 rounded-full blur-[100px]"></div>
        <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="space-y-1">
            <span class="inline-block px-3 py-1 rounded-full bg-white/10 backdrop-blur-md text-[10px] font-bold uppercase tracking-[0.2em] text-emerald-100/90 mb-2 border border-sigma-border">
              System Overview
            </span>
            <h1 class="text-2xl md:text-4xl font-black text-white leading-tight">
              Halo, <span class="text-emerald-300 font-black italic">{{ user.name }}!</span>
            </h1>
            <p class="text-emerald-100/60 text-sm font-medium">
              Akses cepat intelijen data santri Al-Hikmah.
            </p>
          </div>
          
          <div class="flex items-center gap-5 bg-black/20 backdrop-blur-2xl p-5 rounded-3xl border border-sigma-border">
            <div class="text-right">
              <p class="text-[10px] uppercase tracking-widest text-emerald-400 font-black">Database Connection</p>
              <p class="text-sm font-bold text-white">Cloud Active</p>
            </div>
            <div class="relative flex h-3 w-3">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
            </div>
          </div>
        </div>
      </div>

      <!-- Stats Grid (Minimalist) -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group shadow-sm">
          <div class="w-14 h-14 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-400 group-hover:bg-emerald-500 group-hover:text-white transition-all duration-500">
            <Users2 class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-bold text-sigma-muted uppercase tracking-widest">Total Santri</p>
            <h3 class="text-3xl font-black mt-1">{{ stats.total_santri }}</h3>
          </div>
        </div>

        <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group shadow-sm">
          <div class="w-14 h-14 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-400 group-hover:bg-blue-500 group-hover:text-white transition-all duration-500">
            <TrendingUp class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-bold text-sigma-muted uppercase tracking-widest">Growth</p>
            <h3 class="text-3xl font-black mt-1">+12%</h3>
          </div>
        </div>

        <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] hover:bg-sigma-surface-alt transition-all flex items-center gap-5 group shadow-sm">
          <div class="w-14 h-14 bg-orange-500/10 rounded-2xl flex items-center justify-center text-orange-400 group-hover:bg-orange-500 group-hover:text-white transition-all duration-500">
            <ShieldAlert class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-bold text-sigma-muted uppercase tracking-widest">Perizinan</p>
            <h3 class="text-3xl font-black mt-1">{{ stats.perizinan_hari_ini }}</h3>
          </div>
        </div>
      </div>

      <!-- Charts & Content Section -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <!-- Registration Trend Chart -->
        <div class="lg:col-span-2 bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] space-y-6 shadow-sm">
          <div class="flex items-center justify-between">
            <h3 class="text-xl font-bold flex items-center gap-3">
              Trend Pendaftaran
              <span class="text-[10px] bg-emerald-500/10 text-emerald-400 px-2 py-0.5 rounded-full">+8% m/m</span>
            </h3>
            <button class="text-slate-500 hover:text-white transition-colors"><Filter class="w-5 h-5" /></button>
          </div>
          <div class="h-[250px] w-full">
            <VueApexCharts width="100%" height="250" :options="trendChartOptions" :series="trendChartSeries" />
          </div>
        </div>

        <!-- Gender Distribution Chart -->
        <div class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] flex flex-col items-center justify-center space-y-8 shadow-sm">
          <h3 class="text-xl font-bold w-full text-left">Komposisi Gender</h3>
          <div class="relative w-full flex justify-center">
            <VueApexCharts type="donut" width="220" :options="genderChartOptions" :series="genderChartSeries" />
          </div>
          <div class="grid grid-cols-2 gap-4 w-full">
            <div class="bg-sigma-surface-alt p-4 rounded-2xl text-center">
              <p class="text-[10px] text-sigma-muted font-bold uppercase mb-1">Laki-laki</p>
              <p class="text-xl font-black text-emerald-400">{{ stats.chart_gender.L || 0 }}</p>
            </div>
            <div class="bg-sigma-surface-alt p-4 rounded-2xl text-center">
              <p class="text-[10px] text-sigma-muted font-bold uppercase mb-1">Perempuan</p>
              <p class="text-xl font-black text-blue-400">{{ stats.chart_gender.P || 0 }}</p>
            </div>
          </div>
        </div>

        <!-- Table Recent Students -->
        <div class="lg:col-span-3 bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
          <div class="p-8 pb-4 flex justify-between items-center">
            <h3 class="text-xl font-bold">Mahasantri Terbaru</h3>
            <button @click="router.push('/sigmabase/students')" class="text-sm font-bold text-emerald-400 hover:text-emerald-300 transition-all flex items-center gap-2">
              Lihat Semua <TrendingUp class="w-4 h-4" />
            </button>
          </div>
          
          <div class="overflow-x-auto p-4">
            <table class="w-full text-left">
              <thead>
                <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.2em] border-b border-sigma-border">
                  <th class="p-6">Identitas</th>
                  <th class="p-6">Prodi / Kelas / Angkatan</th>
                  <th class="p-6">Domisili</th>
                  <th class="p-6">Status</th>
                  <th class="p-6 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border">
                <tr v-for="student in students" :key="student.id" class="hover:bg-sigma-surface-alt transition-colors group">
                  <td class="p-6">
                    <div class="flex items-center gap-4">
                      <div class="w-10 h-10 rounded-full bg-sigma-surface-alt border border-sigma-border flex items-center justify-center font-bold text-sigma-muted">
                        {{ student.name.charAt(0) }}
                      </div>
                      <div>
                        <div class="font-bold text-sigma-text group-hover:text-emerald-400 transition-colors">{{ student.name }}</div>
                        <div class="text-[10px] font-mono text-sigma-muted">{{ student.nis }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="p-6 text-sm text-sigma-muted">Pondok / {{ student.class }} / {{ student.entry_year }}</td>
                  <td class="p-6 text-sm text-sigma-muted">{{ student.birth_place || 'Jawa Timur' }}</td>
                  <td class="p-6">
                    <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-black uppercase">
                      <span class="w-1 h-1 rounded-full bg-emerald-500"></span> Aktif
                    </span>
                  </td>
                  <td class="p-6 text-right">
                    <button class="p-2.5 hover:bg-sigma-surface-alt text-sigma-muted hover:text-sigma-text rounded-xl transition-all"><Edit class="w-4 h-4" /></button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>

    </main>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>