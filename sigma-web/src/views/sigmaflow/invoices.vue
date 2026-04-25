<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Plus, Search, 
  Trash2, AlertCircle, CheckCircle2, X, Edit3,
  FileText, TrendingUp,
  Filter, ChevronRight, Receipt, Download, 
  ExternalLink, ChevronLeft
} from 'lucide-vue-next'

const closeSuggestions = () => {
  setTimeout(() => {
    showSuggestions.value = false
    if (newInvoice.value.nis && !newInvoice.value.student_id) {
      fetchStudentByNIS()
    }
  }, 200)
}
import axios from 'axios'

const invoices = ref<any[]>([])
const categories = ref<any[]>([])
const searchQuery = ref('')
const statusFilter = ref('All')
const showCreateModal = ref(false)
const isLoading = ref(true)
const isSubmitting = ref(false)
const studentSuggestions = ref<any[]>([])
const showSuggestions = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const editingId = ref<number | null>(null)

const showBulkModal = ref(false)
const classrooms = ref<any[]>([])
const bulkForm = ref({
  payment_category_id: '',
  amount: 0,
  due_date: '',
  target: 'all', // all | classroom
  classroom_id: ''
})

const newInvoice = ref({
  student_id: '',
  nis: '',
  student_name: '',
  guardian_name: '',
  payment_category_id: '',
  amount: 0,
  due_date: '',
  description: ''
})

const searchStudents = async () => {
  if (newInvoice.value.nis.length < 2) {
    studentSuggestions.value = []
    showSuggestions.value = false
    return
  }

  try {
    const res = await axios.get(`/api/v1/admin/base/students?search=${newInvoice.value.nis}`)
    studentSuggestions.value = Array.isArray(res.data) ? res.data.slice(0, 5) : []
    showSuggestions.value = studentSuggestions.value.length > 0
  } catch (err) {
    console.error('Search error', err)
  }
}

const selectStudent = (student: any) => {
  newInvoice.value.nis = student.nis
  newInvoice.value.student_id = student.ID || student.id
  newInvoice.value.student_name = student.name
  newInvoice.value.guardian_name = student.parent_name || student.ParentName || '-'
  showSuggestions.value = false
}

const fetchStudentByNIS = async () => {
  if (!newInvoice.value.nis) return
  
  try {
    const res = await axios.get(`/api/v1/admin/flow/students/nis/${newInvoice.value.nis}`)
    newInvoice.value.student_id = res.data.ID || res.data.id
    newInvoice.value.student_name = res.data.Name || res.data.name
    newInvoice.value.guardian_name = res.data.parent_name || res.data.ParentName || '-'
  } catch (err) {
    console.error('Student not found', err)
    newInvoice.value.student_id = ''
    newInvoice.value.student_name = 'SANTRI TIDAK DITEMUKAN'
    newInvoice.value.guardian_name = '-'
  }
}

const handleCategoryChange = () => {
  const selectedCat = categories.value.find(c => (c.id || c.ID) == newInvoice.value.payment_category_id)
  if (selectedCat) {
    newInvoice.value.amount = selectedCat.Amount || selectedCat.amount || 0
  }
}

const formatRupiah = (val: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val)
}

const summary = computed(() => {
  const totalVal = invoices.value.reduce((acc: number, inv: any) => acc + (inv.amount || 0), 0)
  const paidVal = invoices.value.filter((inv: any) => inv.status === 'Paid').reduce((acc: number, inv: any) => acc + (inv.amount || 0), 0)
  const unpaidVal = totalVal - paidVal
  return {
    total: formatRupiah(totalVal),
    paid: formatRupiah(paidVal),
    unpaid: formatRupiah(unpaidVal),
    count: invoices.value.length
  }
})

