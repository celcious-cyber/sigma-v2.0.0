<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Activity, Users, Bed, Pill, 
  ArrowUpRight, 
  AlertCircle,
  TrendingUp, Clock, Plus,
  FileText, Microscope
} from 'lucide-vue-next'
import axios from 'axios'

const stats = ref<any>({
  today_visits: 0,
  total_visits: 0,
  inpatients: 0,
  low_stock_count: 0,
  bed_capacity: 12,
  top_diseases: []
})

const loading = ref(true)

const fetchStats = async () => {
  try {
    const response = await axios.get('/admin/care/stats')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div class="p-8 space-y-10 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="relative overflow-hidden bg-white border border-rose-100 rounded-[2.5rem] p-10">
      <div class="absolute top-0 right-0 p-12 opacity-10">
        <Activity class="w-40 h-40 text-rose-600 rotate-12" />
      </div>
      
      <div class="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-10">
        <div class="flex items-center gap-6">
          <div class="w-20 h-20 bg-rose-600 rounded-[2rem] flex items-center justify-center text-white shadow-2xl shadow-rose-500/20">
            <Activity class="w-10 h-10" />
          </div>
          <div>
            <h1 class="text-4xl font-black tracking-tight text-slate-800 italic uppercase">Dasbor <span class="text-rose-600">Klinis</span></h1>
            <p class="text-slate-500 font-bold mt-1 tracking-wide flex items-center gap-2">
              <Clock class="w-4 h-4" /> MONITORING KESEHATAN SANTRI REAL-TIME
            </p>
          </div>
        </div>

        <div class="flex flex-wrap gap-4">
          <button class="px-8 py-4 bg-rose-600 text-white rounded-2xl font-black text-sm tracking-widest uppercase hover:bg-rose-700 transition-all shadow-xl shadow-rose-500/20 flex items-center gap-3 active:scale-95">
            <Plus class="w-5 h-5" /> Kunjungan Baru
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Today Visits -->
      <div class="group bg-white p-8 rounded-[2rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-rose-500/5 transition-all relative overflow-hidden">
        <div class="absolute -right-4 -bottom-4 w-24 h-24 bg-rose-50 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-700"></div>
        <div class="relative z-10 flex flex-col gap-4">
          <div class="w-14 h-14 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center">
            <Users class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-black text-slate-400 uppercase tracking-widest mb-1">Kunjungan Hari Ini</p>
            <div class="flex items-end gap-3">
              <h2 class="text-4xl font-black text-slate-800 tracking-tighter">{{ stats.today_visits }}</h2>
              <div class="flex items-center text-emerald-500 font-black text-sm mb-1">
                <ArrowUpRight class="w-4 h-4" /> 12%
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Inpatients -->
      <div class="group bg-white p-8 rounded-[2rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-emerald-500/5 transition-all relative overflow-hidden">
        <div class="absolute -right-4 -bottom-4 w-24 h-24 bg-emerald-50 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-700"></div>
        <div class="relative z-10 flex flex-col gap-4">
          <div class="w-14 h-14 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center">
            <Bed class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-black text-slate-400 uppercase tracking-widest mb-1">Rawat Inap/Observasi</p>
            <div class="flex items-end gap-3">
              <h2 class="text-4xl font-black text-slate-800 tracking-tighter">{{ stats.inpatients }}</h2>
              <span class="text-slate-400 font-bold text-sm mb-1 italic">Bed Terisi</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Pharmacy Alerts -->
      <div class="group bg-white p-8 rounded-[2rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-amber-500/5 transition-all relative overflow-hidden">
        <div class="absolute -right-4 -bottom-4 w-24 h-24 bg-amber-50 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-700"></div>
        <div class="relative z-10 flex flex-col gap-4">
          <div class="w-14 h-14 bg-amber-50 text-amber-600 rounded-2xl flex items-center justify-center">
            <Pill class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-black text-slate-400 uppercase tracking-widest mb-1">Peringatan Stok</p>
            <div class="flex items-end gap-3">
              <h2 class="text-4xl font-black text-slate-800 tracking-tighter">{{ stats.low_stock_count }}</h2>
              <span class="text-amber-600 font-black text-sm mb-1 italic">Item Kritis</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Capacity -->
      <div class="group bg-slate-800 p-8 rounded-[2rem] shadow-2xl relative overflow-hidden">
        <div class="absolute top-0 right-0 w-32 h-32 bg-white/5 rounded-full -mr-16 -mt-16 group-hover:scale-110 transition-transform"></div>
        <div class="relative z-10 flex flex-col gap-4">
          <div class="w-14 h-14 bg-white/10 text-white rounded-2xl flex items-center justify-center">
            <Activity class="w-7 h-7" />
          </div>
          <div>
            <p class="text-xs font-black text-white/50 uppercase tracking-widest mb-1">Kapasitas UKS</p>
            <div class="flex items-end gap-3">
              <h2 class="text-4xl font-black text-white tracking-tighter">{{ Math.round((stats.inpatients / stats.bed_capacity) * 100) }}%</h2>
              <div class="flex-1 h-2 bg-white/10 rounded-full mb-3 ml-2 relative overflow-hidden">
                <div 
                  class="absolute left-0 top-0 h-full bg-rose-500 rounded-full"
                  :style="{ width: `${(stats.inpatients / stats.bed_capacity) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Epidemiologi (Top Diseases) -->
      <div class="lg:col-span-2 bg-white rounded-[2.5rem] border border-slate-100 p-10">
        <div class="flex items-center justify-between mb-10">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center">
              <TrendingUp class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-xl font-black text-slate-800 uppercase italic">Status Epidemiologi</h3>
              <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mt-0.5">Top 5 Keluhan Penyakit</p>
            </div>
          </div>
          <button class="text-rose-600 font-black text-xs uppercase tracking-widest hover:underline">Lihat Detail</button>
        </div>

        <div class="space-y-6">
          <div v-for="(disease, index) in stats.top_diseases" :key="index" class="space-y-2">
            <div class="flex justify-between items-center text-sm">
              <span class="font-black text-slate-700 flex items-center gap-3 uppercase tracking-tighter">
                <span class="w-6 h-6 bg-slate-100 rounded-lg flex items-center justify-center text-[10px]">{{ Number(index) + 1 }}</span>
                {{ disease.name || 'Tidak Terdiagnosis' }}
              </span>
              <span class="font-black text-rose-600 italic">{{ disease.count }} Pasien</span>
            </div>
            <div class="h-3 bg-slate-50 rounded-full relative overflow-hidden">
              <div 
                class="absolute left-0 top-0 h-full bg-rose-500 rounded-full transition-all duration-1000"
                :style="{ width: `${(disease.count / stats.total_visits) * 100}%` }"
              ></div>
            </div>
          </div>

          <div v-if="!stats.top_diseases.length" class="py-20 text-center">
            <div class="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4">
              <Activity class="w-10 h-10 text-slate-200" />
            </div>
            <p class="text-slate-400 font-bold italic">Belum ada data kunjungan hari ini.</p>
          </div>
        </div>
      </div>

      <!-- Quick Actions & Alerts -->
      <div class="space-y-8">
        <div class="bg-rose-600 rounded-[2.5rem] p-8 text-white shadow-2xl shadow-rose-500/20 relative overflow-hidden">
          <div class="absolute -right-10 -bottom-10 opacity-10">
            <AlertCircle class="w-40 h-40" />
          </div>
          <div class="relative z-10 flex flex-col gap-6">
            <h3 class="text-xl font-black italic uppercase tracking-tight">Butuh Perhatian!</h3>
            <div class="space-y-4">
              <div class="flex items-start gap-4 p-4 bg-white/10 rounded-2xl backdrop-blur-md">
                <AlertCircle class="w-6 h-6 mt-1" />
                <div>
                  <p class="text-sm font-black uppercase tracking-tight">Stok Obat Menipis</p>
                  <p class="text-xs font-bold text-white/70 mt-1 uppercase tracking-widest leading-relaxed">Paracetamol 500mg sisa 4 tablet.</p>
                </div>
              </div>
              <div class="flex items-start gap-4 p-4 bg-white/10 rounded-2xl backdrop-blur-md">
                <Users class="w-6 h-6 mt-1" />
                <div>
                  <p class="text-sm font-black uppercase tracking-tight">Jadwal MCU Masal</p>
                  <p class="text-xs font-bold text-white/70 mt-1 uppercase tracking-widest leading-relaxed">Tingkat 1 & 2 - Besok 08:00 WIB.</p>
                </div>
              </div>
            </div>
            <button class="w-full py-4 bg-white text-rose-600 rounded-2xl font-black text-xs uppercase tracking-[0.2em] shadow-xl hover:bg-slate-50 transition-all">
              Buka Manajemen Stok
            </button>
          </div>
        </div>

        <div class="bg-white rounded-[2.5rem] border border-slate-100 p-8">
          <h3 class="text-lg font-black text-slate-800 uppercase italic mb-6">Pintasan Cepat</h3>
          <div class="grid grid-cols-2 gap-4">
            <button class="p-6 bg-slate-50 rounded-2xl hover:bg-rose-50 hover:text-rose-600 transition-all group flex flex-col gap-3">
              <FileText class="w-6 h-6" />
              <span class="text-[10px] font-black uppercase tracking-widest text-left">Buat Surat Sakit</span>
            </button>
            <button class="p-6 bg-slate-50 rounded-2xl hover:bg-rose-50 hover:text-rose-600 transition-all group flex flex-col gap-3">
              <Microscope class="w-6 h-6" />
              <span class="text-[10px] font-black uppercase tracking-widest text-left">Input Hasil MCU</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes pulse-soft {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
.animate-pulse-soft {
  animation: pulse-soft 3s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>
