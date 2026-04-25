<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  CircleDollarSign, ArrowLeft, Activity, TrendingUp, 
  Wallet, Receipt, ArrowUpRight,
  MoreVertical, Download, Plus, Filter, Search,
  Calendar, LayoutDashboard, CreditCard
} from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const isLoading = ref(true)
const stats = ref<any>({
  total_revenue: 0,
  revenue_formatted: 'Rp 0',
  pending_amount: 0,
  pending_formatted: 'Rp 0',
  total_invoices: 0,
  recent_payments: []
})

const fetchStats = async () => {
  isLoading.value = true
  try {
    const res = await axios.get('/api/v1/admin/flow/stats')
    stats.value = res.data
  } catch (err) {
    console.error('Failed to fetch finance stats', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  document.documentElement.classList.remove('dark')
  localStorage.setItem('sigma_theme', 'light')
  fetchStats()
})

const kpis = [
  { label: 'Total Pendapatan', key: 'revenue_formatted', icon: TrendingUp, color: 'text-emerald-600', bg: 'bg-emerald-500/10', border: 'border-emerald-200' },
  { label: 'Tagihan Pending', key: 'pending_formatted', icon: Wallet, color: 'text-amber-600', bg: 'bg-amber-500/10', border: 'border-amber-200' },
  { label: 'Total Faktur', key: 'total_invoices', icon: Receipt, color: 'text-indigo-600', bg: 'bg-indigo-500/10', border: 'border-indigo-200' }
]
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1600px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- HEADER BANNER (Sigmaedu Style) -->
      <section class="relative overflow-hidden rounded-[3rem] bg-white/40 backdrop-blur-2xl p-8 lg:p-12 shadow-2xl border border-white/60">
        <div class="absolute -top-24 -right-24 w-80 h-80 bg-indigo-400/20 rounded-full blur-[100px]"></div>
        <div class="absolute -bottom-24 -left-24 w-60 h-60 bg-emerald-400/10 rounded-full blur-[80px]"></div>
        
        <div class="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-10">
          <div class="flex items-center gap-6">
            <div class="w-20 h-20 bg-indigo-600 rounded-[2rem] flex items-center justify-center text-white shadow-2xl shadow-indigo-500/20">
              <CircleDollarSign class="w-10 h-10" />
            </div>
            <div>
              <div class="flex items-center gap-3 mb-2">
                <span class="px-3 py-1 rounded-full bg-indigo-500/10 text-indigo-600 border border-indigo-500/20 text-[9px] font-black uppercase tracking-widest">
                  Financial Module
                </span>
              </div>
              <h1 class="text-4xl lg:text-5xl font-black text-slate-900 tracking-tighter italic uppercase">
                Sigma<span class="text-indigo-600">flow</span>
              </h1>
              <p class="text-slate-500 text-sm font-medium mt-2 max-w-xl">
                Kelola arus kas, faktur santri, dan sistem pembayaran online secara terpusat dengan transparansi maksimal.
              </p>
            </div>
          </div>

          <div class="flex items-center gap-4">
             <button @click="router.push('/dashboard')" class="px-8 py-3 rounded-2xl bg-white/80 border border-slate-200 text-slate-600 text-[10px] font-black uppercase tracking-widest hover:bg-white transition-all shadow-xl flex items-center gap-2">
                <ArrowLeft class="w-4 h-4" /> Kembali
             </button>
             <button class="px-8 py-3 rounded-2xl bg-[#0F2942] text-white text-[10px] font-black uppercase tracking-widest hover:bg-indigo-600 transition-all shadow-xl flex items-center gap-2">
                <Plus class="w-4 h-4" /> Buat Faktur
             </button>
          </div>
        </div>
      </section>

      <!-- KPI GRID -->
      <section class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div v-for="kpi in kpis" :key="kpi.label" 
             :class="[kpi.bg, kpi.border]"
             class="relative overflow-hidden p-8 rounded-[3rem] border shadow-xl backdrop-blur-md transition-all duration-500 hover:scale-[1.02]">
          <div class="flex items-center justify-between mb-8">
            <div :class="kpi.color" class="w-12 h-12 rounded-2xl bg-white/60 flex items-center justify-center shadow-sm">
              <component :is="kpi.icon" class="w-6 h-6" />
            </div>
            <div class="flex flex-col items-end">
               <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">{{ kpi.label }}</p>
               <span class="text-[9px] font-bold text-emerald-500">+12.5%</span>
            </div>
          </div>
          
          <div>
            <h2 v-if="isLoading" class="text-3xl font-black text-slate-300 italic animate-pulse">Loading...</h2>
            <h2 v-else class="text-3xl font-black text-slate-900 italic tracking-tight">{{ stats[kpi.key] }}</h2>
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-2">Update Terakhir: Hari ini</p>
          </div>
        </div>
      </section>

      <!-- TRANSACTIONS & DETAILS -->
      <section class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Recent Transactions -->
        <div class="lg:col-span-2 bg-white/60 backdrop-blur-2xl rounded-[3rem] p-10 border border-white/60 shadow-2xl">
          <div class="flex items-center justify-between mb-10">
            <div>
              <h3 class="text-lg font-black text-slate-900 italic">Transaksi Terbaru</h3>
              <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Riwayat Pembayaran Masuk</p>
            </div>
            <div class="flex items-center gap-2">
              <button class="p-2 bg-slate-50 text-slate-400 rounded-xl hover:bg-indigo-50 hover:text-indigo-600 transition-all border border-slate-100">
                <Search class="w-4 h-4" />
              </button>
              <button class="p-2 bg-slate-50 text-slate-400 rounded-xl hover:bg-indigo-50 hover:text-indigo-600 transition-all border border-slate-100">
                <Filter class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="text-left border-b border-slate-100">
                  <th class="pb-6 text-[10px] font-black uppercase tracking-widest text-slate-400 px-4">Santri</th>
                  <th class="pb-6 text-[10px] font-black uppercase tracking-widest text-slate-400 px-4">Kategori</th>
                  <th class="pb-6 text-[10px] font-black uppercase tracking-widest text-slate-400 px-4 text-right">Jumlah</th>
                  <th class="pb-6 text-[10px] font-black uppercase tracking-widest text-slate-400 px-4">Status</th>
                  <th class="pb-6 text-[10px] font-black uppercase tracking-widest text-slate-400 px-4 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-50">
                <tr v-if="isLoading" v-for="i in 5" :key="i">
                  <td colspan="5" class="py-6 px-4 animate-pulse">
                    <div class="h-10 bg-slate-100 rounded-2xl w-full"></div>
                  </td>
                </tr>
                <tr v-else v-for="p in stats.recent_payments" :key="p.id" class="group hover:bg-slate-50/50 transition-all">
                  <td class="py-5 px-4">
                    <div class="flex items-center gap-3">
                      <div class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 flex items-center justify-center text-[10px] font-black uppercase">
                        {{ p.student_name.charAt(0) }}
                      </div>
                      <div>
                         <p class="text-xs font-black text-slate-800">{{ p.student_name }}</p>
                         <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">{{ p.date }}</p>
                      </div>
                    </div>
                  </td>
                  <td class="py-5 px-4">
                    <span class="px-3 py-1 rounded-full bg-indigo-50 text-indigo-600 text-[9px] font-black uppercase tracking-widest border border-indigo-100">
                      {{ p.category }}
                    </span>
                  </td>
                  <td class="py-5 px-4 text-right">
                    <p class="text-xs font-black text-emerald-600">{{ p.amount_formatted }}</p>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">{{ p.method }}</p>
                  </td>
                  <td class="py-5 px-4">
                    <div class="flex items-center gap-2">
                       <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
                       <span class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Berhasil</span>
                    </div>
                  </td>
                  <td class="py-5 px-4 text-center">
                    <button class="p-2 text-slate-300 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all">
                      <MoreVertical class="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          
          <div v-if="!stats.recent_payments.length && !isLoading" class="py-20 text-center">
             <div class="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4 border border-slate-100">
                <Activity class="w-8 h-8 text-slate-200" />
             </div>
             <p class="text-xs font-bold text-slate-400 uppercase tracking-widest">Belum ada transaksi terdaftar</p>
          </div>
        </div>

        <!-- Sidebar Actions -->
        <div class="flex flex-col gap-8">
           <!-- Xendit Status Card -->
           <div class="bg-[#0F2942] rounded-[3rem] p-10 text-white shadow-2xl relative overflow-hidden group">
              <div class="absolute -right-10 -bottom-10 w-40 h-40 bg-blue-500/10 rounded-full blur-3xl group-hover:scale-110 transition-transform duration-1000"></div>
              <CreditCard class="absolute right-8 top-8 w-12 h-12 opacity-10" />
              
              <div class="relative z-10">
                 <h4 class="text-[10px] font-black uppercase tracking-widest text-blue-300 mb-8">Payment Gateway</h4>
                 <div class="flex items-center gap-4 mb-8">
                    <div class="w-12 h-12 rounded-2xl bg-white/10 flex items-center justify-center">
                       <Activity class="w-6 h-6 text-emerald-400" />
                    </div>
                    <div>
                       <p class="text-xl font-black italic">Xendit <span class="text-blue-400">Online</span></p>
                       <span class="text-[9px] font-black text-emerald-400 uppercase tracking-widest">System Operational</span>
                    </div>
                 </div>
                 <p class="text-[10px] font-bold text-slate-400 leading-relaxed uppercase tracking-wider">
                    Sistem pembayaran otomatis terintegrasi. Semua transaksi online diproses secara real-time.
                 </p>
              </div>
           </div>

           <!-- Distribution Chart Placeholder -->
           <div class="bg-white/60 backdrop-blur-2xl rounded-[3rem] p-10 border border-white/60 shadow-xl">
              <h4 class="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-8">Distribusi Pendapatan</h4>
              <div class="space-y-6">
                 <div v-for="cat in ['SPP Bulanan', 'Uang Makan', 'Infaq Gedung']" :key="cat" class="space-y-2">
                    <div class="flex justify-between text-[10px] font-black uppercase tracking-widest">
                       <span class="text-slate-600">{{ cat }}</span>
                       <span class="text-indigo-600">65%</span>
                    </div>
                    <div class="h-2 bg-slate-100 rounded-full overflow-hidden">
                       <div class="h-full bg-indigo-500 rounded-full" style="width: 65%"></div>
                    </div>
                 </div>
              </div>
              <button class="w-full mt-10 py-4 rounded-[1.5rem] bg-slate-50 text-[10px] font-black text-slate-400 uppercase tracking-widest hover:bg-indigo-50 hover:text-indigo-600 transition-all border border-slate-100">
                 Detail Laporan
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

::-webkit-scrollbar {
  width: 5px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}
</style>
