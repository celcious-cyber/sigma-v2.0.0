<script setup lang="ts">
import { ref, onMounted, watch, computed, reactive } from 'vue'
import { 
  BookOpen, Save, Calendar as CalendarIcon, 
  CheckCircle2, XCircle, Hash,
  LayoutGrid, List, ChevronRight
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
  date: new Date().toISOString().split('T')[0]
})

// UI State
const tahfidzMap = reactive<Record<number, any>>({})

const tahfidzTypes = [
  { label: 'Sabaq (Hafalan Baru)', value: 'SABAQ', color: 'text-emerald-500 bg-emerald-500/10' },
  { label: 'Sabqi (Murojaah Baru)', value: 'SABQI', color: 'text-blue-500 bg-blue-500/10' },
  { label: 'Manzil (Murojaah Lama)', value: 'MANZIL', color: 'text-amber-500 bg-amber-500/10' },
]

const fetchData = async () => {
  try {
    const clRes = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = clRes.data
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
    // Call load immediately after fetching classrooms
    await loadTahfidzData()
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  }
}

const loadTahfidzData = async () => {
  if (!filters.value.classroom_id) return
  
  isLoading.value = true
  try {
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data

    const params = {
      classroom_id: filters.value.classroom_id,
      date: filters.value.date
    }
    const res = await axios.get('/api/v1/admin/edu/tahfidz', { params })
    
    // Clear existing keys properly for reactive
    Object.keys(tahfidzMap).forEach(key => delete tahfidzMap[Number(key)])

    students.value.forEach(s => {
      const existing = res.data.find((r: any) => r.student_id === s.ID)
      tahfidzMap[s.ID] = existing ? {
        surah_name: existing.surah_name,
        verse_start: existing.verse_start,
        verse_end: existing.verse_end,
        type: existing.type,
        grade: existing.grade,
        remarks: existing.remarks
      } : {
        surah_name: '',
        verse_start: null,
        verse_end: null,
        type: 'SABAQ',
        grade: 'A',
        remarks: ''
      }
    })
  } catch (err) {
    console.error('Gagal memuat data tahfidz:', err)
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
      surah_name: tahfidzMap[s.ID].surah_name,
      verse_start: tahfidzMap[s.ID].verse_start,
      verse_end: tahfidzMap[s.ID].verse_end,
      type: tahfidzMap[s.ID].type,
      grade: tahfidzMap[s.ID].grade,
      remarks: tahfidzMap[s.ID].remarks,
      date: new Date(filters.value.date).toISOString()
    }))

    await axios.post('/api/v1/admin/edu/tahfidz/bulk', payload)
    
    toastMessage.value = 'Setoran Hafalan Berhasil Disimpan!'
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
    
    await loadTahfidzData()
  } catch (err) {
    console.error('Gagal menyimpan tahfidz:', err)
    toastMessage.value = 'Gagal Menyimpan Data.'
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => showToast.value = false, 3000)
  } finally {
    isSaving.value = false
  }
}

