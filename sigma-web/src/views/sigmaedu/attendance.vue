<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  CheckCircle2, XCircle, Calendar as CalendarIcon, 
  Users, Save, List, History, Download, Filter, Search, Trash2, RotateCw
} from 'lucide-vue-next'
import axios from 'axios'
import * as XLSX from 'xlsx'

// State
const classrooms = ref<any[]>([])
const subjects = ref<any[]>([])
const students = ref<any[]>([])
const attendances = ref<any[]>([])
const historyData = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const activeTab = ref('input') // 'input' | 'history'
const searchQuery = ref('')

const filters = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  date: new Date().toISOString().split('T')[0],
  start_date: new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0],
  end_date: new Date().toISOString().split('T')[0]
})

const filteredHistory = computed(() => {
  if (!searchQuery.value) return historyData.value
  const query = searchQuery.value.toLowerCase()
  return historyData.value.filter(h => 
    h.student?.name?.toLowerCase().includes(query) ||
    h.subject?.name?.toLowerCase().includes(query) ||
    h.student?.nis?.toLowerCase().includes(query) ||
    h.classroom?.name?.toLowerCase().includes(query)
  )
})

// UI State
const attendanceMap = ref<Record<number, string>>({}) // studentId -> status

const fetchData = async () => {
  try {
    const [clRes, sbRes] = await Promise.all([
      axios.get('/api/v1/admin/edu/classrooms'),
      axios.get('/api/v1/admin/edu/subjects')
    ])
    classrooms.value = clRes.data
    subjects.value = sbRes.data

    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
  } catch (err) {
    console.error('Gagal mengambil data master:', err)
  }
}

const loadClassData = async () => {
  if (!filters.value.classroom_id || activeTab.value !== 'input') return
  
  isLoading.value = true
  try {
    // 1. Get students in class
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data

    // 2. Get existing attendance for this date/subject
    const params: any = {
      classroom_id: filters.value.classroom_id,
      date: filters.value.date
    }
    if (filters.value.subject_id) params.subject_id = filters.value.subject_id
    
    const atRes = await axios.get('/api/v1/admin/edu/attendance', { params })
    attendances.value = atRes.data

    // 3. Map status
    const newMap: Record<number, string> = {}
    students.value.forEach(s => {
      const existing = attendances.value.find(a => a.student_id === s.ID)
      newMap[s.ID] = existing ? existing.status : 'Hadir'
    })
    attendanceMap.value = newMap

  } catch (err) {
    console.error('Gagal memuat data kelas:', err)
  } finally {
    isLoading.value = false
  }
}

const loadHistoryData = async () => {
  if (!filters.value.classroom_id || activeTab.value !== 'history') return
  
  isLoading.value = true
  try {
    const params: any = {
      classroom_id: filters.value.classroom_id,
      start_date: filters.value.start_date,
      end_date: filters.value.end_date
    }
    if (filters.value.subject_id) params.subject_id = filters.value.subject_id
    
    const res = await axios.get('/api/v1/admin/edu/attendance/report', { params })
    historyData.value = res.data
  } catch (err) {
    console.error('Gagal memuat history:', err)
  } finally {
    isLoading.value = false
  }
}

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

const handleSave = async () => {
  if (!filters.value.classroom_id) return
  
  isSaving.value = true
  try {
    const payload = students.value.map(s => ({
      student_id: s.ID,
      classroom_id: filters.value.classroom_id,
      subject_id: filters.value.subject_id ? Number(filters.value.subject_id) : null,
      date: new Date(filters.value.date).toISOString(),
      status: attendanceMap.value[s.ID]
    }))

    await axios.post('/api/v1/admin/edu/attendance/bulk', payload)
    
    toastMessage.value = 'Presensi Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)

    await loadClassData()
  } catch (err) {
    console.error('Gagal menyimpan presensi:', err)
    toastMessage.value = 'Gagal Menyimpan Presensi.'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isSaving.value = false
  }
}

const setStatusAll = (status: string) => {
  students.value.forEach(s => {
    attendanceMap.value[s.ID] = status
  })
}

const deleteAttendance = async (id: any) => {
  const targetId = id
  if (!targetId) {
    console.error('ID tidak ditemukan:', id)
    return
  }

  if (!confirm('Hapus record presensi ini?')) return
  
  try {
    await axios.delete(`/api/v1/admin/edu/attendance/${targetId}`)
    toastMessage.value = 'Presensi Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    await loadHistoryData()
  } catch (err) {
    console.error('Gagal menghapus presensi:', err)
    toastMessage.value = 'Gagal menghapus data'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  }
}

