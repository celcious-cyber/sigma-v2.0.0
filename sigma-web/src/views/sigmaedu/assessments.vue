<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  GraduationCap, Save, 
  CheckCircle2, XCircle,
  Plus, RotateCw, FileDown,
  Trash2, Search, Filter,
  Edit3, LayoutGrid, List
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
const viewMode = ref<'matrix' | 'summary'>('matrix')
const searchQuery = ref('')
const showPersonalModal = ref(false)
const selectedStudent = ref<any>(null)
const personalForm = ref<Record<string, number | null>>({})
const personalRemarks = ref<Record<string, string>>({})

const bulkForm = ref({
  classroom_id: null as number | null,
  subject_id: null as number | null,
  type: 'TUGAS',
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
  const stored = localStorage.getItem('sigma_assessment_categories')
  if (stored) {
    try {
      assessmentTypes.value = JSON.parse(stored)
    } catch (e) {
      console.error('Gagal memuat kategori:', e)
      assessmentTypes.value = [...STANDARD_CATEGORIES]
    }
  } else {
    assessmentTypes.value = [...STANDARD_CATEGORIES]
  }
}

const saveCategories = () => {
  localStorage.setItem('sigma_assessment_categories', JSON.stringify(assessmentTypes.value))
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

const removeCategory = async (val: string) => {
  if (!confirm(`Hapus kolom "${val}" beserta SELURUH nilainya untuk parameter yang dipilih?`)) return

  try {
    isLoading.value = true
    await axios.delete('/api/v1/admin/edu/assessments', {
      params: {
        ...filters.value,
        type: val
      }
    })
    
    assessmentTypes.value = assessmentTypes.value.filter(t => t.value !== val)
    saveCategories()
    await loadAssessmentData()
    
    toastMessage.value = 'Kolom & Nilai Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal menghapus kolom:', err)
    toastMessage.value = 'Gagal menghapus kolom'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isLoading.value = false
  }
}

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
    const [stRes, asRes] = await Promise.all([
      axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`),
      axios.get('/api/v1/admin/edu/assessments', { params: filters.value })
    ])
    
    students.value = stRes.data
    assessments.value = asRes.data

    // Auto-discover types from data
    const discoveredTypes = [...new Set(assessments.value.map(a => a.type))]
    discoveredTypes.forEach(t => {
      if (!assessmentTypes.value.find(at => at.value === t)) {
        assessmentTypes.value.push({ label: t, value: t })
      }
    })

    const newMatrix: Record<number, Record<string, number | null>> = {}
    const newRemarks: Record<number, string> = {}
    
    students.value.forEach(s => {
      newMatrix[s.ID] = {}
      newRemarks[s.ID] = ''
      
      assessmentTypes.value.forEach(t => {
        const found = assessments.value.find(a => a.student_id === s.ID && a.type === t.value)
        newMatrix[s.ID][t.value] = found ? found.score : null
        if (found?.remarks) newRemarks[s.ID] = found.remarks
      })
    })
    
    scoreMatrix.value = newMatrix
    remarksMap.value = newRemarks

  } catch (err) {
    console.error('Gagal memuat data penilaian:', err)
  } finally {
    isLoading.value = false
  }
}

const openBulkModal = (initialType?: string) => {
  bulkForm.value.classroom_id = filters.value.classroom_id
  bulkForm.value.subject_id = filters.value.subject_id
  
  if (initialType) {
    bulkForm.value.type = initialType
  }
  
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

const openPersonalModal = (student: any) => {
  selectedStudent.value = student
  const scores = scoreMatrix.value[student.ID] || {}
  
  const form: Record<string, number | null> = {}
  const rem: Record<string, string> = {}
  
  assessmentTypes.value.forEach(t => {
    form[t.value] = scores[t.value] ?? null
    rem[t.value] = assessments.value.find(a => a.student_id === student.ID && a.type === t.value)?.remarks || ''
  })
  
  personalForm.value = form
  personalRemarks.value = rem
  showPersonalModal.value = true
}

const savePersonalScores = async () => {
  if (!selectedStudent.value) return
  
  isSaving.value = true
  try {
    const payload = assessmentTypes.value.map(t => ({
      student_id: selectedStudent.value.ID,
      classroom_id: filters.value.classroom_id,
      subject_id: filters.value.subject_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: t.value,
      score: personalForm.value[t.value],
      remarks: personalRemarks.value[t.value],
      grade: calculateGrade(personalForm.value[t.value])
    })).filter(a => a.score !== null)

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    await loadAssessmentData()
    showPersonalModal.value = false
    toastMessage.value = `Nilai ${selectedStudent.value.name} Disimpan!`
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan personal:', err)
  } finally {
    isSaving.value = false
  }
}

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
    toastMessage.value = 'Nilai Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan kolektif:', err)
  } finally {
    isSaving.value = false
  }
}

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

const rekapData = computed(() => {
  const filtered = searchQuery.value 
    ? students.value.filter(s => s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.nis.includes(searchQuery.value))
    : students.value

  return filtered.map(s => {
    const scores = scoreMatrix.value[s.ID] || {}
    const scoreValues = Object.values(scores).filter(v => v !== null) as number[]
    const average = scoreValues.length > 0 ? (scoreValues.reduce((a, b) => a + b, 0) / scoreValues.length).toFixed(1) : '-'
    
    return {
      student: s,
      scores,
      average,
      finalGrade: average !== '-' ? calculateGrade(parseFloat(average)) : '-'
    }
  })
})

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
  if (!grade || grade === '-') return 'bg-sigma-surface-alt text-sigma-muted'
  if (grade.startsWith('A')) return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20'
  if (grade.startsWith('B')) return 'bg-blue-500/10 text-blue-500 border-blue-500/20'
  if (grade.startsWith('C')) return 'bg-amber-500/10 text-amber-500 border-amber-500/20'
  if (grade.startsWith('D')) return 'bg-orange-500/10 text-orange-500 border-orange-500/20'
  if (grade.startsWith('E')) return 'bg-rose-500/10 text-rose-500 border-rose-500/20'
  return 'bg-sigma-surface-alt text-sigma-muted'
}

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

// Watchers
watch([
  () => filters.value.classroom_id, 
  () => filters.value.subject_id, 
  () => filters.value.term,
  () => filters.value.academic_year
], loadAssessmentData)

onMounted(() => {
  loadCategories()
  fetchData()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-[#F8FAFC] dark:bg-sigma-bg transition-colors duration-500">
    <!-- Header Section -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-8">
      <div class="flex items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center text-white shadow-2xl shadow-blue-500/40 transform -rotate-6 hover:rotate-0 transition-transform duration-500">
          <GraduationCap class="w-8 h-8" />
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-sigma-muted mb-1 opacity-60">Akademik • Evaluasi</h2>
          <h1 class="text-3xl font-black italic uppercase tracking-tighter leading-none">
            Gradebook <span class="text-blue-500">Performance</span>
          </h1>
        </div>
      </div>

      <div class="flex items-center gap-4 bg-sigma-surface/50 backdrop-blur-xl p-2 rounded-[2rem] border border-sigma-border shadow-xl">
        <button @click="viewMode = 'matrix'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="viewMode === 'matrix' ? 'bg-blue-500 text-white shadow-lg shadow-blue-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <LayoutGrid class="w-3.5 h-3.5" /> Matrix View
        </button>
        <button @click="viewMode = 'summary'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="viewMode === 'summary' ? 'bg-blue-500 text-white shadow-lg shadow-blue-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <List class="w-3.5 h-3.5" /> Summary View
        </button>
      </div>
    </div>

    <!-- Interactive Filter & Action Bar -->
    <div class="grid grid-cols-1 xl:grid-cols-12 gap-8">
      <!-- Parameters -->
      <div class="xl:col-span-9 bg-sigma-surface border border-sigma-border rounded-[3rem] p-8 shadow-sm flex flex-wrap items-end gap-8">
        <div class="flex-1 min-w-[140px] space-y-3">
          <label class="flex items-center gap-2 text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">
            <Filter class="w-3 h-3 text-blue-500" /> Tahun Ajaran
          </label>
          <select v-model="filters.academic_year" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 appearance-none transition-all">
            <option>2024/2025</option>
            <option>2023/2024</option>
          </select>
        </div>
        
        <div class="flex-1 min-w-[120px] space-y-3">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Semester</label>
          <select v-model="filters.term" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 appearance-none transition-all">
            <option>Ganjil</option>
            <option>Genap</option>
          </select>
        </div>

        <div class="flex-1 min-w-[160px] space-y-3">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Kelas</label>
          <select v-model="filters.classroom_id" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 appearance-none transition-all">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>

        <div class="flex-1 min-w-[200px] space-y-3">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Mata Pelajaran</label>
          <select v-model="filters.subject_id" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-500/50 appearance-none transition-all">
            <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
          </select>
        </div>

        <div class="w-px h-10 bg-sigma-border mx-2"></div>

        <div class="flex items-center gap-3">
          <button @click="loadAssessmentData" 
                  class="p-4 bg-blue-500 text-white rounded-2xl hover:bg-blue-600 transition-all shadow-xl shadow-blue-500/20 active:scale-90"
                  title="Refresh Data">
            <RotateCw class="w-5 h-5" :class="{ 'animate-spin': isLoading }" />
          </button>
          <button @click="exportToExcel" 
                  class="p-4 bg-emerald-500 text-white rounded-2xl hover:bg-emerald-600 transition-all shadow-xl shadow-emerald-500/20 active:scale-90"
                  title="Export Excel">
            <FileDown class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Quick Add Category -->
      <div class="xl:col-span-3 bg-sigma-surface border border-sigma-border rounded-[3rem] p-8 shadow-sm space-y-6">
        <div class="space-y-1">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Quick Add Column</label>
          <div class="flex gap-3">
            <input v-model="newCategory" type="text" placeholder="Tugas ..."
                   @keyup.enter="addCategory"
                   class="flex-1 px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-500/50" />
            <button @click="addCategory" class="p-4 bg-blue-500 text-white rounded-2xl hover:bg-blue-600 transition-all shadow-lg shadow-blue-500/20">
              <Plus class="w-5 h-5" />
            </button>
          </div>
        </div>
        <!-- Custom Tags List (Compact) -->
          <div class="flex flex-wrap gap-1 max-w-[200px]">
             <div v-for="t in assessmentTypes" 
                   :key="t.value"
                   class="px-2 py-1 bg-sigma-surface-alt border border-sigma-border rounded-lg text-[8px] font-black flex items-center gap-2 group">
                {{ t.label }}
                <button @click="removeCategory(t.value)" class="text-rose-500 hover:scale-110 transition-transform">
                  <XCircle class="w-3 h-3" />
                </button>
              </div>
          </div>
      </div>
    </div>

    <!-- Content Area -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[3.5rem] shadow-2xl overflow-hidden flex flex-col min-h-[60vh] relative group/table">
      
      <!-- Table Search & Toolbar -->
      <div class="p-8 border-b border-sigma-border bg-sigma-surface-alt/30 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="relative max-w-md w-full">
          <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-500" />
          <input v-model="searchQuery" type="text" placeholder="Search Santri by name or NIS..."
                 class="w-full pl-14 pr-6 py-4 bg-sigma-surface border border-sigma-border rounded-[2rem] text-sm font-bold outline-none focus:border-blue-500/50 shadow-sm transition-all" />
        </div>
        
        <div class="flex items-center gap-4">
          <button @click="() => openBulkModal()" :disabled="!filters.classroom_id"
                  class="px-10 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-blue-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
            <Edit3 class="w-4 h-4" />
            Input Nilai Kolektif
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-40 space-y-6">
        <div class="relative w-16 h-16">
          <div class="absolute inset-0 border-4 border-blue-500/20 rounded-full"></div>
          <div class="absolute inset-0 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <p class="text-[11px] font-black uppercase tracking-[0.3em] text-sigma-muted animate-pulse">Synchronizing Data...</p>
      </div>

      <!-- MATRIX VIEW -->
      <div v-else-if="viewMode === 'matrix' && students.length > 0" class="flex-1 overflow-auto custom-scrollbar">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/80 sticky top-0 z-20 backdrop-blur-md">
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted min-w-[280px] sticky left-0 bg-sigma-surface-alt/90 z-30">Santri Profile</th>
              <th v-for="t in assessmentTypes" :key="t.value" 
                  class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted min-w-[120px] group/th">
                <div class="flex flex-col items-center gap-2">
                  <span class="group-hover/th:text-blue-500 transition-colors">{{ t.label }}</span>
                  <div class="flex items-center gap-1 opacity-0 group-hover/th:opacity-100 transition-all">
                    <button @click="openBulkModal(t.value)" class="p-1.5 rounded-lg bg-blue-500/10 text-blue-500 hover:bg-blue-500 hover:text-white transition-all scale-75" title="Input Kolektif">
                      <Edit3 class="w-3 h-3" />
                    </button>
                    <button @click="removeCategory(t.value)" 
                            class="p-1.5 rounded-lg bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all scale-75" title="Hapus Kolom">
                      <Trash2 class="w-3 h-3" />
                    </button>
                  </div>
                </div>
              </th>
              <th class="p-8 text-center text-[10px] font-black uppercase tracking-[0.2em] text-blue-500 bg-blue-500/[0.03] min-w-[120px]">Final Avg</th>
              <th class="p-8 text-center text-[10px] font-black uppercase tracking-[0.2em] text-blue-500 bg-blue-500/[0.03] min-w-[100px]">Grade</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border/50">
            <tr v-for="r in rekapData" :key="r.student.ID" class="group hover:bg-blue-500/[0.02] transition-all">
              <td class="p-8 sticky left-0 bg-sigma-surface z-10 group-hover:bg-blue-500/[0.02]">
                <div class="flex items-center gap-5">
                  <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border border border-sigma-border flex items-center justify-center text-xs font-black group-hover:scale-110 group-hover:border-blue-500/50 transition-all shadow-sm">
                    {{ r.student.name.charAt(0) }}
                  </div>
                  <div>
                    <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block mb-0.5 opacity-60">{{ r.student.nis }}</span>
                    <span class="text-sm font-black uppercase tracking-tight italic text-sigma-text group-hover:text-blue-600 transition-colors">{{ r.student.name }}</span>
                  </div>
                </div>
              </td>
              
              <td v-for="t in assessmentTypes" :key="t.value" class="p-6 text-center">
                <button v-if="r.scores[t.value] !== null" 
                        @click="openPersonalModal(r.student)"
                        class="inline-flex flex-col items-center justify-center w-14 h-14 rounded-2xl bg-sigma-surface-alt border border-sigma-border shadow-sm hover:border-blue-500 hover:shadow-md transition-all active:scale-90">
                  <span class="text-sm font-black text-sigma-text">{{ r.scores[t.value] }}</span>
                </button>
                <button v-else 
                        @click="openPersonalModal(r.student)"
                        class="text-[10px] font-black text-sigma-muted opacity-20 hover:opacity-100 transition-all p-4">
                  —
                </button>
              </td>

              <td class="p-8 text-center bg-blue-500/[0.01] group-hover:bg-blue-500/[0.03]">
                <span class="text-lg font-black text-blue-500 tracking-tighter">{{ r.average }}</span>
              </td>
              
              <td class="p-8 text-center bg-blue-500/[0.01] group-hover:bg-blue-500/[0.03]">
                <div class="inline-flex px-5 py-2 rounded-xl text-[11px] font-black italic border shadow-sm transition-transform group-hover:scale-110"
                     :class="getGradeColor(r.finalGrade)">
                  {{ r.finalGrade }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- SUMMARY VIEW -->
      <div v-else-if="viewMode === 'summary' && students.length > 0" class="flex-1 p-8 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 overflow-y-auto custom-scrollbar">
        <div v-for="r in rekapData" :key="r.student.ID" 
             class="bg-sigma-surface-alt/30 border border-sigma-border rounded-[2.5rem] p-8 space-y-6 hover:shadow-2xl hover:shadow-blue-500/10 hover:border-blue-500/30 transition-all group">
           <div class="flex items-center justify-between">
              <div class="flex items-center gap-4">
                 <div class="w-12 h-12 rounded-2xl bg-sigma-surface border border-sigma-border flex items-center justify-center font-black group-hover:scale-110 transition-transform shadow-sm">
                   {{ r.student.name.charAt(0) }}
                 </div>
                 <div>
                   <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ r.student.nis }}</span>
                   <h3 class="text-sm font-black uppercase italic tracking-tight">{{ r.student.name }}</h3>
                 </div>
              </div>
              <div class="flex items-center gap-2">
                 <button @click="openPersonalModal(r.student)" 
                         class="p-2.5 rounded-xl bg-blue-500/10 text-blue-500 hover:bg-blue-500 hover:text-white transition-all active:scale-90"
                         title="Input Personal">
                   <Edit3 class="w-4 h-4" />
                 </button>
                 <button @click="deleteStudentScores(r.student)" 
                         class="p-2.5 rounded-xl bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all active:scale-90"
                         title="Hapus Seluruh Nilai Santri">
                   <Trash2 class="w-4 h-4" />
                 </button>
               </div>
           </div>

           <div class="grid grid-cols-2 gap-4">
              <div class="bg-blue-500/5 rounded-2xl p-4 border border-blue-500/10 text-center">
                 <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block mb-1">Rata-Rata</span>
                 <span class="text-xl font-black text-blue-500">{{ r.average }}</span>
              </div>
              <div class="rounded-2xl p-4 border text-center flex flex-col items-center justify-center" :class="getGradeColor(r.finalGrade)">
                 <span class="text-[8px] font-black uppercase tracking-widest block mb-1 opacity-60">Grade</span>
                 <span class="text-xl font-black italic">{{ r.finalGrade }}</span>
              </div>
           </div>

           <div class="space-y-3">
              <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block ml-1">Detail Skor:</span>
              <div class="flex flex-wrap gap-2">
                 <div v-for="t in assessmentTypes" :key="t.value" 
                      class="px-3 py-1.5 rounded-lg text-[9px] font-black flex items-center gap-2 border"
                      :class="r.scores[t.value] !== null ? 'bg-sigma-surface border-blue-500/20 text-blue-600' : 'bg-sigma-surface-alt/50 border-sigma-border text-sigma-muted opacity-30'">
                   {{ t.label }}: {{ r.scores[t.value] ?? '-' }}
                 </div>
              </div>
           </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="flex-1 flex flex-col items-center justify-center p-40 space-y-8 opacity-40">
        <div class="w-32 h-32 rounded-full bg-sigma-surface-alt border-4 border-dashed border-sigma-border flex items-center justify-center">
          <LayoutGrid class="w-12 h-12 text-sigma-muted" />
        </div>
        <div class="text-center space-y-2">
          <p class="font-black uppercase tracking-[0.4em] text-xs">No Data Synchronized</p>
          <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Please select classroom and subject parameters</p>
        </div>
      </div>
    </div>

    <!-- COLLECTIVE INPUT MODAL -->
    <Transition name="modal">
      <div v-if="showBulkModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl transition-all duration-500" @click="showBulkModal = false"></div>
        <div class="relative w-full max-w-6xl bg-sigma-surface border border-sigma-border rounded-[4rem] shadow-[0_50px_100px_-20px_rgba(0,0,0,0.5)] overflow-hidden animate-in fade-in zoom-in slide-in-from-bottom-12 duration-700">
          
          <div class="p-10 border-b border-sigma-border bg-gradient-to-r from-sigma-surface-alt/50 to-transparent flex items-center justify-between">
            <div class="flex items-center gap-6">
              <div class="w-14 h-14 rounded-[1.5rem] bg-blue-600 flex items-center justify-center text-white shadow-xl shadow-blue-600/30">
                <Edit3 class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Kolektif <span class="text-blue-600">Entry</span></h3>
                <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">Batch grading for multiple students</p>
              </div>
            </div>
            <button @click="showBulkModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-3 bg-sigma-surface-alt/30 border-b border-sigma-border">
             <div class="p-8 space-y-3 border-r border-sigma-border">
                <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">Mata Pelajaran</label>
                <select v-model="bulkForm.subject_id" class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-600/50 shadow-sm transition-all">
                  <option v-for="sub in subjects" :key="sub.ID" :value="sub.ID">{{ sub.name }}</option>
                </select>
             </div>
             <div class="p-8 space-y-3 border-r border-sigma-border">
                <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kategori Penilaian</label>
                <select v-model="bulkForm.type" class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-600/50 shadow-sm transition-all">
                  <option v-for="t in assessmentTypes" :key="t.value" :value="t.value">{{ t.label }}</option>
                </select>
             </div>
             <div class="p-8 space-y-3">
                <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal Evaluasi</label>
                <input v-model="bulkForm.date" type="date" class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-blue-600/50 shadow-sm transition-all" />
             </div>
          </div>

          <div class="p-0 h-[45vh] overflow-y-auto custom-scrollbar bg-sigma-surface">
            <table class="w-full border-collapse">
              <thead class="sticky top-0 z-10 bg-sigma-surface-alt/90 backdrop-blur-md">
                <tr>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pl-12">Santri Information</th>
                  <th class="p-6 text-center text-[9px] font-black uppercase tracking-widest text-sigma-muted w-40">Score (0-100)</th>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pr-12">Evaluation Note</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border/40">
                <tr v-for="s in students" :key="s.ID" class="group hover:bg-blue-500/[0.03] transition-all">
                  <td class="p-6 pl-12">
                    <div class="flex items-center gap-4">
                      <div class="w-10 h-10 rounded-xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center font-black group-hover:border-blue-500/50 transition-all">{{ s.name.charAt(0) }}</div>
                      <div>
                        <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">{{ s.nis }}</span>
                        <span class="text-sm font-black uppercase tracking-tight">{{ s.name }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="p-6 px-10">
                    <input v-model.number="bulkForm.scores[s.ID]" type="number" 
                           class="w-full px-4 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-center font-black text-blue-600 outline-none focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 transition-all text-lg"
                           placeholder="0" />
                  </td>
                  <td class="p-6 pr-12">
                    <input v-model="bulkForm.remarks[s.ID]" type="text"
                           class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-bold outline-none focus:border-blue-600/50 transition-all"
                           placeholder="Teacher's feedback..." />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-between">
            <div class="flex items-center gap-8">
              <div class="flex flex-col">
                <span class="text-[9px] font-black uppercase tracking-widest text-sigma-muted block mb-1">Processed Students</span>
                <div class="flex items-end gap-2 leading-none">
                   <span class="text-3xl font-black text-blue-600">{{ Object.values(bulkForm.scores).filter(v => v !== null).length }}</span>
                   <span class="text-lg font-black text-sigma-muted">/ {{ students.length }}</span>
                </div>
              </div>
            </div>
            <div class="flex gap-4">
              <button @click="showBulkModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all active:scale-95">
                Cancel Session
              </button>
              <button @click="saveBulkScores" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-blue-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
                <Save v-if="!isSaving" class="w-5 h-5" />
                <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Commit Grades
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- PERSONAL INPUT MODAL -->
    <Transition name="modal">
      <div v-if="showPersonalModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl" @click="showPersonalModal = false"></div>
        <div class="relative w-full max-w-2xl bg-sigma-surface border border-sigma-border rounded-[4rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-500">
          
          <div class="p-10 border-b border-sigma-border flex items-center justify-between bg-sigma-surface-alt/30">
            <div class="flex items-center gap-6">
              <div class="w-14 h-14 rounded-[1.5rem] bg-blue-600 flex items-center justify-center text-white shadow-xl shadow-blue-600/30">
                <GraduationCap class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Personal <span class="text-blue-600">Grade</span></h3>
                <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">{{ selectedStudent?.name }}</p>
              </div>
            </div>
            <button @click="showPersonalModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="p-10 space-y-6 max-h-[50vh] overflow-y-auto custom-scrollbar">
             <div v-for="t in assessmentTypes" :key="t.value" class="bg-sigma-surface-alt/50 p-6 rounded-3xl border border-sigma-border space-y-4">
                <div class="flex items-center justify-between">
                   <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ t.label }}</span>
                   <div v-if="personalForm[t.value] !== null" class="px-3 py-1 rounded-lg text-[10px] font-black" :class="getGradeColor(calculateGrade(personalForm[t.value]))">
                      Grade {{ calculateGrade(personalForm[t.value]) }}
                   </div>
                </div>
                <div class="grid grid-cols-4 gap-4">
                   <div class="col-span-1">
                      <input v-model.number="personalForm[t.value]" type="number" 
                             class="w-full px-4 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-center font-black text-blue-600 outline-none focus:border-blue-600 shadow-sm transition-all"
                             placeholder="Score" />
                   </div>
                   <div class="col-span-3">
                      <input v-model="personalRemarks[t.value]" type="text"
                             class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-bold outline-none focus:border-blue-600/30 transition-all"
                             placeholder="Note / Remark..." />
                   </div>
                </div>
             </div>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-end gap-4">
              <button @click="showPersonalModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all">
                Cancel
              </button>
              <button @click="savePersonalScores" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-blue-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
                <Save v-if="!isSaving" class="w-5 h-5" />
                <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Save Changes
              </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Notification Toast -->
    <Transition name="toast">
      <div v-if="showToast" 
           class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[1000] flex items-center gap-5 px-10 py-5 rounded-[2.5rem] border shadow-2xl backdrop-blur-3xl transition-all duration-500"
           :class="toastType === 'success' ? 'bg-emerald-500/90 border-emerald-400/50 text-white shadow-emerald-500/30' : 'bg-rose-500/90 border-rose-400/50 text-white shadow-rose-500/30'">
        <div class="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center">
          <CheckCircle2 v-if="toastType === 'success'" class="w-6 h-6" />
          <XCircle v-else class="w-6 h-6" />
        </div>
        <div class="flex flex-col">
           <span class="text-[8px] font-black uppercase tracking-widest opacity-70">System Message</span>
           <span class="text-xs font-black uppercase tracking-[0.15em]">{{ toastMessage }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1); }
.toast-enter-from { opacity: 0; transform: translate(-50%, 100px) scale(0.8); }
.toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }

.modal-enter-active, .modal-leave-active { transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1); }
.modal-enter-from { opacity: 0; transform: scale(0.9) translateY(40px); }
.modal-leave-to { opacity: 0; transform: scale(0.95); }

.custom-scrollbar::-webkit-scrollbar { width: 6px; height: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.3); }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}

/* Glass effect for active selects */
select:focus {
  background-color: var(--sigma-surface);
  box-shadow: 0 10px 25px -5px rgba(59, 130, 246, 0.1);
}

.animate-in {
  animation-fill-mode: both;
}
</style>