// Quran Surahs List
const quranSurahs = [
  "Al-Fatihah", "Al-Baqarah", "Ali 'Imran", "An-Nisa'", "Al-Ma'idah", "Al-An'am", "Al-A'raf", "Al-Anfal", "At-Tawbah", "Yunus", "Hud", "Yusuf", "Ar-Ra'd", "Ibrahim", "Al-Hijr", "An-Nahl", "Al-Isra'", "Al-Kahf", "Maryam", "Ta-Ha", "Al-Anbiya'", "Al-Hajj", "Al-Mu'minun", "An-Nur", "Al-Furqan", "Ash-Shu'ara'", "An-Naml", "Al-Qasas", "Al-'Ankabut", "Ar-Rum", "Luqman", "As-Sajdah", "Al-Ahzab", "Saba'", "Fatir", "Ya-Sin", "As-Saffat", "Sad", "Az-Zumar", "Ghafir", "Fussilat", "Ash-Shura", "Az-Zukhruf", "Ad-Dukhan", "Al-Jathiyah", "Al-Ahqaf", "Muhammad", "Al-Fath", "Al-Hujurat", "Qaf", "Ad-Dhariyat", "At-Tur", "An-Najm", "Al-Qamar", "Ar-Rahman", "Al-Waqi'ah", "Al-Hadid", "Al-Mujadilah", "Al-Hashr", "Al-Mumtahanah", "As-Saff", "Al-Jumu'ah", "Al-Munafiqun", "At-Taghabun", "At-Talaq", "At-Tahrim", "Al-Mulk", "Al-Qalam", "Al-Haqqah", "Al-Ma'arij", "Nuh", "Al-Jinn", "Al-Muzzammil", "Al-Muddaththir", "Al-Qiyamah", "Al-Insan", "Al-Mursalat", "An-Naba'", "An-Nazi'at", "'Abasa", "At-Takwir", "Al-Infitar", "Al-Mutaffifin", "Al-Inshiqaq", "Al-Buruj", "At-Tariq", "Al-A'la", "Al-Ghashiyah", "Al-Fajr", "Al-Balad", "Ash-Shams", "Al-Layl", "Ad-Duha", "Ash-Sharh", "At-Tin", "Al-'Alaq", "Al-Qadr", "Al-Bayyinah", "Az-Zalzalah", "Al-'Adiyat", "Al-Qari'ah", "At-Takathur", "Al-'Asr", "Al-Humazah", "Al-Fil", "Quraysh", "Al-Ma'un", "Al-Kawthar", "Al-Kafirun", "An-Nasr", "Al-Masad", "Al-Ikhlas", "Al-Falaq", "An-Nas"
]

// Stats for Header
const stats = computed(() => {
  const data = Object.values(tahfidzMap)
  const setoran = data.filter((d: any) => d.surah_name !== '').length
  const totalVerses = data.reduce((acc: number, d: any) => {
    if (d.verse_start && d.verse_end) return acc + (d.verse_end - d.verse_start + 1)
    return acc
  }, 0)
  return { setoran, totalVerses }
})

watch([() => filters.value.classroom_id, () => filters.value.date], loadTahfidzData)
onMounted(fetchData)
</script>

