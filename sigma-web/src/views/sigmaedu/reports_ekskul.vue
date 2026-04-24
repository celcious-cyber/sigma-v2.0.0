<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { 
  Printer, Search, 
  Trophy
} from 'lucide-vue-next'
import axios from 'axios'

// State
const classrooms = ref<any[]>([])
const students = ref<any[]>([])
const isLoading = ref(false)
const selectedStudent = ref<any>(null)
const reportData = ref<any>(null)

const filters = ref({
  classroom_id: null as number | null,
  student_id: null as number | null,
  term: 'Ganjil',
  academic_year: '2024/2025'
})

const fetchData = async () => {
  try {
    const clRes = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = clRes.data
    if (classrooms.value.length > 0 && !filters.value.classroom_id) {
      filters.value.classroom_id = classrooms.value[0].ID
    }
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  }
}

const loadStudents = async () => {
  if (!filters.value.classroom_id) return
  try {
    const stRes = await axios.get(`/api/v1/admin/edu/classrooms/${filters.value.classroom_id}/students`)
    students.value = stRes.data
  } catch (err) {
    console.error('Gagal memuat santri:', err)
  }
}

const generateReport = async () => {
  if (!filters.value.student_id) return
  
  isLoading.value = true
  try {
    // Simulate fetching ekskul data
    reportData.value = [
      { name: 'Pramuka', grade: 'A', remarks: 'Sangat aktif dan memiliki jiwa kepemimpinan.' },
      { name: 'Seni Kaligrafi', grade: 'B', remarks: 'Menunjukkan bakat seni yang baik.' },
      { name: 'Futsal', grade: 'A', remarks: 'Sangat terampil dalam kerja sama tim.' },
    ]
    selectedStudent.value = students.value.find(s => s.ID === filters.value.student_id)
  } finally {
    isLoading.value = false
  }
}

watch(() => filters.value.classroom_id, loadStudents)
onMounted(fetchData)

const printReport = () => window.print()
</script>

<template>
  <div class="p-8 space-y-8 print:p-0">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6 print:hidden">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-rose-500 rounded-full shadow-[0_0_15px_rgba(244,63,94,0.5)]"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Laporan</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Rapot <span class="text-rose-500">Ekstrakurikuler</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button v-if="reportData" @click="printReport"
                class="px-6 py-3 bg-sigma-surface border border-sigma-border text-sigma-text rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-surface-alt transition-all flex items-center gap-2">
          <Printer class="w-4 h-4" /> Cetak Rapor
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[2rem] p-6 shadow-sm flex flex-wrap items-center gap-6 print:hidden">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 flex-1">
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
          <select v-model="filters.classroom_id" 
                  class="w-full px-4 py-2 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Santri</label>
          <select v-model="filters.student_id" 
                  class="w-full px-4 py-2 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold outline-none">
            <option v-for="s in students" :key="s.ID" :value="s.ID">{{ s.name }}</option>
          </select>
        </div>
        <div class="flex items-end">
          <button @click="generateReport" :disabled="isLoading || !filters.student_id"
                  class="w-full px-4 py-2.5 bg-rose-600 text-white rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-rose-700 transition-all shadow-md flex items-center justify-center gap-2">
            <Search class="w-3.5 h-3.5" /> Generate Rapor
          </button>
        </div>
      </div>
    </div>

    <!-- Report View -->
    <div v-if="reportData" class="bg-white text-slate-900 rounded-[3rem] p-12 shadow-2xl border border-slate-200 print:shadow-none print:border-0 print:p-8 space-y-10">
      <!-- Header -->
      <div class="text-center space-y-2 border-b-2 border-rose-500 pb-8">
        <h1 class="text-3xl font-black uppercase tracking-tighter text-rose-900">Laporan Ekstrakurikuler</h1>
        <p class="text-xs font-bold text-slate-400 uppercase tracking-widest">Extracurricular Activity Report</p>
      </div>

      <!-- Info -->
      <div class="grid grid-cols-2 gap-8 text-sm">
        <div class="space-y-2">
          <p class="flex justify-between border-b border-slate-100 pb-2"><span class="text-slate-400 font-bold uppercase text-[10px]">Nama</span> <span class="font-black uppercase">{{ selectedStudent?.name }}</span></p>
          <p class="flex justify-between border-b border-slate-100 pb-2"><span class="text-slate-400 font-bold uppercase text-[10px]">NIS</span> <span class="font-black">{{ selectedStudent?.nis }}</span></p>
        </div>
        <div class="space-y-2">
          <p class="flex justify-between border-b border-slate-100 pb-2"><span class="text-slate-400 font-bold uppercase text-[10px]">Kelas</span> <span class="font-black uppercase">{{ classrooms.find(c => c.ID === filters.classroom_id)?.name }}</span></p>
          <p class="flex justify-between border-b border-slate-100 pb-2"><span class="text-slate-400 font-bold uppercase text-[10px]">Tahun Ajaran</span> <span class="font-black">{{ filters.academic_year }}</span></p>
        </div>
      </div>

      <!-- Table -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 text-rose-600">
          <Trophy class="w-5 h-5" />
          <h2 class="text-sm font-black uppercase tracking-widest">Kegiatan Ekstrakurikuler</h2>
        </div>
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-rose-50">
              <th class="p-4 text-left text-[10px] font-black uppercase border border-slate-200">Nama Kegiatan</th>
              <th class="p-4 text-center text-[10px] font-black uppercase border border-slate-200 w-24">Predikat</th>
              <th class="p-4 text-left text-[10px] font-black uppercase border border-slate-200">Keterangan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="e in reportData" :key="e.name">
              <td class="p-4 text-xs font-black border border-slate-200 uppercase">{{ e.name }}</td>
              <td class="p-4 text-center text-sm font-black border border-slate-200 text-rose-600">{{ e.grade }}</td>
              <td class="p-4 text-xs font-bold text-slate-500 border border-slate-200">{{ e.remarks }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Signatures -->
      <div class="pt-10 grid grid-cols-2 gap-20 text-center">
        <div class="space-y-20">
          <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Pembina Ekskul</p>
          <div class="border-b border-slate-200 w-40 mx-auto"></div>
        </div>
        <div class="space-y-20">
          <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Wali Kelas</p>
          <div class="border-b border-slate-200 w-40 mx-auto"></div>
        </div>
      </div>
    </div>
  </div>
</template>
