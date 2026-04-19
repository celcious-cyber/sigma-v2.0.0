<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Search, Filter, 
  Edit, Trash2, X, Plus, 
  FileSpreadsheet, FileText, Upload, Users, Mail, Phone, Camera, GraduationCap, Download
} from 'lucide-vue-next'
import SigmabaseSidebar from '../../components/SigmabaseSidebar.vue'
import axios from 'axios'
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

// Config API
const API_URL = '/api/v1/admin/base/alumni'

// State
const alumniList = ref<any[]>([])
const loading = ref(true)
const search = ref('')
const isYearDropdownOpen = ref(false)
const selectedYear = ref('')
const showDataMenu = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
  window.addEventListener('click', (e: MouseEvent) => {
    const target = e.target as HTMLElement
    if (!target.closest('.year-dropdown-container')) {
      isYearDropdownOpen.value = false
    }
    if (!target.closest('.data-menu-container')) {
      showDataMenu.value = false
    }
  })
})

const uniqueYears = computed(() => {
  const years = alumniList.value.map(a => a.graduation_year).filter(y => y)
  return [...new Set(years)].sort((a, b) => String(b).localeCompare(String(a)))
})
const showModal = ref(false)
const modalMode = ref('create') // 'create' or 'edit'
const currentId = ref<number | null>(null)

const form = ref({
  name: '',
  gender: 'L',
  nis: '',
  nik: '',
  birth_place: '',
  birth_date: '',
  graduation_year: new Date().getFullYear().toString(),
  batch: '',
  email: '',
  address: '',
  whatsapp: '',
  service_status: 'Tidak Mengabdi',
  photo: ''
})

const photoFile = ref<File | null>(null)
const photoPreview = ref<string | null>(null)

// Methods
const fetchAlumni = async () => {
  loading.value = true
  try {
    const response = await axios.get(API_URL)
    alumniList.value = response.data
  } catch (error) {
    console.error('Error fetching alumni:', error)
  } finally {
    loading.value = false
  }
}

const filteredAlumni = computed(() => {
  return alumniList.value.filter(item => {
    const nameMatch = item.name ? item.name.toLowerCase().includes(search.value.toLowerCase()) : false;
    const yearMatch = item.graduation_year ? String(item.graduation_year).includes(search.value) : false;
    const batchMatch = item.batch ? item.batch.toLowerCase().includes(search.value.toLowerCase()) : false;
    
    if (selectedYear.value && String(item.graduation_year) !== String(selectedYear.value)) {
      return false
    }

    return nameMatch || yearMatch || batchMatch;
  })
})

const handlePhotoUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    photoFile.value = target.files[0]
    photoPreview.value = URL.createObjectURL(target.files[0])
  }
}

const openModal = (mode = 'create', data: any = null) => {
  modalMode.value = mode
  if (mode === 'edit' && data) {
    currentId.value = data.ID || data.id
    form.value = { ...data, birth_date: data.birth_date ? data.birth_date.split('T')[0] : '' }
    photoPreview.value = data.photo || null
  } else {
    currentId.value = null
    form.value = {
      name: '', gender: 'L', nis: '', nik: '', birth_place: '', birth_date: '',
      graduation_year: new Date().getFullYear().toString(), batch: '',
      email: '', address: '', whatsapp: '', service_status: 'Tidak Mengabdi', photo: ''
    }
    photoPreview.value = null
  }
  showModal.value = true
}

const handleSubmit = async () => {
  try {
    let response
    if (modalMode.value === 'create') {
      response = await axios.post(API_URL, form.value)
    } else {
      response = await axios.put(`${API_URL}/${currentId.value}`, form.value)
    }

    // Photo upload if selected
    if (photoFile.value && response.data.ID) {
      const photoFormData = new FormData()
      photoFormData.append('photo', photoFile.value)
      await axios.post(`${API_URL}/${response.data.ID}/photo`, photoFormData)
    }

    showModal.value = false
    fetchAlumni()
  } catch (error: any) {
    console.error('Error saving alumni:', error)
    alert('Gagal menyimpan data alumni: ' + (error.response?.data?.error || error.message))
  }
}

const deleteAlumni = async (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus data alumni ini?')) {
    try {
      await axios.delete(`${API_URL}/${id}`)
      fetchAlumni()
    } catch (error) {
      console.error('Error deleting alumni:', error)
      alert('Gagal menghapus data')
    }
  }
}

