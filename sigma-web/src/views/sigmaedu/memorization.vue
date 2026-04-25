<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  Save, 
  CheckCircle2, XCircle,
  FileText, Search, Plus, Pencil, Trash2
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const subjects = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const searchQuery = ref('')

const filters = ref({
  classroom_id: null as number | null,
  date: new Date().toISOString().split('T')[0]
})

// UI State
const memorizationMatrix = ref<Record<number, any>>({})
const studentHistory = ref<Record<number, any[]>>({})
const loadingHistory = ref<Record<number, boolean>>({})

// Modals
const showPersonalModal = ref(false)
const showBulkModal = ref(false)
const selectedStudent = ref<any>(null)

const personalForm = ref({
  subject_name: '',
  title: '',
  grade: 'A',
  remarks: '',
  date: ''
})

const bulkForm = ref({
  classroom_id: null as number | null,
  date: '',
  subject_name: '',
  entries: {} as Record<number, any>
})

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
    
    // Process current records into a matrix
    const matrix: Record<number, any> = {}
    res.data.forEach((r: any) => {
      matrix[r.student_id] = r
    })
    memorizationMatrix.value = matrix

    // Fetch history for each student
    students.value.forEach(s => fetchStudentHistory(s.ID))
  } catch (err) {
    console.error('Gagal memuat data hafalan:', err)
  } finally {
    isLoading.value = false
  }
}

const fetchStudentHistory = async (studentId: number) => {
  loadingHistory.value[studentId] = true
  try {
    const res = await axios.get(`/api/v1/admin/edu/memorization/history/${studentId}`)
    studentHistory.value[studentId] = res.data
  } catch (err) {
    console.error(`Gagal fetch history ${studentId}:`, err)
  } finally {
    loadingHistory.value[studentId] = false
  }
}

// Rekap Data computed for search
const rekapData = computed(() => {
  let filtered = students.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    filtered = filtered.filter(s => s.name.toLowerCase().includes(q) || s.nis.includes(q))
  }
  return filtered.map(s => ({
    student: s,
    record: memorizationMatrix.value[s.ID] || null
  }))
})

// Modal Actions
const openPersonalModal = (student: any) => {
  selectedStudent.value = student
  const existing = memorizationMatrix.value[student.ID]
  personalForm.value = existing ? {
    subject_name: existing.subject_name,
    title: existing.title,
    grade: existing.grade,
    remarks: existing.remarks,
    date: filters.value.date
  } : {
    subject_name: '',
    title: '',
    grade: 'A',
    remarks: '',
    date: filters.value.date
  }
  showPersonalModal.value = true
}

const savePersonalMemorization = async () => {
  isSaving.value = true
  try {
    const payload = [{
      student_id: selectedStudent.value.ID,
      subject_name: personalForm.value.subject_name,
      title: personalForm.value.title,
      grade: personalForm.value.grade,
      remarks: personalForm.value.remarks,
      date: new Date(personalForm.value.date).toISOString()
    }]
    await axios.post('/api/v1/admin/edu/memorization/bulk', payload)
    await loadMemorizationData()
    showPersonalModal.value = false
    triggerToast('Hafalan Berhasil Disimpan!')
  } catch (err) {
    triggerToast('Gagal Menyimpan Data.', 'error')
  } finally {
    isSaving.value = false
  }
}

const bulkStudents = ref<any[]>([])
const openBulkModal = async () => {
  bulkForm.value.date = filters.value.date
  showBulkModal.value = true
  bulkForm.value.classroom_id = filters.value.classroom_id
  
  if (filters.value.classroom_id) {
    try {
      const res = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
      bulkStudents.value = res.data
      const entries: Record<number, any> = {}
      bulkStudents.value.forEach(s => {
        entries[s.ID] = { subject: '', title: '', grade: 'A', remarks: '' }
      })
      bulkForm.value.entries = entries
    } catch (err) {}
  }
}

watch(() => bulkForm.value.classroom_id, async (newId) => {
  if (!showBulkModal.value || !newId) return
  try {
    const res = await axios.get(`/api/v1/admin/edu/classrooms/${newId}/students`)
    bulkStudents.value = res.data
    const entries: Record<number, any> = {}
    bulkStudents.value.forEach(s => {
      entries[s.ID] = { subject: bulkForm.value.subject_name, title: '', grade: 'A', remarks: '' }
    })
    bulkForm.value.entries = entries
  } catch (err) {}
})

