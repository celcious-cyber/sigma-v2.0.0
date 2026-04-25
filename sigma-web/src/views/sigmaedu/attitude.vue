<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  Save, 
  CheckCircle2, XCircle,
  Plus, RotateCw,
  Trash2, Search, Filter,
  Edit3, LayoutGrid, List,
  Star, Heart
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const assessments = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const showBulkModal = ref(false)
const viewMode = ref<'matrix' | 'summary'>('matrix')
const searchQuery = ref('')
const showPersonalModal = ref(false)
const selectedStudent = ref<any>(null)
const personalForm = ref<Record<string, string>>({})
const personalRemarks = ref<Record<string, string>>({})

const bulkForm = ref({
  classroom_id: null as number | null,
  type: 'KEAKTIFAN',
  date: new Date().toISOString().split('T')[0],
  grades: {} as Record<number, string>,
  remarks: {} as Record<number, string>
})

const filters = ref({
  classroom_id: null as number | null,
  term: 'Ganjil',
  academic_year: '2024/2025'
})

const STANDARD_CHR_CATEGORIES = [
  { label: 'Keaktifan', value: 'KEAKTIFAN' },
  { label: 'Kedisiplinan', value: 'KEDISIPLINAN' },
  { label: 'Kerapihan', value: 'KERAPIHAN' },
  { label: 'Kesantunan', value: 'KESANTUNAN' },
  { label: 'Ketekunan', value: 'KETEKUNAN' },
]

const attitudeTypes = ref([...STANDARD_CHR_CATEGORIES])

const loadCategories = () => {
  const stored = localStorage.getItem('sigma_attitude_categories')
  if (stored) {
    try {
      attitudeTypes.value = JSON.parse(stored)
    } catch (e) {
      console.error('Gagal memuat kategori:', e)
      attitudeTypes.value = [...STANDARD_CHR_CATEGORIES]
    }
  } else {
    attitudeTypes.value = [...STANDARD_CHR_CATEGORIES]
  }
}

const saveCategories = () => {
  localStorage.setItem('sigma_attitude_categories', JSON.stringify(attitudeTypes.value))
}

const newCategory = ref('')
const addCategory = () => {
  if (!newCategory.value) return
  const val = newCategory.value.toUpperCase().replace(/\s+/g, '_')
  if (!attitudeTypes.value.find(t => t.value === val)) {
    attitudeTypes.value.push({ label: newCategory.value, value: val })
    saveCategories()
  }
  newCategory.value = ''
}

const removeCategory = async (val: string) => {
  if (!confirm(`Hapus kolom "${val}" beserta SELURUH nilai karakternya untuk parameter yang dipilih?`)) return

  try {
    isLoading.value = true
    await axios.delete('/api/v1/admin/edu/assessments', {
      params: {
        ...filters.value,
        subject_id: 0,
        type: val
      }
    })
    
    attitudeTypes.value = attitudeTypes.value.filter(t => t.value !== val)
    saveCategories()
    await loadAttitudeData()
    
    toastMessage.value = 'Kolom Karakter Berhasil Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal menghapus kolom:', err)
  } finally {
    isLoading.value = false
  }
}

const fetchData = async () => {
  try {
    const clRes = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = clRes.data
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
    await loadAttitudeData()
  } catch (err) {
    console.error('Gagal mengambil data master:', err)
  }
}

