<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  History, Search, Filter, 
  CheckCircle2, Clock, MoreVertical,
  Download, Calendar, CreditCard
} from 'lucide-vue-next'
import axios from 'axios'

const payments = ref<any[]>([])
const isLoading = ref(true)
const searchQuery = ref('')

const stats = ref({
  total_today: 0,
  total_today_formatted: 'Rp 0',
  success_rate: '0%',
  pending_count: 0
})

const fetchPayments = async () => {
  isLoading.value = true
  try {
    const res = await axios.get('/admin/flow/payments')
    payments.value = Array.isArray(res.data) ? res.data : []
    
    // Calculate actual stats from DB data
    stats.value.total_today = payments.value.reduce((acc, curr) => acc + (curr.amount || 0), 0)
    stats.value.total_today_formatted = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(stats.value.total_today)
    stats.value.pending_count = payments.value.filter(p => p.status === 'PENDING').length
    stats.value.success_rate = payments.value.length > 0 
      ? Math.round((payments.value.filter(p => p.status === 'SUCCESS' || p.status === 'PAID').length / payments.value.length) * 100) + '%'
      : '0%'
  } catch (err) {
    console.error('Failed to fetch payments', err)
    payments.value = []
  } finally {
    isLoading.value = false
  }
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'SUCCESS': return 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
    case 'PENDING': return 'bg-amber-500/10 text-amber-600 border-amber-500/20'
    case 'FAILED': return 'bg-rose-500/10 text-rose-600 border-rose-500/20'
    default: return 'bg-slate-100 text-slate-400 border-slate-200'
  }
}