const saveBulkMemorization = async () => {
  isSaving.value = true
  try {
    const payload = bulkStudents.value.map(s => {
      const e = bulkForm.value.entries[s.ID]
      return {
        student_id: s.ID,
        subject_name: e.subject || bulkForm.value.subject_name,
        title: e.title,
        grade: e.grade,
        remarks: e.remarks,
        date: new Date(bulkForm.value.date).toISOString()
      }
    }).filter(e => e.subject_name !== '' && e.title !== '')

    if (payload.length === 0) {
      triggerToast('Pilih Pelajaran & Judul minimal 1 santri', 'error')
      return
    }

    await axios.post('/api/v1/admin/edu/memorization/bulk', payload)
    filters.value.classroom_id = bulkForm.value.classroom_id
    filters.value.date = bulkForm.value.date
    await loadMemorizationData()
    showBulkModal.value = false
    triggerToast('Hafalan Batch Berhasil Disimpan!')
  } catch (err) {
    triggerToast('Gagal Menyimpan Data.', 'error')
  } finally {
    isSaving.value = false
  }
}

const deleteRecord = async (id: number, studentId: number) => {
  if (!confirm('Hapus record hafalan ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/memorization/${id}`)
    await loadMemorizationData()
    await fetchStudentHistory(studentId)
    triggerToast('Record berhasil dihapus')
  } catch (err) {
    triggerToast('Gagal menghapus record', 'error')
  }
}

// Toast
const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')
const triggerToast = (msg: string, type: string = 'success') => {
  toastMessage.value = msg
  toastType.value = type
  showToast.value = true
  setTimeout(() => showToast.value = false, 3000)
}

const getGradeColor = (grade: string | null) => {
  if (!grade) return 'bg-sigma-border text-sigma-muted'
  switch (grade.toUpperCase()) {
    case 'A': return 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/20'
    case 'B': return 'bg-blue-500/10 text-blue-600 border border-blue-500/20'
    case 'C': return 'bg-amber-500/10 text-amber-600 border border-amber-500/20'
    case 'D': return 'bg-rose-500/10 text-rose-600 border border-rose-500/20'
    default: return 'bg-sigma-border text-sigma-muted'
  }
}

watch([() => filters.value.classroom_id, () => filters.value.date], loadMemorizationData)
onMounted(fetchData)
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-sigma-bg">
    
    <!-- Top Navigation Header -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
      <div class="flex items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-gradient-to-br from-orange-500 to-orange-700 flex items-center justify-center text-white shadow-2xl shadow-orange-500/40 transform -rotate-6 hover:rotate-0 transition-transform duration-500">
          <FileText class="w-8 h-8" />
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-sigma-muted mb-1 opacity-60">Akademik • Hafalan</h2>
          <h1 class="text-3xl font-black italic uppercase tracking-tighter leading-none">
            Hafalan <span class="text-orange-500">Mata Pelajaran</span>
          </h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <select v-model="filters.classroom_id" class="px-5 py-3 bg-white border border-sigma-border rounded-2xl text-[11px] font-bold outline-none focus:border-orange-500 shadow-sm cursor-pointer">
          <option :value="null">Semua Kelas</option>
          <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
        </select>
        <input v-model="filters.date" type="date" class="px-5 py-3 bg-white border border-sigma-border rounded-2xl text-[11px] font-bold outline-none focus:border-orange-500 shadow-sm cursor-pointer" />
        <button @click="openBulkModal"
                class="px-8 py-3 bg-orange-500 text-white rounded-2xl font-black uppercase tracking-wider text-[9px] hover:shadow-xl hover:shadow-orange-500/40 transition-all active:scale-95 flex items-center gap-2">
          <Plus class="w-3.5 h-3.5" /> Input Kolektif
        </button>
      </div>
    </div>

    <!-- Content Area -->
    <div class="flex flex-col relative space-y-8">
      
      <!-- Global Search Bar -->
      <div class="relative w-full">
        <Search class="absolute left-6 top-1/2 -translate-y-1/2 w-6 h-6 text-orange-500" />
        <input v-model="searchQuery" type="text" placeholder="Ketik Nama Santri..."
               class="w-full pl-16 pr-8 py-6 bg-white border border-sigma-border rounded-[2.5rem] text-base font-bold outline-none focus:border-orange-500 shadow-xl transition-all placeholder:opacity-50" />
      </div>

      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-40 space-y-6">
        <div class="w-16 h-16 border-4 border-orange-500/20 border-t-orange-500 rounded-full animate-spin"></div>
        <p class="text-[11px] font-black uppercase tracking-[0.3em] text-sigma-muted">Memuat Data...</p>
      </div>

      <!-- JOURNAL LIST -->
      <div v-if="students.length > 0" class="flex-1 grid grid-cols-1 lg:grid-cols-2 gap-5 overflow-y-auto custom-scrollbar items-start">
        <div v-for="r in rekapData" :key="r.student.ID" 
             class="bg-white border border-sigma-border rounded-[2rem] p-6 space-y-5 hover:shadow-xl transition-all group relative overflow-hidden">
           
           <!-- Card Background Accent -->
           <div class="absolute top-0 right-0 w-80 h-80 bg-orange-500/5 rounded-full -mr-40 -mt-40 blur-3xl group-hover:bg-orange-500/10 transition-colors"></div>

           <!-- Student Info Header -->
           <div class="flex items-center justify-between gap-4 relative z-10">
              <div class="flex items-center gap-4">
                 <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border border border-sigma-border flex items-center justify-center text-base font-black group-hover:scale-110 transition-transform shadow-md">
                   {{ r.student.name.charAt(0) }}
                 </div>
                 <div>
                   <div class="flex items-center gap-2 mb-1">
                      <span class="px-3 py-1 bg-orange-500/10 text-orange-600 rounded-full text-[9px] font-black uppercase tracking-wider">{{ r.student.nis }}</span>
                      <span v-if="r.student.classroom" class="px-3 py-1 bg-blue-500/10 text-blue-600 rounded-full text-[9px] font-black uppercase tracking-wider">{{ r.student.classroom.name }}</span>
                   </div>
                   <h3 class="text-base font-black uppercase italic tracking-tight leading-none">{{ r.student.name }}</h3>
                 </div>
              </div>

              <button @click="openPersonalModal(r.student)" 
                      class="px-5 py-2.5 bg-white border border-sigma-border rounded-xl text-orange-500 font-black uppercase tracking-wider text-[9px] hover:bg-orange-500 hover:text-white transition-all shadow-sm flex items-center gap-2 shrink-0">
                <Plus class="w-3.5 h-3.5" /> Input
              </button>
           </div>

           <!-- Records Table -->
           <div class="relative z-10">
              <div class="flex items-center justify-between mb-3 px-2">
                 <h4 class="text-[9px] font-black uppercase tracking-[0.3em] text-sigma-muted">Riwayat Hafalan</h4>
              </div>

              <div class="overflow-hidden rounded-2xl border border-sigma-border shadow-sm">
                 <table class="w-full border-collapse text-xs">
                    <thead>
                       <tr class="bg-sigma-surface-alt/50">
                          <th class="px-3 py-2.5 text-left text-[9px] font-black uppercase tracking-wider text-sigma-muted">Tanggal</th>
                          <th class="px-3 py-2.5 text-left text-[9px] font-black uppercase tracking-wider text-sigma-muted">Materi</th>
                          <th class="px-3 py-2.5 text-center text-[9px] font-black uppercase tracking-wider text-sigma-muted">Grade</th>
                          <th class="px-3 py-2.5 text-center text-[9px] font-black uppercase tracking-wider text-sigma-muted">Aksi</th>
                       </tr>
                    </thead>
                    <tbody class="divide-y divide-sigma-border/40">
                        <!-- Current Record (if any) -->
                        <tr v-if="r.record" class="bg-orange-400/30">
                           <td class="px-3 py-2">
                              <div class="flex items-center gap-1.5">
                                <span class="w-1.5 h-1.5 rounded-full bg-orange-500 animate-pulse"></span>
                                <span class="text-[10px] font-bold text-orange-700">Hari Ini</span>
                              </div>
                           </td>
                           <td class="px-3 py-2">
                              <div class="flex flex-col">
                                 <span class="text-[9px] font-black text-orange-600 uppercase">{{ r.record.subject_name }}</span>
                                 <span class="text-[11px] font-bold tracking-tight">{{ r.record.title }}</span>
                              </div>
                           </td>
                           <td class="px-3 py-2 text-center">
                              <div class="inline-flex w-7 h-7 rounded-lg items-center justify-center text-[10px] font-black" :class="getGradeColor(r.record.grade)">
                                 {{ r.record.grade }}
                              </div>
                           </td>
                           <td class="px-3 py-2 text-center">
                              <div class="flex items-center justify-center gap-1">
                                <button @click="openPersonalModal(r.student)" class="w-6 h-6 rounded-md bg-blue-500/10 text-blue-500 inline-flex items-center justify-center hover:bg-blue-500 hover:text-white transition-all">
                                  <Pencil class="w-3 h-3" />
                                </button>
                                <button @click="deleteRecord(r.record.ID, r.student.ID)" class="w-6 h-6 rounded-md bg-rose-500/10 text-rose-500 inline-flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all">
                                  <Trash2 class="w-3 h-3" />
                                </button>
                              </div>
                           </td>
                        </tr>

                        <!-- History Records -->
                        <template v-if="studentHistory[r.student.ID]">
                           <tr v-for="log in studentHistory[r.student.ID].filter(l => !r.record || l.ID !== r.record.ID)" :key="log.ID" class="hover:bg-orange-500/[0.02] transition-colors">
                             <td class="px-3 py-2">
                                <span class="text-[10px] font-bold text-sigma-text">{{ new Date(log.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' }) }}</span>
                             </td>
                             <td class="px-3 py-2">
                                <div class="flex flex-col">
                                   <span class="text-[8px] font-black text-sigma-muted uppercase">{{ log.subject_name }}</span>
                                   <span class="text-[10px] font-bold tracking-tight">{{ log.title }}</span>
                                </div>
                             </td>
                             <td class="px-3 py-2 text-center">
                                <div class="inline-flex w-7 h-7 rounded-lg items-center justify-center text-[10px] font-black" :class="getGradeColor(log.grade)">
                                   {{ log.grade }}
                                </div>
                             </td>
                             <td class="px-3 py-2 text-center">
                                <div class="flex items-center justify-center gap-1">
                                  <button @click="deleteRecord(log.ID, r.student.ID)" class="w-6 h-6 rounded-md bg-rose-500/10 text-rose-500 inline-flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all">
                                    <Trash2 class="w-3 h-3" />
                                  </button>
                                </div>
                             </td>
                          </tr>
                        </template>
                        <tr v-else-if="loadingHistory[r.student.ID]">
                           <td colspan="4" class="p-8 text-center">
                              <div class="inline-block w-6 h-6 border-3 border-orange-500/20 border-t-orange-500 rounded-full animate-spin mb-2"></div>
                              <p class="text-[9px] font-black uppercase tracking-wider text-sigma-muted">Loading...</p>
                           </td>
                        </tr>
                        <tr v-else-if="!r.record && (!studentHistory[r.student.ID] || studentHistory[r.student.ID].length === 0)">
                            <td colspan="4" class="p-6 text-center text-sigma-muted opacity-60">
                               <p class="text-[9px] font-black uppercase tracking-wider italic">Belum ada riwayat</p>
                           </td>
                        </tr>
                    </tbody>
                 </table>
              </div>
           </div>
        </div>
      </div>
    </div>

    <!-- PERSONAL MODAL -->
    <Transition name="modal">
      <div v-if="showPersonalModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl" @click="showPersonalModal = false"></div>
        <div class="relative w-full max-w-xl bg-sigma-surface border border-sigma-border rounded-[3rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-300">
           <div class="p-8 border-b border-sigma-border flex items-center justify-between">
              <div class="flex items-center gap-4">
                 <div class="w-12 h-12 rounded-2xl bg-orange-500 flex items-center justify-center text-white shadow-lg">
                    <FileText class="w-6 h-6" />
                 </div>
                 <div>
                    <h3 class="text-xl font-black uppercase italic leading-none">{{ selectedStudent.name }}</h3>
                    <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Input Hafalan Pelajaran</p>
                 </div>
              </div>
              <button @click="showPersonalModal = false" class="text-sigma-muted hover:text-rose-500 transition-colors">
                 <XCircle class="w-6 h-6" />
              </button>
           </div>
           
           <div class="p-8 space-y-6">
              <div class="space-y-2">
                 <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Pelajaran</label>
                 <select v-model="personalForm.subject_name" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-orange-500">
                    <option value="">Pilih...</option>
                    <option v-for="sub in subjects" :key="sub.ID" :value="sub.name">{{ sub.name }}</option>
                 </select>
              </div>
              <div class="space-y-2">
                 <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Judul Hafalan</label>
                 <input v-model="personalForm.title" type="text" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-orange-500" placeholder="Contoh: Bab 1 Sejarah Islam" />
              </div>
              <div class="grid grid-cols-2 gap-6">
                 <div class="space-y-2">
                    <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Nilai / Grade</label>
                    <div class="flex items-center gap-1 bg-sigma-surface-alt p-1 rounded-2xl border border-sigma-border">
                       <button v-for="g in ['A', 'B', 'C', 'D']" :key="g" @click="personalForm.grade = g"
                               class="flex-1 py-3 rounded-xl text-[10px] font-black transition-all"
                               :class="personalForm.grade === g ? 'bg-orange-600 text-white shadow-lg' : 'text-sigma-muted hover:bg-sigma-border'">{{ g }}</button>
                    </div>
                 </div>
                 <div class="space-y-2">
                    <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal</label>
                    <input v-model="personalForm.date" type="date" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-[10px] font-black outline-none" />
                 </div>
              </div>
              <div class="space-y-2">
                 <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Catatan</label>
                 <textarea v-model="personalForm.remarks" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-orange-500" rows="2" placeholder="Keterangan tambahan..."></textarea>
              </div>
           </div>

           <div class="p-8 bg-sigma-surface-alt border-t border-sigma-border">
              <button @click="savePersonalMemorization" :disabled="isSaving"
                      class="w-full py-4 bg-orange-600 text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:bg-orange-700 transition-all flex items-center justify-center gap-3">
                 <Save v-if="!isSaving" class="w-4 h-4" />
                 <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                 Simpan Hafalan
              </button>
           </div>
        </div>
      </div>
    </Transition>

    <!-- BULK MODAL -->
    <Transition name="modal">
      <div v-if="showBulkModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl" @click="showBulkModal = false"></div>
        <div class="relative w-full max-w-6xl bg-sigma-surface border border-sigma-border rounded-[4rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-500">
           <div class="p-10 border-b border-sigma-border flex items-center justify-between bg-sigma-surface-alt/30">
              <div class="flex items-center gap-6">
                 <div class="w-14 h-14 rounded-[1.5rem] bg-orange-600 flex items-center justify-center text-white shadow-xl shadow-orange-600/30">
                    <Plus class="w-7 h-7" />
                 </div>
                 <div class="flex items-center gap-8">
                    <div>
                       <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Hafalan <span class="text-orange-600">Batch</span></h3>
                       <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">Input Kolektif Pelajaran</p>
                    </div>
                    <div class="h-10 w-px bg-sigma-border mx-2"></div>
                    <div class="flex items-center gap-6">
                       <div class="space-y-1">
                          <span class="text-[9px] font-black uppercase text-sigma-muted block ml-1">Pilih Kelas</span>
                          <select v-model="bulkForm.classroom_id" class="px-6 py-3 bg-white border border-sigma-border rounded-xl text-[10px] font-black outline-none focus:border-orange-500 shadow-sm">
                             <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
                          </select>
                       </div>
                       <div class="space-y-1">
                          <span class="text-[9px] font-black uppercase text-sigma-muted block ml-1">Pelajaran (Global)</span>
                          <select v-model="bulkForm.subject_name" class="px-6 py-3 bg-white border border-sigma-border rounded-xl text-[10px] font-black outline-none focus:border-orange-500 shadow-sm">
                             <option value="">Pilih...</option>
                             <option v-for="sub in subjects" :key="sub.ID" :value="sub.name">{{ sub.name }}</option>
                          </select>
                       </div>
                       <div class="space-y-1">
                          <span class="text-[9px] font-black uppercase text-sigma-muted block ml-1">Tanggal</span>
                          <input v-model="bulkForm.date" type="date" class="px-6 py-3 bg-white border border-sigma-border rounded-xl text-[10px] font-black outline-none focus:border-orange-500 shadow-sm" />
                       </div>
                    </div>
                 </div>
              </div>
              <button @click="showBulkModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
                 <XCircle class="w-6 h-6" />
              </button>
           </div>

           <div class="p-0 h-[50vh] overflow-y-auto custom-scrollbar bg-sigma-surface">
              <table v-if="bulkForm.classroom_id" class="w-full border-collapse">
                 <thead class="sticky top-0 z-10 bg-sigma-surface-alt/90 backdrop-blur-md border-b border-sigma-border">
                    <tr>
                       <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pl-12">Santri</th>
                       <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted">Mata Pelajaran</th>
                       <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted">Judul Hafalan</th>
                       <th class="p-6 text-center text-[9px] font-black uppercase tracking-widest text-sigma-muted w-32">Grade</th>
                    </tr>
                 </thead>
                 <tbody class="divide-y divide-sigma-border/40">
                    <tr v-for="s in bulkStudents" :key="s.ID" class="group hover:bg-orange-500/[0.03] transition-all">
                       <td class="p-6 pl-12">
                          <span class="text-sm font-black uppercase tracking-tight">{{ s.name }}</span>
                       </td>
                       <td class="p-6">
                          <select v-model="bulkForm.entries[s.ID].subject" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
                             <option value="">Pilih...</option>
                             <option v-for="sub in subjects" :key="sub.ID" :value="sub.name">{{ sub.name }}</option>
                          </select>
                       </td>
                       <td class="p-6">
                          <input v-model="bulkForm.entries[s.ID].title" type="text" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none" placeholder="Judul..." />
                       </td>
                       <td class="p-6">
                          <div class="flex items-center justify-center gap-1 bg-sigma-surface-alt/50 p-1 rounded-xl border border-sigma-border/30">
                             <button v-for="g in ['A', 'B', 'C', 'D']" :key="g" @click="bulkForm.entries[s.ID].grade = g"
                                     class="w-8 h-8 rounded-lg text-[10px] font-black transition-all"
                                     :class="bulkForm.entries[s.ID].grade === g ? 'bg-orange-600 text-white' : 'text-sigma-muted hover:bg-sigma-border'">{{ g }}</button>
                          </div>
                       </td>
                    </tr>
                 </tbody>
              </table>
           </div>

           <div class="p-10 bg-sigma-surface-alt border-t border-sigma-border flex items-center justify-between">
              <p class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Pastikan semua data sudah benar sebelum menyimpan.</p>
              <button @click="saveBulkMemorization" :disabled="isSaving"
                      class="px-14 py-5 bg-orange-600 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[10px] hover:shadow-2xl hover:shadow-orange-500/40 transition-all flex items-center gap-4 active:scale-95">
                 <Save v-if="!isSaving" class="w-5 h-5" />
                 <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                 Simpan Batch Hafalan
              </button>
           </div>
        </div>
      </div>
    </Transition>

    <!-- TOAST -->
    <Transition name="toast">
      <div v-if="showToast" class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[1000] flex items-center gap-4 px-8 py-4 bg-sigma-text text-white rounded-full shadow-2xl animate-in slide-in-from-bottom duration-300">
        <CheckCircle2 v-if="toastType === 'success'" class="w-5 h-5 text-emerald-400" />
        <XCircle v-else class="w-5 h-5 text-rose-400" />
        <span class="text-xs font-black uppercase tracking-widest">{{ toastMessage }}</span>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(99, 102, 241, 0.1); border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(99, 102, 241, 0.3); }

.modal-enter-active, .modal-leave-active { transition: all 0.3s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; transform: scale(0.95); }

.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 40px); }
</style>
