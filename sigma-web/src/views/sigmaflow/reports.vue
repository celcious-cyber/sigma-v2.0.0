<script setup>
import { ref, onMounted, computed } from 'vue'
import { 
  BarChart3, 
  TrendingUp, 
  TrendingDown, 
  Users, 
  FileText, 
  Download,
  Calendar,
  ChevronRight,
  ArrowUpRight,
  ArrowDownLeft,
  PieChart,
  Search,
  Filter
} from 'lucide-vue-next'
import axios from 'axios'
import { toast } from 'vue3-toastify'

// State
const activeTab = ref('profit_loss')
const loading = ref(false)
const plReport = ref({
  income_categories: {},
  expense_categories: {},
  total_income: 0,
  income_formatted: 'Rp 0',
  total_expense: 0,
  expense_formatted: 'Rp 0',
  net_profit: 0,
  profit_formatted: 'Rp 0',
  month: '',
  year: 2026
})
const arrearsReport = ref([])
const searchArrears = ref('')

const months = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]

const selectedMonth = ref(new Date().getMonth() + 1)
const selectedYear = ref(new Date().getFullYear())

const fetchPLReport = async () => {
  loading.value = true
  try {
    const res = await axios.get(`/api/v1/admin/flow/reports/profit-loss?month=${selectedMonth.value}&year=${selectedYear.value}`)
    plReport.value = res.data
  } catch (err) {
    toast.error('Gagal mengambil laporan laba rugi')
  } finally {
    loading.value = false
  }
}

const fetchArrearsReport = async () => {
  loading.value = true
  try {
    const res = await axios.get('/admin/flow/reports/arrears')
    arrearsReport.value = res.data
  } catch (err) {
    toast.error('Gagal mengambil laporan tunggakan')
  } finally {
    loading.value = false
  }
}

const filteredArrears = computed(() => {
  if (!searchArrears.value) return arrearsReport.value
  const s = searchArrears.value.toLowerCase()
  return arrearsReport.value.filter(a => 
    a.student_name.toLowerCase().includes(s) || 
    a.nis.includes(s) || 
    a.category.toLowerCase().includes(s)
  )
})

const formatCurrency = (val) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(val)
}

