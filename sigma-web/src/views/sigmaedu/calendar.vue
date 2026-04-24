<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Plus, Calendar as CalendarIcon, ArrowLeft, 
  Check, X, Trash2, Pencil, AlertCircle,
  ChevronLeft, ChevronRight, Filter, LayoutGrid, List
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// State
const events = ref<any[]>([])
const isLoading = ref(true)
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingId = ref<number | null>(null)
const viewMode = ref<'grid' | 'list'>('grid')

// Calendar Grid State
const currentMonth = ref(new Date().getMonth())
const currentYear = ref(new Date().getFullYear())
const daysOfWeek = ['Ahad', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']

const form = ref({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
  type: 'Kegiatan', 
  color: '#3b82f6'
})

const eventTypes = [
  { name: 'Kegiatan', color: '#3b82f6', bg: 'bg-blue-500/10', text: 'text-blue-500' },
  { name: 'Libur', color: '#ef4444', bg: 'bg-red-500/10', text: 'text-red-500' },
  { name: 'Ujian', color: '#f59e0b', bg: 'bg-amber-500/10', text: 'text-amber-500' },
  { name: 'Penting', color: '#8b5cf6', bg: 'bg-purple-500/10', text: 'text-purple-500' }
]

const fetchEvents = async () => {
  try {
    const response = await axios.get('/api/v1/admin/edu/calendar')
    events.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data kalender:', err)
  } finally {
    isLoading.value = false
  }
}

// Calendar Logic
const daysInMonth = computed(() => {
  const firstDay = new Date(currentYear.value, currentMonth.value, 1).getDay()
  const totalDays = new Date(currentYear.value, currentMonth.value + 1, 0).getDate()
  
  const days = []
  // Previous month days (padding)
  const prevMonthLastDay = new Date(currentYear.value, currentMonth.value, 0).getDate()
  for (let i = firstDay - 1; i >= 0; i--) {
    days.push({ day: prevMonthLastDay - i, month: currentMonth.value - 1, year: currentYear.value, current: false })
  }
  
  // Current month days
  for (let i = 1; i <= totalDays; i++) {
    days.push({ day: i, month: currentMonth.value, year: currentYear.value, current: true })
  }
  
  // Next month days (padding)
  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    days.push({ day: i, month: currentMonth.value + 1, year: currentYear.value, current: false })
  }
  
  return days
})

const monthName = computed(() => {
  return new Intl.DateTimeFormat('id-ID', { month: 'long' }).format(new Date(currentYear.value, currentMonth.value))
})