const fetchCategories = async () => {
  try {
    const res = await axios.get('/api/v1/admin/flow/categories')
    if (Array.isArray(res.data)) {
      categories.value = res.data.map((c: any) => ({ 
        id: c.ID || c.id, 
        name: c.Name || c.name,
        amount: c.Amount || c.amount 
      }))
    }
  } catch (err) {
    console.error('Failed to fetch categories', err)
    categories.value = [
      { id: 1, name: 'SPP Bulanan', amount: 500000 },
      { id: 2, name: 'Uang Makan', amount: 300000 },
      { id: 3, name: 'Infaq Gedung', amount: 1000000 }
    ]
  }
}

const fetchInvoices = async () => {
  isLoading.value = true
  try {
    const res = await axios.get('/api/v1/admin/flow/invoices')
    invoices.value = Array.isArray(res.data) ? res.data : []
  } catch (err) {
    console.error('Failed to fetch invoices', err)
    invoices.value = []
  } finally {
    isLoading.value = false
  }
}

const openCreateModal = () => {
  modalMode.value = 'create'
  showCreateModal.value = true
  newInvoice.value = { 
    student_id: '', 
    nis: '', 
    student_name: '', 
    guardian_name: '', 
    payment_category_id: '', 
    amount: 0, 
    due_date: '', 
    description: '' 
  }
}

const openEditModal = (inv: any) => {
  modalMode.value = 'edit'
  editingId.value = inv.ID || inv.id
  newInvoice.value = {
    student_id: String(inv.student_id),
    nis: inv.student?.nis || '',
    student_name: inv.student?.name || '',
    guardian_name: inv.student?.parent_name || '',
    payment_category_id: String(inv.payment_category_id),
    amount: inv.amount,
    due_date: inv.due_date ? inv.due_date.split('T')[0] : '',
    description: inv.description || ''
  }
  showCreateModal.value = true
}

const deleteInvoice = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus faktur ini?')) return
  try {
    await axios.delete(`/api/v1/admin/flow/invoices/${id}`)
    await fetchInvoices()
  } catch (err) {
    console.error('Failed to delete invoice', err)
    alert('Gagal menghapus faktur')
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
}

const openPaymentLink = (url: string) => {
  if (url) window.open(url, '_blank')
}

const submitInvoice = async () => {
  if (!newInvoice.value.student_id || !newInvoice.value.payment_category_id || !newInvoice.value.amount) {
    alert('Harap isi semua field yang wajib (ID Santri, Kategori, dan Nominal)!')
    return
  }

  isSubmitting.value = true
  try {
    const payload = {
      student_id: Number(newInvoice.value.student_id),
      payment_category_id: Number(newInvoice.value.payment_category_id),
      amount: Number(newInvoice.value.amount),
      due_date: newInvoice.value.due_date,
      description: newInvoice.value.description
    }
    
    if (modalMode.value === 'create') {
      await axios.post('/api/v1/admin/flow/invoices', payload)
    } else {
      await axios.put(`/api/v1/admin/flow/invoices/${editingId.value}`, payload)
    }

    showCreateModal.value = false
    newInvoice.value = { 
      student_id: '', 
      nis: '', 
      student_name: '', 
      guardian_name: '', 
      payment_category_id: '', 
      amount: 0, 
      due_date: '', 
      description: '' 
    }
    await fetchInvoices()
  } catch (err: any) {
    console.error('Failed to create invoice', err)
    const msg = err?.response?.data?.error || err.message || 'Unknown error'
    alert('Gagal membuat faktur: ' + msg)
  } finally {
    isSubmitting.value = false
  }
}

const handleBulkCategoryChange = () => {
  const selectedCat = categories.value.find(c => (c.id || c.ID) == bulkForm.value.payment_category_id)
  if (selectedCat) {
    bulkForm.value.amount = selectedCat.Amount || selectedCat.amount || 0
  }
}

const fetchClassrooms = async () => {
  try {
    const res = await axios.get('/api/v1/admin/base/classrooms')
    classrooms.value = Array.isArray(res.data) ? res.data : []
  } catch (err) {
    console.error('Failed to fetch classrooms', err)
  }
}