onMounted(() => {
  fetchPLReport()
})

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'profit_loss') fetchPLReport()
  else if (tab === 'arrears') fetchArrearsReport()
}
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1600px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- Header -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-5">
          <div class="p-4 bg-indigo-600 rounded-[2rem] text-white shadow-2xl shadow-indigo-500/20">
            <BarChart3 class="w-8 h-8" />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-3xl font-black italic tracking-tighter uppercase text-slate-900">Laporan <span class="text-indigo-600">Keuangan</span></h1>
              <span class="px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded-full text-[9px] font-black uppercase tracking-widest text-indigo-600">Finance Analytics</span>
            </div>
            <p class="text-slate-500 font-medium text-sm">Analisis mendalam performa finansial yayasan</p>
          </div>
        </div>

        <!-- Tab Switcher -->
        <div class="bg-white p-1.5 rounded-2xl border border-slate-100 shadow-sm flex items-center gap-1">
          <button 
            @click="switchTab('profit_loss')"
            :class="activeTab === 'profit_loss' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-slate-600'"
            class="px-6 py-2.5 rounded-xl font-bold text-xs uppercase tracking-widest transition-all"
          >
            Laba Rugi
          </button>
          <button 
            @click="switchTab('arrears')"
            :class="activeTab === 'arrears' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-slate-600'"
            class="px-6 py-2.5 rounded-xl font-bold text-xs uppercase tracking-widest transition-all"
          >
            Tunggakan
          </button>
        </div>
      </section>

      <!-- Tab: Profit & Loss -->
      <div v-if="activeTab === 'profit_loss'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <!-- P&L Filters -->
        <section class="flex flex-wrap items-center gap-4 bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm">
          <div class="flex items-center gap-2 px-4 py-2.5 bg-slate-50 rounded-xl">
            <Calendar class="w-4 h-4 text-slate-400" />
            <select v-model="selectedMonth" class="bg-transparent border-none text-xs font-bold text-slate-700 outline-none">
              <option v-for="(m, i) in months" :key="m" :value="i + 1">{{ m }}</option>
            </select>
            <select v-model="selectedYear" class="bg-transparent border-none text-xs font-bold text-slate-700 outline-none">
              <option :value="2024">2024</option>
              <option :value="2025">2025</option>
              <option :value="2026">2026</option>
            </select>
          </div>
          <button 
            @click="fetchPLReport"
            class="px-6 py-3 bg-slate-900 text-white rounded-xl font-bold text-xs uppercase tracking-widest hover:bg-slate-800 transition-all"
          >
            Tampilkan Laporan
          </button>
        </section>

        <!-- Summary Widgets -->
        <section class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
            <p class="text-emerald-500 font-black uppercase tracking-widest text-[10px] mb-2">Total Pendapatan</p>
            <h2 class="text-3xl font-black text-slate-900">{{ plReport.income_formatted }}</h2>
            <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:scale-110 transition-transform duration-700">
              <TrendingUp class="w-32 h-32 text-emerald-600" />
            </div>
          </div>
          <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
            <p class="text-rose-500 font-black uppercase tracking-widest text-[10px] mb-2">Total Pengeluaran</p>
            <h2 class="text-3xl font-black text-slate-900">{{ plReport.expense_formatted }}</h2>
            <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:scale-110 transition-transform duration-700">
              <TrendingDown class="w-32 h-32 text-rose-600" />
            </div>
          </div>
          <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
            <p class="text-indigo-500 font-black uppercase tracking-widest text-[10px] mb-2">Laba / Rugi Bersih</p>
            <h2 class="text-3xl font-black text-slate-900">{{ plReport.profit_formatted }}</h2>
            <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:scale-110 transition-transform duration-700">
              <PieChart class="w-32 h-32 text-indigo-600" />
            </div>
          </div>
        </section>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Income Breakdown -->
          <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
            <div class="p-8 border-b border-slate-50 flex items-center justify-between bg-emerald-50/30">
              <h3 class="font-black text-slate-900 uppercase tracking-tighter italic flex items-center gap-2">
                <ArrowUpRight class="w-5 h-5 text-emerald-600" />
                Detail Pendapatan
              </h3>
            </div>
            <div class="p-8 space-y-4">
              <div v-for="(amount, cat) in plReport.income_categories" :key="cat" class="flex items-center justify-between p-4 bg-slate-50 rounded-2xl">
                <span class="text-sm font-bold text-slate-700">{{ cat }}</span>
                <span class="text-sm font-black text-emerald-600">{{ formatCurrency(amount) }}</span>
              </div>
              <div v-if="Object.keys(plReport.income_categories).length === 0" class="py-10 text-center text-slate-400 font-medium italic">
                Tidak ada data pendapatan bulan ini
              </div>
            </div>
          </div>

          <!-- Expense Breakdown -->
          <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
            <div class="p-8 border-b border-slate-50 flex items-center justify-between bg-rose-50/30">
              <h3 class="font-black text-slate-900 uppercase tracking-tighter italic flex items-center gap-2">
                <ArrowDownLeft class="w-5 h-5 text-rose-600" />
                Detail Pengeluaran
              </h3>
            </div>
            <div class="p-8 space-y-4">
              <div v-for="(amount, cat) in plReport.expense_categories" :key="cat" class="flex items-center justify-between p-4 bg-slate-50 rounded-2xl">
                <span class="text-sm font-bold text-slate-700">{{ cat }}</span>
                <span class="text-sm font-black text-rose-600">{{ formatCurrency(amount) }}</span>
              </div>
              <div v-if="Object.keys(plReport.expense_categories).length === 0" class="py-10 text-center text-slate-400 font-medium italic">
                Tidak ada data pengeluaran bulan ini
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab: Arrears -->
      <div v-if="activeTab === 'arrears'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <section class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
          <div class="p-8 border-b border-slate-50 flex flex-col md:flex-row justify-between items-center gap-6">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-rose-50 text-rose-600 rounded-xl">
                <Users class="w-5 h-5" />
              </div>
              <h3 class="font-black text-slate-900 uppercase tracking-tighter italic">Daftar Tunggakan Santri</h3>
            </div>
            <div class="relative w-full md:w-80">
              <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input 
                v-model="searchArrears"
                type="text" 
                placeholder="Cari nama santri / NIS..."
                class="w-full pl-11 pr-5 py-3 bg-slate-50 border-none rounded-xl text-xs font-bold text-slate-600 focus:ring-2 focus:ring-rose-500 transition-all outline-none"
              />
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left">
              <thead>
                <tr class="bg-slate-50/50">
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Santri</th>
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Kategori Tagihan</th>
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Jatuh Tempo</th>
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Total Tagihan</th>
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Sudah Dibayar</th>
                  <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Sisa (Tunggakan)</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-50">
                <tr v-for="a in filteredArrears" :key="a.id" class="hover:bg-slate-50/30 transition-colors">
                  <td class="px-8 py-6">
                    <span class="text-xs font-bold text-slate-900 block mb-0.5">{{ a.student_name }}</span>
                    <span class="text-[10px] text-slate-400 font-medium uppercase tracking-tight">{{ a.nis }}</span>
                  </td>
                  <td class="px-8 py-6">
                    <span class="text-xs font-bold text-slate-700">{{ a.category }}</span>
                  </td>
                  <td class="px-8 py-6">
                    <span class="text-xs font-bold text-slate-700">{{ a.due_date }}</span>
                  </td>
                  <td class="px-8 py-6 text-right">
                    <span class="text-xs font-bold text-slate-900">{{ a.total_formatted }}</span>
                  </td>
                  <td class="px-8 py-6 text-right">
                    <span class="text-xs font-bold text-emerald-600">{{ a.paid_formatted }}</span>
                  </td>
                  <td class="px-8 py-6 text-right">
                    <span class="text-sm font-black text-rose-600 italic tracking-tighter">{{ a.arrears_formatted }}</span>
                  </td>
                </tr>
                <tr v-if="filteredArrears.length === 0">
                  <td colspan="6" class="px-8 py-20 text-center">
                    <div class="flex flex-col items-center gap-3 text-slate-300">
                      <Users class="w-12 h-12" />
                      <p class="font-bold text-sm">Tidak ada tunggakan ditemukan</p>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>

    </main>
  </div>
</template>

<style scoped>
.animate-reveal {
  animation: reveal 0.8s cubic-bezier(0, 0, 0.2, 1) forwards;
}

@keyframes reveal {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