// Export functions
const exportToExcel = () => {
  const data = alumniList.value.map(a => ({
    Nama: a.name,
    Gender: a.gender === 'L' ? 'Laki-laki' : 'Perempuan',
    NIS: a.nis,
    'Tahun Lulus': a.graduation_year,
    Angkatan: a.batch,
    Status: a.service_status,
    Email: a.email,
    WhatsApp: a.whatsapp
  }))
  const ws = XLSX.utils.json_to_sheet(data)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, "Alumni")
  XLSX.writeFile(wb, `Data_Alumni_${new Date().getFullYear()}.xlsx`)
  showDataMenu.value = false
}

const exportToPDF = () => {
  const doc = new jsPDF()
  doc.text("DATA ALUMNI SIGMA", 14, 15)
  const data = alumniList.value.map(a => [
    a.name, a.graduation_year, a.batch, a.service_status, a.whatsapp
  ])
  autoTable(doc, {
    head: [['Nama', 'Tahun', 'Angkatan', 'Status', 'WhatsApp']],
    body: data,
    startY: 20
  })
  doc.save(`Data_Alumni_${new Date().getFullYear()}.pdf`)
  showDataMenu.value = false
}

const downloadTemplate = () => {
  const headers = [
    ['Nama Lengkap', 'L/P', 'NIS', 'NIK', 'Tahun Lulus', 'Angkatan', 'Tempat Lahir', 'Tanggal Lahir', 'Alamat', 'Status Pengabdian', 'Email', 'WhatsApp']
  ]
  const worksheet = XLSX.utils.aoa_to_sheet(headers)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, "Template")
  XLSX.writeFile(workbook, "Template_Import_Alumni.xlsx")
  showDataMenu.value = false
}

const triggerImport = () => {
  fileInput.value?.click()
  showDataMenu.value = false
}