const nextMonth = () => {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

const prevMonth = () => {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

const getEventsForDay = (day: number, month: number, year: number) => {
  return events.value.filter(e => {
    const start = new Date(e.start_date)
    return start.getDate() === day && start.getMonth() === month && start.getFullYear() === year
  })
}

const openModal = (event: any = null, initialDate: string | null = null) => {
  if (event) {
    editingId.value = event.ID
    form.value = { 
      title: event.title, 
      description: event.description || '', 
      start_date: event.start_date ? event.start_date.split('T')[0] : '', 
      end_date: event.end_date ? event.end_date.split('T')[0] : '',
      type: event.type || 'Kegiatan',
      color: event.color || '#3b82f6'
    }
  } else {
    editingId.value = null
    const date = initialDate || new Date().toISOString().split('T')[0]
    form.value = { 
      title: '', 
      description: '', 
      start_date: date, 
      end_date: date, 
      type: 'Kegiatan', 
      color: '#3b82f6' 
    }
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const handleSubmit = async () => {
  isSubmitting.value = true
  try {
    const payload = {
      ...form.value,
      start_date: new Date(form.value.start_date).toISOString(),
      end_date: form.value.end_date ? new Date(form.value.end_date).toISOString() : null
    }

    if (editingId.value) {
      await axios.put(`/api/v1/admin/edu/calendar/${editingId.value}`, payload)
    } else {
      await axios.post('/api/v1/admin/edu/calendar', payload)
    }
    await fetchEvents()
    closeModal()
  } catch (err) {
    alert('Gagal menyimpan data kalender.')
  } finally {
    isSubmitting.value = false
  }
}

const deleteEvent = async (id: number) => {
  if (!confirm('Yakin ingin menghapus agenda ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/calendar/${id}`)
    await fetchEvents()
  } catch (err) {
    alert('Gagal menghapus data.')
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const getEventTypeData = (type: string) => {
  return eventTypes.find(t => t.name === type) || eventTypes[0]
}

onMounted(() => {
  fetchEvents()
})
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-blue-500 rounded-full"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Kurikulum</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Kalender <span class="text-blue-500">Akademik</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <div class="flex bg-sigma-surface p-1 rounded-2xl border border-sigma-border shadow-sm">
          <button @click="viewMode = 'grid'" 
                  :class="['p-2 rounded-xl transition-all', viewMode === 'grid' ? 'bg-blue-500 text-white shadow-lg shadow-blue-500/20' : 'text-sigma-muted hover:text-sigma-text']">
            <LayoutGrid class="w-4 h-4" />
          </button>
          <button @click="viewMode = 'list'" 
                  :class="['p-2 rounded-xl transition-all', viewMode === 'list' ? 'bg-blue-500 text-white shadow-lg shadow-blue-500/20' : 'text-sigma-muted hover:text-sigma-text']">
            <List class="w-4 h-4" />
          </button>
        </div>
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-blue-500/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Calendar View -->
    <div v-if="viewMode === 'grid'" class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm animate-in fade-in zoom-in duration-500">
      <!-- Calendar Header -->
      <div class="p-8 flex items-center justify-between border-b border-sigma-border bg-sigma-surface-alt/20">
        <div class="flex items-center gap-4">
          <h2 class="text-2xl font-black uppercase italic tracking-tighter">{{ monthName }} <span class="text-blue-500 not-italic">{{ currentYear }}</span></h2>
        </div>
        <div class="flex items-center gap-2">
          <button @click="prevMonth" class="p-3 rounded-xl bg-sigma-surface border border-sigma-border hover:border-blue-500/50 transition-all">
            <ChevronLeft class="w-5 h-5" />
          </button>
          <button @click="currentMonth = new Date().getMonth(); currentYear = new Date().getFullYear()" 
                  class="px-6 py-2.5 rounded-xl bg-sigma-surface border border-sigma-border text-[10px] font-black uppercase tracking-widest hover:border-blue-500/50 transition-all">
            Hari Ini
          </button>
          <button @click="nextMonth" class="p-3 rounded-xl bg-sigma-surface border border-sigma-border hover:border-blue-500/50 transition-all">
            <ChevronRight class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Days Header -->
      <div class="grid grid-cols-7 border-b border-sigma-border">
        <div v-for="day in daysOfWeek" :key="day" 
             class="p-4 text-center text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">
          {{ day }}
        </div>
      </div>

      <!-- Calendar Grid -->
      <div class="grid grid-cols-7">
        <div v-for="(date, idx) in daysInMonth" :key="idx" 
             class="min-h-[140px] p-4 border-r border-b border-sigma-border last:border-r-0 transition-all hover:bg-sigma-surface-alt/30 group relative"
             :class="[!date.current ? 'opacity-30 bg-sigma-surface-alt/10' : '']">
          
          <div class="flex justify-between items-start mb-2">
            <span :class="['text-sm font-black', date.current ? 'text-sigma-text' : 'text-sigma-muted']">{{ date.day }}</span>
            <button @click="openModal(null, `${date.year}-${String(date.month + 1).padStart(2, '0')}-${String(date.day).padStart(2, '0')}`)" 
                    class="p-1 rounded-lg hover:bg-blue-500 hover:text-white transition-all opacity-0 group-hover:opacity-100 text-sigma-muted">
              <Plus class="w-3 h-3" />
            </button>
          </div>

          <div class="space-y-1 overflow-y-auto max-h-[100px] custom-scrollbar">
            <div v-for="event in getEventsForDay(date.day, date.month, date.year)" :key="event.ID"
                 @click="openModal(event)"
                 :style="{ backgroundColor: event.color + '20', color: event.color, borderLeft: `3px solid ${event.color}` }"
                 class="px-2 py-1.5 rounded-md text-[9px] font-black uppercase tracking-tight truncate cursor-pointer hover:brightness-95 transition-all">
              {{ event.title }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Agenda View -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div class="lg:col-span-8 space-y-6">
        <div class="flex flex-col md:flex-row gap-4 justify-between items-center bg-sigma-surface p-4 rounded-[2rem] border border-sigma-border shadow-sm">
          <div class="flex items-center gap-2 px-4 py-2 bg-sigma-surface-alt rounded-xl border border-sigma-border">
            <CalendarIcon class="w-4 h-4 text-blue-500" />
            <span class="text-xs font-black uppercase tracking-widest text-sigma-text">Daftar Agenda</span>
          </div>
          <button @click="openModal()" 
                  class="w-full md:w-auto px-8 py-3.5 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center justify-center gap-3">
            <Plus class="w-4 h-4" /> Tambah Agenda
          </button>
        </div>

        <div class="space-y-4">
          <div v-for="event in events" :key="event.ID" 
               class="bg-sigma-surface border border-sigma-border p-6 rounded-[2.5rem] flex items-center gap-6 group hover:border-blue-500/30 transition-all shadow-sm relative overflow-hidden">
            <div class="absolute left-0 top-0 bottom-0 w-2" :style="{ backgroundColor: event.color }"></div>
            <div class="flex flex-col items-center justify-center w-20 h-20 rounded-3xl bg-sigma-surface-alt border border-sigma-border">
              <span class="text-2xl font-black text-sigma-text">{{ new Date(event.start_date).getDate() }}</span>
              <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ new Date(event.start_date).toLocaleDateString('id-ID', { month: 'short' }) }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-3 mb-1">
                <span :class="['px-3 py-1 rounded-full text-[8px] font-black uppercase tracking-widest border', getEventTypeData(event.type).bg, getEventTypeData(event.type).text, 'border-' + getEventTypeData(event.type).text + '/20']">
                  {{ event.type }}
                </span>
                <span v-if="event.end_date" class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest">
                  s/d {{ formatDate(event.end_date) }}
                </span>
              </div>
              <h3 class="text-lg font-black uppercase italic tracking-tighter truncate text-sigma-text group-hover:text-blue-500 transition-colors">
                {{ event.title }}
              </h3>
              <p class="text-xs text-sigma-muted line-clamp-1 mt-1">{{ event.description || 'Tidak ada deskripsi' }}</p>
            </div>
            <div class="flex items-center gap-2">
              <button @click="openModal(event)" class="p-3 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-blue-500 hover:bg-blue-500/5 transition-all border border-sigma-border">
                <Pencil class="w-4 h-4" />
              </button>
              <button @click="deleteEvent(event.ID)" class="p-3 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-red-500 hover:bg-red-500/5 transition-all border border-sigma-border">
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <div class="bg-gradient-to-br from-blue-600 to-indigo-900 p-8 rounded-[3rem] text-white space-y-6 shadow-xl shadow-blue-500/20 relative overflow-hidden">
          <div class="absolute top-0 right-0 -translate-y-1/2 translate-x-1/4 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
          <div class="relative z-10 space-y-4">
            <div class="w-12 h-12 rounded-2xl bg-white/10 flex items-center justify-center">
              <AlertCircle class="w-6 h-6" />
            </div>
            <h3 class="text-xl font-black uppercase italic tracking-tighter">Ringkasan Kalender</h3>
            <p class="text-blue-100/60 text-xs font-medium leading-relaxed uppercase tracking-widest">
              Gunakan kalender akademik untuk mengatur hari libur nasional, jadwal ujian, dan kegiatan besar pesantren lainnya.
            </p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div v-for="type in eventTypes" :key="type.name" class="p-4 rounded-2xl bg-white/5 border border-white/10">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-2 h-2 rounded-full" :style="{ backgroundColor: type.color }"></div>
                <span class="text-[10px] font-black uppercase tracking-widest">{{ type.name }}</span>
              </div>
              <p class="text-lg font-black">{{ events.filter(e => e.type === type.name).length }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Add/Edit -->
    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-xl rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-gradient-to-r from-blue-600 to-blue-900 flex justify-between items-center border-b border-white/10">
          <div>
            <h3 class="text-white font-black uppercase tracking-widest text-xs">{{ editingId ? 'Update' : 'Tambah' }} Agenda</h3>
            <p class="text-blue-100/60 text-[10px] mt-1 font-medium">Manajemen Kalender Akademik Sigmaedu</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-6 max-h-[70vh] overflow-y-auto custom-scrollbar">
          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Judul Agenda</label>
            <input v-model="form.title" type="text" placeholder="e.g. Ujian Akhir Semester" 
                   class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Keterangan</label>
            <textarea v-model="form.description" rows="3" placeholder="Deskripsi singkat kegiatan..." 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight resize-none"></textarea>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Mulai</label>
              <input v-model="form.start_date" type="date" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Selesai (Opsional)</label>
              <input v-model="form.end_date" type="date" 
                     class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Tipe Agenda</label>
              <select v-model="form.type" 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
                <option v-for="t in eventTypes" :key="t.name" :value="t.name">{{ t.name }}</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Label Warna</label>
              <div class="flex items-center gap-2">
                <input v-model="form.color" type="color" 
                       class="w-12 h-13 border-none bg-transparent cursor-pointer" />
                <span class="text-xs font-mono font-bold text-sigma-muted uppercase">{{ form.color }}</span>
              </div>
            </div>
          </div>

          <div class="flex gap-4 pt-4">
            <button @click="closeModal" class="flex-1 px-8 py-4 bg-sigma-surface-alt border border-sigma-border text-sigma-text rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-border transition-all">
              Batal
            </button>
            <button @click="handleSubmit" :disabled="isSubmitting"
                    class="flex-1 px-8 py-4 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 disabled:opacity-50 flex items-center justify-center gap-2">
              <Check v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ editingId ? 'Simpan Perubahan' : 'Tambah Agenda' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }

input[type="color"]::-webkit-color-swatch-wrapper { padding: 0; }
input[type="color"]::-webkit-color-swatch { border: 2px solid rgba(0,0,0,0.1); border-radius: 12px; }
</style>
