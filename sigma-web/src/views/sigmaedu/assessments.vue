<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  GraduationCap, Save, 
  CheckCircle2, XCircle,
  Plus, RotateCw, FileDown,
  Trash2
} from 'lucide-vue-next'
import * as XLSX from 'xlsx'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const subjects = ref<any[]>([])
const students = ref<any[]>([])
const assessments = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const showBulkModal = ref(false)
const bulkForm = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  type: 'LISAN',
  date: new Date().toISOString().split('T')[0],
  scores: {} as Record<number, number | null>,
  remarks: {} as Record<number, string>
})

const filters = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  term: 'Ganjil',
  academic_year: '2024/2025'
})

const STANDARD_CATEGORIES = [
  { label: 'Tugas', value: 'TUGAS' },
  { label: 'PR', value: 'PR' },
  { label: 'MID', value: 'MID' },
  { label: 'Ujian Tulis', value: 'UJIAN_TULIS' },
  { label: 'Ujian Lisan', value: 'UJIAN_LISAN' },
]

const assessmentTypes = ref([...STANDARD_CATEGORIES])

const loadCategories = () => {
  const stored = localStorage.getItem('sigma_custom_categories')
  if (stored) {
    try {
      const custom = JSON.parse(stored)
      assessmentTypes.value = [...STANDARD_CATEGORIES, ...custom]
    } catch (e) {
      console.error('Gagal memuat kategori:', e)
    }
  }
}

const saveCategories = () => {
  const custom = assessmentTypes.value.filter(t => !STANDARD_CATEGORIES.find(s => s.value === t.value))
  if (custom.length > 0) {
    localStorage.setItem('sigma_custom_categories', JSON.stringify(custom))
  } else {
    localStorage.removeItem('sigma_custom_categories')
  }
}

const newCategory = ref('')
const addCategory = () => {
  if (!newCategory.value) return
  const val = newCategory.value.toUpperCase().replace(/\s+/g, '_')
  if (!assessmentTypes.value.find(t => t.value === val)) {
    assessmentTypes.value.push({ label: newCategory.value, value: val })
    saveCategories()
  }
  newCategory.value = ''
}

const removeCategory = (val: string) => {
  if (STANDARD_CATEGORIES.find(s => s.value === val)) return
  assessmentTypes.value = assessmentTypes.value.filter(t => t.value !== val)
  saveCategories()
}

// Load immediately
loadCategories()

const deleteStudentScores = async (student: any) => {
  if (!confirm(`Hapus semua nilai untuk ${student.name}?`)) return
  
  isSaving.value = true
  try {
    const payload = assessmentTypes.value.map(t => ({
      student_id: student.ID,
      classroom_id: filters.value.classroom_id,
      subject_id: filters.value.subject_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: t.value,
      score: null,
      remarks: '',
      grade: ''
    }))

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    await loadAssessmentData()
    
    toastMessage.value = 'Data Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal menghapus nilai:', err)
  } finally {
    isSaving.value = false
  }
}

const openBulkModal = (initialType?: string) => {
  bulkForm.value.classroom_id = filters.value.classroom_id
  bulkForm.value.subject_id = filters.value.subject_id
  
  if (initialType) {
    bulkForm.value.type = initialType
  }
  
  // Initialize scores with current matrix data if available
  const newScores: Record<number, number | null> = {}
  const newRemarks: Record<number, string> = {}
  
  students.value.forEach(s => {
    newScores[s.ID] = scoreMatrix.value[s.ID]?.[bulkForm.value.type] || null
    newRemarks[s.ID] = remarksMap.value[s.ID] || ''
  })
  
  bulkForm.value.scores = newScores
  bulkForm.value.remarks = newRemarks
  showBulkModal.value = true
}

// Watch for type change in bulk modal to reload scores
watch(() => bulkForm.value.type, (newType) => {
  students.value.forEach(s => {
    bulkForm.value.scores[s.ID] = scoreMatrix.value[s.ID]?.[newType] || null
    // Also try to find remarks for this type if possible, 
    // though the current backend might share remarks across types or store per type.
    // Given the current service logic, remarks are tied to (student, subject, type).
    // Let's check the remarksMap structure.
  })
})

