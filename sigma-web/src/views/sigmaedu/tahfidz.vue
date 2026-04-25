<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { 
  BookOpen, Save, 
  CheckCircle2, XCircle,
  Plus, 
  Search,
  LayoutGrid,
  Pencil, Trash2
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const records = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const showBulkModal = ref(false)
const searchQuery = ref('')
const showPersonalModal = ref(false)
const selectedStudent = ref<any>(null)
const personalEntries = ref<Record<string, { surah: string, start: number | null, end: number | null, grade: string, remarks: string }>>({})
const studentHistory = ref<Record<number, any[]>>({})
const loadingHistory = ref<Record<number, boolean>>({})

const filters = ref({
  classroom_id: null as number | null,
  date: new Date().toISOString().split('T')[0]
})

const bulkForm = ref({
  type: 'SABAQ',
  classroom_id: null as number | null,
  date: new Date().toISOString().split('T')[0],
  entries: {} as Record<number, { surah: string, start: number | null, end: number | null, grade: string, remarks: string }>
})

const TAHFIDZ_TYPES = [
  { label: 'Sabaq (Hafalan Baru)', value: 'SABAQ' },
  { label: 'Sabqi (Murojaah Baru)', value: 'SABQI' },
  { label: 'Manzil (Murojaah Lama)', value: 'MANZIL' },
]

const quranSurahs = [
  "Al-Fatihah", "Al-Baqarah", "Ali 'Imran", "An-Nisa'", "Al-Ma'idah", "Al-An'am", "Al-A'raf", "Al-Anfal", "At-Tawbah", "Yunus", "Hud", "Yusuf", "Ar-Ra'd", "Ibrahim", "Al-Hijr", "An-Nahl", "Al-Isra'", "Al-Kahf", "Maryam", "Ta-Ha", "Al-Anbiya'", "Al-Hajj", "Al-Mu'minun", "An-Nur", "Al-Furqan", "Ash-Shu'ara'", "An-Naml", "Al-Qasas", "Al-'Ankabut", "Ar-Rum", "Luqman", "As-Sajdah", "Al-Ahzab", "Saba'", "Fatir", "Ya-Sin", "As-Saffat", "Sad", "Az-Zumar", "Ghafir", "Fussilat", "Ash-Shura", "Az-Zukhruf", "Ad-Dukhan", "Al-Jathiyah", "Al-Ahqaf", "Muhammad", "Al-Fath", "Al-Hujurat", "Qaf", "Ad-Dhariyat", "At-Tur", "An-Najm", "Al-Qamar", "Ar-Rahman", "Al-Waqi'ah", "Al-Hadid", "Al-Mujadilah", "Al-Hashr", "Al-Mumtahanah", "As-Saff", "Al-Jumu'ah", "Al-Munafiqun", "At-Taghabun", "At-Talaq", "At-Tahrim", "Al-Mulk", "Al-Qalam", "Al-Haqqah", "Al-Ma'arij", "Nuh", "Al-Jinn", "Al-Muzzammil", "Al-Muddaththir", "Al-Qiyamah", "Al-Insan", "Al-Mursalat", "An-Naba'", "An-Nazi'at", "'Abasa", "At-Takwir", "Al-Infitar", "Al-Mutaffifin", "Al-Inshiqaq", "Al-Buruj", "At-Tariq", "Al-A'la", "Al-Ghashiyah", "Al-Fajr", "Al-Balad", "Ash-Shams", "Al-Layl", "Ad-Duha", "Ash-Sharh", "At-Tin", "Al-'Alaq", "Al-Qadr", "Al-Bayyinah", "Az-Zalzalah", "Al-'Adiyat", "Al-Qari'ah", "At-Takathur", "Al-'Asr", "Al-Humazah", "Al-Fil", "Quraysh", "Al-Ma'un", "Al-Kawthar", "Al-Kafirun", "An-Nasr", "Al-Masad", "Al-Ikhlas", "Al-Falaq", "An-Nas"
]

const fetchData = async () => {
  try {
    const clRes = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = clRes.data
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
    await loadTahfidzData()
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  }
}

const loadTahfidzData = async () => {
  isLoading.value = true
  try {
    let stRes;
    if (filters.value.classroom_id) {
       stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    } else {
       stRes = await axios.get('/api/v1/admin/base/students')
    }
    
    const tfRes = await axios.get('/api/v1/admin/edu/tahfidz', { params: { classroom_id: filters.value.classroom_id, date: filters.value.date } })
    
    students.value = stRes.data
    records.value = tfRes.data

    const newMatrix: Record<number, Record<string, any>> = {}
    students.value.forEach(s => {
      newMatrix[s.ID] = {}
      TAHFIDZ_TYPES.forEach(t => {
        const found = records.value.find(r => r.student_id === s.ID && r.type === t.value)
        newMatrix[s.ID][t.value] = found || null
      })
    })
    tahfidzMatrix.value = newMatrix
  } catch (err) {
    console.error('Gagal memuat data tahfidz:', err)
  } finally {
    isLoading.value = false
  }
}

const tahfidzMatrix = ref<Record<number, Record<string, any>>>({})

const openBulkModal = async () => {
  bulkForm.value.date = filters.value.date
  showBulkModal.value = true
  
  // Set classroom_id after showBulkModal so the watcher picks it up
  bulkForm.value.classroom_id = filters.value.classroom_id
  
  // Also fetch students directly for initial open
  if (filters.value.classroom_id) {
    try {
      const res = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
      bulkStudents.value = res.data
      
      const newEntries: Record<number, any> = {}
      bulkStudents.value.forEach((s: any) => {
        const existing = tahfidzMatrix.value[s.ID]?.[bulkForm.value.type]
        newEntries[s.ID] = existing ? {
          surah: existing.surah_name,
          start: existing.verse_start,
          end: existing.verse_end,
          grade: existing.grade,
          remarks: existing.remarks
        } : {
          surah: '',
          start: null,
          end: null,
          grade: 'A',
          remarks: ''
        }
      })
      bulkForm.value.entries = newEntries
    } catch (err) {
      console.error('Gagal fetch santri:', err)
    }
  }
}

const bulkStudents = ref<any[]>([])

// Single watcher for class changes inside the bulk modal
watch(() => bulkForm.value.classroom_id, async (newClassId) => {
  if (!showBulkModal.value || !newClassId) return
  try {
    const res = await axios.get(`/api/v1/admin/edu/classrooms/${newClassId}/students`)
    bulkStudents.value = res.data
    
    const newEntries: Record<number, any> = {}
    bulkStudents.value.forEach((s: any) => {
      newEntries[s.ID] = {
        surah: '',
        start: null,
        end: null,
        grade: 'A',
        remarks: ''
      }
    })
    bulkForm.value.entries = newEntries
  } catch (err) {
    console.error('Gagal fetch santri untuk bulk:', err)
  }
})

const openPersonalModal = (student: any) => {
  selectedStudent.value = student
  const m = tahfidzMatrix.value[student.ID] || {}
  
  const entries: Record<string, any> = {}
  TAHFIDZ_TYPES.forEach(t => {
    const existing = m[t.value]
    entries[t.value] = existing ? {
      surah: existing.surah_name,
      start: existing.verse_start,
      end: existing.verse_end,
      grade: existing.grade,
      remarks: existing.remarks
    } : {
      surah: '',
      start: null,
      end: null,
      grade: 'A',
      remarks: ''
    }
  })
  
  personalEntries.value = entries
  showPersonalModal.value = true
}
const fetchStudentHistory = async (studentID: number) => {
  if (loadingHistory.value[studentID]) return
  loadingHistory.value[studentID] = true
  try {
    const res = await axios.get(`/api/v1/admin/edu/tahfidz/history/${studentID}`)
    studentHistory.value[studentID] = res.data
  } catch (err) {
    console.error('Gagal ambil history:', err)
  } finally {
    loadingHistory.value[studentID] = false
  }
}

const savePersonalTahfidz = async () => {
  if (!selectedStudent.value) return
  
  isSaving.value = true
  try {
    const payload = TAHFIDZ_TYPES.map(t => {
      const entry = personalEntries.value[t.value]
      return {
        student_id: selectedStudent.value.ID,
        date: new Date(filters.value.date).toISOString(),
        type: t.value,
        surah_name: entry.surah,
        verse_start: entry.start,
        verse_end: entry.end,
        grade: entry.grade,
        remarks: entry.remarks
      }
    }).filter(e => e.surah_name !== '')

    await axios.post('/api/v1/admin/edu/tahfidz/bulk', payload)
    await loadTahfidzData()
    await fetchStudentHistory(selectedStudent.value.ID)
    showPersonalModal.value = false
    toastMessage.value = `Hafalan ${selectedStudent.value.name} Disimpan!`
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan personal tahfidz:', err)
  } finally {
    isSaving.value = false
  }
}

const saveBulkTahfidz = async () => {
  if (!bulkForm.value.classroom_id) return
  
  isSaving.value = true
  try {
    const payload = bulkStudents.value.map(s => {
      const entry = bulkForm.value.entries[s.ID]
      return {
        student_id: s.ID,
        date: new Date(bulkForm.value.date).toISOString(),
        type: bulkForm.value.type,
        surah_name: entry.surah,
        verse_start: entry.start,
        verse_end: entry.end,
        grade: entry.grade,
        remarks: entry.remarks
      }
    }).filter(e => e.surah_name !== '')

    await axios.post('/api/v1/admin/edu/tahfidz/bulk', payload)
    
    // Sync main page filter with the class we just saved
    filters.value.classroom_id = bulkForm.value.classroom_id
    filters.value.date = bulkForm.value.date
    
    await loadTahfidzData()
    
    // Refresh history for students in the payload
    payload.forEach(item => {
      fetchStudentHistory(item.student_id)
    })

    showBulkModal.value = false
    toastMessage.value = 'Hafalan Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal simpan tahfidz:', err)
  } finally {
    isSaving.value = false
  }
}


