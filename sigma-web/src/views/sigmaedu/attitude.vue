<script setup lang="ts">
import { ref, onMounted, watch, computed, reactive } from 'vue'
import { 
  Heart, Save, 
  CheckCircle2, XCircle, Star,
  LayoutGrid, List, Sparkles, Clock, BookOpen,
  MessageSquare
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const activeTab = ref('input') // 'input' or 'rekap'

const filters = ref({
  classroom_id: null as number | null,
  term: 'Ganjil',
  academic_year: '2024/2025'
})

const attitudeCategories = [
  { id: 'participation', label: 'Keaktifan', icon: MessageSquare, color: 'text-blue-500' },
  { id: 'discipline', label: 'Kedisiplinan', icon: Clock, color: 'text-emerald-500' },
  { id: 'neatness', label: 'Kerapihan', icon: Sparkles, color: 'text-amber-500' },
  { id: 'etiquette', label: 'Kesantunan', icon: Heart, color: 'text-rose-500' },
  { id: 'diligence', label: 'Ketekunan', icon: BookOpen, color: 'text-indigo-500' },
]

// UI State
const attitudeMap = reactive<Record<number, Record<string, string>>>({})

const fetchData = async () => {
  try {
    const clRes = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = clRes.data
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
    await loadAttitudeData()
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  }
}

const loadAttitudeData = async () => {
  if (!filters.value.classroom_id) return
  
  isLoading.value = true
  try {
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data

    const params = {
      classroom_id: filters.value.classroom_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year,
      type: 'ATTITUDE'
    }
    const asRes = await axios.get('/api/v1/admin/edu/assessments', { params })
    
    // Clear existing keys properly for reactive
    Object.keys(attitudeMap).forEach(key => delete attitudeMap[Number(key)])

    students.value.forEach(s => {
      attitudeMap[s.ID] = {}
      attitudeCategories.forEach(cat => {
        const existing = asRes.data.find((a: any) => a.student_id === s.ID && a.remarks === cat.id)
        attitudeMap[s.ID][cat.id] = existing ? existing.grade : 'B'
      })
    })
  } catch (err) {
    console.error('Gagal memuat data karakter:', err)
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
    const payload: any[] = []
    students.value.forEach(s => {
      attitudeCategories.forEach(cat => {
        payload.push({
          student_id: s.ID,
          classroom_id: filters.value.classroom_id,
          term: filters.value.term,
          academic_year: filters.value.academic_year,
          type: 'ATTITUDE',
          grade: attitudeMap[s.ID][cat.id],
          remarks: cat.id
        })
      })
    })

    await axios.post('/api/v1/admin/edu/assessments/bulk', payload)
    
    toastMessage.value = 'Nilai Karakter Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    
    await loadAttitudeData()
  } catch (err) {
    console.error('Gagal menyimpan nilai karakter:', err)
    toastMessage.value = 'Gagal Menyimpan Nilai.'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isSaving.value = false
  }
}

const updateGrade = (studentId: number, categoryId: string, grade: string) => {
  if (!attitudeMap[studentId]) {
    attitudeMap[studentId] = {}
  }
  attitudeMap[studentId][categoryId] = grade
}

const getGradeColor = (grade: string) => {
  switch (grade) {
    case 'A': return 'text-emerald-500 bg-emerald-500/10 border-emerald-500/20'
    case 'B': return 'text-blue-500 bg-blue-500/10 border-blue-500/20'
    case 'C': return 'text-amber-500 bg-amber-500/10 border-amber-500/20'
    case 'D': return 'text-rose-500 bg-rose-500/10 border-rose-500/20'
    default: return 'text-sigma-muted'
  }
}

watch([() => filters.value.classroom_id, () => filters.value.term], loadAttitudeData)
onMounted(fetchData)

// Aggregation for Rekap
const rekapData = computed(() => {
  if (students.value.length === 0) return []
  
  return students.value.map(s => {
    const map = attitudeMap[s.ID] || {}
    
    // Count frequencies only for active categories
    const counts: Record<string, number> = { 'A': 0, 'B': 0, 'C': 0, 'D': 0 }
    let totalCount = 0
    
    attitudeCategories.forEach(cat => {
      const g = map[cat.id]
      if (g && counts[g] !== undefined) {
        counts[g]++
        totalCount++
      }
    })
    
    // Default if no data
    if (totalCount === 0) return { student: s, grades: map, summary: '-' }
    
    // Determine summary (Most frequent grade)
    let summary = 'B'
    let maxCount = -1
    
    // Priority: A > B > C > D (higher grade wins on ties)
    const priority = ['A', 'B', 'C', 'D']
    priority.forEach(g => {
      if (counts[g] > maxCount) {
        maxCount = counts[g]
        summary = g
      }
    })

    return {
      student: s,
      grades: map,
      summary
    }
  })
})
</script>

<template>
  <div class="attitude-page">
    <div class="p-8 space-y-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="h-10 w-1 bg-rose-500 rounded-full shadow-[0_0_15px_rgba(244,63,94,0.5)]"></div>
          <div>
            <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Penilaian</h2>
            <h1 class="text-2xl font-black italic uppercase tracking-tight">Nilai <span class="text-rose-500">Karakter & Adab</span></h1>
          </div>
        </div>

      <div class="flex items-center gap-3">
        <!-- Tabs -->
        <div class="bg-sigma-surface-alt border border-sigma-border p-1 rounded-2xl flex items-center gap-1 mr-4">
          <button @click="activeTab = 'input'"
                  class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                  :class="activeTab === 'input' ? 'bg-rose-600 text-white shadow-lg shadow-rose-500/20' : 'text-sigma-muted hover:text-sigma-text'">
            <LayoutGrid class="w-3.5 h-3.5" /> Input
          </button>
          <button @click="activeTab = 'rekap'"
                  class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                  :class="activeTab === 'rekap' ? 'bg-rose-600 text-white shadow-lg shadow-rose-500/20' : 'text-sigma-muted hover:text-sigma-text'">
            <List class="w-3.5 h-3.5" /> Rekap
          </button>
        </div>

        <button v-if="activeTab === 'input'" @click="handleSave" :disabled="isSaving || students.length === 0"
                class="px-8 py-3 bg-rose-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-rose-700 transition-all shadow-lg shadow-rose-500/20 disabled:opacity-50 flex items-center gap-3 active:scale-95">
          <Save v-if="!isSaving" class="w-4 h-4" />
          <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
          Simpan Nilai
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center gap-8">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 flex-1">
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tahun Ajaran</label>
          <select v-model="filters.academic_year" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none appearance-none">
            <option>2024/2025</option><option>2023/2024</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Semester</label>
          <select v-model="filters.term" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none appearance-none">
            <option>Ganjil</option><option>Genap</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
          <select v-model="filters.classroom_id" class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none appearance-none">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Main Table -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm min-h-[60vh]">
      <div v-if="isLoading" class="flex flex-col items-center justify-center p-40 space-y-4">
        <div class="w-10 h-10 border-4 border-rose-500/20 border-t-rose-500 rounded-full animate-spin"></div>
        <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Memuat Data Karakter...</p>
      </div>

      <!-- Input Mode -->
      <table v-else-if="activeTab === 'input' && students.length > 0" class="w-full border-collapse">
        <thead>
          <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
            <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Santri</th>
            <th v-for="cat in attitudeCategories" :key="cat.id" class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">
              <div class="flex flex-col items-center gap-1">
                <component :is="cat.icon" class="w-4 h-4 mb-1" :class="cat.color" />
                {{ cat.label }}
              </div>
            </th>
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
            <td v-for="cat in attitudeCategories" :key="cat.id" class="p-6 text-center">
              <div class="flex items-center justify-center gap-1 bg-sigma-surface-alt/50 p-1.5 rounded-2xl border border-sigma-border/50 w-fit mx-auto">
                <button v-for="grade in ['A', 'B', 'C', 'D']" 
                        :key="grade"
                        @click="updateGrade(s.ID, cat.id, grade)"
                        class="w-9 h-9 rounded-xl text-xs font-black transition-all border flex items-center justify-center"
                        :class="attitudeMap[s.ID][cat.id] === grade 
                          ? (grade === 'A' ? 'bg-emerald-500 text-white border-emerald-500 shadow-lg shadow-emerald-500/20' : 
                             grade === 'B' ? 'bg-blue-600 text-white border-blue-600 shadow-lg shadow-blue-500/20' : 
                             grade === 'C' ? 'bg-amber-500 text-white border-amber-500 shadow-lg shadow-amber-500/20' : 
                             'bg-rose-500 text-white border-rose-500 shadow-lg shadow-rose-500/20')
                          : 'bg-transparent text-sigma-muted border-transparent hover:border-sigma-border hover:text-sigma-text'">
                  {{ grade }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Rekap Mode -->
      <table v-else-if="activeTab === 'rekap' && students.length > 0" class="w-full border-collapse">
        <thead>
          <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
            <th class="p-6 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Santri</th>
            <th v-for="cat in attitudeCategories" :key="cat.id" class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">
              {{ cat.label }}
            </th>
            <th class="p-6 text-center text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted bg-rose-500/5">Kesimpulan</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rekapData" :key="row.student.ID" class="border-b border-sigma-border last:border-0 hover:bg-sigma-surface-alt/10 transition-all">
            <td class="p-6">
              <span class="text-sm font-black uppercase tracking-tighter italic">{{ row.student.name }}</span>
            </td>
            <td v-for="cat in attitudeCategories" :key="cat.id" class="p-6 text-center">
              <span class="inline-flex items-center justify-center w-8 h-8 rounded-lg text-xs font-black border" 
                    :class="getGradeColor(row.grades[cat.id])">
                {{ row.grades[cat.id] }}
              </span>
            </td>
            <td class="p-6 text-center bg-rose-500/5">
              <span class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-black border"
                    :class="getGradeColor(row.summary)">
                <Star class="w-3.5 h-3.5" /> {{ row.summary }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="flex flex-col items-center justify-center p-40 space-y-4 opacity-30">
        <Heart class="w-16 h-16" />
        <p class="font-black uppercase tracking-widest text-[10px]">Pilih kelas untuk memulai penilaian karakter</p>
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
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }
</style>

<style scoped>
.toast-enter-active, .toast-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px) scale(0.9);
}
</style>
