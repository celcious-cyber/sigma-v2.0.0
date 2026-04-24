<script setup lang="ts">
import { ref, onMounted, watch, computed, reactive } from 'vue'
import { 
  FileText, Save, Calendar as CalendarIcon, 
  CheckCircle2, XCircle,
  LayoutGrid, List, BookOpen
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const subjects = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const activeTab = ref('input') // 'input' or 'rekap'

const filters = ref({
  classroom_id: null as number | null,
  date: new Date().toISOString().split('T')[0]
})

// UI State
const memorizationMap = reactive<Record<number, any>>({})

const fetchData = async () => {
  try {
    const [clRes, subRes] = await Promise.all([
      axios.get('/api/v1/admin/edu/classrooms'),
      axios.get('/api/v1/admin/edu/subjects')
    ])
    classrooms.value = clRes.data
    subjects.value = subRes.data
    
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
    await loadMemorizationData()
  } catch (err) {
    console.error('Gagal mengambil data awal:', err)
  }
}

const loadMemorizationData = async () => {
  if (!filters.value.classroom_id) return
  
  isLoading.value = true
  try {
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data

    const params = {
      classroom_id: filters.value.classroom_id,
      date: filters.value.date
    }
    const res = await axios.get('/api/v1/admin/edu/memorization', { params })
    
    // Clear existing keys
    Object.keys(memorizationMap).forEach(key => delete memorizationMap[Number(key)])

    students.value.forEach(s => {
      const existing = res.data.find((r: any) => r.student_id === s.ID)
      memorizationMap[s.ID] = existing ? {
        subject_name: existing.subject_name,
        title: existing.title,
        grade: existing.grade,
        remarks: existing.remarks
      } : {
        subject_name: '',
        title: '',
        grade: 'A',
        remarks: ''
      }
    })
  } catch (err) {
    console.error('Gagal memuat data hafalan:', err)
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
      teacher_id: 0,
      subject_name: memorizationMap[s.ID].subject_name,
      title: memorizationMap[s.ID].title,
      grade: memorizationMap[s.ID].grade,
      remarks: memorizationMap[s.ID].remarks,
      date: new Date(filters.value.date).toISOString()
    }))

    await axios.post('/api/v1/admin/edu/memorization/bulk', payload)
    
    toastMessage.value = 'Hafalan Pelajaran Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    
    await loadMemorizationData()
  } catch (err) {
    console.error('Gagal menyimpan hafalan:', err)
    toastMessage.value = 'Gagal Menyimpan Data.'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isSaving.value = false
  }
}

// Stats for Header
const stats = computed(() => {
  const data = Object.values(memorizationMap)
  const total = data.filter((d: any) => d.subject_name !== '' || d.title !== '').length
  return { total }
})

watch([() => filters.value.classroom_id, () => filters.value.date], loadMemorizationData)
onMounted(fetchData)
</script>

