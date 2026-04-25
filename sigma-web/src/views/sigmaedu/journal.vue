<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  CheckCircle2, XCircle, 
  BookOpen, Save, History, Search, Trash2, RotateCw, Clock, Plus,
  ChevronRight, ChevronLeft, Layout, ClipboardList, AlertCircle, Sparkles,
  Eye, Pencil
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const subjects = ref<any[]>([])
const studyHours = ref<any[]>([])
const students = ref<any[]>([])
const journals = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const activeTab = ref('input') // 'input' | 'history'
const currentStep = ref(1)
const searchQuery = ref('')

const filters = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  date: new Date().toISOString().split('T')[0]
})

const form = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  study_hour_id: null as number | null,
  date: new Date().toISOString().split('T')[0],
  competence: '',
  main_material: '',
  learning_goal: '',
  activities: '',
  media: '',
  student_notes: '',
  obstacles: '',
  absence_records: [] as any[], // { student_id, name, status: 'S'|'I'|'A' }
  status: 'Hadir',
  remarks: ''
})

// Auto-save logic
const saveToLocal = () => {
  localStorage.setItem('sigma_journal_draft', JSON.stringify(form.value))
}

const loadFromLocal = () => {
  const draft = localStorage.getItem('sigma_journal_draft')
  if (draft) {
    try {
      const parsed = JSON.parse(draft)
      // Only restore if date is today (optional, but safer)
      form.value = { ...form.value, ...parsed }
    } catch (e) {
      console.error('Failed to parse draft')
    }
  }
}

// Watch for changes to auto-save
watch(form, saveToLocal, { deep: true })

const fetchData = async () => {
  try {
    const [clRes, sbRes, hrRes] = await Promise.all([
      axios.get('/api/v1/admin/edu/classrooms'),
      axios.get('/api/v1/admin/edu/subjects'),
      axios.get('/api/v1/admin/edu/hours')
    ])
    classrooms.value = clRes.data
    subjects.value = sbRes.data
    studyHours.value = hrRes.data

    if (classrooms.value.length > 0 && !form.value.classroom_id) {
      form.value.classroom_id = classrooms.value[0].ID
    }
    if (subjects.value.length > 0 && !form.value.subject_id) {
      form.value.subject_id = subjects.value[0].ID
    }
    if (studyHours.value.length > 0 && !form.value.study_hour_id) {
      form.value.study_hour_id = studyHours.value[0].ID
    }
    
    loadFromLocal()
  } catch (err) {
    console.error('Gagal mengambil data master:', err)
  }
}

const loadStudents = async () => {
  if (!form.value.classroom_id) return
  try {
    const res = await axios.get(`/api/v1/admin/edu/classrooms/${form.value.classroom_id}/students`)
    students.value = res.data
    // Initialize absence records if empty
    if (form.value.absence_records.length === 0) {
        // We don't pre-fill everyone, only those with specific status
    }
  } catch (err) {
    console.error('Gagal memuat siswa:', err)
  }
}

const toggleAbsence = (student: any, status: 'H' | 'S' | 'I' | 'A' | 'P') => {
    const idx = form.value.absence_records.findIndex(r => r.student_id === student.ID)
    if (idx > -1) {
        if (form.value.absence_records[idx].status === status) {
            form.value.absence_records.splice(idx, 1) // Remove if same
        } else {
            form.value.absence_records[idx].status = status // Change if different
        }
    } else {
        form.value.absence_records.push({
            student_id: student.ID,
            name: student.name,
            status: status
        })
    }
}

const getAbsenceStatus = (studentID: number) => {
    return form.value.absence_records.find(r => r.student_id === studentID)?.status || null
}

const loadJournals = async () => {
  isLoading.value = true
  try {
    const params: any = {
      date: filters.value.date,
      classroom_id: filters.value.classroom_id,
      subject_id: filters.value.subject_id
    }
    const res = await axios.get('/api/v1/admin/edu/journals', { params })
    journals.value = res.data
  } catch (err) {
    console.error('Gagal memuat jurnal:', err)
  } finally {
    isLoading.value = false
  }
}

