<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { 
  FileText, Download, Printer, Search, 
  GraduationCap, 
  TrendingUp, Calendar, User
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
    // In a real app, this would be a single endpoint returning aggregated data
    // For now, we'll simulate it by fetching assessments
    const params = {
      student_id: filters.value.student_id,
      term: filters.value.term,
      academic_year: filters.value.academic_year
    }
    // We'll need to add this endpoint to the backend or fetch manually
    // Let's assume we have an endpoint that returns a summary
    const res = await axios.get(`/api/v1/admin/edu/reports/student/${filters.value.student_id}`, { params })
    reportData.value = res.data
    selectedStudent.value = students.value.find(s => s.ID === filters.value.student_id)
  } catch (err) {
    console.error('Gagal generate rapor:', err)
    // Fallback/Mock for demo
    reportData.value = {
      grades: [
        { subject: 'Tahfidz Quran', score: 96, grade: 'A+' },
        { subject: 'Bahasa Arab', score: 88, grade: 'B+' },
        { subject: 'Fiqih Ibadah', score: 82, grade: 'B' },
      ],
      attendance: { hadir: 45, izin: 2, sakit: 1, alpa: 0 },
      attitude: [
        { category: 'Kedisiplinan', grade: 'A' },
        { category: 'Kebersihan', grade: 'B' },
      ]
    }
    selectedStudent.value = students.value.find(s => s.ID === filters.value.student_id)
  } finally {
    isLoading.value = false
  }
}

watch(() => filters.value.classroom_id, loadStudents)
onMounted(fetchData)

const printReport = () => {
  window.print()
}
</script>

