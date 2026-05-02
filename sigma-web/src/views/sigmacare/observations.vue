<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Bed, Clock, 
  Activity, Thermometer, Pill,
  Plus,
  MoreVertical, X,
  Save
} from 'lucide-vue-next'
import axios from 'axios'

const activePatients = ref<any[]>([])
const loading = ref(true)
const showCareModal = ref(false)
const selectedPatient = ref<any>(null)

const careUpdate = ref({
  temperature: 36.5,
  blood_pressure: '120/80',
  notes: '',
  meds: ''
})

const fetchObservations = async () => {
  try {
    const response = await axios.get('/admin/care/observations')
    activePatients.value = response.data
  } catch (error) {
    console.error('Failed to fetch observations:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchObservations()
})

const openCareUpdate = (patient: any) => {
  selectedPatient.value = patient
  careUpdate.value = {
    temperature: patient.temperature || 36.5,
    blood_pressure: patient.blood_pressure || '120/80',
    notes: patient.notes || '',
    meds: patient.meds || ''
  }
  showCareModal.value = true
}

const saveCareUpdate = async () => {
  try {
    const payload = {
      temperature: parseFloat(careUpdate.value.temperature.toString()),
      blood_pressure: careUpdate.value.blood_pressure,
      notes: careUpdate.value.notes,
      meds: careUpdate.value.meds
    }
    await axios.patch(`/admin/care/observations/${selectedPatient.value.ID}`, payload)
    showCareModal.value = false
    fetchObservations()
    alert('Kondisi santri berhasil diperbarui!')
  } catch (error) {
    console.error('Update error:', error)
    alert('Gagal memperbarui kondisi')
  }
}

const dischargePatient = async (id: number) => {
  if (!confirm('Apakah santri sudah diperbolehkan keluar dari ruang observasi?')) return
  try {
    await axios.delete(`/admin/care/observations/${id}`)
    fetchObservations()
    alert('Santri berhasil diperbolehkan keluar.')
  } catch (error) {
    alert('Gagal memproses kepulangan santri')
  }
}

const calculateDuration = (checkIn: string) => {
  const start = new Date(checkIn)
  const now = new Date()
  const diffHours = Math.floor((now.getTime() - start.getTime()) / (1000 * 60 * 60))
  if (diffHours < 24) return `${diffHours} Jam`
  const days = Math.floor(diffHours / 24)
  return `${days} Hari`
}
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Observasi <span class="text-rose-600">Klinis</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <Bed class="w-4 h-4" /> Monitoring Santri Rawat Inap & Istirahat UKS
        </p>
      </div>
      
      <div class="flex items-center gap-4 bg-white px-6 py-3 rounded-2xl border border-slate-100 shadow-sm">
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 bg-rose-500 rounded-full animate-pulse"></div>
          <span class="text-[10px] font-black uppercase text-slate-400 tracking-widest">Kapasitas Bed:</span>
        </div>
        <span class="text-sm font-black text-slate-800">{{ activePatients.length }} / 12 Terisi</span>
      </div>
    </div>

    <!-- Active Patients Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div v-for="patient in activePatients" :key="patient.ID" class="group bg-white rounded-[2.5rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-rose-500/5 transition-all overflow-hidden flex flex-col">
        <!-- Card Header -->
        <div class="p-8 border-b border-slate-50 flex items-start justify-between">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center font-black text-xl">
              {{ patient.visit?.student?.name?.charAt(0) }}
            </div>
            <div>
              <h3 class="font-black text-slate-800 uppercase tracking-tight">{{ patient.visit?.student?.name }}</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">{{ patient.visit?.student?.classroom?.name || 'Kelas N/A' }}</p>
            </div>
          </div>
          <button class="p-2 text-slate-300 hover:text-rose-600 transition-colors">
            <MoreVertical class="w-5 h-5" />
          </button>
        </div>

        <!-- Monitoring Info -->
        <div class="p-8 space-y-6 flex-1">
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Diagnosis</p>
              <p class="text-xs font-black text-slate-700 uppercase leading-relaxed">{{ patient.visit?.diagnosis }}</p>
            </div>
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Durasi</p>
              <p class="text-xs font-black text-rose-600 uppercase">{{ calculateDuration(patient.check_in) }}</p>
            </div>
          </div>

          <div class="space-y-4">
            <div class="flex items-center justify-between text-xs">
              <span class="font-bold text-slate-400 uppercase tracking-widest flex items-center gap-2">
                <Clock class="w-3.5 h-3.5" /> Check-in
              </span>
              <span class="font-black text-slate-600">{{ new Date(patient.check_in).toLocaleString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} WIB</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="font-bold text-slate-400 uppercase tracking-widest flex items-center gap-2">
                <Bed class="w-3.5 h-3.5" /> Nomor Bed
              </span>
              <span class="font-black text-rose-600 italic uppercase">{{ patient.bed?.name || 'B-' + patient.bed_id }}</span>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="p-6 bg-slate-50/50 border-t border-slate-50 flex gap-3">
          <button @click="openCareUpdate(patient)" class="flex-1 py-3 bg-white border border-slate-200 text-rose-600 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-rose-600 hover:text-white transition-all shadow-sm">
            Update Kondisi
          </button>
          <button @click="dischargePatient(patient.ID)" class="px-4 py-3 bg-emerald-600 text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-500/10">
            Selesai
          </button>
        </div>
      </div>

      <!-- Add New Patient Placeholder -->
      <button class="group border-2 border-dashed border-slate-200 rounded-[2.5rem] p-10 flex flex-col items-center justify-center gap-4 hover:border-rose-500 hover:bg-rose-50/30 transition-all text-slate-400 hover:text-rose-600">
        <div class="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center group-hover:bg-rose-100 transition-all">
          <Plus class="w-10 h-10" />
        </div>
        <p class="font-black uppercase tracking-[0.2em] text-xs">Tambah Pasien Istirahat</p>
      </button>
    </div>

    <!-- Modal Update Condition -->
    <div v-if="showCareModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 bg-rose-600 text-white flex items-center justify-between">
          <div>
            <h3 class="text-xl font-black italic uppercase">Update Monitoring</h3>
            <p class="text-xs font-bold text-white/70 uppercase tracking-widest">Observasi Berkala Santri</p>
          </div>
          <button @click="showCareModal = false" class="p-2 hover:bg-white/10 rounded-xl transition-colors">
            <X class="w-6 h-6" />
          </button>
        </div>

        <div class="p-10 space-y-6">
          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 flex items-center gap-2">
                <Thermometer class="w-3 h-3" /> Suhu Tubuh (°C)
              </label>
              <input v-model="careUpdate.temperature" type="text" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-lg font-black focus:border-rose-500 outline-none" />
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 flex items-center gap-2">
                <Activity class="w-3 h-3" /> Tekanan Darah
              </label>
              <input v-model="careUpdate.blood_pressure" type="text" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-lg font-black focus:border-rose-500 outline-none" />
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 flex items-center gap-2">
              <Pill class="w-3 h-3" /> Terapi Obat Berjalan
            </label>
            <textarea v-model="careUpdate.meds" rows="2" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold resize-none focus:border-rose-500 outline-none"></textarea>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 flex items-center gap-2">
              <FileText class="w-3 h-3" /> Catatan Perkembangan
            </label>
            <textarea v-model="careUpdate.notes" rows="3" placeholder="Bagaimana kondisi santri sekarang?" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold resize-none focus:border-rose-500 outline-none"></textarea>
          </div>
        </div>

        <div class="p-8 bg-slate-50 border-t border-slate-100 flex gap-4">
          <button @click="showCareModal = false" class="flex-1 py-4 bg-white border border-slate-200 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-100 transition-all">Batal</button>
          <button @click="saveCareUpdate" class="flex-[2] py-4 bg-rose-600 text-white rounded-2xl font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-rose-500/20 hover:bg-rose-700 transition-all flex items-center justify-center gap-3">
            <Save class="w-4 h-4" /> Simpan Observasi
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