const copyLastWeek = async () => {
    if (!form.value.subject_id) return
    try {
        const res = await axios.get('/api/v1/admin/edu/journals/last', {
            params: { subject_id: form.value.subject_id }
        })
        const last = res.data
        if (last) {
            form.value.competence = last.competence
            form.value.main_material = last.main_material
            form.value.learning_goal = last.learning_goal
            form.value.media = last.media
            
            toastMessage.value = 'Data pekan lalu berhasil disalin!'
            toastType.value = 'success'
            showToast.value = true
            setTimeout(() => showToast.value = false, 3000)
        }
    } catch (err) {
        toastMessage.value = 'Belum ada data pekan lalu.'
        toastType.value = 'info'
        showToast.value = true
        setTimeout(() => showToast.value = false, 3000)
    }
}

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

const editingId = ref<number | null>(null)
const viewingJournal = ref<any>(null)
const showViewModal = ref(false)

const editJournal = (journal: any) => {
  editingId.value = journal.ID
  form.value = {
    classroom_id: journal.classroom_id,
    subject_id: journal.subject_id,
    study_hour_id: journal.study_hour_id,
    date: journal.date.split('T')[0],
    competence: journal.competence || '',
    main_material: journal.main_material || '',
    learning_goal: journal.learning_goal || '',
    activities: journal.activities || '',
    media: journal.media || '',
    student_notes: journal.student_notes || '',
    obstacles: journal.obstacles || '',
    absence_records: journal.absence_records ? JSON.parse(journal.absence_records) : [],
    status: journal.status || 'Hadir',
    remarks: journal.remarks || ''
  }
  activeTab.value = 'input'
  currentStep.value = 1
}

const openViewModal = (journal: any) => {
    viewingJournal.value = journal
    showViewModal.value = true
}

const handleSave = async () => {
  if (!form.value.classroom_id || !form.value.subject_id || !form.value.study_hour_id) {
    toastMessage.value = 'Mohon lengkapi identitas pelajaran'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    return
  }
  
  isSaving.value = true
  try {
    const payload = {
      ...form.value,
      absence_records: JSON.stringify(form.value.absence_records),
      date: new Date(form.value.date).toISOString()
    }

    if (editingId.value) {
        await axios.put(`/api/v1/admin/edu/journals/${editingId.value}`, payload)
    } else {
        await axios.post('/api/v1/admin/edu/journals', payload)
    }
    
    toastMessage.value = editingId.value ? 'Jurnal Berhasil Diperbarui!' : 'Jurnal Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)

    // Clear local storage on success
    localStorage.removeItem('sigma_journal_draft')
    
    // Reset form
    currentStep.value = 1
    editingId.value = null
    form.value = {
        ...form.value,
        competence: '',
        main_material: '',
        learning_goal: '',
        activities: '',
        media: '',
        student_notes: '',
        obstacles: '',
        absence_records: [],
        remarks: ''
    }
    
    if (activeTab.value === 'history') await loadJournals()
  } catch (err: any) {
    toastMessage.value = err.response?.data?.error || 'Gagal Menyimpan Jurnal.'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isSaving.value = false
  }
}

const deleteJournal = async (id: any) => {
  if (!confirm('Hapus jurnal mengajar ini?')) return
  
  try {
    await axios.delete(`/api/v1/admin/edu/journals/${id}`)
    toastMessage.value = 'Jurnal Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    await loadJournals()
  } catch (err) {
    console.error('Gagal menghapus jurnal:', err)
  }
}

const filteredJournals = computed(() => {
  if (!searchQuery.value) return journals.value
  const query = searchQuery.value.toLowerCase()
  return journals.value.filter(j => 
    j.main_material?.toLowerCase().includes(query) ||
    j.subject?.name?.toLowerCase().includes(query) ||
    j.teacher?.name?.toLowerCase().includes(query)
  )
})

const getParsedAbsence = (records: any) => {
    if (!records) return []
    if (typeof records === 'string') {
        try {
            return JSON.parse(records)
        } catch (e) {
            console.error('Failed to parse absence records:', e)
            return []
        }
    }
    return Array.isArray(records) ? records : []
}