<template>
  <div class="memorization-page">
    <div class="p-8 space-y-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="h-10 w-1 bg-indigo-500 rounded-full shadow-[0_0_15px_rgba(99,102,241,0.5)]"></div>
          <div>
            <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Hafalan</h2>
            <h1 class="text-2xl font-black italic uppercase tracking-tight">Hafalan <span class="text-indigo-500">Mata Pelajaran</span></h1>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <div class="bg-sigma-surface-alt border border-sigma-border p-1 rounded-2xl flex items-center gap-1 mr-4">
            <button @click="activeTab = 'input'"
                    class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                    :class="activeTab === 'input' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-500/20' : 'text-sigma-muted hover:text-sigma-text'">
              <LayoutGrid class="w-3.5 h-3.5" /> Input
            </button>
            <button @click="activeTab = 'rekap'"
                    class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                    :class="activeTab === 'rekap' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-500/20' : 'text-sigma-muted hover:text-sigma-text'">
              <List class="w-3.5 h-3.5" /> Rekap
            </button>
          </div>

          <button v-if="activeTab === 'input'" @click="handleSave" :disabled="isSaving || students.length === 0"
                  class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/20 disabled:opacity-50 flex items-center gap-3 active:scale-95">
            <Save v-if="!isSaving" class="w-4 h-4" />
            <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            Simpan Hafalan
          </button>
        </div>
      </div>

      <!-- Quick Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] shadow-sm flex items-center gap-4">
          <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center text-indigo-500">
            <FileText class="w-6 h-6" />
          </div>
          <div>
            <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">Total Input</span>
            <span class="text-xl font-black">{{ stats.total }} <span class="text-[10px] text-sigma-muted">Santri</span></span>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center gap-8">
        <div class="flex items-center gap-6 flex-1">
          <div class="space-y-1 min-w-[200px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal Hafalan</label>
            <div class="relative">
              <CalendarIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-indigo-500" />
              <input v-model="filters.date" type="date" 
                     class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-indigo-500/50" />
            </div>
          </div>
          <div class="space-y-1 min-w-[200px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Kelas</label>
            <select v-model="filters.classroom_id" 
                    class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-indigo-500/50">
              <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Main Table -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm">
        <div v-if="isLoading" class="flex flex-col items-center justify-center p-40 space-y-4">
          <div class="w-10 h-10 border-4 border-indigo-500/20 border-t-indigo-500 rounded-full animate-spin"></div>
          <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Memuat Data Hafalan...</p>
        </div>

        <table v-else-if="students.length > 0" class="w-full border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
              <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Santri</th>
              <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Detail Hafalan</th>
              <th class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted w-32">Nilai</th>
              <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Catatan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in students" :key="s.ID" 
                class="border-b border-sigma-border last:border-0 group hover:bg-sigma-surface-alt/10 transition-all">
              <td class="p-6">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-2xl bg-sigma-surface border border-sigma-border flex items-center justify-center text-xs font-black group-hover:scale-110 transition-transform shadow-sm">
                    {{ s.name.charAt(0) }}
                  </div>
                  <div>
                    <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ s.nis }}</span>
                    <span class="text-sm font-black uppercase tracking-tighter italic text-sigma-text">{{ s.name }}</span>
                  </div>
                </div>
              </td>
              <td class="p-6">
                <div class="flex items-center gap-4">
                  <div class="flex flex-col gap-2 flex-1">
                    <div class="relative">
                      <BookOpen class="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-sigma-muted" />
                      <select v-model="memorizationMap[s.ID].subject_name"
                              class="w-full pl-9 pr-4 py-2.5 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-indigo-500/50">
                        <option value="">Pilih Pelajaran...</option>
                        <option v-for="sub in subjects" :key="sub.ID" :value="sub.name">{{ sub.name }}</option>
                      </select>
                    </div>
                    <div class="relative">
                      <FileText class="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-sigma-muted" />
                      <input v-model="memorizationMap[s.ID].title" type="text"
                             class="w-full pl-9 pr-4 py-2.5 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold focus:border-indigo-500/50 outline-none"
                             placeholder="Judul Hafalan..." />
                    </div>
                  </div>
                </div>
              </td>
              <td class="p-6 text-center">
                <div class="flex items-center justify-center gap-1 bg-sigma-surface-alt/50 p-1 rounded-2xl border border-sigma-border/50">
                  <button v-for="grade in ['A', 'B', 'C', 'D']" 
                          :key="grade"
                          @click="memorizationMap[s.ID].grade = grade"
                          class="w-8 h-8 rounded-xl text-[10px] font-black transition-all border flex items-center justify-center"
                          :class="memorizationMap[s.ID].grade === grade 
                            ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-500/20' 
                            : 'bg-transparent text-sigma-muted border-transparent hover:border-sigma-border hover:text-sigma-text'">
                    {{ grade }}
                  </button>
                </div>
              </td>
              <td class="p-6">
                <input v-model="memorizationMap[s.ID].remarks" type="text"
                       class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-indigo-500/50"
                       placeholder="Catatan..." />
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-else class="flex flex-col items-center justify-center p-40 space-y-4 opacity-30">
          <FileText class="w-16 h-16" />
          <p class="font-black uppercase tracking-widest text-[10px]">Pilih kelas untuk memulai pencatatan hafalan</p>
        </div>
      </div>

      <!-- Notification Toast -->
      <Transition name="toast">
        <div v-if="showToast" 
             class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[200] flex items-center gap-4 px-8 py-4 rounded-[2rem] border shadow-2xl backdrop-blur-2xl transition-all"
             :class="toastType === 'success' ? 'bg-indigo-500/90 border-indigo-400/50 text-white shadow-indigo-500/20' : 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/20'">
          <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
            <CheckCircle2 v-if="toastType === 'success'" class="w-5 h-5" />
            <XCircle v-else class="w-5 h-5" />
          </div>
          <span class="text-xs font-black uppercase tracking-widest">{{ toastMessage }}</span>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }
</style>