const loadAttitudeData = async () => {
  if (!filters.value.classroom_id) return
  
  isLoading.value = true
  try {
    const [stRes, asRes] = await Promise.all([
      axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`),
      axios.get('/api/v1/admin/edu/assessments', { params: { ...filters.value, subject_id: 0 } })
    ])
    
    students.value = stRes.data
    assessments.value = asRes.data

    // Auto-discover attitude types from data if they match some pattern or just any subject-less assessment
    const discoveredTypes = [...new Set(assessments.value.map(a => a.type))]
    discoveredTypes.forEach(t => {
      if (!attitudeTypes.value.find(at => at.value === t)) {
        attitudeTypes.value.push({ label: t, value: t })
      }
    })

    const newMatrix: Record<number, Record<string, string>> = {}
    students.value.forEach(s => {
      newMatrix[s.ID] = {}
      attitudeTypes.value.forEach(t => {
        const found = assessments.value.find(a => a.student_id === s.ID && a.type === t.value)
        newMatrix[s.ID][t.value] = found ? found.grade : ''
      })
    })
    
    gradeMatrix.value = newMatrix
  } catch (err) {
    console.error('Gagal memuat data karakter:', err)
  } finally {
    isLoading.value = false
  }
}

const gradeMatrix = ref<Record<number, Record<string, string>>>({})

const openBulkModal = (initialType?: string) => {
  bulkForm.value.classroom_id = filters.value.classroom_id
  if (initialType) bulkForm.value.type = initialType
  
  const newGrades: Record<number, string> = {}
  const newRemarks: Record<number, string> = {}
  
  students.value.forEach(s => {
    newGrades[s.ID] = gradeMatrix.value[s.ID]?.[bulkForm.value.type] || 'B'
    const found = assessments.value.find(a => a.student_id === s.ID && a.type === bulkForm.value.type)
    newRemarks[s.ID] = found?.remarks || ''
  })
  
  bulkForm.value.grades = newGrades
  bulkForm.value.remarks = newRemarks
  showBulkModal.value = true
}

const openPersonalModal = (student: any) => {
  selectedStudent.value = student
  const grades = gradeMatrix.value[student.ID] || {}
  
  const form: Record<string, string> = {}
  const rem: Record<string, string> = {}
  
  attitudeTypes.value.forEach(t => {
    form[t.value] = grades[t.value] || ''
    rem[t.value] = assessments.value.find(a => a.student_id === student.ID && a.type === t.value)?.remarks || ''
  })
  
  personalForm.value = form
  personalRemarks.value = rem
  showPersonalModal.value = true
}

const savePersonalGrades = async () => {
  if (!selectedStudent.value) return
  
  isSaving.value = true
  try {
    const payload = attitudeTypes.value.map(t => ({
      student_id: selectedStudent.value.ID,
      classroom_id: filters.value.classroom_id,
      subject_id: null,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: t.value,
      grade: personalForm.value[t.value],
      remarks: personalRemarks.value[t.value],
      score: null
    })).filter(a => a.grade !== '')

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    await loadAttitudeData()
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

const saveBulkGrades = async () => {
  if (!bulkForm.value.classroom_id) return
  
  isSaving.value = true
  try {
    const payload = students.value.map(s => ({
      student_id: s.ID,
      classroom_id: bulkForm.value.classroom_id,
      subject_id: null,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: bulkForm.value.type,
      grade: bulkForm.value.grades[s.ID],
      remarks: bulkForm.value.remarks[s.ID],
      score: null
    })).filter(a => a.grade !== '')

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    await loadAttitudeData()
    showBulkModal.value = false
    toastMessage.value = 'Nilai Karakter Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan kolektif:', err)
  } finally {
    isSaving.value = false
  }
}

const deleteStudentAttitudes = async (student: any) => {
  if (!confirm(`Hapus seluruh nilai karakter untuk ${student.name} di parameter ini?`)) return
  
  isSaving.value = true
  try {
    const payload = attitudeTypes.value.map(t => ({
      student_id: student.ID,
      classroom_id: filters.value.classroom_id,
      subject_id: null,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: t.value,
      grade: '',
      remarks: '',
      score: null
    }))

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    await loadAttitudeData()
    toastMessage.value = 'Data Karakter Santri Dihapus'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal menghapus nilai karakter:', err)
  } finally {
    isSaving.value = false
  }
}

const rekapData = computed(() => {
  const filtered = searchQuery.value 
    ? students.value.filter(s => s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.nis.includes(searchQuery.value))
    : students.value

  return filtered.map(s => {
    const grades = gradeMatrix.value[s.ID] || {}
    const activeGrades = Object.values(grades).filter(g => g !== '')
    
    // Summary logic: Most frequent grade
    let summary = '-'
    if (activeGrades.length > 0) {
      const counts: Record<string, number> = { 'A': 0, 'B': 0, 'C': 0, 'D': 0 }
      activeGrades.forEach(g => { if (counts[g] !== undefined) counts[g]++ })
      
      let max = -1
      const priority = ['A', 'B', 'C', 'D']
      priority.forEach(p => {
        if (counts[p] > max) {
          max = counts[p]
          summary = p
        }
      })
    }
    
    return {
      student: s,
      grades,
      summary
    }
  })
})

const getGradeColor = (grade: string) => {
  if (!grade || grade === '-') return 'bg-sigma-surface-alt text-sigma-muted'
  if (grade === 'A') return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20'
  if (grade === 'B') return 'bg-blue-500/10 text-blue-500 border-blue-500/20'
  if (grade === 'C') return 'bg-amber-500/10 text-amber-500 border-amber-500/20'
  if (grade === 'D') return 'bg-rose-500/10 text-rose-500 border-rose-500/20'
  return 'bg-sigma-surface-alt text-sigma-muted'
}

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

watch([() => filters.value.classroom_id, () => filters.value.term], loadAttitudeData)

onMounted(() => {
  loadCategories()
  fetchData()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-[#FDFCFB] dark:bg-sigma-bg transition-colors duration-500">
    <!-- Header Section -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-8">
      <div class="flex items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-gradient-to-br from-rose-500 to-rose-700 flex items-center justify-center text-white shadow-2xl shadow-rose-500/40 transform -rotate-6 hover:rotate-0 transition-transform duration-500">
          <Heart class="w-8 h-8" />
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-sigma-muted mb-1 opacity-60">Akademik • Adab & Karakter</h2>
          <h1 class="text-3xl font-black italic uppercase tracking-tighter leading-none">
            Attitude <span class="text-rose-500">Analytics</span>
          </h1>
        </div>
      </div>

      <div class="flex items-center gap-4 bg-sigma-surface/50 backdrop-blur-xl p-2 rounded-[2rem] border border-sigma-border shadow-xl">
        <button @click="viewMode = 'matrix'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="viewMode === 'matrix' ? 'bg-rose-500 text-white shadow-lg shadow-rose-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <LayoutGrid class="w-3.5 h-3.5" /> Matrix View
        </button>
        <button @click="viewMode = 'summary'" 
                class="flex items-center gap-2 px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="viewMode === 'summary' ? 'bg-rose-500 text-white shadow-lg shadow-rose-500/20' : 'text-sigma-muted hover:text-sigma-text'">
          <List class="w-3.5 h-3.5" /> Summary View
        </button>
      </div>
    </div>

    <!-- Filter & Action Bar -->
    <div class="grid grid-cols-1 xl:grid-cols-12 gap-8">
      <div class="xl:col-span-9 bg-sigma-surface border border-sigma-border rounded-[3rem] p-8 shadow-sm flex flex-wrap items-end gap-8">
        <div class="flex-1 min-w-[140px] space-y-3">
          <label class="flex items-center gap-2 text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1">
            <Filter class="w-3 h-3 text-rose-500" /> Tahun Ajaran
          </label>
          <select v-model="filters.academic_year" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-rose-500/50 appearance-none transition-all">
            <option>2024/2025</option><option>2023/2024</option>
          </select>
        </div>
        
        <div class="flex-1 min-w-[120px] space-y-3">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Semester</label>
          <select v-model="filters.term" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-rose-500/50 appearance-none transition-all">
            <option>Ganjil</option><option>Genap</option>
          </select>
        </div>

        <div class="flex-1 min-w-[160px] space-y-3">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Kelas</label>
          <select v-model="filters.classroom_id" class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-rose-500/50 appearance-none transition-all">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>

        <button @click="loadAttitudeData" 
                class="p-4 bg-rose-500 text-white rounded-2xl hover:bg-rose-600 transition-all shadow-xl shadow-rose-500/20 active:scale-90"
                title="Refresh Data">
          <RotateCw class="w-5 h-5" :class="{ 'animate-spin': isLoading }" />
        </button>
      </div>

      <!-- Quick Add Category -->
      <div class="xl:col-span-3 bg-sigma-surface border border-sigma-border rounded-[3rem] p-8 shadow-sm space-y-6">
        <div class="space-y-1">
          <label class="text-[9px] font-black uppercase tracking-widest text-sigma-muted ml-1 block">Quick Add Column</label>
          <div class="flex gap-3">
            <input v-model="newCategory" type="text" placeholder="Keaktifan..."
                   @keyup.enter="addCategory"
                   class="flex-1 px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-xs font-black outline-none focus:border-rose-500/50" />
            <button @click="addCategory" class="p-4 bg-rose-500 text-white rounded-2xl hover:bg-rose-600 transition-all shadow-lg shadow-rose-500/20">
              <Plus class="w-5 h-5" />
            </button>
          </div>
        </div>
        <div class="flex flex-wrap gap-1">
           <div v-for="t in attitudeTypes" :key="t.value"
                 class="px-2 py-1 bg-sigma-surface-alt border border-sigma-border rounded-lg text-[8px] font-black flex items-center gap-2 group animate-in zoom-in">
              {{ t.label }}
              <button @click="removeCategory(t.value)" class="text-rose-500 hover:scale-125 transition-transform">
                <XCircle class="w-3 h-3" />
              </button>
            </div>
        </div>
      </div>
    </div>

    <!-- Content Area -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[3.5rem] shadow-2xl overflow-hidden flex flex-col min-h-[60vh] relative">
      
      <!-- Table Search & Toolbar -->
      <div class="p-8 border-b border-sigma-border bg-sigma-surface-alt/30 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="relative max-w-md w-full">
          <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-5 h-5 text-rose-500" />
          <input v-model="searchQuery" type="text" placeholder="Search Santri..."
                 class="w-full pl-14 pr-6 py-4 bg-sigma-surface border border-sigma-border rounded-[2rem] text-sm font-bold outline-none focus:border-rose-500/50 shadow-sm transition-all" />
        </div>
        
        <div class="flex items-center gap-4">
          <button @click="() => openBulkModal()" :disabled="!filters.classroom_id"
                  class="px-10 py-4 bg-gradient-to-r from-rose-600 to-rose-700 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-rose-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
            <Edit3 class="w-4 h-4" /> Input Nilai Kolektif
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-40 space-y-6">
        <div class="w-16 h-16 border-4 border-rose-500/20 border-t-rose-500 rounded-full animate-spin"></div>
        <p class="text-[11px] font-black uppercase tracking-[0.3em] text-sigma-muted">Synchronizing Character Data...</p>
      </div>

      <!-- MATRIX VIEW -->
      <div v-else-if="viewMode === 'matrix' && students.length > 0" class="flex-1 overflow-auto custom-scrollbar">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/80 sticky top-0 z-20 backdrop-blur-md">
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted min-w-[280px] sticky left-0 bg-sigma-surface-alt/90 z-30">Santri Profile</th>
              <th v-for="t in attitudeTypes" :key="t.value" 
                  class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted min-w-[120px] group/th">
                <div class="flex flex-col items-center gap-2">
                  <span class="group-hover/th:text-rose-500 transition-colors">{{ t.label }}</span>
                  <div class="flex items-center gap-1 opacity-0 group-hover/th:opacity-100 transition-all">
                    <button @click="openBulkModal(t.value)" class="p-1.5 rounded-lg bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all scale-75">
                      <Edit3 class="w-3 h-3" />
                    </button>
                    <button @click="removeCategory(t.value)" class="p-1.5 rounded-lg bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all scale-75">
                      <Trash2 class="w-3 h-3" />
                    </button>
                  </div>
                </div>
              </th>
              <th class="p-8 text-center text-[10px] font-black uppercase tracking-[0.2em] text-rose-500 bg-rose-500/[0.03] min-w-[100px]">Summary</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border/50">
            <tr v-for="r in rekapData" :key="r.student.ID" class="group hover:bg-rose-500/[0.02] transition-all">
              <td class="p-8 sticky left-0 bg-sigma-surface z-10 group-hover:bg-rose-500/[0.02]">
                <div class="flex items-center gap-5">
                  <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border border border-sigma-border flex items-center justify-center text-xs font-black group-hover:scale-110 group-hover:border-rose-500/50 transition-all">
                    {{ r.student.name.charAt(0) }}
                  </div>
                  <div>
                    <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block mb-0.5 opacity-60">{{ r.student.nis }}</span>
                    <span class="text-sm font-black uppercase tracking-tight italic text-sigma-text group-hover:text-rose-600 transition-colors">{{ r.student.name }}</span>
                  </div>
                </div>
              </td>
              
              <td v-for="t in attitudeTypes" :key="t.value" class="p-6 text-center">
                <button v-if="r.grades[t.value]" 
                     @click="openPersonalModal(r.student)"
                     class="inline-flex flex-col items-center justify-center w-12 h-12 rounded-xl text-sm font-black border transition-all"
                     :class="getGradeColor(r.grades[t.value])">
                  {{ r.grades[t.value] }}
                </button>
                <button v-else 
                        @click="openPersonalModal(r.student)"
                        class="text-[10px] font-black text-sigma-muted opacity-20 hover:opacity-100 transition-all p-4">
                  —
                </button>
              </td>

              <td class="p-8 text-center bg-rose-500/[0.01] group-hover:bg-rose-500/[0.03]">
                <div class="inline-flex items-center gap-2 px-5 py-2 rounded-xl text-[11px] font-black italic border shadow-sm transition-transform group-hover:scale-110"
                     :class="getGradeColor(r.summary)">
                  <Star class="w-3.5 h-3.5" /> {{ r.summary }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- SUMMARY VIEW -->
      <div v-else-if="viewMode === 'summary' && students.length > 0" class="flex-1 p-8 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 overflow-y-auto custom-scrollbar">
        <div v-for="r in rekapData" :key="r.student.ID" 
             class="bg-sigma-surface-alt/30 border border-sigma-border rounded-[2.5rem] p-8 space-y-6 hover:shadow-2xl hover:shadow-rose-500/10 hover:border-rose-500/30 transition-all group">
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
                        class="p-2.5 rounded-xl bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all active:scale-90"
                        title="Input Personal">
                  <Edit3 class="w-4 h-4" />
                </button>
                <button @click="deleteStudentAttitudes(r.student)" 
                        class="p-2.5 rounded-xl bg-rose-500/10 text-rose-500 hover:bg-rose-500 hover:text-white transition-all active:scale-90"
                        title="Hapus Nilai Karakter Santri">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
           </div>

           <div class="bg-rose-500/5 rounded-2xl p-6 border border-rose-500/10 text-center flex flex-col items-center justify-center" :class="getGradeColor(r.summary)">
              <span class="text-[8px] font-black uppercase tracking-widest block mb-1 opacity-60">Attitude Summary</span>
              <div class="flex items-center gap-2">
                 <Star class="w-5 h-5" />
                 <span class="text-2xl font-black italic">{{ r.summary }}</span>
              </div>
           </div>

           <div class="space-y-3">
              <div class="flex flex-wrap gap-2">
                 <div v-for="t in attitudeTypes" :key="t.value" 
                      class="px-3 py-1.5 rounded-lg text-[9px] font-black flex items-center gap-2 border"
                      :class="r.grades[t.value] ? getGradeColor(r.grades[t.value]) : 'bg-sigma-surface-alt/50 border-sigma-border text-sigma-muted opacity-30'">
                   {{ t.label }}: {{ r.grades[t.value] || '—' }}
                 </div>
              </div>
           </div>
        </div>
      </div>
    </div>

    <!-- COLLECTIVE INPUT MODAL -->
    <Transition name="modal">
      <div v-if="showBulkModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl" @click="showBulkModal = false"></div>
        <div class="relative w-full max-w-5xl bg-sigma-surface border border-sigma-border rounded-[4rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-500">
          
          <div class="p-10 border-b border-sigma-border flex items-center justify-between">
            <div class="flex items-center gap-6">
              <div class="w-14 h-14 rounded-[1.5rem] bg-rose-600 flex items-center justify-center text-white shadow-xl shadow-rose-600/30">
                <Star class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Batch <span class="text-rose-600">Grading</span></h3>
                <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">Karakter • {{ bulkForm.type }}</p>
              </div>
            </div>
            <button @click="showBulkModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="p-0 h-[45vh] overflow-y-auto custom-scrollbar bg-sigma-surface">
            <table class="w-full border-collapse">
              <thead class="sticky top-0 z-10 bg-sigma-surface-alt/90 backdrop-blur-md border-b border-sigma-border">
                <tr>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pl-12">Santri</th>
                  <th class="p-6 text-center text-[9px] font-black uppercase tracking-widest text-sigma-muted">Grade Select</th>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pr-12">Catatan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border/40">
                <tr v-for="s in students" :key="s.ID" class="group hover:bg-rose-500/[0.03] transition-all">
                  <td class="p-6 pl-12">
                    <span class="text-sm font-black uppercase tracking-tight">{{ s.name }}</span>
                  </td>
                  <td class="p-6">
                    <div class="flex items-center justify-center gap-2">
                       <button v-for="g in ['A', 'B', 'C', 'D']" :key="g"
                               @click="bulkForm.grades[s.ID] = g"
                               class="w-12 h-12 rounded-2xl text-xs font-black border transition-all"
                               :class="bulkForm.grades[s.ID] === g ? getGradeColor(g) + ' border-rose-500 shadow-md' : 'bg-sigma-surface-alt border-transparent text-sigma-muted hover:border-sigma-border'">
                         {{ g }}
                       </button>
                    </div>
                  </td>
                  <td class="p-6 pr-12">
                    <input v-model="bulkForm.remarks[s.ID]" type="text"
                           class="w-full px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-bold outline-none focus:border-rose-600/50 transition-all"
                           placeholder="Teacher's feedback..." />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-end gap-4">
              <button @click="showBulkModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all active:scale-95">
                Cancel
              </button>
              <button @click="saveBulkGrades" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-rose-600 to-rose-700 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-rose-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
                <Save v-if="!isSaving" class="w-5 h-5" />
                <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Commit Grades
              </button>
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
              <div class="w-14 h-14 rounded-[1.5rem] bg-rose-600 flex items-center justify-center text-white shadow-xl shadow-rose-600/30">
                <Star class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Personal <span class="text-rose-600">Adab</span></h3>
                <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">{{ selectedStudent?.name }}</p>
              </div>
            </div>
            <button @click="showPersonalModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="p-10 space-y-6 max-h-[50vh] overflow-y-auto custom-scrollbar">
             <div v-for="t in attitudeTypes" :key="t.value" class="bg-sigma-surface-alt/50 p-6 rounded-3xl border border-sigma-border space-y-4">
                <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ t.label }}</span>
                <div class="grid grid-cols-2 gap-4">
                   <div class="flex items-center gap-1.5">
                      <button v-for="g in ['A', 'B', 'C', 'D']" :key="g"
                               @click="personalForm[t.value] = g"
                               class="w-10 h-10 rounded-xl text-[10px] font-black border transition-all"
                               :class="personalForm[t.value] === g ? getGradeColor(g) + ' border-rose-500 shadow-md' : 'bg-sigma-surface border-transparent text-sigma-muted hover:border-sigma-border'">
                         {{ g }}
                      </button>
                   </div>
                   <input v-model="personalRemarks[t.value]" type="text"
                          class="w-full px-6 py-3 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-bold outline-none focus:border-rose-600/30 transition-all"
                          placeholder="Catatan..." />
                </div>
             </div>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-end gap-4">
              <button @click="showPersonalModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all">
                Cancel
              </button>
              <button @click="savePersonalGrades" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-rose-600 to-rose-700 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-rose-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
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
           class="fixed bottom-10 left-1/2 -translate-x-1/2 z-[1000] flex items-center gap-5 px-10 py-5 rounded-[2.5rem] border shadow-2xl backdrop-blur-3xl transition-all"
           :class="toastType === 'success' ? 'bg-emerald-500/90 border-emerald-400/50 text-white' : 'bg-rose-500/90 border-rose-400/50 text-white'">
        <CheckCircle2 v-if="toastType === 'success'" class="w-6 h-6" />
        <XCircle v-else class="w-6 h-6" />
        <span class="text-xs font-black uppercase tracking-[0.15em]">{{ toastMessage }}</span>
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
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(244, 63, 94, 0.1); border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(244, 63, 94, 0.3); }

.animate-in { animation-fill-mode: both; }
</style>
