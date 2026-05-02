<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  History, User, Search, 
  Stethoscope,
  Activity, Pill, FileText,
  TrendingUp, ArrowLeft, Filter,
  ExternalLink, HeartPulse
} from 'lucide-vue-next'
import axios from 'axios'

const students = ref<any[]>([])
const selectedStudent = ref<any>(null)
const medicalHistory = ref<any[]>([])
const loading = ref(false)

const fetchStudents = async () => {
  try {
    const response = await axios.get('/admin/base/students')
    students.value = response.data.data || response.data
  } catch (error) {
    console.error('Failed to fetch students:', error)
  }
}

const selectStudent = async (student: any) => {
  selectedStudent.value = student
  loading.value = true
  try {
    // In a real app, fetch visits for this specific student
    const response = await axios.get('/admin/care/visits')
    medicalHistory.value = response.data.filter((v: any) => v.student_id === student.ID)
  } catch (error) {
    console.error('Failed to fetch history:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStudents()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div class="flex items-center gap-4">
        <button v-if="selectedStudent" @click="selectedStudent = null" class="p-3 bg-white border border-slate-200 rounded-xl hover:bg-rose-50 hover:text-rose-600 transition-all">
          <ArrowLeft class="w-5 h-5" />
        </button>
        <div>
          <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Riwayat <span class="text-rose-600">Terintegrasi</span></h1>
          <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
            <History class="w-4 h-4" /> Rekam Medis Elektronik (RME) Personal Santri
          </p>
        </div>
      </div>
      
      <div v-if="!selectedStudent" class="relative group">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-rose-600 transition-colors" />
        <input 
          type="text" 
          placeholder="Cari Nama Santri..." 
          class="pl-11 pr-6 py-3.5 bg-white border border-slate-200 rounded-2xl text-sm font-bold focus:outline-none focus:ring-4 focus:ring-rose-500/10 focus:border-rose-500 transition-all w-64 md:w-80 shadow-sm"
        />
      </div>
    </div>

    <!-- Student Selection Grid -->
    <div v-if="!selectedStudent" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <button 
        v-for="s in students" 
        :key="s.ID" 
        @click="selectStudent(s)"
        class="bg-white p-8 rounded-[2rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-rose-500/5 transition-all text-left flex flex-col items-center group"
      >
        <div class="w-20 h-20 bg-slate-50 text-slate-300 rounded-[2rem] flex items-center justify-center mb-4 group-hover:bg-rose-50 group-hover:text-rose-600 transition-all">
          <User class="w-10 h-10" />
        </div>
        <h3 class="font-black text-slate-800 uppercase text-center tracking-tight leading-tight mb-1">{{ s.full_name }}</h3>
        <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ s.classroom?.name || 'Kelas N/A' }}</p>
        <div class="mt-6 flex gap-2">
          <span class="px-3 py-1 bg-slate-100 text-slate-500 rounded-lg text-[9px] font-black uppercase tracking-widest">NISN: {{ s.nisn }}</span>
        </div>
      </button>
    </div>

    <!-- Patient Profile View -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <!-- Left Column: Patient Profile -->
      <div class="space-y-8">
        <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm p-10 flex flex-col items-center text-center relative overflow-hidden">
          <div class="absolute -right-10 -top-10 w-40 h-40 bg-rose-50 rounded-full opacity-30"></div>
          
          <div class="w-28 h-28 bg-rose-600 rounded-[2.5rem] flex items-center justify-center text-white text-4xl font-black mb-6 shadow-2xl shadow-rose-500/30">
            {{ selectedStudent.full_name.charAt(0) }}
          </div>
          <h2 class="text-2xl font-black text-slate-800 italic uppercase tracking-tighter">{{ selectedStudent.full_name }}</h2>
          <p class="text-xs font-bold text-slate-400 uppercase tracking-[0.2em] mt-2 italic">Data Santri Terverifikasi</p>
          
          <div class="w-full grid grid-cols-2 gap-4 mt-10">
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Tinggi Badan</p>
              <p class="text-sm font-black text-slate-800">172 cm</p>
            </div>
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Berat Badan</p>
              <p class="text-sm font-black text-slate-800">65 kg</p>
            </div>
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Golongan Darah</p>
              <p class="text-sm font-black text-rose-600">B+</p>
            </div>
            <div class="bg-slate-50 p-4 rounded-2xl border border-slate-100">
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Alergi</p>
              <p class="text-sm font-black text-emerald-600 uppercase">Debu</p>
            </div>
          </div>
        </div>

        <div class="bg-slate-800 rounded-[2.5rem] p-10 text-white shadow-2xl shadow-slate-800/20 relative overflow-hidden">
          <div class="absolute -right-4 -bottom-4 opacity-10">
            <HeartPulse class="w-32 h-32" />
          </div>
          <h3 class="text-lg font-black italic uppercase mb-6 flex items-center gap-3">
            <TrendingUp class="w-5 h-5 text-rose-500" /> Skor Kesejahteraan
          </h3>
          <div class="space-y-6">
            <div class="flex justify-between items-end">
              <span class="text-3xl font-black tracking-tighter">B+ <span class="text-xs text-white/50 tracking-normal italic">(85%)</span></span>
              <span class="text-[10px] font-black uppercase text-emerald-400 tracking-widest">Sangat Sehat</span>
            </div>
            <div class="h-2 bg-white/10 rounded-full overflow-hidden">
              <div class="h-full bg-emerald-500 w-[85%] rounded-full shadow-[0_0_10px_rgba(16,185,129,0.5)]"></div>
            </div>
            <p class="text-xs font-bold text-white/40 uppercase tracking-widest leading-relaxed italic">
              Berdasarkan kunjungan & riwayat penyakit dalam 12 bulan terakhir.
            </p>
          </div>
        </div>
      </div>

      <!-- Right Column: Medical Timeline -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm p-10">
          <div class="flex items-center justify-between mb-10">
            <h3 class="text-xl font-black text-slate-800 uppercase italic">Timeline Rekam Medis</h3>
            <button class="p-2.5 bg-slate-50 text-slate-400 rounded-xl hover:text-rose-600 transition-colors">
              <Filter class="w-4 h-4" />
            </button>
          </div>

          <div class="relative pl-10 space-y-12 before:absolute before:left-[19px] before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-100">
            <div v-for="item in medicalHistory" :key="item.ID" class="relative group">
              <!-- Marker -->
              <div class="absolute -left-[31px] top-1 w-10 h-10 bg-white border-4 border-rose-50 text-rose-600 rounded-2xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-all group-hover:border-rose-100 group-hover:shadow-rose-500/10">
                <component :is="item.status === 'Resting' ? Activity : Stethoscope" class="w-4 h-4" />
              </div>
              
              <div class="bg-slate-50/50 group-hover:bg-rose-50/50 rounded-3xl p-8 transition-all border border-transparent group-hover:border-rose-100">
                <div class="flex justify-between items-start mb-4">
                  <div>
                    <span class="text-[10px] font-black text-rose-500 uppercase tracking-widest mb-1 block">{{ new Date(item.visit_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) }}</span>
                    <h4 class="text-lg font-black text-slate-800 uppercase tracking-tight">{{ item.diagnosis }}</h4>
                  </div>
                  <span class="px-3 py-1 bg-white border border-slate-200 rounded-lg text-[9px] font-black text-slate-400 uppercase tracking-widest italic group-hover:border-rose-200 group-hover:text-rose-600">{{ item.icd10_code || 'ICD-10' }}</span>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
                  <div class="flex items-start gap-4">
                    <div class="w-8 h-8 bg-white rounded-lg flex items-center justify-center shadow-sm">
                      <FileText class="w-4 h-4 text-slate-400" />
                    </div>
                    <div>
                      <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Anamnesis</p>
                      <p class="text-xs font-bold text-slate-600 italic leading-relaxed mt-1">{{ item.anamnesis }}</p>
                    </div>
                  </div>
                  <div class="flex items-start gap-4">
                    <div class="w-8 h-8 bg-white rounded-lg flex items-center justify-center shadow-sm">
                      <Pill class="w-4 h-4 text-slate-400" />
                    </div>
                    <div>
                      <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Tindakan/Terapi</p>
                      <p class="text-xs font-bold text-slate-600 italic leading-relaxed mt-1">{{ item.treatment }}</p>
                    </div>
                  </div>
                </div>

                <div v-if="item.status === 'Resting'" class="mt-6 pt-6 border-t border-slate-100 flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <span class="px-3 py-1 bg-amber-50 text-amber-600 rounded-lg text-[9px] font-black uppercase tracking-widest">Rawat Inap 2 Hari</span>
                    <button class="text-[10px] font-black text-rose-600 uppercase tracking-widest hover:underline flex items-center gap-2">
                      <FileText class="w-3 h-3" /> Lihat Surat Sakit
                    </button>
                  </div>
                  <button class="text-slate-300 hover:text-rose-600 transition-colors">
                    <ExternalLink class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <div v-if="!medicalHistory.length" class="py-20 text-center bg-slate-50/50 rounded-3xl border-2 border-dashed border-slate-200">
              <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center mx-auto mb-4 shadow-sm text-slate-300">
                <History class="w-8 h-8" />
              </div>
              <p class="text-slate-400 font-bold italic uppercase tracking-widest text-xs leading-relaxed">
                Belum ada riwayat penyakit<br/>yang tercatat.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