const submitBulkInvoices = async () => {
  if (!bulkForm.value.payment_category_id || !bulkForm.value.due_date) {
    alert('Harap isi kategori dan tanggal jatuh tempo')
    return
  }
  
  if (bulkForm.value.target === 'classroom' && !bulkForm.value.classroom_id) {
    alert('Harap pilih kelas')
    return
  }

  isSubmitting.value = true
  try {
    const res = await axios.post('/api/v1/admin/flow/invoices/bulk', {
      payment_category_id: Number(bulkForm.value.payment_category_id),
      amount: bulkForm.value.amount,
      due_date: bulkForm.value.due_date,
      target: bulkForm.value.target,
      classroom_id: bulkForm.value.classroom_id ? Number(bulkForm.value.classroom_id) : 0
    })
    
    alert(res.data.message)
    showBulkModal.value = false
    await fetchInvoices()
  } catch (err: any) {
    console.error('Bulk error', err)
    alert('Gagal membuat penagihan massal: ' + (err.response?.data?.error || err.message))
  } finally {
    isSubmitting.value = false
  }
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'Paid': return 'status-paid'
    case 'Unpaid': return 'status-unpaid'
    case 'Partial': return 'status-partial'
    default: return 'status-default'
  }
}

const getStatusDot = (status: string) => {
  switch (status) {
    case 'Paid': return 'dot-paid'
    case 'Unpaid': return 'dot-unpaid'
    case 'Partial': return 'dot-partial'
    default: return 'dot-default'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'Paid': return 'Lunas'
    case 'Unpaid': return 'Belum Bayar'
    case 'Partial': return 'Cicilan'
    default: return status
  }
}

