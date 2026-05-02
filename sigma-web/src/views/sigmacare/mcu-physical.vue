<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { 
  Activity, Scale, Ruler,
  Plus, Search, User,
  Save, X, CheckCircle2
} from 'lucide-vue-next'
import axios from 'axios'

const students = ref<any[]>([])
const mcuRecords = ref<any[]>([])
const loading = ref(true)
const showAddModal = ref(false)
const searchQuery = ref('')
const studentFilter = ref('')

const filteredRecords = computed(() => {
  if (!searchQuery.value) return mcuRecords.value
  const q = searchQuery.value.toLowerCase()
  return mcuRecords.value.filter(r => 
    r.student?.name?.toLowerCase().includes(q) ||
    r.student?.nisn?.toLowerCase().includes(q)
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
  if (!newMCU.value.student_id) return null
  return students.value.find(s => s.ID === newMCU.value.student_id)
})

const selectStudent = (student: any) => {
  newMCU.value.student_id = student.ID
  studentFilter.value = ''
}

const clearSelectedStudent = () => {
  newMCU.value.student_id = ''
}

const newMCU = ref({
  student_id: '',
  height: 0,
  weight: 0,
  vision_left: '6/6',
  vision_right: '6/6',
  blood_type: '-',
  allergies: '-'
})

const bmi = computed(() => {
  if (newMCU.value.height > 0 && newMCU.value.weight > 0) {
    const hMeter = newMCU.value.height / 100
    return (newMCU.value.weight / (hMeter * hMeter)).toFixed(1)
  }
  return '0.0'
})

const bmiCategory = computed(() => {
  const val = parseFloat(bmi.value)
  if (val === 0) return { label: 'N/A', color: 'text-slate-400' }
  if (val < 18.5) return { label: 'Underweight', color: 'text-amber-500' }
  if (val < 25) return { label: 'Normal', color: 'text-emerald-500' }
  if (val < 30) return { label: 'Overweight', color: 'text-orange-500' }
  return { label: 'Obese', color: 'text-rose-600' }
})


const fetchStudents = async () => {
  try {
    const response = await axios.get('/admin/base/students')
    students.value = response.data.data || response.data
  } catch (error) {
    console.error('Failed to fetch students:', error)
  }
}

const fetchMCURecords = async () => {
  try {
    const response = await axios.get('/admin/care/mcu')
    mcuRecords.value = response.data
  } catch (error) {
    console.error('Failed to fetch MCU records:', error)
  } finally {
    loading.value = false
  }
}

const submitMCU = async () => {
  try {
    await axios.post('/admin/care/mcu', newMCU.value)
    showAddModal.value = false
    alert('Data Antropometri Berhasil Disimpan')
    newMCU.value = { student_id: '', height: 0, weight: 0, vision_left: '6/6', vision_right: '6/6', blood_type: '-', allergies: '-' }
    fetchMCURecords()
  } catch (error) {
    alert('Gagal menyimpan data MCU')
  }
}

onMounted(() => {
  fetchStudents()
  fetchMCURecords()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Prosedur <span class="text-rose-600">Antropometri</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <Activity class="w-4 h-4" /> Pemeriksaan Fisik Rutin & Monitoring Pertumbuhan
        </p>
      </div>
      
      <button 
        @click="showAddModal = true"
        class="px-8 py-4 bg-rose-600 text-white rounded-2xl font-black shadow-lg shadow-rose-500/20 hover:bg-rose-700 active:scale-95 transition-all flex items-center gap-3"
      >
        <Plus class="w-5 h-5" /> <span class="text-xs uppercase tracking-widest">Input Hasil Pemeriksaan</span>
      </button>
    </div>

    <!-- BMI Reference Card -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <div v-for="cat in [
        { l: 'Underweight', r: '< 18.5', c: 'bg-amber-500' },
        { l: 'Normal', r: '18.5 - 24.9', c: 'bg-emerald-500' },
        { l: 'Overweight', r: '25 - 29.9', c: 'bg-orange-500' },
        { l: 'Obese', r: '> 30', c: 'bg-rose-600' }
      ]" :key="cat.l" class="bg-white p-4 rounded-2xl border border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div :class="`w-2 h-8 ${cat.c} rounded-full`"></div>
          <span class="text-[10px] font-black uppercase text-slate-400 tracking-widest">{{ cat.l }}</span>
        </div>
        <span class="text-xs font-black text-slate-700 italic">{{ cat.r }}</span>
      </div>
    </div>

    <!-- Student List for MCU -->
    <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
      <div class="p-8 border-b border-slate-50 flex items-center justify-between">
        <h2 class="text-lg font-black text-slate-800 uppercase italic">Status Pertumbuhan Santri</h2>
        <div class="relative">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari Nama/Kelas..." 
            class="pl-11 pr-6 py-2.5 bg-slate-50 border border-slate-100 rounded-xl text-xs font-bold focus:outline-none focus:ring-2 focus:ring-rose-500/20" 
          />
        </div>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Santri</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Tinggi (cm)</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Berat (kg)</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">IMT / BMI</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Visus (L/R)</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Update Terakhir</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="record in filteredRecords" :key="record.ID" class="hover:bg-slate-50 transition-colors">
              <td class="px-8 py-6">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-slate-400">
                    <User class="w-5 h-5" />
                  </div>
                  <div>
                    <p class="font-black text-slate-800 uppercase tracking-tight">{{ record.student?.name }}</p>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">NISN: {{ record.student?.nisn }}</p>
                  </div>
                </div>
              </td>
              <td class="px-8 py-6 font-bold text-slate-600">{{ record.height }}</td>
              <td class="px-8 py-6 font-bold text-slate-600">{{ record.weight }}</td>
              <td class="px-8 py-6">
                <div class="flex items-center gap-2">
                  <span class="font-black text-slate-800">{{ record.bmi.toFixed(1) }}</span>
                  <span 
                    class="px-2 py-0.5 rounded text-[9px] font-black uppercase"
                    :class="record.bmi_category === 'Normal' ? 'bg-emerald-50 text-emerald-600' : 'bg-amber-50 text-amber-600'"
                  >
                    {{ record.bmi_category }}
                  </span>
                </div>
              </td>
              <td class="px-8 py-6">
                <span class="font-black text-slate-600 tracking-widest">{{ record.vision_left }} | {{ record.vision_right }}</span>
              </td>
              <td class="px-8 py-6 text-slate-400 text-xs font-bold italic">
                {{ new Date(record.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) }}
              </td>
            </tr>
            <tr v-if="!mcuRecords.length && !loading">
              <td colspan="6" class="px-8 py-20 text-center">
                <p class="text-slate-400 font-bold italic">Belum ada data pemeriksaan rutin.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Input MCU -->
    <div v-if="showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-2xl rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 flex items-center justify-between bg-rose-600 text-white font-black italic uppercase tracking-tight">
          Hasil Pemeriksaan Fisik
          <button @click="showAddModal = false" class="p-2 hover:bg-white/10 rounded-xl transition-colors">
            <X class="w-6 h-6" />
          </button>
        </div>

        <div class="p-10 space-y-8">
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

          <!-- Physical Metrics -->
          <div class="grid grid-cols-2 gap-8">
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Tinggi Badan (cm)</label>
              <div class="relative">
                <Ruler class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                <input v-model.number="newMCU.height" type="number" class="w-full pl-12 pr-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-lg font-black focus:border-rose-500 outline-none" />
              </div>
            </div>
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Berat Badan (kg)</label>
              <div class="relative">
                <Scale class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                <input v-model.number="newMCU.weight" type="number" class="w-full pl-12 pr-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-lg font-black focus:border-rose-500 outline-none" />
              </div>
            </div>
          </div>

          <!-- BMI Calculator View -->
          <div class="bg-slate-50 rounded-3xl p-8 flex items-center justify-between border-2 border-dashed border-slate-200">
            <div>
              <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Kalkulasi IMT (BMI) Otomatis</p>
              <h4 :class="`text-4xl font-black tracking-tighter ${bmiCategory.color}`">{{ bmi }}</h4>
            </div>
            <div class="text-right">
              <span :class="`px-6 py-2 rounded-full text-xs font-black uppercase tracking-widest ${bmiCategory.color.replace('text-', 'bg-').replace('-600', '-100').replace('-500', '-100')} border`">
                {{ bmiCategory.label }}
              </span>
            </div>
          </div>

          <!-- Vision -->
          <div class="grid grid-cols-2 gap-8">
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Visus Mata Kiri</label>
              <input v-model="newMCU.vision_left" type="text" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
            </div>
            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Visus Mata Kanan</label>
              <input v-model="newMCU.vision_right" type="text" class="w-full px-6 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold" />
            </div>
          </div>
        </div>

        <div class="p-8 bg-slate-50 border-t border-slate-100">
          <button @click="submitMCU" class="w-full py-5 bg-rose-600 text-white rounded-2xl font-black text-xs uppercase tracking-[0.2em] shadow-xl hover:bg-rose-700 transition-all flex items-center justify-center gap-3">
            <Save class="w-5 h-5" /> Simpan Hasil MCU
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