onMounted(fetchPayments)
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1400px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- Header -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-5">
          <div class="p-4 bg-indigo-600 rounded-[2rem] text-white shadow-2xl shadow-indigo-500/20">
            <History class="w-8 h-8" />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-3xl font-black italic tracking-tighter uppercase text-slate-900">Riwayat <span class="text-indigo-600">Pembayaran</span></h1>
              <span class="px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded-full text-[9px] font-black uppercase tracking-widest text-indigo-600">Transaksi Real-time</span>
            </div>
            <p class="text-xs font-bold text-slate-400 uppercase tracking-[0.2em]">Pantau semua dana masuk &amp; status pembayaran</p>
          </div>
        </div>
        
        <div class="flex items-center gap-3">
          <button class="flex items-center gap-3 px-6 py-3 bg-white border border-slate-200 text-slate-600 rounded-2xl shadow-sm hover:bg-slate-50 transition-all active:scale-95">
            <Download class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-widest">Export Laporan</span>
          </button>
        </div>
      </section>

      <!-- Stats Grid -->
      <section class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-indigo-50 flex items-center justify-center text-indigo-600">
              <CreditCard class="w-6 h-6" />
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Masuk Hari Ini</p>
              <h3 class="text-xl font-black text-slate-900 italic">{{ stats.total_today_formatted }}</h3>
            </div>
          </div>
        </div>
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-emerald-50 flex items-center justify-center text-emerald-600">
              <CheckCircle2 class="w-6 h-6" />
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Success Rate</p>
              <h3 class="text-xl font-black text-slate-900 italic">{{ stats.success_rate }}</h3>
            </div>
          </div>
        </div>
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-amber-50 flex items-center justify-center text-amber-600">
              <Clock class="w-6 h-6" />
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Pending</p>
              <h3 class="text-xl font-black text-slate-900 italic">{{ stats.pending_count }} Transaksi</h3>
            </div>
          </div>
        </div>
        <div class="bg-indigo-600 p-6 rounded-[2.5rem] shadow-2xl shadow-indigo-500/20 text-white flex items-center justify-between">
          <div>
            <p class="text-[9px] font-black text-indigo-200 uppercase tracking-widest">Status Gateway</p>
            <h3 class="text-lg font-bold italic tracking-tight leading-tight">Xendit: Operational</h3>
          </div>
          <div class="w-3 h-3 rounded-full bg-emerald-400 animate-pulse"></div>
        </div>
      </section>

      <!-- Table Section -->
      <section class="bg-white/60 backdrop-blur-2xl rounded-[3rem] border border-white/60 shadow-2xl overflow-hidden">
        <div class="p-8 border-b border-slate-100 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="relative flex-1 max-w-md">
            <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input v-model="searchQuery" type="text" placeholder="Cari santri atau ID transaksi..." 
                   class="w-full pl-14 pr-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-xs font-bold focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all" />
          </div>
          <div class="flex items-center gap-3">
            <button class="flex items-center gap-2 px-5 py-3 bg-slate-50 border border-slate-100 rounded-xl text-[10px] font-black uppercase tracking-widest text-slate-500 hover:bg-indigo-50 hover:text-indigo-600 transition-all">
              <Filter class="w-3.5 h-3.5" /> Filter
            </button>
            <button class="flex items-center gap-2 px-5 py-3 bg-slate-50 border border-slate-100 rounded-xl text-[10px] font-black uppercase tracking-widest text-slate-500 hover:bg-indigo-50 hover:text-indigo-600 transition-all">
              <Calendar class="w-3.5 h-3.5" /> Periode
            </button>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="text-left bg-slate-50/50">
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400">ID Transaksi & Tanggal</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400">Santri</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400">Kategori</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400">Metode</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400 text-right">Nominal</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400 text-center">Status</th>
                <th class="py-6 px-8 text-[10px] font-black uppercase tracking-widest text-slate-400"></th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-if="isLoading" v-for="i in 5" :key="i">
                <td colspan="7" class="py-8 px-8"><div class="h-12 bg-slate-50 animate-pulse rounded-2xl w-full"></div></td>
              </tr>
              <tr v-else v-for="p in payments" :key="p.id" class="group hover:bg-indigo-50/30 transition-all duration-300">
                <td class="py-6 px-8">
                  <div class="flex flex-col">
                    <span class="text-[10px] font-black text-slate-900 uppercase tracking-tight">TRX-{{ p.id.toString().padStart(6, '0') }}</span>
                    <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ p.date }}</span>
                  </div>
                </td>
                <td class="py-6 px-8">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 flex items-center justify-center text-[10px] font-black uppercase group-hover:bg-white transition-colors border border-transparent group-hover:border-indigo-100">
                      {{ p.student_name.charAt(0) }}
                    </div>
                    <span class="text-xs font-black text-slate-800">{{ p.student_name }}</span>
                  </div>
                </td>
                <td class="py-6 px-8">
                  <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">{{ p.category }}</span>
                </td>
                <td class="py-6 px-8">
                  <div class="flex items-center gap-2">
                    <CreditCard class="w-3.5 h-3.5 text-slate-400" />
                    <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">{{ p.method }}</span>
                  </div>
                </td>
                <td class="py-6 px-8 text-right">
                  <span class="text-sm font-black text-slate-900 italic tracking-tight">{{ p.amount_formatted }}</span>
                </td>
                <td class="py-6 px-8">
                  <div class="flex justify-center">
                    <span :class="getStatusClass(p.status)" class="px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest border flex items-center gap-2">
                      <div class="w-1.5 h-1.5 rounded-full" :class="p.status === 'SUCCESS' ? 'bg-emerald-500' : p.status === 'PENDING' ? 'bg-amber-500' : 'bg-rose-500'"></div>
                      {{ p.status }}
                    </span>
                  </div>
                </td>
                <td class="py-6 px-8 text-right">
                  <button class="p-2 text-slate-300 hover:text-indigo-600 hover:bg-white rounded-xl transition-all shadow-sm border border-transparent hover:border-indigo-100">
                    <MoreVertical class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty State -->
        <div v-if="!isLoading && payments.length === 0" class="py-32 flex flex-col items-center justify-center text-slate-300">
          <History class="w-20 h-20 mb-6 opacity-10" />
          <p class="text-sm font-black uppercase tracking-[0.2em]">Tidak ada riwayat transaksi</p>
        </div>

        <!-- Pagination -->
        <div class="p-8 border-t border-slate-50 flex items-center justify-between">
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Menampilkan {{ payments.length }} dari {{ payments.length }} transaksi</p>
          <div class="flex items-center gap-2">
            <button disabled class="p-2 bg-slate-50 rounded-xl text-slate-300 border border-slate-100 cursor-not-allowed">
              <History class="w-4 h-4 rotate-180" />
            </button>
            <button disabled class="p-2 bg-slate-50 rounded-xl text-slate-300 border border-slate-100 cursor-not-allowed">
              <History class="w-4 h-4" />
            </button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.animate-reveal {
  animation: reveal 0.8s cubic-bezier(0, 0, 0.2, 1) forwards;
}
@keyframes reveal {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