const saveBulkScores = async () => {
  if (!bulkForm.value.classroom_id || !bulkForm.value.subject_id) return
  
  isSaving.value = true
  try {
    const payload = students.value.map(s => ({
      student_id: s.ID,
      classroom_id: bulkForm.value.classroom_id,
      subject_id: bulkForm.value.subject_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: bulkForm.value.type,
      score: bulkForm.value.scores[s.ID],
      remarks: bulkForm.value.remarks[s.ID],
      grade: calculateGrade(bulkForm.value.scores[s.ID] || null)
    })).filter(a => a.score !== null)

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    
    await loadAssessmentData()
    showBulkModal.value = false
    toastMessage.value = 'Nilai Kolektif Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan kolektif:', err)
  } finally {
    isSaving.value = false
  }
}

// UI State
// Map: [studentID][type] = score
const scoreMatrix = ref<Record<number, Record<string, number | null>>>({})
const remarksMap = ref<Record<number, string>>({})

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
    if (subjects.value.length > 0 && !filters.value.subject_id) {
      filters.value.subject_id = subjects.value[0].ID
    }
    await loadAssessmentData()
  } catch (err) {
    console.error('Gagal mengambil data master:', err)
  }
}

const loadAssessmentData = async () => {
  if (!filters.value.classroom_id || !filters.value.subject_id) return
  
  isLoading.value = true
  try {
    // 1. Get students in class
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data

    // 2. Get existing assessments
    const params = {
      classroom_id: filters.value.classroom_id,
      subject_id: filters.value.subject_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year
    }
    const asRes = await axios.get('/api/v1/admin/edu/assessments', { params })
    assessments.value = asRes.data

    // Auto-discover types from data
    const discoveredTypes = [...new Set(assessments.value.map(a => a.type))]
    discoveredTypes.forEach(t => {
      if (!assessmentTypes.value.find(at => at.value === t)) {
        assessmentTypes.value.push({ label: t, value: t })
      }
    })

    const newRemarks: Record<number, string> = {}
    
    students.value.forEach(s => {
      scoreMatrix.value[s.ID] = scoreMatrix.value[s.ID] || {}
      newRemarks[s.ID] = ''
      
      assessmentTypes.value.forEach(t => {
        const found = assessments.value.find(a => a.student_id === s.ID && a.type === t.value)
        scoreMatrix.value[s.ID][t.value] = found ? found.score : null
        if (found?.remarks) newRemarks[s.ID] = found.remarks
      })
    })
    
    remarksMap.value = newRemarks

  } catch (err) {
    console.error('Gagal memuat data penilaian:', err)
  } finally {
    isLoading.value = false
  }
}

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

const calculateGrade = (score: number | null) => {
  if (score === null || score === undefined || isNaN(score)) return ''
  if (score >= 95) return 'A+'
  if (score >= 90) return 'A'
  if (score >= 85) return 'B+'
  if (score >= 80) return 'B'
  if (score >= 75) return 'C+'
  if (score >= 70) return 'C'
  if (score >= 60) return 'D'
  return 'E'
}

const exportToExcel = () => {
  const data = rekapData.value.map(r => {
    const row: any = {
      NIS: r.student.nis,
      Nama: r.student.name,
    }
    assessmentTypes.value.forEach(t => {
      row[t.label] = r.scores[t.value] || '-'
    })
    row['Rata-rata'] = r.average
    row['Grade'] = r.finalGrade
    return row
  })
  
  const ws = XLSX.utils.json_to_sheet(data)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, "Nilai Pelajaran")
  XLSX.writeFile(wb, `Nilai_Pelajaran_${filters.value.classroom_id}.xlsx`)
}

const getGradeColor = (grade: string) => {
  if (!grade) return 'text-sigma-muted'
  if (grade.startsWith('A')) return 'text-emerald-500'
  if (grade.startsWith('B')) return 'text-blue-500'
  if (grade.startsWith('C')) return 'text-amber-500'
  if (grade.startsWith('D')) return 'text-orange-500'
  if (grade.startsWith('E')) return 'text-rose-500'
  return 'text-sigma-muted'
}