const exportHistory = () => {
  if (historyData.value.length === 0) return
  
  const data = historyData.value.map(h => ({
    Tanggal: new Date(h.date).toLocaleDateString('id-ID'),
    Santri: h.student?.name,
    NIS: h.student?.nis,
    Mata_Pelajaran: h.subject?.name || 'Harian',
    Status: h.status
  }))
  
  const ws = XLSX.utils.json_to_sheet(data)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, "History Presensi")
  XLSX.writeFile(wb, `History_Presensi_${filters.value.classroom_id}.xlsx`)
}

// Consolidated Watchers
watch(filters, () => {
  if (activeTab.value === 'input') loadClassData()
  else loadHistoryData()
}, { deep: true })

watch(activeTab, (val) => {
  if (val === 'input') loadClassData()
  else loadHistoryData()
})

const refreshData = () => {
  if (activeTab.value === 'input') loadClassData()
  else loadHistoryData()
}

onMounted(fetchData)

// Helpers
const getStatusColor = (status: string) => {
  switch (status) {
    case 'Hadir': return 'bg-emerald-500/10 text-emerald-500'
    case 'Izin': return 'bg-blue-500/10 text-blue-500'
    case 'Sakit': return 'bg-amber-500/10 text-amber-500'
    case 'Alpa': return 'bg-rose-500/10 text-rose-500'
    default: return 'bg-sigma-surface-alt text-sigma-muted'
  }
}