const handleImport = async (event: any) => {
  const file = event.target.files[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = async (e) => {
    const data = new Uint8Array(e.target?.result as ArrayBuffer)
    const workbook = XLSX.read(data, { type: 'array' })
    const firstSheetName = workbook.SheetNames[0]
    const worksheet = workbook.Sheets[firstSheetName]
    const rawData = XLSX.utils.sheet_to_json(worksheet)
    
    const formattedData = rawData.map((row: any) => ({
      name: row['Nama Lengkap'] || row['Nama'],
      gender: row['L/P'] || row['Jenis Kelamin'] || row['Gender'] || 'L',
      nis: String(row['NIS'] || ''),
      nik: String(row['NIK'] || ''),
      graduation_year: String(row['Tahun Lulus'] || ''),
      batch: String(row['Angkatan'] || ''),
      birth_place: row['Tempat Lahir'] || '',
      birth_date: row['Tanggal Lahir'] || '',
      address: row['Alamat'] || '',
      service_status: row['Status Pengabdian'] || 'Tidak Mengabdi',
      email: row['Email'] || '',
      whatsapp: String(row['WhatsApp'] || '')
    }))

    try {
      loading.value = true
      await Promise.all(formattedData.map((dataItem: any) => axios.post(API_URL, dataItem)))
      await fetchAlumni()
      alert('Import berhasil!')
    } catch (err: any) {
      console.error('Import gagal:', err)
      alert('Gagal mengimpor data. Pastikan NIS/NIK tidak duplikat.')
    } finally {
      loading.value = false
      if (fileInput.value) fileInput.value.value = ''
    }
  }
  reader.readAsArrayBuffer(file)
}

onMounted(fetchAlumni)
</script>

<template>
  <div class="flex h-screen bg-[#020617] text-slate-200 overflow-hidden font-sans selection:bg-sigma-emerald/30">
    <SigmabaseSidebar activeItem="Data Alumni" />

    <main class="flex-1 overflow-y-auto custom-scrollbar flex flex-col">
      <!-- Sticky Header Container -->
      <div class="sticky top-0 z-40 bg-[#020617]/80 backdrop-blur-xl border-b border-slate-800/60 px-8 py-6 space-y-6 shadow-sm">
        <!-- Page Title & Stats -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-400">
              <GraduationCap class="w-6 h-6" />
            </div>
            <div>
              <h2 class="text-3xl font-black tracking-tight italic">DATABASE <span class="text-sigma-emerald not-italic uppercase text-2xl">ALUMNI</span></h2>
              <p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-1 hidden md:block">Database Lulusan & Pengabdian</p>
            </div>
          </div>

          <div class="text-[10px] text-slate-400 font-bold uppercase tracking-[0.2em] px-4 py-2 bg-slate-900 rounded-full border border-slate-800 hidden lg:block">
            Verified: {{ filteredAlumni.length }} alumni
          </div>
        </div>

        <div class="flex flex-col lg:flex-row gap-4 items-center">
          <!-- Search Input -->
          <div class="relative flex-1 w-full lg:w-auto">
            <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
            <input 
              v-model="search"
              type="text" 
              placeholder="Cari nama alumni, angkatan, atau tahun lulus..."
              class="w-full pl-14 pr-6 py-4 rounded-2xl bg-slate-900 border border-slate-800 focus:border-sigma-emerald outline-none transition-all placeholder:text-slate-500 text-sm font-medium"
            >
          </div>
          
          <!-- Consolidated Utilities -->
          <div class="flex flex-wrap items-center gap-3 w-full lg:w-auto">
            <!-- Year Filter -->
            <div class="relative year-dropdown-container">
              <button 
                @click.stop="isYearDropdownOpen = !isYearDropdownOpen"
                class="flex items-center gap-4 pl-6 pr-10 py-4 bg-slate-900 border border-slate-800 rounded-2xl text-xs font-bold uppercase tracking-widest hover:border-sigma-emerald/50 transition-all shadow-sm text-slate-300 min-w-[180px] justify-between group"
              >
                <div class="flex items-center gap-3">
                  <Filter class="w-4 h-4 text-slate-500 group-hover:text-sigma-emerald transition-colors" />
                  <span>{{ selectedYear ? 'Tahun ' + selectedYear : 'Semua Tahun' }}</span>
                </div>
                <div class="w-1 h-1 rounded-full bg-slate-500"></div>
              </button>

              <div v-if="isYearDropdownOpen" class="absolute top-full left-0 mt-3 w-full bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl z-50 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                <div class="max-h-60 overflow-y-auto custom-scrollbar">
                  <button @click="selectedYear = ''; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-slate-800" :class="selectedYear === '' ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-slate-400'">Semua Tahun</button>
                  <button v-for="year in uniqueYears" :key="year" @click="selectedYear = year as string; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-slate-800 last:border-0" :class="selectedYear === year ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-slate-400'">Tahun {{ year }}</button>
                </div>
              </div>
            </div>

            <input type="file" ref="fileInput" class="hidden" accept=".xlsx, .xls" @change="handleImport" />
            
            <!-- Data Operations Dropdown -->
            <div class="relative data-menu-container">
              <button 
                @click.stop="showDataMenu = !showDataMenu"
                class="flex items-center gap-2 px-6 py-4 bg-slate-900 border border-slate-800 rounded-2xl hover:bg-slate-800 transition-all text-xs font-bold uppercase tracking-widest text-slate-400 hover:text-slate-200 shadow-sm"
              >
                <Download class="w-4 h-4 text-sigma-emerald" /> Export/Import
              </button>
              
              <div v-if="showDataMenu" class="absolute right-0 mt-3 w-56 bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl z-50 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                <div class="p-2 border-b border-slate-800 bg-slate-800/30">
                  <button @click="downloadTemplate" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-sigma-emerald hover:bg-sigma-emerald/10 rounded-xl transition-all"><Download class="w-4 h-4" /> Template Excel</button>
                  <button @click="triggerImport" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-slate-400 hover:bg-slate-800 hover:text-slate-200 rounded-xl transition-all mt-1"><Upload class="w-4 h-4" /> Import Data</button>
                </div>
                <div class="p-2">
                  <button @click="exportToExcel" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-slate-400 hover:bg-slate-800 hover:text-sigma-emerald rounded-xl transition-all"><FileSpreadsheet class="w-4 h-4" /> Export Excel</button>
                  <button @click="exportToPDF" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-slate-400 hover:bg-slate-800 hover:text-rose-500 rounded-xl transition-all mt-1"><FileText class="w-4 h-4" /> Laporan PDF</button>
                </div>
              </div>
            </div>

            <button @click="openModal('create')" class="flex items-center gap-2 px-8 py-4 bg-sigma-emerald text-white rounded-2xl font-bold uppercase tracking-widest text-xs hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20">
              <Plus class="w-4 h-4" /> Tambah
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content Scrolling Area -->
      <div class="p-8 space-y-10">
        <!-- Alumni Table -->
        <div class="bg-slate-900/40 border border-slate-800/60 rounded-[2.5rem] overflow-hidden backdrop-blur-xl shadow-sm">
          <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-800/20 border-b border-slate-800/60">
                <th class="px-8 py-5 text-xs font-bold text-slate-500 uppercase tracking-widest">Identitas Alumni</th>
                <th class="px-8 py-5 text-xs font-bold text-slate-500 uppercase tracking-widest text-center">Tahun Lulus</th>
                <th class="px-8 py-5 text-xs font-bold text-slate-500 uppercase tracking-widest text-center">Angkatan</th>
                <th class="px-8 py-5 text-xs font-bold text-slate-500 uppercase tracking-widest">Status Pengabdian</th>
                <th class="px-8 py-5 text-xs font-bold text-slate-500 uppercase tracking-widest text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/40">
              <tr v-for="alumni in filteredAlumni" :key="alumni.ID || alumni.id" class="hover:bg-slate-800/20 transition-colors group">
                <td class="px-8 py-5">
                  <div class="flex items-center gap-5">
                    <div class="w-14 h-14 rounded-xl overflow-hidden bg-slate-800 ring-4 ring-slate-800/30">
                      <img v-if="alumni.photo" :src="alumni.photo" class="w-full h-full object-cover" />
                      <div v-else class="w-full h-full flex items-center justify-center text-slate-600">
                        <Users class="w-6 h-6" />
                      </div>
                    </div>
                    <div class="space-y-1">
                      <div class="flex items-center gap-2">
                        <span class="font-bold text-white tracking-tight text-lg">{{ alumni.name }}</span>
                        <span class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-wider" 
                              :class="alumni.gender === 'L' ? 'bg-blue-500/10 text-blue-500' : 'bg-pink-500/10 text-pink-500'">
                          {{ alumni.gender }}
                        </span>
                        <span v-if="alumni.nik" class="px-2 py-0.5 rounded bg-sigma-emerald/10 text-sigma-emerald text-[9px] font-black uppercase tracking-tighter border border-sigma-emerald/20">
                          NIK: {{ alumni.nik }}
                        </span>
                      </div>
                      <div class="flex items-center gap-4 text-xs font-medium text-slate-400">
                        <span class="flex items-center gap-1.5"><Mail class="w-3.5 h-3.5 text-sigma-emerald" /> {{ alumni.email || '-' }}</span>
                        <span class="flex items-center gap-1.5"><Phone class="w-3.5 h-3.5 text-emerald-400" /> {{ alumni.whatsapp || '-' }}</span>
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-5">
                   <div class="flex flex-col items-center">
                    <span class="px-4 py-1.5 bg-slate-800/50 rounded-full text-sm font-bold text-white border border-slate-700/30">
                      {{ alumni.graduation_year }}
                    </span>
                  </div>
                </td>
                <td class="px-8 py-5">
                  <div class="flex flex-col items-center">
                    <span class="text-xs font-black text-sigma-emerald tracking-widest uppercase">{{ alumni.batch || 'GENERAL' }}</span>
                  </div>
                </td>
                <td class="px-8 py-5">
                  <div class="flex flex-col">
                    <span class="text-sm font-bold" 
                          :class="alumni.service_status === 'Mengabdi' ? 'text-emerald-400' : alumni.service_status === 'Tidak Mengabdi' ? 'text-slate-400' : 'text-amber-400'">
                      {{ alumni.service_status }}
                    </span>
                    <span class="text-[10px] text-slate-500 uppercase tracking-tighter">{{ alumni.address?.substring(0, 30) }}...</span>
                  </div>
                </td>
                <td class="px-8 py-5 text-center">
                  <div class="flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button @click="openModal('edit', alumni)" class="w-8 h-8 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center hover:bg-emerald-500 hover:text-white transition-colors" title="Edit">
                      <Edit class="w-4 h-4" />
                    </button>
                    <button @click="deleteAlumni(alumni.ID || alumni.id)" class="w-8 h-8 rounded-lg bg-rose-500/10 text-rose-400 flex items-center justify-center hover:bg-rose-500 hover:text-white transition-colors" title="Delete">
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      </div>
    </main>

    <!-- Modal Alumni -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-6 backdrop-blur-sm bg-black/60">
      <div class="bg-slate-900 border border-slate-800 w-full max-w-4xl rounded-3xl overflow-hidden shadow-2xl">
        <div class="px-8 py-6 bg-slate-800/50 border-b border-slate-800 flex items-center justify-between">
          <h2 class="text-2xl font-black text-white uppercase italic italic-sigma">
            {{ modalMode === 'create' ? 'TAMBAH' : 'EDIT' }} <span class="text-sigma-emerald">ALUMNI</span>
          </h2>
          <button @click="showModal = false" class="p-2 hover:bg-slate-700/50 rounded-lg transition-colors text-slate-400">
            <X class="w-6 h-6" />
          </button>
        </div>

        <div class="p-8 max-h-[75vh] overflow-y-auto custom-scrollbar">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <!-- Identity Section -->
            <div class="space-y-6">
              <h3 class="text-xs font-black text-sigma-emerald tracking-widest uppercase mb-4 py-1 border-b border-sigma-emerald/20">Identitas Pribadi</h3>
              
              <div class="space-y-2">
                <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Nama Lengkap</label>
                <input v-model="form.name" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="Nama Alumni...">
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Jenis Kelamin</label>
                  <select v-model="form.gender" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all">
                    <option value="L">Laki-laki</option>
                    <option value="P">Perempuan</option>
                  </select>
                </div>
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">NIK</label>
                  <input v-model="form.nik" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="NIK Nasional...">
                </div>
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">NIS</label>
                  <input v-model="form.nis" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="NIS...">
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Tempat Lahir</label>
                  <input v-model="form.birth_place" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="Kota...">
                </div>
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Tanggal Lahir</label>
                  <input v-model="form.birth_date" type="date" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all">
                </div>
              </div>

              <div class="space-y-2">
                <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Alamat Domisili</label>
                <textarea v-model="form.address" rows="3" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all resize-none" placeholder="Alamat lengkap..."></textarea>
              </div>
            </div>

            <!-- Academic & Contact Section -->
            <div class="space-y-6">
              <h3 class="text-xs font-black text-sigma-emerald tracking-widest uppercase mb-4 py-1 border-b border-sigma-emerald/20">Lulusan & Kontak</h3>
              
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Tahun Lulus</label>
                  <input v-model="form.graduation_year" type="number" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all">
                </div>
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Angkatan (Batch)</label>
                  <input v-model="form.batch" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="Nama Angkatan...">
                </div>
              </div>

              <div class="space-y-2">
                <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Status Pengabdian</label>
                <select v-model="form.service_status" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all">
                  <option value="Tidak Mengabdi">Tidak Mengabdi</option>
                  <option value="Mengabdi">Mengabdi</option>
                  <option value="Lainnya">Lainnya</option>
                </select>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Email</label>
                  <input v-model="form.email" type="email" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="example@mail.com">
                </div>
                <div class="space-y-2">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">WhatsApp</label>
                  <input v-model="form.whatsapp" type="text" class="w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-xl focus:border-sigma-emerald outline-none transition-all" placeholder="08...">
                </div>
              </div>

              <!-- Photo Upload -->
              <div class="space-y-2">
                <label class="text-xs font-bold text-slate-400 uppercase tracking-wider ml-1">Foto Alumni</label>
                <div class="flex items-center gap-4 p-4 bg-slate-950 border border-dashed border-slate-800 rounded-2xl">
                  <div class="w-20 h-20 rounded-xl bg-slate-900 overflow-hidden flex-shrink-0 border border-slate-800">
                    <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover">
                    <div v-else class="w-full h-full flex items-center justify-center text-slate-700">
                      <Camera class="w-8 h-8" />
                    </div>
                  </div>
                  <div class="space-y-1">
                    <label class="cursor-pointer px-4 py-1.5 bg-sigma-emerald/10 hover:bg-sigma-emerald/20 text-sigma-emerald rounded-lg text-xs font-bold border border-sigma-emerald/20 transition-all inline-block">
                      <Upload class="w-3.5 h-3.5 inline mr-1" /> UPLOAD FOTO
                      <input type="file" @change="handlePhotoUpload" accept="image/*" class="hidden">
                    </label>
                    <p class="text-[10px] text-slate-500">Maksimal 2MB (JPG/PNG)</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="px-8 py-6 bg-slate-800/30 border-t border-slate-800 flex justify-end gap-3">
          <button @click="showModal = false" class="px-6 py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl font-bold text-sm transition-all border border-slate-700/50">
            BATAL
          </button>
          <button @click="handleSubmit" class="px-8 py-2.5 bg-gradient-to-r from-sigma-emerald to-emerald-600 hover:from-emerald-500 hover:to-emerald-500 text-white rounded-xl shadow-lg shadow-sigma-emerald/20 font-black text-sm transition-all">
            SIMPAN DATA
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.italic-sigma {
  transform: skewX(-10deg);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #1e293b;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #334155;
}
</style>