const filteredInvoices = computed(() => {
  if (!Array.isArray(invoices.value)) return []
  return invoices.value.filter((inv: any) => {
    const sName = String(inv.student_name || '')
    const iNum = String(inv.invoice_number || '')
    const q = searchQuery.value.toLowerCase()
    const matchesSearch = sName.toLowerCase().includes(q) || iNum.toLowerCase().includes(q)
    const matchesStatus = statusFilter.value === 'All' || inv.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const filterLabels: Record<string, string> = {
  'All': 'Semua',
  'Paid': 'Lunas',
  'Unpaid': 'Belum Bayar',
  'Partial': 'Cicilan'
}

onMounted(async () => {
  await fetchCategories()
  await fetchInvoices()
  await fetchClassrooms()
})
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] font-sans pb-20">
    <main class="max-w-[1600px] mx-auto p-8 lg:p-10 flex flex-col gap-8 animate-reveal">
      
      <!-- Header -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-5">
          <div class="p-4 bg-indigo-600 rounded-[2rem] text-white shadow-2xl shadow-indigo-500/20">
            <Receipt class="w-8 h-8" />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-3xl font-black italic tracking-tighter uppercase text-slate-900">Data <span class="text-indigo-600">Faktur</span></h1>
              <span class="px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded-full text-[9px] font-black uppercase tracking-widest text-indigo-600">Manajemen Penagihan</span>
            </div>
            <p class="text-xs font-bold text-slate-400 uppercase tracking-[0.2em]">Daftar faktur santri &amp; status pembayaran</p>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <button class="flex items-center gap-3 px-6 py-3 bg-white border border-slate-200 hover:bg-slate-50 rounded-2xl transition-all shadow-xl shadow-slate-200/50">
            <Download class="w-4 h-4 text-slate-400" />
            <span class="text-[10px] font-black uppercase tracking-widest text-slate-600">Export PDF</span>
          </button>
          <button @click="showBulkModal = true" class="flex items-center gap-3 px-6 py-3 bg-white border border-slate-200 hover:bg-slate-50 rounded-2xl transition-all shadow-xl shadow-slate-200/50">
            <Plus class="w-4 h-4 text-indigo-500" />
            <span class="text-[10px] font-black uppercase tracking-widest text-slate-600">Tagih Massal</span>
          </button>
          <button @click="openCreateModal" class="flex items-center gap-3 px-6 py-3 bg-[#0F2942] hover:bg-indigo-600 text-white rounded-2xl shadow-xl shadow-indigo-600/20 transition-all active:scale-95">
            <Plus class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-widest">Buat Faktur Baru</span>
          </button>
        </div>
      </section>

      <!-- Summary Cards -->
      <section class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div class="bg-white/40 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-white/60 shadow-xl group hover:border-indigo-500/30 transition-all">
          <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2">Total Penagihan</p>
          <h3 class="text-2xl font-black text-slate-900 italic">{{ summary.total }}</h3>
          <div class="flex items-center gap-2 mt-3">
            <div class="w-1.5 h-1.5 rounded-full bg-indigo-500"></div>
            <span class="text-[9px] font-bold text-slate-500 uppercase">{{ summary.count }} Faktur Terbit</span>
          </div>
        </div>
        <div class="bg-emerald-50/50 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-emerald-100 shadow-xl group hover:border-emerald-500/30 transition-all">
          <p class="text-[9px] font-black text-emerald-600 uppercase tracking-widest mb-2">Dana Diterima</p>
          <h3 class="text-2xl font-black text-emerald-600 italic">{{ summary.paid }}</h3>
          <div class="flex items-center gap-2 mt-3">
            <TrendingUp class="w-3 h-3 text-emerald-500" />
            <span class="text-[9px] font-bold text-emerald-600 uppercase tracking-widest">Lunas</span>
          </div>
        </div>
        <div class="bg-rose-50/50 backdrop-blur-2xl p-6 rounded-[2.5rem] border border-rose-100 shadow-xl group hover:border-rose-500/30 transition-all">
          <p class="text-[9px] font-black text-rose-600 uppercase tracking-widest mb-2">Sisa Tunggakan</p>
          <h3 class="text-2xl font-black text-rose-600 italic">{{ summary.unpaid }}</h3>
          <div class="flex items-center gap-2 mt-3">
            <AlertCircle class="w-3 h-3 text-rose-500" />
            <span class="text-[9px] font-bold text-rose-600 uppercase tracking-widest">Belum Terbayar</span>
          </div>
        </div>
        <div class="bg-indigo-600 p-6 rounded-[2.5rem] shadow-2xl shadow-indigo-500/20 relative overflow-hidden group">
           <div class="absolute -right-6 -bottom-6 w-24 h-24 bg-white/10 rounded-full blur-2xl group-hover:scale-110 transition-transform duration-1000"></div>
           <p class="text-[9px] font-black text-indigo-200 uppercase tracking-widest mb-2">Efisiensi</p>
           <h3 class="text-3xl font-black text-white italic">84<span class="text-sm">%</span></h3>
           <p class="text-[9px] font-bold text-indigo-100 mt-3 uppercase tracking-widest">Kolektibilitas Dana</p>
        </div>
      </section>

      <!-- Search & Filters -->
      <section class="bg-white/40 backdrop-blur-2xl rounded-[2.5rem] p-6 border border-white/60 shadow-xl flex flex-col md:flex-row gap-4 items-center">
        <div class="relative flex-1 w-full">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input v-model="searchQuery" type="text" placeholder="Cari nama santri atau nomor faktur..." class="w-full pl-12 pr-4 py-3 bg-white/60 border border-slate-100 rounded-2xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all" />
        </div>
        <div class="flex items-center gap-3 w-full md:w-auto">
          <div class="flex bg-slate-100/50 p-1 rounded-2xl border border-slate-200/50">
            <button v-for="s in ['All', 'Paid', 'Unpaid', 'Partial']" :key="s"
                    @click="statusFilter = s"
                    :class="statusFilter === s ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-400 hover:text-slate-600'"
                    class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">
              {{ filterLabels[s] }}
            </button>
          </div>
          <button class="p-3 bg-white border border-slate-200 rounded-2xl text-slate-400 hover:text-indigo-600 transition-all shadow-sm">
            <Filter class="w-4 h-4" />
          </button>
        </div>
      </section>

      <!-- Table Section -->
      <section class="bg-white/40 backdrop-blur-2xl rounded-[3rem] border border-white/60 shadow-2xl overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
              <tr class="border-b border-slate-100">
                <th class="py-8 px-10 text-[10px] font-black uppercase tracking-widest text-slate-400">Nomor Faktur</th>
                <th class="py-8 px-6 text-[10px] font-black uppercase tracking-widest text-slate-400">Nama Santri</th>
                <th class="py-8 px-6 text-[10px] font-black uppercase tracking-widest text-slate-400">Kategori</th>
                <th class="py-8 px-6 text-[10px] font-black uppercase tracking-widest text-slate-400 text-right">Jumlah</th>
                <th class="py-8 px-6 text-[10px] font-black uppercase tracking-widest text-slate-400">Jatuh Tempo</th>
                <th class="py-8 px-6 text-[10px] font-black uppercase tracking-widest text-slate-400 text-center">Status</th>
                <th class="py-8 px-10 text-[10px] font-black uppercase tracking-widest text-slate-400 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <!-- Loading skeleton -->
              <template v-if="isLoading">
                <tr v-for="i in 5" :key="'skel-' + i">
                  <td colspan="7" class="py-10 px-10">
                    <div class="h-6 bg-slate-100 rounded-lg animate-pulse w-full"></div>
                  </td>
                </tr>
              </template>
              <!-- Data rows -->
              <template v-else>
                <tr v-for="inv in filteredInvoices" :key="inv.id" class="group hover:bg-white/60 transition-all duration-300">
                  <td class="py-6 px-10">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 rounded-lg bg-indigo-50 flex items-center justify-center">
                        <FileText class="w-4 h-4 text-indigo-500" />
                      </div>
                      <span class="text-xs font-black text-slate-900 uppercase">{{ inv.invoice_number }}</span>
                    </div>
                  </td>
                  <td class="py-6 px-6">
                    <p class="text-xs font-black text-slate-800">{{ inv.student_name }}</p>
                  </td>
                  <td class="py-6 px-6">
                    <span class="px-3 py-1 rounded-full bg-slate-100 text-slate-500 text-[9px] font-black uppercase tracking-widest">{{ inv.category }}</span>
                  </td>
                  <td class="py-6 px-6 text-right">
                    <p class="text-xs font-black text-slate-900 italic">{{ inv.amount_formatted }}</p>
                  </td>
                  <td class="py-6 px-6">
                    <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ inv.due_date }}</p>
                  </td>
                  <td class="py-6 px-6">
                    <div class="flex justify-center">
                      <span :class="getStatusClass(inv.status)" class="px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest border flex items-center gap-2">
                        <div class="w-1.5 h-1.5 rounded-full" :class="getStatusDot(inv.status)"></div>
                        {{ getStatusLabel(inv.status) }}
                      </span>
                    </div>
                  </td>
                  <td class="py-6 px-10 text-right">
                    <div class="flex items-center justify-end gap-2">
                      <button v-if="inv.payment_url" @click="openPaymentLink(inv.payment_url)" title="Bayar Sekarang" class="p-2 text-indigo-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all">
                        <ExternalLink class="w-4 h-4" />
                      </button>
                      <button @click="openEditModal(inv)" title="Edit Faktur" class="p-2 text-slate-300 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all">
                        <Edit3 class="w-4 h-4" />
                      </button>
                      <button @click="deleteInvoice(inv.id || inv.ID)" title="Hapus Faktur" class="p-2 text-slate-300 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all">
                        <Trash2 class="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>

        <!-- Empty state -->
        <div v-if="!isLoading && filteredInvoices.length === 0" class="py-20 flex flex-col items-center justify-center text-slate-300">
           <AlertCircle class="w-16 h-16 mb-4 opacity-20" />
           <p class="text-xs font-black uppercase tracking-widest">Tidak ada faktur ditemukan</p>
        </div>

        <!-- Pagination -->
        <div class="p-8 border-t border-slate-100 flex items-center justify-between">
           <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Menampilkan {{ filteredInvoices.length }} dari {{ invoices.length }} faktur</p>
           <div class="flex items-center gap-2">
              <button class="p-2 border border-slate-200 rounded-xl text-slate-400 hover:bg-slate-50 transition-all">
                <ChevronLeft class="w-4 h-4" />
              </button>
              <button class="p-2 border border-slate-200 rounded-xl text-slate-400 hover:bg-slate-50 transition-all">
                <ChevronRight class="w-4 h-4" />
              </button>
           </div>
        </div>
      </section>
    </main>

    <!-- Create Invoice Modal -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="modal-overlay" @click.self="closeCreateModal">
        <div class="modal-container">
          <!-- Modal Header -->
          <div class="p-10 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div>
              <h3 class="text-2xl font-black text-slate-900 italic uppercase">
                {{ modalMode === 'create' ? 'Buat' : 'Edit' }} <span class="text-indigo-600">Faktur Baru</span>
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Input data penagihan manual</p>
            </div>
            <button @click="closeCreateModal" class="p-3 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-2xl transition-all border border-slate-100 shadow-sm">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Modal Body -->
          <div class="p-10 space-y-8">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2 relative">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">NIS atau Nama Santri</label>
                <div class="relative">
                  <Search class="absolute left-6 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                  <input v-model="newInvoice.nis" @input="searchStudents" @blur="closeSuggestions" type="text" placeholder="Ketik NIS atau Nama..." class="w-full pl-14 pr-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold" />
                </div>
                
                <!-- Suggestions Dropdown -->
                <div v-if="showSuggestions" class="absolute z-[100] left-0 right-0 top-full mt-2 bg-white border border-slate-100 rounded-3xl shadow-2xl overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
                  <div v-for="s in studentSuggestions" :key="s.id" 
                       @click="selectStudent(s)"
                       class="p-4 hover:bg-indigo-50 cursor-pointer flex items-center justify-between border-b border-slate-50 last:border-0 transition-colors">
                    <div>
                      <p class="text-[10px] font-black text-slate-900 uppercase">{{ s.name }}</p>
                      <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">NIS: {{ s.nis }} • Kelas: {{ s.class || '-' }}</p>
                    </div>
                    <ChevronRight class="w-3 h-3 text-slate-300" />
                  </div>
                </div>
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kategori Pembayaran</label>
                <select v-model="newInvoice.payment_category_id" @change="handleCategoryChange" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 appearance-none transition-all font-bold">
                  <option value="">Pilih Kategori...</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
            </div>

            <!-- Auto-filled Student Info -->
            <div v-if="newInvoice.student_name" class="p-6 bg-indigo-50/50 rounded-[2rem] border border-indigo-100 grid grid-cols-2 gap-4 animate-in fade-in slide-in-from-top-2 duration-300">
              <div>
                <p class="text-[9px] font-black text-indigo-400 uppercase tracking-widest mb-1">Nama Santri</p>
                <p class="text-xs font-black text-slate-800 uppercase">{{ newInvoice.student_name }}</p>
              </div>
              <div>
                <p class="text-[9px] font-black text-indigo-400 uppercase tracking-widest mb-1">Nama Wali</p>
                <p class="text-xs font-black text-slate-800 uppercase">{{ newInvoice.guardian_name }}</p>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nominal (Otomatis)</label>
                <div class="relative">
                  <input :value="formatRupiah(newInvoice.amount)" type="text" readonly class="w-full px-6 py-4 bg-slate-100 border border-slate-200 rounded-[1.5rem] text-sm font-black italic text-indigo-600 text-lg cursor-not-allowed" />
                  <div class="absolute right-6 top-1/2 -translate-y-1/2 p-2 bg-white rounded-full border border-slate-100 shadow-sm">
                    <CheckCircle2 class="w-4 h-4 text-emerald-500" />
                  </div>
                </div>
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Jatuh Tempo</label>
                <input v-model="newInvoice.due_date" type="date" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all" />
              </div>
            </div>

            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Keterangan</label>
              <textarea v-model="newInvoice.description" rows="3" placeholder="Contoh: SPP Bulan Mei 2026" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all resize-none"></textarea>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="p-10 bg-slate-50/50 border-t border-slate-100 flex items-center justify-end gap-4">
            <button @click="closeCreateModal" class="px-8 py-4 rounded-2xl text-[10px] font-black text-slate-400 uppercase tracking-widest hover:bg-slate-100 transition-all">Batal</button>
            <button @click="submitInvoice" :disabled="isSubmitting" class="px-10 py-4 rounded-2xl bg-indigo-600 text-white text-[10px] font-black uppercase tracking-widest shadow-xl shadow-indigo-600/20 hover:bg-indigo-700 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3">
              <CheckCircle2 v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ isSubmitting ? 'Memproses...' : (modalMode === 'create' ? 'Terbitkan Faktur' : 'Simpan Perubahan') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Bulk Invoice Modal -->
      <div v-if="showBulkModal" class="modal-overlay" @click.self="showBulkModal = false">
        <div class="modal-container">
          <!-- Modal Header -->
          <div class="p-10 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div>
              <h3 class="text-2xl font-black text-slate-900 italic uppercase">Tagih <span class="text-indigo-600">Massal</span></h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Buat faktur untuk banyak santri sekaligus</p>
            </div>
            <button @click="showBulkModal = false" class="p-3 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-2xl transition-all border border-slate-100 shadow-sm">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Modal Body -->
          <div class="p-10 space-y-8">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kategori Pembayaran</label>
                <select v-model="bulkForm.payment_category_id" @change="handleBulkCategoryChange" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold">
                  <option value="">Pilih Kategori</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nominal (Rp)</label>
                <input v-model="bulkForm.amount" type="number" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold" />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Target Santri</label>
                <select v-model="bulkForm.target" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold">
                  <option value="all">Semua Santri Aktif</option>
                  <option value="classroom">Berdasarkan Kelas</option>
                </select>
              </div>
              <div v-if="bulkForm.target === 'classroom'" class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Pilih Kelas</label>
                <select v-model="bulkForm.classroom_id" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold">
                  <option value="">Pilih Kelas</option>
                  <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
                </select>
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Jatuh Tempo</label>
                <input v-model="bulkForm.due_date" type="date" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-[1.5rem] text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all font-bold" />
              </div>
            </div>
            
            <div class="p-6 bg-amber-50 rounded-3xl border border-amber-100 flex gap-4">
              <AlertCircle class="w-6 h-6 text-amber-500 shrink-0" />
              <p class="text-[11px] font-medium text-amber-700 leading-relaxed">
                <strong class="uppercase block mb-1">Penting:</strong>
                Tindakan ini akan membuat faktur secara otomatis untuk banyak santri sekaligus. Pastikan data kategori dan nominal sudah benar sebelum melanjutkan.
              </p>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="p-10 bg-slate-50/50 border-t border-slate-100 flex items-center justify-end gap-4">
            <button @click="showBulkModal = false" class="px-8 py-4 rounded-2xl text-[10px] font-black text-slate-400 uppercase tracking-widest hover:bg-slate-100 transition-all">Batal</button>
            <button @click="submitBulkInvoices" :disabled="isSubmitting" class="px-10 py-4 rounded-2xl bg-[#0F2942] text-white text-[10px] font-black uppercase tracking-widest shadow-xl shadow-indigo-600/20 hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3">
              <CheckCircle2 v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ isSubmitting ? 'Memproses...' : 'Proses Tagihan Massal' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
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

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(12px);
}
.modal-container {
  background: white;
  border-radius: 3rem;
  width: 100%;
  max-width: 42rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  animation: modal-in 0.3s ease-out;
}
@keyframes modal-in {
  from { opacity: 0; transform: scale(0.95) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.status-paid { background: #ecfdf5; color: #059669; border-color: #d1fae5; }
.status-unpaid { background: #fff1f2; color: #e11d48; border-color: #ffe4e6; }
.status-partial { background: #fffbeb; color: #d97706; border-color: #fef3c7; }
.status-default { background: #f8fafc; color: #475569; border-color: #e2e8f0; }

.dot-paid { background: #10b981; }
.dot-unpaid { background: #f43f5e; }
.dot-partial { background: #f59e0b; }
.dot-default { background: #94a3b8; }

::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
</style>