const stats = computed(() => {
  const counts = { Hadir: 0, Izin: 0, Sakit: 0, Alpa: 0 }
  Object.values(attendanceMap.value).forEach(v => {
    if (counts.hasOwnProperty(v)) counts[v as keyof typeof counts]++
  })
  return counts
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-emerald-500 rounded-full shadow-[0_0_15px_rgba(16,185,129,0.5)]"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Presensi</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Manajemen <span class="text-emerald-500">Kehadiran Santri</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3 bg-sigma-surface p-1.5 rounded-2xl border border-sigma-border shadow-sm">
        <button @click="activeTab = 'input'"
                class="px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                :class="activeTab === 'input' ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <List class="w-3.5 h-3.5" /> Input Presensi
        </button>
        <button @click="activeTab = 'history'"
                class="px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                :class="activeTab === 'history' ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <History class="w-3.5 h-3.5" /> Data Absensi
        </button>
      </div>
    </div>

    <!-- MAIN VIEW: RECORDING -->
    <div v-if="activeTab === 'input'" class="space-y-8 animate-in fade-in duration-500">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <!-- Filter Panel -->
        <div class="lg:col-span-4 space-y-6">
          <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-8 shadow-sm space-y-8 sticky top-8">
            <div class="space-y-6">
              <div class="flex items-center gap-3 text-emerald-500 mb-2">
                <CalendarIcon class="w-5 h-5" />
                <span class="text-xs font-black uppercase tracking-widest">Parameter Presensi</span>
              </div>

              <!-- Date Picker -->
              <div class="space-y-2">
                <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Tanggal</label>
                <input v-model="filters.date" type="date" 
                       class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight" />
              </div>

              <!-- Class Picker -->
              <div class="space-y-2">
                <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Kelas</label>
                <select v-model="filters.classroom_id" 
                        class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight appearance-none">
                  <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
                </select>
              </div>

              <!-- Subject Picker -->
              <div class="space-y-2">
                <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Mata Pelajaran (Opsional)</label>
                <select v-model="filters.subject_id" 
                        class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight appearance-none">
                  <option :value="null">Presensi Harian</option>
                  <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
                </select>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="pt-6 border-t border-sigma-border space-y-3">
              <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted mb-2 block">Setel Semua Status:</span>
              <div class="grid grid-cols-2 gap-2">
                <button @click="setStatusAll('Hadir')" class="py-2.5 rounded-xl border border-emerald-500/20 text-emerald-500 text-[9px] font-black uppercase tracking-widest hover:bg-emerald-500 hover:text-white transition-all">Hadir Semua</button>
                <button @click="setStatusAll('Izin')" class="py-2.5 rounded-xl border border-blue-500/20 text-blue-500 text-[9px] font-black uppercase tracking-widest hover:bg-blue-500 hover:text-white transition-all">Izin Semua</button>
              </div>
            </div>

            <button @click="handleSave" :disabled="isSaving || students.length === 0"
                    class="w-full py-4 bg-emerald-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-emerald-700 transition-all shadow-xl shadow-emerald-500/20 disabled:opacity-50 flex items-center justify-center gap-3 active:scale-95">
              <Save v-if="!isSaving" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              Simpan Presensi
            </button>
          </div>
        </div>

        <!-- Student List Panel -->
        <div class="lg:col-span-8 space-y-6">
          <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm">
             <div class="p-8 border-b border-sigma-border bg-sigma-surface-alt/30 flex justify-between items-center">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 rounded-2xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                    <Users class="w-5 h-5" />
                  </div>
                  <div>
                    <h3 class="text-sm font-black italic uppercase tracking-tight">Daftar Santri</h3>
                    <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Total: {{ students.length }} Santri</p>
                  </div>
                </div>
                <div class="flex gap-4">
                   <div class="text-center px-4">
                      <span class="text-[7px] font-black uppercase text-sigma-muted block">Hadir</span>
                      <span class="text-xs font-black text-emerald-500">{{ stats.Hadir }}</span>
                   </div>
                   <div class="text-center px-4">
                      <span class="text-[7px] font-black uppercase text-sigma-muted block">Izin/Sakit</span>
                      <span class="text-xs font-black text-blue-500">{{ stats.Izin + stats.Sakit }}</span>
                   </div>
                   <div class="text-center px-4">
                      <span class="text-[7px] font-black uppercase text-sigma-muted block">Alpha</span>
                      <span class="text-xs font-black text-rose-500">{{ stats.Alpa }}</span>
                   </div>
                </div>
             </div>

            <div v-if="isLoading" class="flex flex-col items-center justify-center p-40 space-y-4">
              <div class="w-10 h-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin"></div>
              <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Memuat Daftar Santri...</p>
            </div>

            <table v-else-if="students.length > 0" class="w-full border-collapse">
              <thead>
                <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
                  <th class="px-8 py-5 text-left text-[8px] font-black uppercase tracking-widest text-sigma-muted">No</th>
                  <th class="px-8 py-5 text-left text-[8px] font-black uppercase tracking-widest text-sigma-muted">Santri</th>
                  <th class="px-8 py-5 text-center text-[8px] font-black uppercase tracking-widest text-sigma-muted">Status Kehadiran</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border/50">
                <tr v-for="(s, index) in students" :key="s.ID" 
                    class="group hover:bg-emerald-500/[0.02] transition-all">
                  <td class="px-8 py-5 w-16">
                    <span class="text-xs font-black text-sigma-muted">{{ index + 1 }}</span>
                  </td>
                  <td class="px-8 py-5">
                    <div class="flex items-center gap-4">
                      <div class="w-10 h-10 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border flex items-center justify-center text-sigma-text font-black text-xs border border-sigma-border group-hover:scale-110 transition-transform shadow-sm">
                        {{ s.name.charAt(0) }}
                      </div>
                      <div>
                        <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ s.nis }}</span>
                        <span class="text-xs font-black uppercase tracking-tighter italic text-sigma-text">{{ s.name }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="px-8 py-5">
                    <div class="flex items-center justify-center gap-2">
                      <button v-for="status in ['Hadir', 'Izin', 'Sakit', 'Alpa']"
                              :key="status"
                              @click="attendanceMap[s.ID] = status"
                              class="px-4 py-2 rounded-xl text-[9px] font-black uppercase tracking-widest transition-all border"
                              :class="attendanceMap[s.ID] === status ? getStatusColor(status) + ' border-transparent' : 'bg-sigma-surface-alt text-sigma-muted border-transparent hover:border-sigma-border'">
                        {{ status }}
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>

            <div v-else class="flex flex-col items-center justify-center p-40 space-y-4 opacity-30">
              <Users class="w-16 h-16 text-sigma-muted" />
              <p class="font-black uppercase tracking-widest text-[10px]">Pilih kelas untuk memulai presensi</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MAIN VIEW: HISTORY (LIST DATA ABSENSI) -->
    <div v-else class="space-y-8 animate-in fade-in duration-500">
      <!-- Search & Filters -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center justify-between gap-8">
        <div class="flex items-center gap-6">
          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
            <select v-model="filters.classroom_id" 
                    class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50 appearance-none">
              <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
            </select>
          </div>
          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mapel</label>
            <select v-model="filters.subject_id" 
                    class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50 appearance-none">
              <option :value="null">Semua / Harian</option>
              <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
            </select>
          </div>
          <div class="space-y-1 min-w-[120px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Dari</label>
            <input v-model="filters.start_date" type="date" 
                   class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50" />
          </div>
          <div class="space-y-1 min-w-[120px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Sampai</label>
            <input v-model="filters.end_date" type="date" 
                   class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50" />
          </div>
          <div class="flex items-center gap-2">
            <div class="space-y-1 min-w-[200px]">
              <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pencarian</label>
              <div class="relative">
                <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-emerald-500" />
                <input v-model="searchQuery" type="text" placeholder="Santri/NIS..."
                       class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50" />
              </div>
            </div>
            <div class="pt-4">
              <button @click="refreshData" 
                      class="p-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-emerald-500 hover:bg-emerald-500 hover:text-white transition-all shadow-sm"
                      title="Refresh Data">
                <RotateCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
              </button>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <button @click="loadHistoryData" class="px-6 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-[10px] font-black uppercase tracking-widest text-sigma-muted hover:text-sigma-text transition-all flex items-center gap-2">
            <Filter class="w-3.5 h-3.5" /> Terapkan Filter
          </button>
          <button @click="exportHistory" class="px-8 py-3 bg-emerald-600 text-white rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-500/20 flex items-center gap-2 active:scale-95">
            <Download class="w-3.5 h-3.5" /> Export Excel
          </button>
        </div>
      </div>

      <!-- History Table -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Tanggal</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Santri</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted">Kelas & Mapel</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted text-center">Status</th>
              <th class="px-8 py-5 text-[8px] font-black uppercase tracking-widest text-sigma-muted text-right">Info</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border/50">
            <tr v-if="isLoading" class="animate-pulse">
              <td colspan="5" class="px-8 py-20 text-center text-[10px] font-black uppercase tracking-widest text-sigma-muted">
                Mengambil Riwayat...
              </td>
            </tr>
            <tr v-else-if="filteredHistory.length > 0" v-for="h in filteredHistory" :key="h.ID"
                class="hover:bg-emerald-500/[0.02] transition-all group">
              <td class="px-8 py-5">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                    <CalendarIcon class="w-4 h-4" />
                  </div>
                  <span class="text-xs font-black">{{ formatDate(h.date) }}</span>
                </div>
              </td>
              <td class="px-8 py-5">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border flex items-center justify-center text-sigma-text font-black text-xs border border-sigma-border group-hover:scale-110 transition-transform">
                    {{ h.student?.name?.charAt(0) }}
                  </div>
                  <div>
                    <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ h.student?.nis }}</span>
                    <span class="text-xs font-black uppercase tracking-tighter italic text-sigma-text">{{ h.student?.name }}</span>
                  </div>
                </div>
              </td>
              <td class="px-8 py-5">
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-sigma-text">{{ h.classroom?.name || '-' }}</span>
                  <span class="text-[9px] font-black text-emerald-500 uppercase tracking-tighter">{{ h.subject?.name || 'Presensi Harian' }}</span>
                </div>
              </td>
              <td class="px-8 py-5 text-center">
                <span class="px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest inline-block shadow-sm"
                      :class="getStatusColor(h.status)">
                  {{ h.status }}
                </span>
              </td>
              <td class="px-8 py-5 text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-all">
                   <button @click="deleteAttendance(h.ID || h.id)"
                           class="p-2 rounded-xl bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all shadow-sm"
                           title="Hapus Presensi">
                     <Trash2 class="w-4 h-4" />
                   </button>
                </div>
              </td>
            </tr>
            <tr v-else>
              <td colspan="5" class="px-8 py-20 text-center opacity-30">
                <Users class="w-12 h-12 mx-auto mb-4 text-sigma-muted" />
                <p class="font-black uppercase tracking-widest text-[10px]">Tidak ada data absensi ditemukan</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- Notification Toast -->
  <Transition name="toast">
    <div v-if="showToast" 
         class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[200] flex items-center gap-4 px-8 py-4 rounded-[2rem] border shadow-2xl backdrop-blur-2xl transition-all"
         :class="toastType === 'success' ? 'bg-emerald-500/90 border-emerald-400/50 text-white shadow-emerald-500/20' : 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/20'">
      <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
        <CheckCircle2 v-if="toastType === 'success'" class="w-5 h-5" />
        <XCircle v-else class="w-5 h-5" />
      </div>
      <span class="text-xs font-black uppercase tracking-widest">{{ toastMessage }}</span>
    </div>
  </Transition>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px) scale(0.9);
}

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(16, 185, 129, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(16, 185, 129, 0.2); }

input[type="date"]::-webkit-calendar-picker-indicator {
  filter: invert(0.5);
  cursor: pointer;
}
</style>