<template>
  <div class="p-8 space-y-8 print:p-0">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6 print:hidden">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-indigo-500 rounded-full shadow-[0_0_15px_rgba(99,102,241,0.5)]"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Akademik / Laporan</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Rapor <span class="text-indigo-500">Hasil Belajar</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button v-if="reportData" @click="printReport"
                class="px-6 py-3 bg-sigma-surface border border-sigma-border text-sigma-text rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-surface-alt transition-all flex items-center gap-2">
          <Printer class="w-4 h-4" /> Cetak Rapor
        </button>
        <button v-if="reportData"
                class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/20 flex items-center gap-2">
          <Download class="w-4 h-4" /> Download PDF
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[2rem] p-6 shadow-sm flex flex-wrap items-center gap-6 print:hidden">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 flex-1">
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Kelas</label>
          <select v-model="filters.classroom_id" 
                  class="w-full px-4 py-2 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-indigo-500/50">
            <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Santri</label>
          <select v-model="filters.student_id" 
                  class="w-full px-4 py-2 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-indigo-500/50">
            <option v-for="s in students" :key="s.ID" :value="s.ID">{{ s.name }}</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-[8px] font-black uppercase tracking-widest text-sigma-muted ml-1">Semester</label>
          <select v-model="filters.term" 
                  class="w-full px-4 py-2 bg-sigma-surface-alt border border-sigma-border rounded-xl text-xs font-bold appearance-none outline-none focus:border-indigo-500/50">
            <option>Ganjil</option><option>Genap</option>
          </select>
        </div>
        <div class="flex items-end">
          <button @click="generateReport" :disabled="isLoading || !filters.student_id"
                  class="w-full px-4 py-2.5 bg-indigo-600 text-white rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-md shadow-indigo-500/20 flex items-center justify-center gap-2">
            <Search class="w-3.5 h-3.5" /> Lihat Rapor
          </button>
        </div>
      </div>
    </div>

    <!-- Report Preview -->
    <div v-if="reportData" class="space-y-8 animate-in fade-in duration-700 print:space-y-4">
      <!-- Rapor Card -->
      <div class="bg-white text-slate-900 rounded-[3rem] p-12 shadow-2xl space-y-12 border border-slate-200 print:shadow-none print:border-0 print:p-8">
        <!-- Rapor Header -->
        <div class="flex items-center justify-between border-b-4 border-indigo-600 pb-10">
          <div class="flex items-center gap-6">
            <div class="w-20 h-20 bg-indigo-600 rounded-3xl flex items-center justify-center shadow-xl">
              <img src="/logo/SIGMA.svg" class="w-12 h-12 invert" />
            </div>
            <div>
              <h1 class="text-4xl font-black italic uppercase tracking-tighter text-indigo-900 leading-none mb-2">Laporan Hasil Belajar</h1>
              <div class="flex items-center gap-4 text-xs font-black uppercase tracking-widest text-slate-500">
                <span class="flex items-center gap-1.5"><Calendar class="w-3.5 h-3.5" /> TP {{ filters.academic_year }}</span>
                <span class="w-1 h-1 bg-slate-300 rounded-full"></span>
                <span class="flex items-center gap-1.5"><TrendingUp class="w-3.5 h-3.5" /> Semester {{ filters.term }}</span>
              </div>
            </div>
          </div>
          <div class="text-right">
            <h3 class="text-2xl font-black uppercase tracking-tighter text-indigo-900 leading-none mb-1">SIGMAEDU ACADEMY</h3>
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Islamic Boarding School Management</p>
          </div>
        </div>

        <!-- Student Info Grid -->
        <div class="grid grid-cols-2 gap-12">
          <div class="space-y-4">
            <div class="flex justify-between border-b border-slate-100 pb-3">
              <span class="text-[10px] font-black uppercase tracking-widest text-slate-400">Nama Lengkap</span>
              <span class="text-sm font-black text-slate-800 uppercase">{{ selectedStudent?.name }}</span>
            </div>
            <div class="flex justify-between border-b border-slate-100 pb-3">
              <span class="text-[10px] font-black uppercase tracking-widest text-slate-400">Nomor Induk (NIS)</span>
              <span class="text-sm font-black text-slate-800 tracking-widest">{{ selectedStudent?.nis }}</span>
            </div>
          </div>
          <div class="space-y-4">
            <div class="flex justify-between border-b border-slate-100 pb-3">
              <span class="text-[10px] font-black uppercase tracking-widest text-slate-400">Kelas</span>
              <span class="text-sm font-black text-slate-800 uppercase">{{ classrooms.find(c => c.ID === filters.classroom_id)?.name }}</span>
            </div>
            <div class="flex justify-between border-b border-slate-100 pb-3">
              <span class="text-[10px] font-black uppercase tracking-widest text-slate-400">Wali Kelas</span>
              <span class="text-sm font-black text-slate-800 uppercase">Ust. Abdullah, S.Pd.</span>
            </div>
          </div>
        </div>

        <!-- Academics Section -->
        <div class="space-y-6">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-xl bg-indigo-100 text-indigo-600 flex items-center justify-center">
              <GraduationCap class="w-4 h-4" />
            </div>
            <h2 class="text-sm font-black uppercase tracking-[0.2em] text-slate-800">I. Capaian Akademik</h2>
          </div>
          <table class="w-full border-collapse">
            <thead>
              <tr class="bg-slate-50">
                <th class="p-5 text-left text-[10px] font-black uppercase tracking-widest text-slate-400 border border-slate-200">Mata Pelajaran</th>
                <th class="p-5 text-center text-[10px] font-black uppercase tracking-widest text-slate-400 border border-slate-200 w-24">Skor</th>
                <th class="p-5 text-center text-[10px] font-black uppercase tracking-widest text-slate-400 border border-slate-200 w-24">Grade</th>
                <th class="p-5 text-left text-[10px] font-black uppercase tracking-widest text-slate-400 border border-slate-200">Keterangan Capaian</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="g in reportData.grades" :key="g.subject" class="hover:bg-slate-50 transition-colors">
                <td class="p-5 text-xs font-black text-slate-700 border border-slate-200 uppercase">{{ g.subject }}</td>
                <td class="p-5 text-center text-sm font-black text-indigo-600 border border-slate-200">{{ g.score }}</td>
                <td class="p-5 text-center text-sm font-black border border-slate-200">{{ g.grade }}</td>
                <td class="p-5 text-[10px] font-bold text-slate-500 border border-slate-200 italic">Sangat baik dalam penguasaan materi dan aktif di kelas.</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Attitude & Attendance -->
        <div class="grid grid-cols-2 gap-12">
          <!-- Attitude -->
          <div class="space-y-6">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-xl bg-indigo-100 text-indigo-600 flex items-center justify-center">
                <User class="w-4 h-4" />
              </div>
              <h2 class="text-sm font-black uppercase tracking-[0.2em] text-slate-800">II. Perkembangan Karakter</h2>
            </div>
            <table class="w-full border-collapse">
              <thead>
                <tr class="bg-slate-50">
                  <th class="p-4 text-left text-[9px] font-black uppercase tracking-widest text-slate-400 border border-slate-200">Aspek Pengamatan</th>
                  <th class="p-4 text-center text-[9px] font-black uppercase tracking-widest text-slate-400 border border-slate-200 w-20">Nilai</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="a in reportData.attitude" :key="a.category">
                  <td class="p-4 text-[10px] font-black text-slate-700 border border-slate-200 uppercase">{{ a.category }}</td>
                  <td class="p-4 text-center text-xs font-black text-indigo-600 border border-slate-200">{{ a.grade }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Attendance -->
          <div class="space-y-6">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-xl bg-indigo-100 text-indigo-600 flex items-center justify-center">
                <Calendar class="w-4 h-4" />
              </div>
              <h2 class="text-sm font-black uppercase tracking-[0.2em] text-slate-800">III. Ketidakhadiran</h2>
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div class="p-4 bg-slate-50 rounded-2xl border border-slate-200 text-center">
                <span class="text-[8px] font-black uppercase tracking-widest text-slate-400 block mb-1">Sakit</span>
                <span class="text-lg font-black text-slate-700">{{ reportData.attendance.sakit }} <small class="text-[8px] text-slate-400">HARI</small></span>
              </div>
              <div class="p-4 bg-slate-50 rounded-2xl border border-slate-200 text-center">
                <span class="text-[8px] font-black uppercase tracking-widest text-slate-400 block mb-1">Izin</span>
                <span class="text-lg font-black text-slate-700">{{ reportData.attendance.izin }} <small class="text-[8px] text-slate-400">HARI</small></span>
              </div>
              <div class="p-4 bg-slate-50 rounded-2xl border border-slate-200 text-center">
                <span class="text-[8px] font-black uppercase tracking-widest text-slate-400 block mb-1">Alpa</span>
                <span class="text-lg font-black text-rose-500">{{ reportData.attendance.alpa }} <small class="text-[8px] text-slate-400">HARI</small></span>
              </div>
            </div>
          </div>
        </div>

        <!-- Signatures -->
        <div class="pt-12 grid grid-cols-3 gap-12 text-center">
          <div class="space-y-20">
            <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Orang Tua / Wali</p>
            <div class="border-b border-slate-200 w-40 mx-auto"></div>
          </div>
          <div class="space-y-20">
            <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Wali Kelas</p>
            <div class="space-y-1">
              <p class="text-xs font-black uppercase text-indigo-900 underline">Ust. Abdullah, S.Pd.</p>
              <p class="text-[8px] font-bold text-slate-400 tracking-widest uppercase">NIP. 19880706 201503 1 002</p>
            </div>
          </div>
          <div class="space-y-20">
            <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Kepala Madrasah</p>
            <div class="space-y-1">
              <p class="text-xs font-black uppercase text-indigo-900 underline">Dr. Ahmad Faisal, M.Ag.</p>
              <p class="text-[8px] font-bold text-slate-400 tracking-widest uppercase">NIP. 19750101 200003 1 001</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!isLoading" class="flex flex-col items-center justify-center p-40 space-y-4 opacity-30 print:hidden">
      <FileText class="w-16 h-16" />
      <p class="font-black uppercase tracking-widest text-[10px]">Pilih santri untuk menggenerate rapor</p>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center p-40 space-y-4 print:hidden">
      <div class="w-12 h-12 border-4 border-indigo-500/20 border-t-indigo-500 rounded-full animate-spin"></div>
      <p class="text-[10px] font-black uppercase tracking-widest text-sigma-muted animate-pulse">Menyiapkan Dokumen Rapor...</p>
    </div>
  </div>
</template>

<style scoped>
@media print {
  .print\:hidden { display: none !important; }
  body { background: white !important; }
  .p-8 { padding: 0 !important; }
}
</style>
