<script setup>
import { ref, onMounted, computed } from 'vue'
import { 
  TrendingUp, 
  TrendingDown, 
  Wallet, 
  Plus, 
  Search, 
  Filter,
  Calendar,
  MoreVertical,
  Trash2,
  ArrowUpRight,
  ArrowDownLeft,
  DollarSign,
  FileText,
  History
} from 'lucide-vue-next'
import axios from 'axios'
import { toast } from 'vue3-toastify'

// State
const cashFlows = ref([])
const summary = ref({
  total_income: 0,
  income_formatted: 'Rp 0',
  total_expense: 0,
  expense_formatted: 'Rp 0',
  balance: 0,
  balance_formatted: 'Rp 0'
})
const loading = ref(false)
const syncing = ref(false)
const showAddModal = ref(false)

const syncData = async () => {
  syncing.value = true
  try {
    const res = await axios.post('/admin/flow/cashflow/sync')
    toast.success(res.data.message)
    await fetchCashFlows()
    await fetchSummary()
  } catch (err) {
    toast.error('Gagal menyinkronkan data')
  } finally {
    syncing.value = false
  }
}

// Filters
const filters = ref({
  type: '',
  start_date: '',
  end_date: ''
})

// Form
const form = ref({
  type: 'expense',
  category: '',
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  description: '',
  admin_name: ''
})

const loadInvoiceSettings = () => {
  const saved = localStorage.getItem('sigma_invoice_settings')
  if (saved) {
    const settings = JSON.parse(saved)
    form.value.admin_name = settings.adminName
  }
}

const fetchCashFlows = async () => {
  loading.ref = true
  try {
    const params = new URLSearchParams()
    if (filters.value.type) params.append('type', filters.value.type)
    if (filters.value.start_date) params.append('start_date', filters.value.start_date)
    if (filters.value.end_date) params.append('end_date', filters.value.end_date)
    
    const res = await axios.get(`/api/v1/admin/flow/cashflow?${params.toString()}`)
    cashFlows.value = res.data
  } catch (err) {
    toast.error('Gagal mengambil data arus kas')
  } finally {
    loading.value = false
  }
}

