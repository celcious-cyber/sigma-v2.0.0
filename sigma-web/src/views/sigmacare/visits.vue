<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Plus, Search, Filter, 
  Stethoscope, Calendar, User,
  CheckCircle2, AlertCircle, X,
  FileText, Activity, ArrowUpRight
} from 'lucide-vue-next'
import axios from 'axios'

const visits = ref<any[]>([])
const students = ref<any[]>([])
const loading = ref(true)
const showAddModal = ref(false)
const stats = ref<any>({
  total_visits: 0,
  inpatients: 0,
  recovery_rate: 0
})

const newVisit = ref({
  student_id: '',
  anamnesis: '',
  diagnosis: '',
  icd10_code: '',
  treatment: '',
  status: 'Treated'
})

const fetchVisits = async () => {
  try {
    const response = await axios.get('/admin/care/visits')
    visits.value = response.data
  } catch (error) {
    console.error('Failed to fetch visits:', error)
  } finally {
    loading.value = false
  }
}

const fetchStudents = async () => {
  try {
    const response = await axios.get('/admin/base/students')
    console.log('Fetched students:', response.data)
    students.value = response.data.data || response.data
  } catch (error) {
    console.error('Failed to fetch students:', error)
  }
}

const fetchStats = async () => {
  try {
    const response = await axios.get('/admin/care/stats')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const submitVisit = async () => {
  try {
    await axios.post('/admin/care/visits', newVisit.value)
    showAddModal.value = false
    newVisit.value = { student_id: '', anamnesis: '', diagnosis: '', icd10_code: '', treatment: '', status: 'Treated' }
    fetchVisits()
    fetchStats()
  } catch (error) {
    alert('Gagal menyimpan data kunjungan')
  }
}

const createCertificate = async (visit: any) => {
  try {
    const cert = {
      visit_id: visit.ID,
      type: 'SAKIT',
      start_date: new Date(),
      end_date: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000), // 3 days
      digital_sign: 'SIGNED_BY_MEDIC'
    }
    await axios.post('/admin/care/certificates', cert)
    alert('Surat Sakit berhasil dibuat!')
  } catch (error) {
    alert('Gagal membuat surat sakit')
  }
}

const startObservation = async (visit: any) => {
  try {
    const obs = {
      visit_id: visit.ID,
      bed_id: 1, // Default bed
      notes: 'Observasi rutin',
      status: 'Active'
    }
    await axios.post('/admin/care/observations', obs)
    fetchStats()
    alert('Santri berhasil didaftarkan ke ruang observasi!')
  } catch (error) {
    alert('Gagal memulai observasi. Pastikan bed tersedia.')
  }
}

const createReferral = async (visit: any) => {
  try {
    const cert = {
      visit_id: visit.ID,
      type: 'RUJUKAN',
      start_date: new Date(),
      end_date: new Date(),
      digital_sign: 'SIGNED_BY_MEDIC',
      notes: 'Dirujuk ke RS terdekat untuk penanganan lebih lanjut'
    }
    await axios.post('/admin/care/certificates', cert)
    alert('Surat Rujukan berhasil dibuat!')
  } catch (error) {
    alert('Gagal membuat surat rujukan')
  }
}

const searchQuery = ref('')
const studentFilter = ref('')

const filteredVisits = computed(() => {
  if (!searchQuery.value) return visits.value
  const q = searchQuery.value.toLowerCase()
  return visits.value.filter(v => 
    v.student?.name?.toLowerCase().includes(q) || 
    v.diagnosis?.toLowerCase().includes(q) ||
    v.icd10_code?.toLowerCase().includes(q)
  )
})

const filteredStudents = computed(() => {
  if (!studentFilter.value) return students.value
  const q = studentFilter.value.toLowerCase()
  return students.value.filter(s => 
    s.name?.toLowerCase().includes(q) || 
    s.nisn?.toLowerCase().includes(q)
  )
})

const selectedStudentInfo = computed(() => {
  if (!newVisit.value.student_id) return null
  return students.value.find(s => s.ID === newVisit.value.student_id)
})

const selectStudent = (student: any) => {
  newVisit.value.student_id = student.ID
  studentFilter.value = ''
}

const clearSelectedStudent = () => {
  newVisit.value.student_id = ''
}

onMounted(() => {
  fetchVisits()
  fetchStudents()
  fetchStats()
})

const getStatusColor = (status: string) => {
  switch (status) {
    case 'Treated': return 'bg-emerald-50 text-emerald-600 border-emerald-100'
    case 'Resting': return 'bg-amber-50 text-amber-600 border-amber-100'
    case 'Referred': return 'bg-rose-50 text-rose-600 border-rose-100'
    default: return 'bg-slate-50 text-slate-600 border-slate-100'
  }
}
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Anamnesis & <span class="text-rose-600">Diagnosis</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <Stethoscope class="w-4 h-4" /> Pencatatan Kunjungan Pasien Poliklinik
        </p>
      </div>
      
      <div class="flex items-center gap-3">
        <div class="relative group">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-rose-600 transition-colors" />
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari Santri atau Diagnosis..." 
            class="pl-11 pr-6 py-3.5 bg-white border border-slate-200 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all w-64 md:w-80 shadow-sm"
          />
        </div>
        <button 
          @click="showAddModal = true"
          class="p-3.5 bg-rose-600 text-white rounded-2xl font-black shadow-lg shadow-rose-500/20 hover:bg-rose-700 active:scale-95 transition-all flex items-center gap-3"
        >
          <Plus class="w-5 h-5" /> <span class="hidden md:inline text-xs uppercase tracking-widest">Catat Kunjungan</span>
        </button>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white p-6 rounded-[2rem] border border-slate-100 flex items-center gap-6 shadow-sm">
        <div class="w-14 h-14 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center">
          <Calendar class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Bulan Ini</p>
          <p class="text-2xl font-black text-slate-800 tracking-tighter">{{ stats.total_visits }} Kunjungan</p>
        </div>
      </div>
      <div class="bg-white p-6 rounded-[2rem] border border-slate-100 flex items-center gap-6 shadow-sm">
        <div class="w-14 h-14 bg-amber-50 text-amber-600 rounded-2xl flex items-center justify-center">
          <AlertCircle class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Pasien Observasi</p>
          <p class="text-2xl font-black text-slate-800 tracking-tighter">{{ stats.inpatients }} Aktif</p>
        </div>
      </div>
      <div class="bg-white p-6 rounded-[2rem] border border-slate-100 flex items-center gap-6 shadow-sm">
        <div class="w-14 h-14 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center">
          <CheckCircle2 class="w-7 h-7" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tingkat Kesembuhan</p>
          <p class="text-2xl font-black text-slate-800 tracking-tighter">{{ stats.recovery_rate.toFixed(1) }}%</p>
        </div>
      </div>
    </div>

    <!-- Table Section -->
    <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
      <div class="p-8 border-b border-slate-50 flex items-center justify-between">
        <h2 class="text-lg font-black text-slate-800 uppercase italic">Riwayat Kunjungan Terakhir</h2>
        <div class="flex items-center gap-2">
          <button class="p-2.5 bg-slate-50 text-slate-400 rounded-xl hover:text-rose-600 transition-colors">
            <Filter class="w-4 h-4" />
          </button>
        </div>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Santri</th>
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Tanggal / Jam</th>
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Diagnosis</th>
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Status</th>
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Tindakan</th>
              <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="visit in filteredVisits" :key="visit.ID" class="hover:bg-rose-50/30 transition-colors group">
              <td class="px-8 py-6">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-slate-400">
                    <User class="w-5 h-5" />
                  </div>
                  <div>
                    <p class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ visit.student?.name }}</p>
                    <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">NISN: {{ visit.student?.nisn }}</p>
                  </div>
                </div>
              </td>
              <td class="px-8 py-6">
                <p class="font-bold text-slate-600 text-sm italic">{{ new Date(visit.visit_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) }}</p>
                <p class="text-[10px] font-black text-slate-400 uppercase mt-0.5 tracking-tighter">{{ new Date(visit.visit_date).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} WIB</p>
              </td>
              <td class="px-8 py-6">
                <div class="flex items-center gap-2">
                  <span class="px-2 py-0.5 bg-slate-100 text-slate-500 rounded text-[9px] font-black uppercase">{{ visit.icd10_code || 'N/A' }}</span>
                  <p class="font-black text-slate-700 text-sm uppercase tracking-tight">{{ visit.diagnosis }}</p>
                </div>
              </td>
              <td class="px-8 py-6">
                <span 
                  class="px-4 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border"
                  :class="getStatusColor(visit.status)"
                >
                  {{ visit.status }}
                </span>
              </td>
              <td class="px-8 py-6">
                <p class="text-xs font-bold text-slate-500 italic max-w-[200px] truncate">{{ visit.treatment }}</p>
              </td>
              <td class="px-8 py-6 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    @click="createCertificate(visit)"
                    title="Buat Surat Sakit" 
                    class="p-2.5 bg-amber-50 text-amber-600 rounded-xl hover:bg-amber-100 transition-all active:scale-90 group-hover:shadow-sm"
                  >
                    <FileText class="w-4 h-4" />
                  </button>
                  <button 
                    @click="startObservation(visit)"
                    title="Mulai Observasi" 
                    class="p-2.5 bg-emerald-50 text-emerald-600 rounded-xl hover:bg-emerald-100 transition-all active:scale-90 group-hover:shadow-sm"
                  >
                    <Activity class="w-4 h-4" />
                  </button>
                  <button 
                    @click="createReferral(visit)"
                    title="Buat Rujukan" 
                    class="p-2.5 bg-rose-50 text-rose-600 rounded-xl hover:bg-rose-100 transition-all active:scale-90 group-hover:shadow-sm"
                  >
                    <ArrowUpRight class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!visits.length && !loading">
              <td colspan="6" class="px-8 py-20 text-center">
                <div class="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Stethoscope class="w-8 h-8 text-slate-200" />
                </div>
                <p class="text-slate-400 font-bold italic">Belum ada data kunjungan hari ini.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Add Visit -->
    <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm animate-in fade-in duration-300">
      <div class="bg-white w-full max-w-2xl rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 flex items-center justify-between bg-rose-600 text-white">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-white/20 rounded-2xl flex items-center justify-center">
              <Plus class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-xl font-black italic uppercase">Catat Kunjungan Baru</h3>
              <p class="text-xs font-bold text-white/70 uppercase tracking-widest">Lengkapi Detail Rekam Medis</p>
            </div>
          </div>
          <button @click="showAddModal = false" class="p-2 hover:bg-white/10 rounded-xl transition-colors">
            <X class="w-6 h-6" />
          </button>
        </div>

        <div class="p-10 space-y-8 max-h-[70vh] overflow-y-auto custom-scrollbar">
          <!-- Student Search (Invoice Style) -->
          <div class="space-y-4">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Cari Nama atau NISN Santri</label>
            <div class="relative group">
              <div class="absolute left-6 top-1/2 -translate-y-1/2 w-10 h-10 bg-rose-50 rounded-xl flex items-center justify-center">
                <Search class="w-4 h-4 text-rose-600" />
              </div>
              <input 
                v-model="studentFilter"
                type="text" 
                placeholder="Ketik Nama atau NISN..." 
                class="w-full pl-20 pr-8 py-5 bg-slate-50 border border-slate-100 rounded-[2rem] text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all"
              />
            </div>

            <!-- Student Results -->
            <div v-if="studentFilter && filteredStudents.length" class="bg-white rounded-3xl border border-slate-100 shadow-xl overflow-hidden divide-y divide-slate-50 max-h-60 overflow-y-auto custom-scrollbar">
              <button 
                v-for="s in filteredStudents" 
                :key="s.ID"
                @click="selectStudent(s)"
                class="w-full px-8 py-4 text-left hover:bg-rose-50 transition-colors flex items-center justify-between group"
              >
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-slate-50 rounded-xl flex items-center justify-center group-hover:bg-white transition-colors">
                    <User class="w-5 h-5 text-slate-400 group-hover:text-rose-600" />
                  </div>
                  <div>
                    <p class="font-black text-slate-800 uppercase text-xs">{{ s.name }}</p>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">NISN: {{ s.nisn }}</p>
                  </div>
                </div>
                <div class="px-3 py-1 rounded-lg bg-slate-100 text-[9px] font-black text-slate-400 uppercase group-hover:bg-rose-600 group-hover:text-white transition-all">Pilih</div>
              </button>
            </div>

            <!-- Selected Student Chip -->
            <div v-if="selectedStudentInfo" class="p-6 bg-emerald-50 rounded-3xl border border-emerald-100 flex items-center justify-between animate-in zoom-in duration-300">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 bg-emerald-500 text-white rounded-2xl flex items-center justify-center shadow-lg shadow-emerald-500/20">
                  <CheckCircle2 class="w-6 h-6" />
                </div>
                <div>
                  <p class="text-[10px] font-black text-emerald-600 uppercase tracking-widest leading-none mb-1">Santri Terpilih</p>
                  <p class="font-black text-slate-800 uppercase leading-tight">{{ selectedStudentInfo.name }}</p>
                  <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mt-0.5">NIS: {{ selectedStudentInfo.nis }}</p>
                </div>
              </div>
              <button @click="clearSelectedStudent" class="p-2 text-slate-300 hover:text-rose-500 transition-colors">
                <X class="w-5 h-5" />
              </button>
            </div>
            
            <p v-if="studentFilter && !filteredStudents.length" class="text-[10px] text-slate-400 font-bold italic text-center py-4 bg-slate-50 rounded-2xl">Nama tidak ditemukan.</p>
          </div>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kode ICD-10 (Opsional)</label>
              <input 
                v-model="newVisit.icd10_code"
                type="text" 
                placeholder="Misal: J00"
                class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all"
              />
            </div>
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Status Kunjungan</label>
              <select 
                v-model="newVisit.status"
                class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all"
              >
                <option value="Treated">Ditangani (Langsung Pulang)</option>
                <option value="Resting">Observasi (Istirahat di UKS)</option>
                <option value="Referred">Rujuk (RS/Klinik Luar)</option>
              </select>
            </div>
          </div>

          <div class="space-y-3">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Anamnesis (Keluhan Utama)</label>
            <textarea 
              v-model="newVisit.anamnesis"
              rows="3"
              placeholder="Ceritakan keluhan santri secara detail..."
              class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all resize-none"
            ></textarea>
          </div>

          <div class="space-y-3">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Diagnosis & Tindakan Medis</label>
            <textarea 
              v-model="newVisit.treatment"
              rows="3"
              placeholder="Tindakan yang dilakukan dan obat yang diberikan..."
              class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all resize-none"
            ></textarea>
          </div>
        </div>

        <div class="p-8 bg-slate-50 border-t border-slate-100 flex gap-4">
          <button 
            @click="showAddModal = false"
            class="flex-1 py-4 bg-white border border-slate-200 text-slate-500 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-100 transition-all"
          >
            Batal
          </button>
          <button 
            @click="submitVisit"
            class="flex-[2] py-4 bg-rose-600 text-white rounded-2xl font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-rose-500/20 hover:bg-rose-700 active:scale-95 transition-all"
          >
            Simpan Rekam Medis
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(225, 29, 72, 0.1); border-radius: 10px; }
</style>