<template>
  <div class="tahfidz-page">
    <div class="p-8 space-y-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="h-10 w-1 bg-emerald-500 rounded-full shadow-[0_0_15px_rgba(16,185,129,0.5)]"></div>
          <div>
            <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Tahfidz</h2>
            <h1 class="text-2xl font-black italic uppercase tracking-tight">Setoran <span class="text-emerald-500">Hafalan Quran</span></h1>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <!-- Tabs -->
          <div class="bg-sigma-surface-alt border border-sigma-border p-1 rounded-2xl flex items-center gap-1 mr-4">
            <button @click="activeTab = 'input'"
                    class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                    :class="activeTab === 'input' ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
              <LayoutGrid class="w-3.5 h-3.5" /> Input
            </button>
            <button @click="activeTab = 'rekap'"
                    class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2"
                    :class="activeTab === 'rekap' ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-500/20' : 'text-sigma-muted hover:text-sigma-text'">
              <List class="w-3.5 h-3.5" /> Rekap
            </button>
          </div>

          <button v-if="activeTab === 'input'" @click="handleSave" :disabled="isSaving || students.length === 0"
                  class="px-8 py-3 bg-emerald-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-500/20 disabled:opacity-50 flex items-center gap-3 active:scale-95">
            <Save v-if="!isSaving" class="w-4 h-4" />
            <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            Simpan Setoran
          </button>
        </div>
      </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
          <BookOpen class="w-6 h-6" />
        </div>
        <div>
          <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">Total Setoran</span>
          <span class="text-xl font-black">{{ stats.setoran }} <span class="text-[10px] text-sigma-muted">Santri</span></span>
        </div>
      </div>
      <div class="bg-sigma-surface border border-sigma-border p-6 rounded-[2rem] shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-500">
          <Hash class="w-6 h-6" />
        </div>
        <div>
          <span class="text-[8px] font-black uppercase tracking-widest text-sigma-muted block">Total Ayat</span>
          <span class="text-xl font-black">{{ stats.totalVerses }} <span class="text-[10px] text-sigma-muted">Ayat</span></span>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] p-6 shadow-sm flex flex-wrap items-center gap-8">
      <div class="flex items-center gap-6 flex-1">
        <div class="space-y-1 min-w-[200px]">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Tanggal Setoran</label>
          <div class="relative">
            <CalendarIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-emerald-500" />
            <input v-model="filters.date" type="date" 
                   class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50" />
          </div>
        </div>
        <div class="space-y-1 min-w-[200px]">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Pilih Kelas</label>
          <select v-model="filters.classroom_id" 
                  class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-emerald-500/50">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Main Table -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm">
      <div v-if="isLoading" class="flex flex-col items-center justify-center p-40 space-y-4">
        <div class="w-10 h-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin"></div>
        <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Memuat Data Tahfidz...</p>
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
                <!-- Jenis Setoran -->
                <select v-model="tahfidzMap[s.ID].type" 
                        class="px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-[10px] font-black uppercase tracking-widest outline-none focus:border-emerald-500/50 w-48">
                  <option v-for="t in tahfidzTypes" :key="t.value" :value="t.value">{{ t.label }}</option>
                </select>
                
                <ChevronRight class="w-4 h-4 text-sigma-muted" />

                <!-- Surah & Ayat -->
                <div class="flex items-center gap-2 flex-1">
                  <div class="relative flex-1">
                    <BookOpen class="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-sigma-muted" />
                    <select v-model="tahfidzMap[s.ID].surah_name"
                            class="w-full pl-9 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-emerald-500/50">
                      <option value="">Pilih Surah...</option>
                      <option v-for="surah in quranSurahs" :key="surah" :value="surah">{{ surah }}</option>
                    </select>
                  </div>
                  <div class="flex items-center gap-1">
                    <input v-model.number="tahfidzMap[s.ID].verse_start" type="number"
                           class="w-16 px-2 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold text-center outline-none focus:border-emerald-500/50"
                           placeholder="Awl" />
                    <span class="text-sigma-muted font-bold">-</span>
                    <input v-model.number="tahfidzMap[s.ID].verse_end" type="number"
                           class="w-16 px-2 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold text-center outline-none focus:border-emerald-500/50"
                           placeholder="Akh" />
                  </div>
                </div>
              </div>
            </td>
            <td class="p-6 text-center">
              <div class="flex items-center justify-center gap-1 bg-sigma-surface-alt/50 p-1 rounded-2xl border border-sigma-border/50">
                <button v-for="grade in ['A', 'B', 'C', 'D']" 
                        :key="grade"
                        @click="tahfidzMap[s.ID].grade = grade"
                        class="w-8 h-8 rounded-xl text-[10px] font-black transition-all border flex items-center justify-center"
                        :class="tahfidzMap[s.ID].grade === grade 
                          ? 'bg-emerald-600 text-white border-emerald-600 shadow-lg shadow-emerald-500/20' 
                          : 'bg-transparent text-sigma-muted border-transparent hover:border-sigma-border hover:text-sigma-text'">
                  {{ grade }}
                </button>
              </div>
            </td>
            <td class="p-6">
              <input v-model="tahfidzMap[s.ID].remarks" type="text"
                     class="w-full px-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none focus:border-emerald-500/50"
                     placeholder="Catatan..." />
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="flex flex-col items-center justify-center p-40 space-y-4 opacity-30">
        <BookOpen class="w-16 h-16" />
        <p class="font-black uppercase tracking-widest text-[10px]">Pilih kelas untuk memulai pencatatan tahfidz</p>
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
  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px) scale(0.9); }
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
</style>