// Watchers
watch([
  () => filters.value.classroom_id, 
  () => filters.value.subject_id, 
  () => filters.value.term
], loadAssessmentData)

onMounted(fetchData)

const rekapData = computed(() => {
  return students.value.map(s => {
    const row: any = {
      student: s,
      scores: scoreMatrix.value[s.ID] || {}
    }
    
    // Calculate Final (Simple Average)
    const scores = Object.values(row.scores).filter(v => v !== null) as number[]
    row.average = scores.length > 0 ? (scores.reduce((a: number, b: number) => a + b, 0) / scores.length).toFixed(1) : '-'
    row.finalGrade = row.average !== '-' ? calculateGrade(parseFloat(row.average)) : '-'
    
    return row
  })
})
</script>

<template>
  <div class="p-8 space-y-8">
    <div class="assessment-page space-y-8">
    <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="h-10 w-1 bg-blue-500 rounded-full shadow-[0_0_15px_rgba(59,130,246,0.5)]"></div>
          <div>
            <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Penilaian</h2>
            <h1 class="text-2xl font-black italic uppercase tracking-tight">Rekap <span class="text-blue-500">Nilai Pelajaran</span></h1>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button @click="() => openBulkModal()" :disabled="!filters.classroom_id"
                  class="px-8 py-3 bg-blue-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-lg shadow-blue-500/20 active:scale-95 disabled:opacity-50 flex items-center gap-3">
            <Plus class="w-4 h-4" />
            Input Nilai Kolektif
          </button>
        </div>
      </div>

    <div class="flex flex-col gap-6">
      <!-- Horizontal Filter Bar -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center justify-between gap-8">
        <div class="flex items-center gap-6">
          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tahun Ajaran</label>
            <select v-model="filters.academic_year" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
              <option class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">2024/2025</option><option class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">2023/2024</option>
            </select>
          </div>
          <div class="space-y-1 min-w-[120px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Semester</label>
            <select v-model="filters.term" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
              <option class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">Ganjil</option><option class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">Genap</option>
            </select>
          </div>
          <div class="space-y-1 min-w-[150px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
            <select v-model="filters.classroom_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
              <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID" class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">{{ cls.name }}</option>
            </select>
          </div>
          <div class="space-y-1 min-w-[180px]">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mata Pelajaran</label>
            <select v-model="filters.subject_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
              <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID" class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">{{ sub.name }}</option>
            </select>
          </div>

          <!-- Quick Add Session -->
          <div class="space-y-1">
            <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tambah Sesi (Tugas 2, dll)</label>
            <div class="flex gap-2">
              <input v-model="newCategory" type="text" placeholder="Tugas 2..."
                     @keyup.enter="addCategory"
                     class="w-[120px] px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-blue-500/50" />
              <button @click="addCategory" class="p-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-all shadow-lg shadow-blue-500/20 active:scale-95">
                <Plus class="w-4 h-4" />
              </button>
            </div>
          </div>
          
          <!-- Custom Tags List (Compact) -->
          <div class="flex flex-wrap gap-1 max-w-[200px]">
             <div v-for="t in assessmentTypes.filter(cat => !STANDARD_CATEGORIES.find(s => s.value === cat.value))" 
                   :key="t.value"
                   class="px-2 py-1 bg-sigma-surface-alt border border-sigma-border rounded-lg text-[8px] font-black flex items-center gap-2 group">
                {{ t.label }}
                <button @click="removeCategory(t.value)" class="text-rose-500 hover:scale-110 transition-transform">
                  <XCircle class="w-3 h-3" />
                </button>
              </div>
          </div>
        </div>
        

        <div class="flex items-center gap-3">
          <button @click="exportToExcel" 
                  class="p-4 bg-sigma-surface-alt border border-sigma-border rounded-xl text-emerald-500 hover:bg-emerald-500 hover:text-white transition-all shadow-sm"
                  title="Export Excel">
            <FileDown class="w-5 h-5" />
          </button>
          <button @click="loadAssessmentData" 
                  class="p-4 bg-sigma-surface-alt border border-sigma-border rounded-xl text-blue-500 hover:bg-blue-500 hover:text-white transition-all shadow-sm"
                  title="Refresh Data">
            <RotateCw class="w-5 h-5" :class="{ 'animate-spin': isLoading }" />
          </button>
        </div>
      </div>

      <!-- Main Content -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] shadow-sm min-h-[70vh] overflow-hidden flex flex-col">
        
        <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-40 space-y-4">
          <div class="w-10 h-10 border-4 border-blue-500/20 border-t-blue-500 rounded-full animate-spin"></div>
          <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Memuat Data Penilaian...</p>
        </div>

        <!-- Summary Mode -->
        <div class="flex-1 overflow-x-auto">
          <table v-if="students.length > 0" class="w-full border-collapse">
            <thead>
              <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border sticky top-0 z-10 backdrop-blur-md">
                <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted min-w-[200px]">Santri</th>
                <th class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Jumlah Nilai</th>
                <th class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted bg-blue-500/5">Rata-rata</th>
                <th class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted bg-blue-500/5">Grade</th>
                <th class="p-6 text-right text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in students" :key="s.ID" class="border-b border-sigma-border last:border-0 group hover:bg-sigma-surface-alt/20 transition-all">
                <td class="p-6">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-[10px] font-black group-hover:scale-110 transition-transform group-hover:border-blue-500/50">{{ s.name.charAt(0) }}</div>
                    <div>
                      <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ s.nis }}</span>
                      <span class="text-sm font-black uppercase tracking-tighter italic">{{ s.name }}</span>
                    </div>
                  </div>
                </td>
                <td class="p-6">
                  <div class="flex flex-wrap gap-1 justify-center">
                    <template v-for="t in assessmentTypes" :key="t.value">
                      <button v-if="scoreMatrix[s.ID]?.[t.value] !== null && scoreMatrix[s.ID]?.[t.value] !== undefined"
                              @click="openBulkModal(t.value)"
                              class="px-2 py-1 bg-blue-500/10 border border-blue-500/20 rounded-md text-[8px] font-black text-blue-600 hover:bg-blue-600 hover:text-white transition-all">
                        {{ t.label }}: {{ scoreMatrix[s.ID][t.value] }}
                      </button>
                    </template>
                    <span v-if="Object.values(scoreMatrix[s.ID] || {}).filter(v => v !== null).length === 0" 
                          class="text-[10px] font-black text-sigma-muted opacity-30 italic">Belum Ada Nilai</span>
                  </div>
                </td>
                <td class="p-6 text-center font-black text-blue-500 bg-blue-500/5">
                  {{ rekapData.find(r => r.student.ID === s.ID)?.average }}
                </td>
                <td class="p-6 text-center font-black italic bg-blue-500/5" :class="getGradeColor(rekapData.find(r => r.student.ID === s.ID)?.finalGrade || '')">
                  {{ rekapData.find(r => r.student.ID === s.ID)?.finalGrade }}
                </td>
                <td class="p-6 text-right">
                  <div class="flex items-center justify-end gap-2">
                    <button @click="deleteStudentScores(s)"
                            class="p-3 bg-rose-500/10 text-rose-600 border border-rose-400/20 rounded-xl hover:bg-rose-600 hover:text-white transition-all shadow-sm"
                            title="Hapus Semua Nilai">
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-else class="flex-1 flex flex-col items-center justify-center p-40 space-y-4 opacity-30">
            <GraduationCap class="w-16 h-16" />
            <p class="font-black uppercase tracking-widest text-[10px]">Pilih kelas dan mapel</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Bulk Input Modal -->
    <Transition name="modal">
      <div v-if="showBulkModal" class="fixed inset-0 z-[300] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/80 backdrop-blur-xl" @click="showBulkModal = false"></div>
        <div class="relative w-full max-w-5xl bg-sigma-surface border border-sigma-border rounded-[3rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-300">
          <!-- Modal Header -->
          <div class="p-8 border-b border-sigma-border flex items-center justify-between bg-sigma-surface-alt/50">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-blue-500 flex items-center justify-center text-white">
                <Plus class="w-6 h-6" />
              </div>
              <div>
                <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted block">Input Nilai Kolektif</span>
                <h3 class="text-xl font-black uppercase tracking-tighter italic">BANYAK SANTRI SEKALIGUS</h3>
              </div>
            </div>
            <button @click="showBulkModal = false" class="p-3 bg-sigma-surface border border-sigma-border rounded-xl hover:bg-rose-500 hover:text-white transition-all">
              <XCircle class="w-5 h-5" />
            </button>
          </div>

          <!-- Modal Configuration Bar -->
          <div class="p-6 bg-sigma-surface-alt border-b border-sigma-border flex items-center gap-6">
            <div class="flex-1 space-y-1">
              <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mata Pelajaran</label>
              <select v-model="bulkForm.subject_id" class="w-full px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none">
                <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID" class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">{{ sub.name }}</option>
              </select>
            </div>
            <div class="flex-1 space-y-1">
              <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Jenis Penilaian</label>
              <select v-model="bulkForm.type" class="w-full px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none">
                <option v-for="t in assessmentTypes" :key="t.value" :value="t.value" class="bg-white text-slate-900 dark:bg-slate-900 dark:text-slate-100">{{ t.label }}</option>
              </select>
            </div>
            <div class="flex-1 space-y-1">
              <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal</label>
              <input v-model="bulkForm.date" type="date" class="w-full px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none" />
            </div>
          </div>

          <!-- Modal Body: Student List -->
          <div class="p-0 max-h-[50vh] overflow-y-auto">
            <table class="w-full border-collapse">
              <thead class="sticky top-0 z-10 bg-sigma-surface shadow-sm">
                <tr>
                  <th class="p-4 text-left text-[8px] font-black uppercase tracking-widest text-sigma-muted pl-8">Santri</th>
                  <th class="p-4 text-center text-[8px] font-black uppercase tracking-widest text-sigma-muted w-32">Skor (0-100)</th>
                  <th class="p-4 text-left text-[8px] font-black uppercase tracking-widest text-sigma-muted pr-8">Catatan Singkat</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="s in students" :key="s.ID" class="border-b border-sigma-border last:border-0 hover:bg-sigma-surface-alt/30 transition-all">
                  <td class="p-4 pl-8">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 rounded-xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center text-[10px] font-black">{{ s.name.charAt(0) }}</div>
                      <span class="text-xs font-black uppercase tracking-tighter">{{ s.name }}</span>
                    </div>
                  </td>
                  <td class="p-4">
                    <input v-model.number="bulkForm.scores[s.ID]" type="number" 
                           class="w-full px-3 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-center font-black text-blue-600 outline-none focus:border-blue-500/50"
                           placeholder="0" />
                  </td>
                  <td class="p-4 pr-8">
                    <input v-model="bulkForm.remarks[s.ID]" type="text"
                           class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-[10px] font-bold outline-none focus:border-blue-500/50"
                           placeholder="Evaluasi..." />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Modal Footer -->
          <div class="p-8 border-t border-sigma-border flex items-center justify-between bg-sigma-surface-alt/30">
            <div class="flex items-center gap-6">
              <div class="text-center">
                <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">Santri Terisi</span>
                <span class="text-xl font-black text-blue-500">{{ Object.values(bulkForm.scores).filter(v => v !== null).length }} / {{ students.length }}</span>
              </div>
            </div>
            <div class="flex gap-4">
              <button @click="showBulkModal = false" class="px-8 py-4 bg-sigma-surface border border-sigma-border rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-rose-500 hover:text-white transition-all active:scale-95">
                Batal
              </button>
              <button @click="saveBulkScores" :disabled="isSaving"
                      class="px-12 py-4 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[11px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center gap-3 active:scale-95 disabled:opacity-50">
                <Save v-if="!isSaving" class="w-5 h-5" />
                <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Simpan Kolektif
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

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
  </div>
</div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }
.modal-enter-active, .modal-leave-active { transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.modal-enter-from, .modal-leave-to { opacity: 0; transform: scale(0.95); }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
</style>