watch(() => form.value.classroom_id, loadStudents)
watch(activeTab, (val) => {
  if (val === 'history') loadJournals()
})

onMounted(fetchData)

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
  <div class="p-8 space-y-8 min-h-screen bg-sigma-bg">
    <!-- Header Section -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
      <div class="flex items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-gradient-to-br from-emerald-500 to-emerald-700 flex items-center justify-center text-white shadow-2xl shadow-emerald-500/40 transform -rotate-6 hover:rotate-0 transition-transform duration-500">
          <BookOpen class="w-8 h-8" />
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-sigma-muted mb-1 opacity-60">Akademik • Jurnal</h2>
          <h1 class="text-3xl font-black italic uppercase tracking-tighter leading-none">
            Teaching <span class="text-emerald-500">Journal</span>
          </h1>
        </div>
      </div>

      <div class="flex items-center gap-4 bg-sigma-surface/50 backdrop-blur-xl p-2 rounded-[2rem] border border-sigma-border shadow-xl">
        <button @click="activeTab = 'input'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="activeTab === 'input' ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <Plus class="w-3.5 h-3.5" /> Input Jurnal
        </button>
        <button @click="activeTab = 'history'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="activeTab === 'history' ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <History class="w-3.5 h-3.5" /> Riwayat
        </button>
      </div>
    </div>

    <!-- MAIN VIEW: INPUT (STEPPER) -->
    <div v-if="activeTab === 'input'" class="animate-in fade-in zoom-in duration-500 max-w-5xl mx-auto w-full">
      
      <!-- Stepper Progress -->
      <div class="flex items-center justify-between mb-12 relative px-4">
        <div class="absolute top-1/2 left-0 w-full h-1 bg-sigma-border -translate-y-1/2 z-0"></div>
        <div class="absolute top-1/2 left-0 h-1 bg-emerald-500 -translate-y-1/2 z-0 transition-all duration-500"
             :style="{ width: ((currentStep - 1) / 2) * 100 + '%' }"></div>
        
        <div v-for="s in 3" :key="s" class="relative z-10 flex flex-col items-center">
            <div class="w-12 h-12 rounded-full border-4 flex items-center justify-center font-black transition-all duration-500"
                 :class="currentStep >= s ? 'bg-emerald-500 border-sigma-bg text-white shadow-lg scale-110' : 'bg-sigma-bg border-sigma-border text-sigma-muted'">
                {{ s }}
            </div>
            <span class="absolute -bottom-8 whitespace-nowrap text-[9px] font-black uppercase tracking-widest text-sigma-muted"
                  :class="{ 'text-emerald-500': currentStep >= s }">
                {{ s === 1 ? 'Identitas' : s === 2 ? 'Materi & Proses' : 'Evaluasi' }}
            </span>
        </div>
      </div>

      <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] shadow-2xl overflow-hidden">
          <!-- Step 1: Identitas -->
          <div v-if="currentStep === 1" class="p-10 space-y-10 animate-in slide-in-from-right duration-500">
             <div class="flex items-center justify-between">
                <div class="flex items-center gap-4">
                   <div class="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                      <Layout class="w-5 h-5" />
                   </div>
                   <h3 class="text-lg font-black uppercase italic">Identitas Pelajaran</h3>
                </div>
                <button @click="copyLastWeek" class="flex items-center gap-2 px-5 py-2.5 bg-emerald-500/10 text-emerald-600 rounded-xl text-[9px] font-black uppercase tracking-widest hover:bg-emerald-500 hover:text-white transition-all shadow-sm">
                   <Sparkles class="w-3.5 h-3.5" /> Copy Last Week
                </button>
             </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div class="space-y-3">
                   <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Hari / Tanggal</label>
                   <input v-model="form.date" type="date" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50" />
                </div>
                <div class="space-y-3">
                   <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Jam Pelajaran</label>
                   <select v-model="form.study_hour_id" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50 appearance-none">
                      <option v-for="h in studyHours" :key="h.ID" :value="h.ID">{{ h.name }} ({{ h.start_time }} - {{ h.end_time }})</option>
                   </select>
                </div>
                <div class="space-y-3">
                   <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas & Rombel</label>
                   <select v-model="form.classroom_id" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50 appearance-none">
                      <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
                   </select>
                </div>
                <div class="space-y-3">
                   <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mata Pelajaran</label>
                   <select v-model="form.subject_id" class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50 appearance-none">
                      <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
                   </select>
                </div>
             </div>
          </div>

          <!-- Step 2: Materi & Proses -->
          <div v-if="currentStep === 2" class="p-10 space-y-10 animate-in slide-in-from-right duration-500">
             <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                   <ClipboardList class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-black uppercase italic">Inti Materi & Proses</h3>
             </div>

             <div class="space-y-8">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                   <div class="space-y-3">
                      <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kompetensi Dasar (KD)</label>
                      <input v-model="form.competence" type="text" placeholder="Misal: KD 3.1 Memahami..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50" />
                   </div>
                   <div class="space-y-3">
                      <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Materi Pokok</label>
                      <input v-model="form.main_material" type="text" placeholder="Judul bab atau sub-bab..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50" />
                   </div>
                </div>

                <div class="space-y-3">
                   <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tujuan Belajar</label>
                   <textarea v-model="form.learning_goal" rows="3" placeholder="Apa yang ingin dicapai pada pertemuan ini..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-[1.5rem] text-xs font-black outline-none focus:border-emerald-500/50 resize-none"></textarea>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div class="space-y-3">
                       <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Ringkasan Kegiatan</label>
                       <textarea v-model="form.activities" rows="4" placeholder="Poin-poin aktivitas selama di kelas..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-[1.5rem] text-xs font-black outline-none focus:border-emerald-500/50 resize-none"></textarea>
                    </div>
                    <div class="space-y-3">
                       <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Media / Alat Peraga</label>
                       <textarea v-model="form.media" rows="4" placeholder="Papan tulis, LCD, Alat Peraga, dll..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-[1.5rem] text-xs font-black outline-none focus:border-emerald-500/50 resize-none"></textarea>
                    </div>
                </div>
             </div>
          </div>

          <!-- Step 3: Evaluasi & Absensi -->
          <div v-if="currentStep === 3" class="p-10 space-y-10 animate-in slide-in-from-right duration-500">
             <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                   <AlertCircle class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-black uppercase italic">Evaluasi & Rekap Absensi</h3>
             </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                 <div class="space-y-8">
                    <div class="space-y-3">
                       <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Catatan Siswa</label>
                       <textarea v-model="form.student_notes" rows="3" placeholder="Catatan khusus siswa aktif atau bermasalah..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-[1.5rem] text-xs font-black outline-none focus:border-emerald-500/50 resize-none"></textarea>
                    </div>
                    <div class="space-y-3">
                       <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Hambatan</label>
                       <textarea v-model="form.obstacles" rows="3" placeholder="Gangguan (Listrik mati, kelas bising, dll)..." class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-[1.5rem] text-xs font-black outline-none focus:border-emerald-500/50 resize-none"></textarea>
                    </div>
                 </div>

                 <div class="space-y-4">
                    <label class="text-[10px] font-black uppercase tracking-widest text-sigma-muted ml-1">Rekap Ketidakhadiran Siswa</label>
                    <div class="bg-sigma-surface-alt border border-sigma-border rounded-[2rem] overflow-hidden max-h-[300px] overflow-y-auto custom-scrollbar">
                        <div v-for="s in students" :key="s.ID" class="p-4 flex items-center justify-between border-b border-sigma-border/50 last:border-0 hover:bg-sigma-surface transition-colors">
                            <div>
                                <p class="text-[10px] font-black uppercase tracking-tight">{{ s.name }}</p>
                                <p class="text-[8px] text-sigma-muted font-bold">{{ s.nis }}</p>
                            </div>
                            <div class="flex gap-1">
                                <button v-for="st in ['H', 'S', 'I', 'A', 'P']" :key="st"
                                        @click="toggleAbsence(s, st as any)"
                                        class="w-7 h-7 rounded-lg text-[9px] font-black flex items-center justify-center transition-all"
                                        :class="getAbsenceStatus(s.ID) === st ? 
                                            (st === 'H' ? 'bg-emerald-500 text-white' : 
                                             st === 'S' ? 'bg-amber-500 text-white' : 
                                             st === 'I' ? 'bg-blue-500 text-white' : 
                                             st === 'A' ? 'bg-rose-500 text-white' : 
                                             'bg-purple-500 text-white') : 
                                            'bg-sigma-surface border border-sigma-border text-sigma-muted hover:border-emerald-500/30'">
                                    {{ st }}
                                </button>
                            </div>
                        </div>
                    </div>
                    <p class="text-[8px] text-sigma-muted italic text-center">*Klik H (Hadir), S (Sakit), I (Izin), A (Alpha), P (Piket). Klik ulang untuk batal.</p>
                 </div>
             </div>
          </div>

          <!-- Navigation Footer -->
          <div class="p-10 bg-sigma-surface-alt/30 border-t border-sigma-border flex items-center justify-between">
              <button @click="currentStep--" :disabled="currentStep === 1"
                      class="flex items-center gap-2 px-8 py-4 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all disabled:opacity-0"
                      :class="'text-sigma-muted hover:text-sigma-text'">
                <ChevronLeft class="w-4 h-4" /> Back
              </button>

              <button v-if="currentStep < 3" @click="currentStep++"
                      class="flex items-center gap-2 px-10 py-4 bg-emerald-500 text-white rounded-2xl text-[10px] font-black uppercase tracking-widest hover:bg-emerald-600 transition-all shadow-xl shadow-emerald-500/20 active:scale-95">
                Next <ChevronRight class="w-4 h-4" />
              </button>

              <button v-else @click="handleSave" :disabled="isSaving"
                      class="flex items-center gap-3 px-12 py-4 bg-emerald-600 text-white rounded-2xl text-[10px] font-black uppercase tracking-widest hover:bg-emerald-700 transition-all shadow-xl shadow-emerald-500/30 active:scale-95 disabled:opacity-50">
                <div v-if="isSaving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                <Save v-else class="w-4 h-4" />
                Finish & Save Journal
              </button>
          </div>
      </div>
    </div>

    <!-- MAIN VIEW: HISTORY -->
    <div v-else class="space-y-8 animate-in fade-in duration-500">
      <!-- Filter Bar -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-8 shadow-sm flex flex-wrap items-end gap-8">
        <div class="flex-1 min-w-[200px] space-y-3">
           <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">Cari Materi / Guru</label>
           <div class="relative">
              <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-emerald-500" />
              <input v-model="searchQuery" type="text" placeholder="Ketik..."
                     class="w-full pl-12 pr-4 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50" />
           </div>
        </div>
        <div class="w-[180px] space-y-3">
           <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal</label>
           <input v-model="filters.date" type="date" 
                  class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-emerald-500/50" />
        </div>
        <button @click="loadJournals" 
                class="p-4 bg-emerald-500 text-white rounded-2xl hover:bg-emerald-600 transition-all shadow-xl shadow-emerald-500/20 active:scale-90">
          <RotateCw class="w-5 h-5" :class="{ 'animate-spin': isLoading }" />
        </button>
      </div>

      <!-- Journal List -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center py-40 gap-4">
         <div class="w-12 h-12 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin"></div>
         <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Synchronizing Data...</p>
      </div>

      <div v-else-if="filteredJournals.length > 0" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div v-for="j in filteredJournals" :key="j.ID" 
             class="bg-sigma-surface border border-sigma-border rounded-[3rem] p-8 hover:shadow-2xl hover:border-emerald-500/20 transition-all group relative overflow-hidden">
           
           <div class="flex justify-between items-start mb-8">
             <div class="flex items-center gap-5">
               <div class="w-14 h-14 rounded-2xl bg-emerald-500/10 flex items-center justify-center text-emerald-500 shadow-sm group-hover:scale-110 transition-transform">
                 <BookOpen class="w-7 h-7" />
               </div>
               <div>
                 <h3 class="text-base font-black uppercase italic tracking-tight text-sigma-text">{{ j.subject?.name }}</h3>
                 <div class="flex items-center gap-3 text-[9px] font-black uppercase tracking-widest text-sigma-muted mt-1">
                    <span class="px-2 py-0.5 bg-sigma-surface-alt rounded border border-sigma-border">{{ j.classroom?.name }}</span>
                    <span>•</span>
                    <span class="flex items-center gap-1"><Clock class="w-3 h-3" /> {{ j.study_hour?.name }}</span>
                 </div>
               </div>
             </div>
             <div class="flex items-center gap-2">
                <button @click="openViewModal(j)" class="p-3 rounded-xl bg-emerald-500/10 text-emerald-500 opacity-0 group-hover:opacity-100 transition-all hover:bg-emerald-500 hover:text-white">
                    <Eye class="w-4 h-4" />
                </button>
                <button @click="editJournal(j)" class="p-3 rounded-xl bg-blue-500/10 text-blue-500 opacity-0 group-hover:opacity-100 transition-all hover:bg-blue-500 hover:text-white">
                    <Pencil class="w-4 h-4" />
                </button>
                <button @click="deleteJournal(j.ID)" class="p-3 rounded-xl bg-rose-500/10 text-rose-500 opacity-0 group-hover:opacity-100 transition-all hover:bg-rose-500 hover:text-white">
                    <Trash2 class="w-4 h-4" />
                </button>
             </div>
           </div>

           <div class="space-y-6 mb-8">
              <div v-if="j.main_material" class="space-y-2">
                 <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Materi Pokok</p>
                 <p class="text-xs font-black text-sigma-text">{{ j.main_material }}</p>
              </div>
              <div v-if="j.learning_goal" class="space-y-2">
                 <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Tujuan & Ringkasan</p>
                 <p class="text-[11px] font-bold text-sigma-muted leading-relaxed line-clamp-3">{{ j.learning_goal }}</p>
              </div>
           </div>

           <div class="flex items-center justify-between pt-6 border-t border-sigma-border/50">
              <div class="flex items-center gap-4">
                 <div class="w-10 h-10 rounded-full bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-[11px] font-black">
                    {{ j.teacher?.name?.charAt(0) }}
                 </div>
                 <div>
                    <span class="text-[10px] font-black uppercase block text-sigma-text">{{ j.teacher?.name }}</span>
                    <span class="text-[8px] font-black text-sigma-muted uppercase">{{ formatDate(j.date) }}</span>
                 </div>
              </div>
              
              <div v-if="j.absence_records" class="flex -space-x-2">
                  <div v-for="(rec, idx) in JSON.parse(j.absence_records).slice(0, 3)" :key="idx"
                       class="w-8 h-8 rounded-full border-2 border-sigma-bg flex items-center justify-center text-[8px] font-black text-white shadow-sm"
                       :class="rec.status === 'S' ? 'bg-amber-500' : rec.status === 'I' ? 'bg-blue-500' : 'bg-rose-500'"
                       :title="rec.name">
                      {{ rec.status }}
                  </div>
                  <div v-if="JSON.parse(j.absence_records).length > 3" class="w-8 h-8 rounded-full border-2 border-sigma-bg bg-sigma-surface-alt flex items-center justify-center text-[8px] font-black text-sigma-muted">
                      +{{ JSON.parse(j.absence_records).length - 3 }}
                  </div>
              </div>
           </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-40 opacity-30">
        <BookOpen class="w-16 h-16 text-sigma-muted mb-4" />
        <p class="font-black uppercase tracking-widest text-[10px]">Belum ada jurnal untuk filter ini</p>
      </div>
    </div>

    <!-- Modal View Detail -->
    <Transition name="toast">
      <div v-if="showViewModal" class="fixed inset-0 z-[300] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/80 backdrop-blur-sm" @click="showViewModal = false"></div>
        <div v-if="viewingJournal" class="relative bg-sigma-surface border border-sigma-border w-full max-w-2xl rounded-[3rem] shadow-2xl overflow-hidden animate-in zoom-in duration-300">
            <div class="p-8 border-b border-sigma-border flex items-center justify-between bg-emerald-500/5">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-emerald-500 flex items-center justify-center text-white">
                        <BookOpen class="w-6 h-6" />
                    </div>
                    <div>
                        <h3 class="text-lg font-black uppercase italic">{{ viewingJournal.subject?.name }}</h3>
                        <p class="text-[9px] font-black uppercase tracking-widest text-sigma-muted">{{ formatDate(viewingJournal.date) }} • {{ viewingJournal.classroom?.name }}</p>
                    </div>
                </div>
                <button @click="showViewModal = false" class="w-10 h-10 rounded-full bg-sigma-surface border border-sigma-border flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all">
                    <XCircle class="w-5 h-5" />
                </button>
            </div>
            
            <div class="p-10 max-h-[70vh] overflow-y-auto custom-scrollbar space-y-10">
                <div class="grid grid-cols-2 gap-8">
                    <div class="space-y-2">
                        <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Kompetensi Dasar</p>
                        <p class="text-xs font-black text-sigma-text">{{ viewingJournal.competence || '-' }}</p>
                    </div>
                    <div class="space-y-2">
                        <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Materi Pokok</p>
                        <p class="text-xs font-black text-sigma-text">{{ viewingJournal.main_material || '-' }}</p>
                    </div>
                </div>

                <div class="space-y-2">
                    <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Tujuan Belajar</p>
                    <p class="text-[11px] font-bold text-sigma-muted leading-relaxed">{{ viewingJournal.learning_goal || '-' }}</p>
                </div>

                <div class="space-y-2">
                    <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Ringkasan Kegiatan</p>
                    <p class="text-[11px] font-bold text-sigma-muted leading-relaxed">{{ viewingJournal.activities || '-' }}</p>
                </div>

                <div class="space-y-2">
                    <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Media / Alat</p>
                    <p class="text-[11px] font-bold text-sigma-muted leading-relaxed">{{ viewingJournal.media || '-' }}</p>
                </div>

                <div class="pt-6 border-t border-sigma-border">
                    <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted mb-4">Absensi & Catatan</p>
                    <div class="flex flex-wrap gap-2 mb-6">
                        <div v-for="(rec, idx) in getParsedAbsence(viewingJournal.absence_records)" :key="idx"
                             class="flex items-center gap-2 px-3 py-1.5 rounded-xl border border-sigma-border bg-sigma-surface-alt">
                            <span class="w-5 h-5 rounded-lg flex items-center justify-center text-[8px] font-black text-white"
                                  :class="rec.status === 'H' ? 'bg-emerald-500' : rec.status === 'S' ? 'bg-amber-500' : rec.status === 'I' ? 'bg-blue-500' : rec.status === 'A' ? 'bg-rose-500' : 'bg-purple-500'">
                                {{ rec.status }}
                            </span>
                            <span class="text-[9px] font-black uppercase">{{ rec.name }}</span>
                        </div>
                        <p v-if="getParsedAbsence(viewingJournal.absence_records).length === 0" class="text-[10px] italic text-sigma-muted">Tidak ada data absensi khusus.</p>
                    </div>
                    <div class="grid grid-cols-2 gap-8">
                        <div class="space-y-2">
                            <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Catatan Siswa</p>
                            <p class="text-xs font-black text-sigma-text">{{ viewingJournal.student_notes || '-' }}</p>
                        </div>
                        <div class="space-y-2">
                            <p class="text-[8px] font-black uppercase tracking-widest text-sigma-muted">Hambatan</p>
                            <p class="text-xs font-black text-sigma-text">{{ viewingJournal.obstacles || '-' }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
      </div>
    </Transition>

    <!-- Toast Notification -->
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
  </div>
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
</style>