const deleteRecord = async (recordID: number, studentID: number) => {
  if (!confirm('Yakin ingin menghapus data hafalan ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/tahfidz/${recordID}`)
    await loadTahfidzData()
    await fetchStudentHistory(studentID)
    toastMessage.value = 'Data hafalan berhasil dihapus!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } catch (err) {
    console.error('Gagal hapus record:', err)
  }
}

watch(searchQuery, (newVal) => {
  if (newVal.length > 2) {
    const matched = students.value.filter(s => s.name.toLowerCase().includes(newVal.toLowerCase()))
    if (matched.length <= 3) {
      matched.forEach(s => {
        if (!studentHistory.value[s.ID]) fetchStudentHistory(s.ID)
      })
    }
  }
})

const rekapData = computed(() => {
  const filtered = searchQuery.value 
    ? students.value.filter(s => s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.nis.includes(searchQuery.value))
    : students.value

  return filtered.map(s => {
    const m = tahfidzMatrix.value[s.ID] || {}
    let totalAyat = 0
    TAHFIDZ_TYPES.forEach(t => {
      const rec = m[t.value]
      if (rec && rec.verse_start && rec.verse_end) {
        totalAyat += (rec.verse_end - rec.verse_start + 1)
      }
    })
    
    return {
      student: s,
      records: m,
      totalAyat
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

watch([() => filters.value.classroom_id, () => filters.value.date], loadTahfidzData)

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-[#F7FCF9] dark:bg-sigma-bg transition-colors duration-500">
    <!-- Header Section -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
      <div class="flex items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-gradient-to-br from-emerald-500 to-emerald-700 flex items-center justify-center text-white shadow-2xl shadow-emerald-500/40 transform -rotate-6 hover:rotate-0 transition-transform duration-500">
          <BookOpen class="w-8 h-8" />
        </div>
        <div>
          <h2 class="text-[10px] font-black uppercase tracking-[0.4em] text-sigma-muted mb-1 opacity-60">Akademik • Tahfidz Quran</h2>
          <h1 class="text-3xl font-black italic uppercase tracking-tighter leading-none">
            Tahfidz <span class="text-emerald-500">Journal</span>
          </h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <select v-model="filters.classroom_id" class="px-5 py-3 bg-white border border-sigma-border rounded-2xl text-[11px] font-bold outline-none focus:border-emerald-500 shadow-sm cursor-pointer">
          <option :value="null">Semua Kelas</option>
          <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
        </select>
        <input v-model="filters.date" type="date" class="px-5 py-3 bg-white border border-sigma-border rounded-2xl text-[11px] font-bold outline-none focus:border-emerald-500 shadow-sm cursor-pointer" />
        <button @click="() => openBulkModal()"
                class="px-8 py-3 bg-emerald-500 text-white rounded-2xl font-black uppercase tracking-wider text-[9px] hover:shadow-xl hover:shadow-emerald-500/40 transition-all active:scale-95 flex items-center gap-2">
          <Plus class="w-3.5 h-3.5" /> Input Kolektif
        </button>
      </div>
    </div>

    <!-- Content Area -->
    <div class="flex flex-col relative space-y-8">
      
      <!-- Global Search Bar -->
      <div class="relative w-full">
        <Search class="absolute left-6 top-1/2 -translate-y-1/2 w-6 h-6 text-emerald-500" />
        <input v-model="searchQuery" type="text" placeholder="Ketik Nama Santri untuk melihat Jurnal..."
               class="w-full pl-16 pr-8 py-6 bg-white border border-sigma-border rounded-[2.5rem] text-base font-bold outline-none focus:border-emerald-500 shadow-xl transition-all placeholder:opacity-50" />
      </div>

      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-40 space-y-6">
        <div class="w-16 h-16 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin"></div>
        <p class="text-[11px] font-black uppercase tracking-[0.3em] text-sigma-muted">Synchronizing Journal Data...</p>
      </div>

      <!-- JOURNAL LIST (FORMERLY SUMMARY) -->
      <div v-if="students.length > 0" class="flex-1 grid grid-cols-1 lg:grid-cols-2 gap-5 overflow-y-auto custom-scrollbar items-start">
        <div v-for="r in rekapData" :key="r.student.ID" 
             class="bg-white border border-sigma-border rounded-[2rem] p-6 space-y-5 hover:shadow-xl transition-all group relative overflow-hidden">
           
           <!-- Card Background Accent -->
           <div class="absolute top-0 right-0 w-80 h-80 bg-emerald-500/5 rounded-full -mr-40 -mt-40 blur-3xl group-hover:bg-emerald-500/10 transition-colors"></div>

           <div class="flex items-center justify-between gap-4 relative z-10">
              <div class="flex items-center gap-4">
                 <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-sigma-surface-alt to-sigma-border border border-sigma-border flex items-center justify-center text-base font-black group-hover:scale-110 transition-transform shadow-md">
                   {{ r.student.name.charAt(0) }}
                 </div>
                 <div>
                   <div class="flex items-center gap-2 mb-1">
                      <span class="px-3 py-1 bg-emerald-500/10 text-emerald-600 rounded-full text-[9px] font-black uppercase tracking-wider">{{ r.student.nis }}</span>
                      <span v-if="r.student.classroom" class="px-3 py-1 bg-blue-500/10 text-blue-600 rounded-full text-[9px] font-black uppercase tracking-wider">{{ r.student.classroom.name }}</span>
                   </div>
                   <h3 class="text-base font-black uppercase italic tracking-tight leading-none">{{ r.student.name }}</h3>
                 </div>
              </div>

              <button @click="openPersonalModal(r.student)" 
                      class="px-5 py-2.5 bg-white border border-sigma-border rounded-xl text-emerald-500 font-black uppercase tracking-wider text-[9px] hover:bg-emerald-500 hover:text-white transition-all shadow-sm flex items-center gap-2 shrink-0">
                <Plus class="w-3.5 h-3.5" /> Input
              </button>
           </div>

           <!-- Records Table (Journal Mode) -->
           <div class="relative z-10">
              <div class="flex items-center justify-between mb-3 px-2">
                 <h4 class="text-[9px] font-black uppercase tracking-[0.3em] text-sigma-muted">Riwayat Hafalan</h4>
                 <div v-if="searchQuery" class="text-[9px] font-black uppercase tracking-wider text-emerald-600">
                    "{{ searchQuery }}"
                 </div>
              </div>

              <div class="overflow-hidden rounded-2xl border border-sigma-border shadow-sm">
                 <table class="w-full border-collapse text-xs">
                    <thead>
                       <tr class="bg-sigma-surface-alt/50">
                          <th class="px-3 py-2.5 text-left text-[9px] font-black uppercase tracking-wider text-sigma-muted">Tanggal</th>
                          <th class="px-3 py-2.5 text-left text-[9px] font-black uppercase tracking-wider text-sigma-muted">Surah</th>
                          <th class="px-3 py-2.5 text-center text-[9px] font-black uppercase tracking-wider text-sigma-muted">Ayat</th>
                          <th class="px-3 py-2.5 text-left text-[9px] font-black uppercase tracking-wider text-sigma-muted">Jenis</th>
                          <th class="px-3 py-2.5 text-center text-[9px] font-black uppercase tracking-wider text-sigma-muted">Grade</th>
                          <th class="px-3 py-2.5 text-center text-[9px] font-black uppercase tracking-wider text-sigma-muted">Aksi</th>
                       </tr>
                    </thead>
                    <tbody class="divide-y divide-sigma-border/40">
                        <!-- Show Today's Records first if available -->
                        <template v-for="(rec, type) in r.records" :key="type">
                          <tr v-if="rec && rec.surah_name" class="bg-emerald-50/30">
                           <td class="px-3 py-2">
                              <div class="flex items-center gap-1.5">
                                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                                <span class="text-[10px] font-bold text-emerald-700">Hari Ini</span>
                              </div>
                           </td>
                           <td class="px-3 py-2">
                              <div class="flex items-center gap-2">
                                 <BookOpen class="w-3 h-3 text-emerald-500" />
                                 <span class="text-[11px] font-bold tracking-tight">{{ rec.surah_name }}</span>
                              </div>
                           </td>
                           <td class="px-3 py-2 text-center">
                              <span class="px-2 py-0.5 bg-white border border-sigma-border rounded text-[10px] font-bold">{{ rec.verse_start }}-{{ rec.verse_end }}</span>
                           </td>
                           <td class="px-3 py-2">
                              <span class="text-[9px] font-bold uppercase tracking-wider text-emerald-600">
                                 {{ rec.type === 'SABAQ' ? 'Sabaq' : 'Murojaah' }}
                              </span>
                           </td>
                           <td class="px-3 py-2 text-center">
                              <div class="inline-flex w-7 h-7 rounded-lg items-center justify-center text-[10px] font-black" :class="getGradeColor(rec.grade)">
                                 {{ rec.grade }}
                              </div>
                           </td>
                           <td class="px-3 py-2 text-center">
                              <button @click="openPersonalModal(r.student)" class="w-6 h-6 rounded-md bg-blue-500/10 text-blue-500 inline-flex items-center justify-center hover:bg-blue-500 hover:text-white transition-all" title="Edit">
                                <Pencil class="w-3 h-3" />
                              </button>
                           </td>
                        </tr>
                        </template>

                        <template v-if="studentHistory[r.student.ID]">
                           <tr v-for="log in studentHistory[r.student.ID].filter(l => new Date(l.date).toDateString() !== new Date().toDateString())" :key="log.ID" class="hover:bg-emerald-500/[0.02] transition-colors">
                             <td class="px-3 py-2">
                                <span class="text-[10px] font-bold text-sigma-text">{{ new Date(log.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) }}</span>
                             </td>
                             <td class="px-3 py-2">
                                <div class="flex items-center gap-2">
                                   <BookOpen class="w-3 h-3 text-emerald-500" />
                                   <span class="text-[11px] font-bold tracking-tight">{{ log.surah_name }}</span>
                                </div>
                             </td>
                             <td class="px-3 py-2 text-center">
                                <span class="px-2 py-0.5 bg-white border border-sigma-border rounded text-[10px] font-bold">{{ log.verse_start }}-{{ log.verse_end }}</span>
                             </td>
                             <td class="px-3 py-2">
                                <span class="text-[9px] font-bold uppercase tracking-wider"
                                      :class="log.type === 'SABAQ' ? 'text-emerald-600' : 'text-blue-600'">
                                   {{ log.type === 'SABAQ' ? 'Sabaq' : 'Murojaah' }}
                                </span>
                             </td>
                             <td class="px-3 py-2 text-center">
                                <div class="inline-flex w-7 h-7 rounded-lg items-center justify-center text-[10px] font-black" :class="getGradeColor(log.grade)">
                                   {{ log.grade }}
                                </div>
                             </td>
                             <td class="px-3 py-2 text-center">
                                <div class="flex items-center justify-center gap-1">
                                  <button @click="openPersonalModal(r.student)" class="w-6 h-6 rounded-md bg-blue-500/10 text-blue-500 inline-flex items-center justify-center hover:bg-blue-500 hover:text-white transition-all" title="Edit">
                                    <Pencil class="w-3 h-3" />
                                  </button>
                                  <button @click="deleteRecord(log.ID, r.student.ID)" class="w-6 h-6 rounded-md bg-rose-500/10 text-rose-500 inline-flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all" title="Hapus">
                                    <Trash2 class="w-3 h-3" />
                                  </button>
                                </div>
                             </td>
                          </tr>
                        </template>
                        <tr v-else-if="loadingHistory[r.student.ID]">
                           <td colspan="6" class="p-8 text-center">
                              <div class="inline-block w-6 h-6 border-3 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin mb-2"></div>
                              <p class="text-[9px] font-black uppercase tracking-wider text-sigma-muted">Loading...</p>
                           </td>
                        </tr>
                         <tr v-else-if="Object.keys(r.records).length === 0">
                            <td colspan="6" class="p-6 text-center text-sigma-muted opacity-60">
                               <div class="flex flex-col items-center gap-1">
                                 <Search class="w-4 h-4 opacity-20" />
                                 <p class="text-[9px] font-black uppercase tracking-wider">Cari nama untuk riwayat</p>
                              </div>
                           </td>
                        </tr>
                    </tbody>
                 </table>
              </div>
           </div>
        </div>
      </div>

    </div>

    <!-- COLLECTIVE INPUT MODAL -->
    <Transition name="modal">
      <div v-if="showBulkModal" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-sigma-bg/60 backdrop-blur-2xl" @click="showBulkModal = false"></div>
        <div class="relative w-full max-w-6xl bg-sigma-surface border border-sigma-border rounded-[4rem] shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-500">
          
          <div class="p-10 border-b border-sigma-border flex items-center justify-between bg-sigma-surface-alt/30">
            <div class="flex items-center gap-6">
              <div class="w-14 h-14 rounded-[1.5rem] bg-emerald-600 flex items-center justify-center text-white shadow-xl shadow-emerald-600/30">
                <Plus class="w-7 h-7" />
              </div>
              <div class="flex items-center gap-8">
                <div>
                  <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Setoran <span class="text-emerald-600">Batch</span></h3>
                  <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">Input Kolektif Per Kelas</p>
                </div>
                
                <div class="h-10 w-px bg-sigma-border mx-2"></div>

                <!-- Selection Inside Modal -->
                <div class="flex flex-wrap items-center gap-x-14 gap-y-6">
                   <div class="space-y-3">
                      <span class="text-[10px] font-black uppercase tracking-[0.25em] text-sigma-muted ml-1 opacity-60 block">Pilih Kelas</span>
                      <select v-model="bulkForm.classroom_id" class="px-8 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-[11px] font-black outline-none focus:border-emerald-500 shadow-md transition-all min-w-[160px] cursor-pointer hover:bg-white">
                         <option :value="null">Pilih...</option>
                         <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
                      </select>
                   </div>
                   <div class="space-y-3">
                      <span class="text-[10px] font-black uppercase tracking-[0.25em] text-sigma-muted ml-1 opacity-60 block">Tanggal</span>
                      <input v-model="bulkForm.date" type="date" class="px-8 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-[11px] font-black outline-none focus:border-emerald-500 shadow-md transition-all cursor-pointer hover:bg-white" />
                   </div>
                   <div class="space-y-3">
                      <span class="text-[10px] font-black uppercase tracking-[0.25em] text-sigma-muted ml-1 opacity-60 block">Tipe Setoran</span>
                      <select v-model="bulkForm.type" class="px-8 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl text-[11px] font-black outline-none focus:border-emerald-500 shadow-md transition-all min-w-[200px] cursor-pointer hover:bg-white">
                         <option v-for="t in TAHFIDZ_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
                      </select>
                   </div>
                </div>
              </div>
            </div>
            <button @click="showBulkModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="p-0 h-[55vh] overflow-y-auto custom-scrollbar bg-sigma-surface">
            <div v-if="!bulkForm.classroom_id" class="h-full flex flex-col items-center justify-center space-y-4 opacity-40">
               <LayoutGrid class="w-16 h-16 text-emerald-500" />
               <p class="text-sm font-black uppercase tracking-widest">Silakan pilih kelas terlebih dahulu</p>
            </div>
            <table v-else class="w-full border-collapse">
              <thead class="sticky top-0 z-10 bg-sigma-surface-alt/90 backdrop-blur-md border-b border-sigma-border">
                <tr>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pl-12">Santri</th>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted">Hafalan (Surah & Ayat)</th>
                  <th class="p-6 text-center text-[9px] font-black uppercase tracking-widest text-sigma-muted w-32">Grade</th>
                  <th class="p-6 text-left text-[9px] font-black uppercase tracking-widest text-sigma-muted pr-12">Catatan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border/40">
                <tr v-for="s in bulkStudents" :key="s.ID" class="group hover:bg-emerald-500/[0.03] transition-all">
                  <td class="p-6 pl-12">
                    <span class="text-sm font-black uppercase tracking-tight">{{ s.name }}</span>
                  </td>
                  <td class="p-6">
                    <div v-if="bulkForm.entries[s.ID]" class="flex items-center gap-3">
                       <select v-model="bulkForm.entries[s.ID].surah" class="flex-1 px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50">
                          <option value="">Pilih Surah...</option>
                          <option v-for="surah in quranSurahs" :key="surah" :value="surah">{{ surah }}</option>
                       </select>
                       <div class="flex items-center gap-1">
                          <input v-model.number="bulkForm.entries[s.ID].start" type="number" placeholder="Mulai" class="w-16 px-2 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold text-center" />
                          <span class="text-sigma-muted">-</span>
                          <input v-model.number="bulkForm.entries[s.ID].end" type="number" placeholder="Akhir" class="w-16 px-2 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold text-center" />
                       </div>
                    </div>
                  </td>
                  <td class="p-6">
                    <div v-if="bulkForm.entries[s.ID]" class="flex items-center justify-center gap-1 bg-sigma-surface-alt/50 p-1 rounded-xl border border-sigma-border/30">
                       <button v-for="g in ['A', 'B', 'C', 'D']" :key="g"
                               @click="bulkForm.entries[s.ID].grade = g"
                               class="w-8 h-8 rounded-lg text-[10px] font-black transition-all"
                               :class="bulkForm.entries[s.ID].grade === g ? 'bg-emerald-600 text-white' : 'text-sigma-muted hover:bg-sigma-border'">
                         {{ g }}
                       </button>
                    </div>
                  </td>
                  <td class="p-6 pr-12">
                    <input v-if="bulkForm.entries[s.ID]" v-model="bulkForm.entries[s.ID].remarks" type="text"
                           class="w-full px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none"
                           placeholder="Feedback..." />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-end gap-4">
              <button @click="showBulkModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all">
                Cancel
              </button>
              <button @click="saveBulkTahfidz" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-emerald-600 to-emerald-700 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-emerald-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
                <Save v-if="!isSaving" class="w-5 h-5" />
                <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Commit Setoran
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
              <div class="w-14 h-14 rounded-[1.5rem] bg-emerald-600 flex items-center justify-center text-white shadow-xl shadow-emerald-600/30">
                <BookOpen class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black uppercase tracking-tighter italic leading-none mb-1">Personal <span class="text-emerald-600">Tahfidz</span></h3>
                <p class="text-[10px] font-black uppercase tracking-[0.3em] text-sigma-muted">{{ selectedStudent?.name }}</p>
              </div>
            </div>
            <button @click="showPersonalModal = false" class="w-12 h-12 bg-sigma-surface-alt border border-sigma-border rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white transition-all transform hover:rotate-90">
              <XCircle class="w-6 h-6" />
            </button>
          </div>

          <div class="p-10 space-y-6 max-h-[50vh] overflow-y-auto custom-scrollbar">
             <div v-for="t in TAHFIDZ_TYPES" :key="t.value" class="bg-sigma-surface-alt/50 p-6 rounded-3xl border border-sigma-border space-y-4">
                <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ t.label }}</span>
                <div class="grid grid-cols-1 gap-4">
                   <div class="flex items-center gap-3">
                      <select v-model="personalEntries[t.value].surah" class="flex-1 px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50">
                         <option value="">Pilih Surah...</option>
                         <option v-for="surah in quranSurahs" :key="surah" :value="surah">{{ surah }}</option>
                      </select>
                      <div class="flex items-center gap-1">
                         <input v-model.number="personalEntries[t.value].start" type="number" placeholder="Mulai" class="w-16 px-2 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold text-center" />
                         <span class="text-sigma-muted">-</span>
                         <input v-model.number="personalEntries[t.value].end" type="number" placeholder="Akhir" class="w-16 px-2 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold text-center" />
                      </div>
                   </div>
                   <div class="flex items-center justify-between gap-4">
                      <div class="flex items-center gap-1 bg-sigma-surface p-1 rounded-xl border border-sigma-border/30">
                         <button v-for="g in ['A', 'B', 'C', 'D']" :key="g"
                                  @click="personalEntries[t.value].grade = g"
                                  class="w-8 h-8 rounded-lg text-[10px] font-black transition-all"
                                  :class="personalEntries[t.value].grade === g ? 'bg-emerald-600 text-white' : 'text-sigma-muted hover:bg-sigma-border'">
                           {{ g }}
                         </button>
                      </div>
                      <input v-model="personalEntries[t.value].remarks" type="text"
                             class="flex-1 px-4 py-3 bg-sigma-surface border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50"
                             placeholder="Catatan..." />
                   </div>
                </div>
             </div>
          </div>

          <div class="p-10 border-t border-sigma-border bg-sigma-surface-alt/30 flex items-center justify-end gap-4">
              <button @click="showPersonalModal = false" class="px-10 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] font-black uppercase tracking-widest text-[11px] hover:bg-rose-500 hover:text-white transition-all">
                Cancel
              </button>
              <button @click="savePersonalTahfidz" :disabled="isSaving"
                      class="px-14 py-5 bg-gradient-to-r from-emerald-600 to-emerald-700 text-white rounded-[2rem] font-black uppercase tracking-[0.2em] text-[11px] hover:shadow-2xl hover:shadow-emerald-500/40 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-3">
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
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(16, 185, 129, 0.1); border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(16, 185, 129, 0.3); }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
.animate-in { animation-fill-mode: both; }
</style>
