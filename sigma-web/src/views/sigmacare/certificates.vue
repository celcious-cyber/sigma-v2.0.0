<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  FileText, CheckCircle2,
  Printer, Download, Send, 
  Plus, Search, User,
  MoreVertical, ShieldCheck
} from 'lucide-vue-next'
import axios from 'axios'

const loading = ref(true)
const showIssueModal = ref(false)

const certificates = ref<any[]>([])
const visits = ref<any[]>([])

const fetchCertificates = async () => {
  try {
    const response = await axios.get('/admin/care/visits')
    visits.value = response.data
    // For now, let's say all resting/referred visits can have certificates
    certificates.value = response.data.filter((v: any) => v.status !== 'Treated')
  } catch (error) {
    console.error('Failed to fetch certificates data:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCertificates()
})
</script>

<template>
  <div class="p-8 space-y-8 min-h-screen bg-slate-50/50">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <h1 class="text-3xl font-black tracking-tight text-slate-800 italic uppercase">Sertifikasi <span class="text-rose-600">Medis</span></h1>
        <p class="text-slate-500 font-bold mt-1 tracking-wide uppercase text-xs flex items-center gap-2">
          <FileText class="w-4 h-4" /> Penerbitan Surat Keterangan Sakit & Sinkronisasi Absensi
        </p>
      </div>
      
      <button 
        @click="showIssueModal = true"
        class="px-8 py-4 bg-rose-600 text-white rounded-2xl font-black shadow-lg shadow-rose-500/20 hover:bg-rose-700 active:scale-95 transition-all flex items-center gap-3"
      >
        <Plus class="w-5 h-5" /> <span class="text-xs uppercase tracking-widest">Terbitkan Surat Baru</span>
      </button>
    </div>

    <!-- Active Certificates Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div v-for="i in 2" :key="i" class="group bg-white rounded-[2.5rem] border border-slate-100 shadow-sm hover:shadow-xl hover:shadow-rose-500/5 transition-all p-8 relative overflow-hidden">
        <div class="absolute -right-10 -top-10 w-40 h-40 bg-rose-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-700 flex items-center justify-center">
          <FileText class="w-20 h-20 text-rose-100 rotate-12" />
        </div>
        
        <div class="relative z-10 space-y-6">
          <div class="flex justify-between items-start">
            <div class="space-y-1">
              <span class="px-3 py-1 bg-rose-50 text-rose-600 rounded-lg text-[10px] font-black uppercase tracking-widest border border-rose-100">SURAT KETERANGAN SAKIT</span>
              <h3 class="text-xl font-black text-slate-800 uppercase tracking-tighter mt-2">Muhammad Akmal</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-[0.2em]">REF: UKS/SK/2026/04/082</p>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 bg-rose-500 rounded-full animate-pulse"></div>
              <span class="text-[10px] font-black text-rose-600 uppercase tracking-widest">Aktif</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-6 py-6 border-y border-dashed border-slate-100">
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Mulai Istirahat</p>
              <p class="text-sm font-black text-slate-700">28 Apr 2026</p>
            </div>
            <div>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Berakhir Pada</p>
              <p class="text-sm font-black text-slate-700">30 Apr 2026</p>
            </div>
          </div>

          <div class="flex items-center justify-between pt-2">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-emerald-50 text-emerald-600 rounded-lg flex items-center justify-center">
                <ShieldCheck class="w-4 h-4" />
              </div>
              <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest italic">Tersinkron ke Absensi Akademik</span>
            </div>
            <div class="flex gap-2">
              <button class="p-3 bg-slate-50 text-slate-400 rounded-xl hover:bg-rose-50 hover:text-rose-600 transition-all shadow-sm">
                <Printer class="w-4 h-4" />
              </button>
              <button class="p-3 bg-slate-50 text-slate-400 rounded-xl hover:bg-rose-50 hover:text-rose-600 transition-all shadow-sm">
                <Download class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- History Table -->
    <div class="bg-white rounded-[2.5rem] border border-slate-100 shadow-sm overflow-hidden">
      <div class="p-8 border-b border-slate-50 flex items-center justify-between">
        <h2 class="text-lg font-black text-slate-800 uppercase italic">Arsip Sertifikasi Medis</h2>
        <div class="flex items-center gap-4">
          <div class="relative">
            <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input type="text" placeholder="Cari No. Referensi..." class="pl-11 pr-6 py-2.5 bg-slate-50 border border-slate-100 rounded-xl text-xs font-bold" />
          </div>
        </div>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Nomor Referensi</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Santri</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Jenis Surat</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Periode</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest">Status</th>
              <th class="px-8 py-5 font-black text-slate-400 uppercase tracking-widest text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50 text-slate-600">
            <tr v-for="cert in certificates" :key="cert.ID" class="hover:bg-slate-50 transition-colors group">
              <td class="px-8 py-6 font-black text-xs tracking-widest uppercase">UKS/SK/{{ new Date(cert.CreatedAt).getFullYear() }}/{{ (new Date(cert.CreatedAt).getMonth()+1).toString().padStart(2, '0') }}/{{ cert.ID.toString().padStart(3, '0') }}</td>
              <td class="px-8 py-6">
                <div class="flex items-center gap-3 font-black text-slate-700 uppercase tracking-tight">
                  <User class="w-4 h-4 text-slate-400" /> {{ cert.student?.full_name }}
                </div>
              </td>
              <td class="px-8 py-6">
                <span class="px-3 py-1 bg-slate-100 text-slate-500 rounded-lg text-[10px] font-black uppercase italic">{{ cert.status }}</span>
              </td>
              <td class="px-8 py-6 font-bold italic">{{ new Date(cert.visit_date).toLocaleDateString('id-ID') }} - Selesai</td>
              <td class="px-8 py-6">
                <span class="text-emerald-500 font-black text-[10px] uppercase tracking-widest flex items-center gap-2">
                  <CheckCircle2 class="w-4 h-4" /> Aktif
                </span>
              </td>
              <td class="px-8 py-6 text-right">
                <button class="p-2 text-slate-300 hover:text-rose-600 transition-colors">
                  <MoreVertical class="w-5 h-5" />
                </button>
              </td>
            </tr>
            <tr v-if="!certificates.length && !loading">
              <td colspan="6" class="px-8 py-20 text-center">
                <p class="text-slate-400 font-bold italic">Belum ada arsip sertifikasi medis.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Issue Certificate -->
    <div v-if="showIssueModal" class="fixed inset-0 z-[100] flex items-center justify-center p-6 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300">
        <div class="p-8 border-b border-slate-100 bg-rose-600 text-white flex items-center gap-4">
          <div class="w-12 h-12 bg-white/20 rounded-2xl flex items-center justify-center">
            <FileText class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-xl font-black italic uppercase">Terbitkan Surat</h3>
            <p class="text-xs font-bold text-white/70 uppercase tracking-widest">Validasi Izin Medis Santri</p>
          </div>
        </div>

        <div class="p-10 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 italic">Pilih Riwayat Kunjungan</label>
            <select class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:border-rose-500 outline-none">
              <option value="">Cari berdasarkan kunjungan terakhir...</option>
              <option>Kunjungan Akmal - Demam Tinggi (28 Apr)</option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 italic">Mulai Tanggal</label>
              <input type="date" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:border-rose-500 outline-none" />
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 italic">Sampai Tanggal</label>
              <input type="date" class="w-full px-5 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-sm font-bold focus:border-rose-500 outline-none" />
            </div>
          </div>

          <div class="bg-emerald-50 p-6 rounded-3xl border-2 border-dashed border-emerald-100 flex items-center gap-4">
            <ShieldCheck class="w-8 h-8 text-emerald-500" />
            <p class="text-[10px] font-bold text-emerald-700 uppercase tracking-wide leading-relaxed">
              Menerbitkan surat ini akan otomatis mengubah status kehadiran santri menjadi <span class="font-black italic">"SAKIT"</span> pada modul Absensi Akademik.
            </p>
          </div>
        </div>

        <div class="p-8 bg-slate-50 border-t border-slate-100 flex gap-4">
          <button @click="showIssueModal = false" class="flex-1 py-4 bg-white border border-slate-200 rounded-2xl font-black text-xs uppercase tracking-widest">Batal</button>
          <button class="flex-[2] py-4 bg-rose-600 text-white rounded-2xl font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-rose-500/20 flex items-center justify-center gap-3">
            <Send class="w-4 h-4" /> Proses Penerbitan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