const fetchSummary = async () => {
  try {
    const res = await axios.get('/admin/flow/cashflow/summary')
    summary.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const handleSubmit = async () => {
  if (!form.value.category || !form.value.amount || !form.value.description) {
    toast.warning('Mohon lengkapi semua data')
    return
  }

  try {
    await axios.post('/admin/flow/cashflow', form.value)
    toast.success('Transaksi berhasil dicatat')
    showAddModal.value = false
    resetForm()
    await fetchCashFlows()
    await fetchSummary()
  } catch (err) {
    toast.error('Gagal mencatat transaksi')
  }
}

const deleteEntry = async (id) => {
  if (!confirm('Hapus transaksi ini?')) return
  try {
    await axios.delete(`/api/v1/admin/flow/cashflow/${id}`)
    toast.success('Transaksi dihapus')
    await fetchCashFlows()
    await fetchSummary()
  } catch (err) {
    toast.error('Gagal menghapus transaksi')
  }
}

const resetForm = () => {
  form.value = {
    type: 'expense',
    category: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    description: '',
    admin_name: form.value.admin_name // Preserve admin name
  }
}

const formatCurrency = (val) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(val)
}

onMounted(() => {
  loadInvoiceSettings()
  fetchCashFlows()
  fetchSummary()
})
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1600px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- Header -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-5">
          <div class="p-4 bg-emerald-600 rounded-[2rem] text-white shadow-2xl shadow-emerald-500/20">
            <Wallet class="w-8 h-8" />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-3xl font-black italic tracking-tighter uppercase text-slate-900">Arus <span class="text-emerald-600">Kas</span></h1>
              <span class="px-3 py-1 bg-emerald-500/10 border border-emerald-500/20 rounded-full text-[9px] font-black uppercase tracking-widest text-emerald-600">Laporan Keuangan</span>
            </div>
            <p class="text-slate-500 font-medium text-sm">Kelola pemasukan dan pengeluaran operasional yayasan</p>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button 
            @click="syncData"
            class="flex items-center gap-2 px-6 py-3.5 bg-white text-slate-600 border border-slate-200 rounded-2xl font-bold text-sm shadow-sm hover:bg-slate-50 transition-all"
            :disabled="syncing"
          >
            <History v-if="!syncing" class="w-5 h-5" />
            <div v-else class="w-5 h-5 border-2 border-slate-300 border-t-slate-600 rounded-full animate-spin"></div>
            {{ syncing ? 'Menyinkronkan...' : 'Sinkronisasi Data' }}
          </button>
          <button 
            @click="showAddModal = true"
            class="flex items-center gap-2 px-6 py-3.5 bg-slate-900 text-white rounded-2xl font-bold text-sm shadow-xl shadow-slate-900/10 hover:bg-slate-800 transition-all hover:-translate-y-1 active:translate-y-0"
          >
            <Plus class="w-5 h-5" />
            Catat Transaksi
          </button>
        </div>
      </section>

      <!-- Summary Cards -->
      <section class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Balance -->
        <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
          <div class="absolute top-0 right-0 p-8 opacity-5 group-hover:scale-110 transition-transform duration-700">
            <Wallet class="w-24 h-24 text-slate-900" />
          </div>
          <p class="text-slate-400 font-black uppercase tracking-widest text-[10px] mb-2">Saldo Saat Ini</p>
          <h2 class="text-3xl font-black text-slate-900 mb-2">{{ summary.balance_formatted }}</h2>
          <div class="flex items-center gap-2">
            <span class="px-2 py-0.5 bg-slate-100 rounded-full text-[9px] font-bold text-slate-500 uppercase tracking-tight">Total Kas Tersedia</span>
          </div>
        </div>

        <!-- Income -->
        <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
          <div class="absolute top-0 right-0 p-8 opacity-5 group-hover:scale-110 transition-transform duration-700">
            <TrendingUp class="w-24 h-24 text-emerald-600" />
          </div>
          <p class="text-emerald-500 font-black uppercase tracking-widest text-[10px] mb-2">Total Pemasukan</p>
          <h2 class="text-3xl font-black text-slate-900 mb-2">{{ summary.income_formatted }}</h2>
          <div class="flex items-center gap-2">
            <ArrowUpRight class="w-4 h-4 text-emerald-500" />
            <span class="text-[11px] font-bold text-emerald-600 uppercase tracking-tight">Uang Masuk</span>
          </div>
        </div>

        <!-- Expense -->
        <div class="bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm relative overflow-hidden group">
          <div class="absolute top-0 right-0 p-8 opacity-5 group-hover:scale-110 transition-transform duration-700">
            <TrendingDown class="w-24 h-24 text-rose-600" />
          </div>
          <p class="text-rose-500 font-black uppercase tracking-widest text-[10px] mb-2">Total Pengeluaran</p>
          <h2 class="text-3xl font-black text-slate-900 mb-2">{{ summary.expense_formatted }}</h2>
          <div class="flex items-center gap-2">
            <ArrowDownLeft class="w-4 h-4 text-rose-500" />
            <span class="text-[11px] font-bold text-rose-600 uppercase tracking-tight">Uang Keluar</span>
          </div>
        </div>
      </section>

      <!-- Filters & Table -->
      <section class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
        <div class="p-8 border-b border-slate-50 flex flex-col md:flex-row justify-between items-center gap-6">
          <div class="flex items-center gap-3">
            <div class="p-2.5 bg-emerald-50 text-emerald-600 rounded-xl">
              <Filter class="w-5 h-5" />
            </div>
            <h3 class="font-black text-slate-900 uppercase tracking-tighter italic">Riwayat Transaksi</h3>
          </div>

          <div class="flex flex-wrap items-center gap-4">
            <select 
              v-model="filters.type"
              @change="fetchCashFlows"
              class="px-4 py-2.5 bg-slate-50 border-none rounded-xl text-xs font-bold text-slate-600 focus:ring-2 focus:ring-emerald-500 transition-all outline-none"
            >
              <option value="">Semua Tipe</option>
              <option value="income">Pemasukan</option>
              <option value="expense">Pengeluaran</option>
            </select>

            <div class="flex items-center gap-2 bg-slate-50 rounded-xl px-3 py-2 border border-transparent focus-within:border-emerald-500 transition-all">
              <Calendar class="w-4 h-4 text-slate-400" />
              <input 
                type="date" 
                v-model="filters.start_date"
                @change="fetchCashFlows"
                class="bg-transparent border-none text-[11px] font-bold text-slate-600 outline-none p-0"
              />
              <span class="text-slate-300 mx-1">-</span>
              <input 
                type="date" 
                v-model="filters.end_date"
                @change="fetchCashFlows"
                class="bg-transparent border-none text-[11px] font-bold text-slate-600 outline-none p-0"
              />
            </div>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
              <tr class="bg-slate-50/50">
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Tipe</th>
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Kategori</th>
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Deskripsi</th>
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Tanggal</th>
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Nominal</th>
                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr 
                v-for="entry in cashFlows" 
                :key="entry.id"
                class="hover:bg-slate-50/30 transition-colors group"
              >
                <td class="px-8 py-6">
                  <div class="flex items-center gap-3">
                    <div 
                      :class="entry.type === 'income' ? 'bg-emerald-100 text-emerald-600' : 'bg-rose-100 text-rose-600'"
                      class="p-2 rounded-lg"
                    >
                      <ArrowUpRight v-if="entry.type === 'income'" class="w-4 h-4" />
                      <ArrowDownLeft v-else class="w-4 h-4" />
                    </div>
                    <span 
                      :class="entry.type === 'income' ? 'text-emerald-600' : 'text-rose-600'"
                      class="text-[10px] font-black uppercase tracking-widest"
                    >
                      {{ entry.type === 'income' ? 'Masuk' : 'Keluar' }}
                    </span>
                  </div>
                </td>
                <td class="px-8 py-6">
                  <span class="text-xs font-bold text-slate-900 block mb-0.5">{{ entry.category }}</span>
                  <span class="text-[10px] text-slate-400 font-medium uppercase tracking-tight">Kategori</span>
                </td>
                <td class="px-8 py-6">
                  <span class="text-xs text-slate-600 font-medium line-clamp-1 max-w-xs">{{ entry.description }}</span>
                </td>
                <td class="px-8 py-6">
                  <span class="text-xs font-bold text-slate-900 block mb-0.5">{{ entry.date }}</span>
                  <span class="text-[10px] text-slate-400 font-medium uppercase tracking-tight">Waktu Catat</span>
                </td>
                <td class="px-8 py-6 text-right">
                  <span 
                    :class="entry.type === 'income' ? 'text-emerald-600' : 'text-rose-600'"
                    class="text-sm font-black italic tracking-tighter"
                  >
                    {{ entry.type === 'income' ? '+' : '-' }} {{ entry.amount_formatted }}
                  </span>
                </td>
                <td class="px-8 py-6 text-center">
                  <button 
                    @click="deleteEntry(entry.id)"
                    class="p-2 text-slate-300 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </td>
              </tr>
              <tr v-if="cashFlows.length === 0">
                <td colspan="6" class="px-8 py-20 text-center">
                  <div class="flex flex-col items-center gap-3">
                    <div class="p-5 bg-slate-50 text-slate-300 rounded-full">
                      <FileText class="w-12 h-12" />
                    </div>
                    <p class="text-slate-400 font-bold text-sm">Belum ada transaksi tercatat</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>

    <!-- Modal Catat Transaksi -->
    <div 
      v-if="showAddModal"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-md animate-fade-in"
      @click.self="showAddModal = false"
    >
      <div class="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-zoom-in">
        <div class="p-8 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
          <div class="flex items-center gap-3">
            <div class="p-2.5 bg-slate-900 text-white rounded-xl">
              <Plus class="w-5 h-5" />
            </div>
            <h3 class="text-xl font-black italic tracking-tighter uppercase text-slate-900">Catat <span class="text-emerald-600">Transaksi</span></h3>
          </div>
          <button @click="showAddModal = false" class="p-2 hover:bg-white rounded-full transition-colors">
            <Plus class="w-6 h-6 text-slate-400 rotate-45" />
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="p-8 space-y-6">
          <!-- Type Toggle -->
          <div class="grid grid-cols-2 gap-4">
            <button 
              type="button"
              @click="form.type = 'income'"
              :class="form.type === 'income' ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-500/20' : 'bg-slate-50 text-slate-400'"
              class="py-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2"
            >
              <ArrowUpRight class="w-4 h-4" />
              Pemasukan
            </button>
            <button 
              type="button"
              @click="form.type = 'expense'"
              :class="form.type === 'expense' ? 'bg-rose-600 text-white shadow-lg shadow-rose-500/20' : 'bg-slate-50 text-slate-400'"
              class="py-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2"
            >
              <ArrowDownLeft class="w-4 h-4" />
              Pengeluaran
            </button>
          </div>

          <div class="space-y-4">
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kategori Transaksi</label>
              <input 
                v-model="form.category"
                type="text"
                placeholder="Contoh: Belanja Dapur, Gaji Guru, Donasi..."
                class="w-full px-5 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold text-slate-900 focus:ring-2 focus:ring-emerald-500 transition-all outline-none"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nominal (Rp)</label>
                <div class="relative">
                  <div class="absolute left-5 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-xs">Rp</div>
                  <input 
                    v-model.number="form.amount"
                    type="number"
                    class="w-full pl-12 pr-5 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold text-slate-900 focus:ring-2 focus:ring-emerald-500 transition-all outline-none"
                  />
                </div>
              </div>
              <div class="space-y-1.5">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Tanggal</label>
                <input 
                  v-model="form.date"
                  type="date"
                  class="w-full px-5 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold text-slate-900 focus:ring-2 focus:ring-emerald-500 transition-all outline-none"
                />
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Deskripsi / Keterangan</label>
              <textarea 
                v-model="form.description"
                rows="3"
                placeholder="Tulis detail transaksi di sini..."
                class="w-full px-5 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold text-slate-900 focus:ring-2 focus:ring-emerald-500 transition-all outline-none resize-none"
              ></textarea>
            </div>
          </div>

          <button 
            type="submit"
            class="w-full py-4.5 bg-slate-900 text-white rounded-[1.5rem] font-black text-xs uppercase tracking-widest shadow-2xl shadow-slate-900/20 hover:bg-slate-800 transition-all hover:-translate-y-1 active:translate-y-0"
          >
            Simpan Transaksi
          </button>
        </form>
      </div>
    </div>
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

.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

.animate-zoom-in {
  animation: zoomIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes zoomIn {
  from { opacity: 0; transform: scale(0.9) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
</style>
