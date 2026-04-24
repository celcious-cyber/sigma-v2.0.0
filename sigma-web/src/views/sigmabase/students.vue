<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Search, Filter, 
  Download, Edit, Trash2, X, Plus, GraduationCap, Camera, 
  FileSpreadsheet, FileText, Upload
} from 'lucide-vue-next'
import axios from 'axios'
import * as XLSX from 'xlsx'
import jsPDF from 'jspdf'
import 'jspdf-autotable'

// Config API
const API_URL = '/api/v1/admin/base/students'

// State
const isModalOpen = ref(false)
const isEditing = ref(false)
const currentId = ref<number | null>(null)
const searchQuery = ref('')
const isLoading = ref(false)
const isYearDropdownOpen = ref(false)
const selectedYear = ref('')
const showDataMenu = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

// Handle click outside to close dropdowns
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

// Computed: Unique Years for Filter
const uniqueYears = computed(() => {
  const years = students.value.map(s => s.entry_year).filter(y => y)
  return [...new Set(years)].sort((a, b) => b.localeCompare(a))
})

// Form State
const form = ref({
  name: '',
  gender: 'L',
  birth_place: '',
  birth_date: '',
  address: '',
  parent_name: '',
  parent_phone: '',
  nis: '',
  nik: '',
  nisn: '',
  class: '',
  entry_year: ''
})

// Graduation State
const showGradModal = ref(false)
const gradStudentId = ref<number | null>(null)
const gradForm = ref({
  graduation_year: new Date().getFullYear().toString(),
  batch: '',
  email: '',
  whatsapp: '',
  service_status: 'Tidak Mengabdi',
  other_status: ''
})
const gradPhotoFile = ref<File | null>(null)
const gradPhotoPreview = ref<string | null>(null)

// Database Data
const students = ref<any[]>([])

const fetchStudents = async () => {
  try {
    isLoading.value = true
    const response = await axios.get(API_URL)
    students.value = response.data
  } catch (err: any) {
    console.error('Gagal mengambil data:', err)
  } finally {
    isLoading.value = false
  }
}

// Client-side Search & Year Filtering
const filteredStudents = computed(() => {
  let result = students.value
  
  // Filter by Search Query
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(s => 
      s.name.toLowerCase().includes(q) || 
      s.nis.toLowerCase().includes(q)
    )
  }
  
  // Filter by Selected Year
  if (selectedYear.value) {
    result = result.filter(s => s.entry_year === selectedYear.value)
  }
  
  return result
})

const openNewModal = () => {
  isEditing.value = false
  currentId.value = null
  form.value = {
    name: '',
    gender: 'L',
    birth_place: '',
    birth_date: '',
    address: '',
    parent_name: '',
    parent_phone: '',
    nis: '',
    nik: '',
    nisn: '',
    class: '',
    entry_year: ''
  }
  isModalOpen.value = true
}

const handleEdit = (student: any) => {
  isEditing.value = true
  currentId.value = student.ID
  form.value = {
    name: student.name,
    gender: student.gender || 'L',
    birth_place: student.birth_place,
    birth_date: student.birth_date ? student.birth_date.split('T')[0] : '',
    address: student.address,
    parent_name: student.parent_name,
    parent_phone: student.parent_phone,
    nis: student.nis,
    nik: student.nik || '',
    nisn: student.nisn,
    class: student.class,
    entry_year: student.entry_year
  }
  isModalOpen.value = true
}

const handleDelete = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus data santri ini?')) return
  
  try {
    isLoading.value = true
    await axios.delete(`${API_URL}/${id}`)
    await fetchStudents()
  } catch (err) {
    alert('Gagal menghapus data.')
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const openGradModal = (student: any) => {
  gradStudentId.value = student.ID
  gradForm.value = {
    graduation_year: new Date().getFullYear().toString(),
    batch: student.entry_year || '',
    email: '',
    whatsapp: student.parent_phone || '',
    service_status: 'Tidak Mengabdi',
    other_status: ''
  }
  gradPhotoPreview.value = null
  gradPhotoFile.value = null
  showGradModal.value = true
}

const handleGradPhoto = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    gradPhotoFile.value = target.files[0]
    gradPhotoPreview.value = URL.createObjectURL(target.files[0])
  }
}

const handleGraduation = async () => {
  try {
    isLoading.value = true
    
    // Prepare data
    const finalStatus = gradForm.value.service_status === 'Lainnya' 
      ? gradForm.value.other_status 
      : gradForm.value.service_status

    const payload = {
      ...gradForm.value,
      service_status: finalStatus
    }

    const response = await axios.post(`${API_URL}/${gradStudentId.value}/graduate`, payload)
    
    // Upload photo to the NEW Alumni record if exist
    if (gradPhotoFile.value && response.data.alumni_id) {
      const photoFormData = new FormData()
      photoFormData.append('photo', gradPhotoFile.value)
      // Call the alumni photo upload endpoint (not students!)
      await axios.post(`/api/v1/admin/base/alumni/${response.data.alumni_id}/photo`, photoFormData)
    }
    
    alert('Santri berhasil diluluskan menjadi Alumni!')
    showGradModal.value = false
    fetchStudents()
  } catch (err: any) {
    alert('Gagal memproses kelulusan: ' + (err.response?.data?.error || err.message))
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchStudents()
})

// Export & Import Logic
const exportToExcel = () => {
  const data = students.value.map(s => ({
    'Nama Lengkap': s.name,
    'L/P': s.gender,
    'NIS': s.nis,
    'NISN': s.nisn,
    'Kelas': s.class,
    'Angkatan': s.entry_year,
    'Tempat Lahir': s.birth_place,
    'Tanggal Lahir': s.birth_date ? s.birth_date.split('T')[0] : '',
    'Alamat': s.address,
    'Orang Tua': s.parent_name,
    'Telepon': s.parent_phone
  }))
  
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, "Data Santri")
  XLSX.writeFile(workbook, `Data_Santri_${new Date().toLocaleDateString()}.xlsx`)
  showDataMenu.value = false
}

const exportToPDF = () => {
  const doc = new jsPDF()
  doc.text("Laporan Data Santri Al-Hikmah", 14, 15)
  
  const tableData: any[][] = students.value.map((s: any) => [
    s.name, s.gender, s.nis, s.class, s.entry_year, s.parent_name, s.parent_phone
  ]);
  
  (doc as any).autoTable({
    startY: 20,
    head: [['Nama', 'L/P', 'NIS', 'Kelas', 'Angkatan', 'Orang Tua', 'Telepon']],
    body: tableData,
    theme: 'grid',
    headStyles: { fillColor: [16, 185, 129] }
  })
  
  doc.save(`Laporan_Santri_${new Date().toLocaleDateString()}.pdf`)
  showDataMenu.value = false
}

const downloadTemplate = () => {
  const headers = [
    ['Nama Lengkap', 'L/P', 'NIS', 'NISN', 'Angkatan', 'Tempat Lahir', 'Tanggal Lahir', 'Alamat', 'Orang Tua', 'Telepon']
  ]
  
  const worksheet = XLSX.utils.aoa_to_sheet(headers)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, "Template")
  XLSX.writeFile(workbook, "Template_Import_Santri.xlsx")
}

const triggerImport = () => {
  fileInput.value?.click()
}

const handleImport = async (event: any) => {
  const file = event.target.files[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = async (e) => {
    const data = new Uint8Array(e.target?.result as ArrayBuffer)
    const workbook = XLSX.read(data, { type: 'array' })
    const worksheet = workbook.Sheets[workbook.SheetNames[0]]
    const jsonData = XLSX.utils.sheet_to_json(worksheet)
    
    // Map Excel headers to StudentDTO
    const formattedData = jsonData.map((row: any) => ({
      name: row['Nama Lengkap'] || row['Nama'],
      gender: row['L/P'] || row['Jenis Kelamin'] || row['Gender'] || 'L',
      nis: String(row['NIS'] || ''),
      nisn: String(row['NISN'] || ''),
      class: row['Kelas'] || '',
      entry_year: String(row['Angkatan'] || row['Tahun'] || ''),
      birth_place: row['Tempat Lahir'] || '',
      birth_date: row['Tanggal Lahir'] || '',
      address: row['Alamat'] || '',
      parent_name: row['Orang Tua'] || row['Wali'],
      parent_phone: String(row['Telepon'] || row['WhatsApp'] || '')
    }))

    try {
      isLoading.value = true
      await axios.post(`${API_URL}/bulk`, formattedData)
      await fetchStudents()
      alert('Import berhasil!')
    } catch (err: any) {
      console.error('Import gagal:', err)
      alert('Gagal mengimpor data. Periksa format file Anda.')
    } finally {
      isLoading.value = false
    }
  }
  reader.readAsArrayBuffer(file)
}

const handleSubmit = async () => {
  try {
    isLoading.value = true
    if (isEditing.value && currentId.value) {
      await axios.put(`${API_URL}/${currentId.value}`, form.value)
    } else {
      await axios.post(API_URL, form.value)
    }
    
    // Refresh data dan tutup modal
    await fetchStudents()
    isModalOpen.value = false
  } catch (err: any) {
    alert('Gagal menyimpan data. Pastikan NIS/NISN belum terdaftar.')
    console.error(err)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="flex flex-col">
    <!-- Sticky Header Container -->
    <div class="sticky top-0 z-40 bg-sigma-app/80 backdrop-blur-xl border-b border-sigma-border px-8 py-6 space-y-6 shadow-sm">
      <!-- Page Title & Stats -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-400">
            <GraduationCap class="w-6 h-6" />
          </div>
          <div>
            <h2 class="text-3xl font-black tracking-tight italic">DATABASE <span class="text-sigma-emerald not-italic uppercase text-2xl">SANTRI</span></h2>
            <p class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest mt-1 hidden md:block">Al-Hikmah Academic Intelligence</p>
          </div>
        </div>

        <div class="text-[10px] text-sigma-muted font-bold uppercase tracking-[0.2em] px-4 py-2 bg-sigma-surface-alt rounded-full border border-sigma-border hidden lg:block">
          Verified: {{ filteredStudents.length }} santri active
        </div>
      </div>

      <div class="flex flex-col lg:flex-row gap-4 items-center">
        <!-- Search Input -->
        <div class="relative flex-1 w-full lg:w-auto">
          <input type="file" ref="fileInput" class="hidden" accept=".xlsx, .xls" @change="handleImport" />
          <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted" />
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari Nama atau NIS Santri..." 
            class="w-full pl-14 pr-6 py-4 rounded-2xl bg-sigma-surface border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/50 text-sm font-medium"
          />
        </div>
        
        <!-- Consolidated Utilities -->
        <div class="flex flex-wrap items-center gap-3 w-full lg:w-auto">
          <!-- Year Filter -->
          <div class="relative year-dropdown-container">
            <button 
              @click.stop="isYearDropdownOpen = !isYearDropdownOpen"
              class="flex items-center gap-4 pl-6 pr-10 py-4 bg-sigma-surface border border-sigma-border rounded-2xl text-xs font-bold uppercase tracking-widest hover:border-sigma-emerald/50 transition-all shadow-sm text-sigma-text min-w-[180px] justify-between group"
            >
              <div class="flex items-center gap-3">
                <Filter class="w-4 h-4 text-sigma-muted group-hover:text-sigma-emerald transition-colors" />
                <span>{{ selectedYear ? 'Angk. ' + selectedYear : 'Semua Angkatan' }}</span>
              </div>
              <div class="w-1 h-1 rounded-full bg-sigma-muted"></div>
            </button>

            <div v-if="isYearDropdownOpen" class="absolute top-full left-0 mt-3 w-full bg-sigma-surface dark:bg-slate-900 border border-sigma-border rounded-2xl shadow-2xl z-50 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
              <div class="max-h-60 overflow-y-auto custom-scrollbar">
                <button @click="selectedYear = ''; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-sigma-border" :class="selectedYear === '' ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-sigma-muted'">Semua Angkatan</button>
                <button v-for="year in uniqueYears" :key="year" @click="selectedYear = year; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-sigma-border last:border-0" :class="selectedYear === year ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-sigma-muted'">Angk. {{ year }}</button>
              </div>
            </div>
          </div>

          <!-- Data Operations (Combined) -->
          <div class="relative data-menu-container">
            <button 
              @click.stop="showDataMenu = !showDataMenu"
              class="flex items-center gap-2 px-6 py-4 bg-sigma-surface border border-sigma-border rounded-2xl hover:bg-sigma-surface-alt transition-all text-xs font-bold uppercase tracking-widest text-sigma-muted hover:text-sigma-text shadow-sm"
            >
              <Download class="w-4 h-4 text-sigma-emerald" /> Export/Import
            </button>
            
            <div v-if="showDataMenu" class="absolute right-0 mt-3 w-56 bg-sigma-surface dark:bg-slate-900 border border-sigma-border rounded-2xl shadow-2xl z-50 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
              <div class="p-2 border-b border-sigma-border bg-sigma-surface-alt/30">
                <button @click="downloadTemplate" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-sigma-emerald hover:bg-sigma-emerald/10 rounded-xl transition-all"><Download class="w-4 h-4" /> Template Excel</button>
                <button @click="triggerImport" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-sigma-muted hover:bg-sigma-surface-alt hover:text-sigma-text rounded-xl transition-all mt-1"><Upload class="w-4 h-4" /> Import Data</button>
              </div>
              <div class="p-2">
                <button @click="exportToExcel" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-sigma-muted hover:bg-sigma-surface-alt hover:text-sigma-emerald rounded-xl transition-all"><FileSpreadsheet class="w-4 h-4" /> Export Excel</button>
                <button @click="exportToPDF" class="w-full flex items-center gap-3 px-4 py-3 text-[10px] font-bold uppercase tracking-widest text-sigma-muted hover:bg-sigma-surface-alt hover:text-red-400 rounded-xl transition-all mt-1"><FileText class="w-4 h-4" /> Laporan PDF</button>
              </div>
            </div>
          </div>

          <!-- Registration Button -->
          <button @click="openNewModal" class="flex items-center gap-2 px-8 py-4 bg-sigma-emerald text-white rounded-2xl font-bold uppercase tracking-widest text-xs hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20">
            <Plus class="w-4 h-4" /> Registrasi
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content Scrolling Area -->
    <div class="p-8 space-y-10">
      <!-- Table Visualization -->
      <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.3em] bg-sigma-surface-alt">
              <th class="p-8 border-b border-sigma-border">Santri / NIS</th>
              <th class="p-8 border-b border-sigma-border text-center">Gender</th>
              <th class="p-8 border-b border-sigma-border text-center">Kelas</th>
              <th class="p-8 border-b border-sigma-border text-center">Angkatan</th>
              <th class="p-8 border-b border-sigma-border">Informasi Wali</th>
              <th class="p-8 border-b border-sigma-border">Status</th>
              <th class="p-8 border-b border-sigma-border text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border">
            <tr v-for="student in filteredStudents" :key="student.ID" class="hover:bg-sigma-surface-alt transition-all group">
              <td class="p-8">
                <div class="flex items-center gap-5">
                  <div class="w-12 h-12 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center font-black text-sigma-emerald text-xl shadow-sm group-hover:scale-110 transition-transform">
                    {{ student.name.charAt(0) }}
                  </div>
                  <div>
                    <div class="font-bold text-sigma-text text-lg group-hover:text-sigma-emerald transition-colors">{{ student.name }}</div>
                    <div class="text-[11px] font-mono text-sigma-muted flex items-center gap-2 mt-1">
                      <span class="px-1.5 py-0.5 rounded-md bg-sigma-surface border border-sigma-border" title="NIS">{{ student.nis }}</span>
                      <span v-if="student.nik" class="px-1.5 py-0.5 rounded-md bg-sigma-emerald/10 border border-sigma-emerald/20 text-sigma-emerald" title="NIK">{{ student.nik }}</span>
                      <span class="text-sigma-border">|</span>
                      <span title="NISN">{{ student.nisn }}</span>
                    </div>
                  </div>
                </div>
              </td>
              <td class="p-8 text-center text-xs font-bold text-sigma-muted">
                <span :class="student.gender === 'P' ? 'text-pink-400 font-black' : 'text-blue-400 font-black'">
                  {{ student.gender }}
                </span>
              </td>
              <td class="p-8 text-center">
                <span class="inline-block px-4 py-2 rounded-xl bg-sigma-surface-alt border border-sigma-border font-black text-sm text-sigma-muted group-hover:text-sigma-text transition-all">
                  {{ student.class }}
                </span>
              </td>
              <td class="p-8 text-center">
                <span class="inline-block px-4 py-2 rounded-xl bg-sigma-surface-alt border border-sigma-border font-black text-xs text-sigma-muted group-hover:text-sigma-text transition-all">
                  {{ student.entry_year }}
                </span>
              </td>
              <td class="p-8">
                <div class="text-sm font-bold text-sigma-text/80">{{ student.parent_name }}</div>
                <div class="text-[10px] text-sigma-muted mt-1 uppercase tracking-widest font-bold">{{ student.parent_phone }}</div>
              </td>
              <td class="p-8">
                <span class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-sigma-emerald/10 text-sigma-emerald text-[10px] font-black uppercase tracking-widest border border-sigma-emerald/5">
                  <span class="w-1.5 h-1.5 rounded-full bg-sigma-emerald animate-pulse"></span> Aktif
                </span>
              </td>
              <td class="p-8 text-right">
                <div class="flex justify-end gap-3 translate-x-4 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-500">
                  <button @click="openGradModal(student)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-sigma-emerald text-sigma-muted hover:text-white rounded-xl transition-all shadow-md hover:shadow-sigma-emerald/20" title="Luluskan"><GraduationCap class="w-4 h-4" /></button>
                  <button @click="handleEdit(student)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-sigma-emerald text-sigma-muted hover:text-white rounded-xl transition-all shadow-md hover:shadow-sigma-emerald/20"><Edit class="w-4 h-4" /></button>
                  <button @click="handleDelete(student.ID)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-red-500 text-sigma-muted hover:text-white rounded-xl transition-all shadow-md hover:shadow-red-500/20"><Trash2 class="w-4 h-4" /></button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="filteredStudents.length === 0" class="p-20 text-center space-y-4">
          <div class="w-20 h-20 bg-sigma-surface-alt rounded-full flex items-center justify-center mx-auto">
            <Search class="w-10 h-10 text-sigma-muted/20" />
          </div>
          <p class="text-sigma-muted font-bold uppercase tracking-widest text-sm">Tidak ada santri yang cocok dengan pencarian Anda.</p>
        </div>
      </div>
    </div>
  </div>

    <!-- Modal Registration (Premium Redesign) -->
    <Transition name="fade">
      <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-sigma-app/80 backdrop-blur-md" @click="isModalOpen = false"></div>
        
        <div class="bg-sigma-surface w-full max-w-3xl rounded-[3rem] border border-sigma-border relative z-10 overflow-hidden shadow-2xl animate-in zoom-in-95 duration-500">
          <!-- Modal Header -->
          <div class="p-10 pb-6 border-b border-sigma-border flex justify-between items-start bg-gradient-to-br from-sigma-surface-alt to-transparent">
            <div>
              <div class="flex items-center gap-3 mb-3">
                <span class="px-3 py-1 rounded-full bg-sigma-emerald text-white text-[10px] font-black uppercase tracking-[0.2em]">{{ isEditing ? 'Update' : 'New Enrollment' }}</span>
              </div>
              <h3 class="text-3xl font-black text-sigma-text">{{ isEditing ? 'Edit Profil Santri' : 'Registrasi Santri Baru' }}</h3>
              <p class="text-sigma-muted text-sm mt-1 uppercase font-bold tracking-widest">Sigmabase Module / Database Central</p>
            </div>
            <button @click="isModalOpen = false" class="p-3 bg-sigma-surface-alt hover:bg-sigma-emerald/10 rounded-full transition-all text-sigma-muted hover:text-sigma-emerald">
              <X class="w-6 h-6" />
            </button>
          </div>

          <!-- Modal Content -->
          <div class="p-10 max-h-[60vh] overflow-y-auto custom-scrollbar">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
              <!-- Section: Personal Info -->
              <div class="md:col-span-2 flex items-center gap-4 mb-2">
                <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Biodata Lengkap</span>
                <div class="h-[1px] flex-1 bg-sigma-border"></div>
              </div>

              <div class="space-y-2 md:col-span-2 flex gap-10 items-center bg-sigma-surface-alt p-4 rounded-2xl border border-sigma-border">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Jenis Kelamin</label>
                <div class="flex gap-6">
                  <label class="flex items-center gap-2 cursor-pointer group">
                    <input type="radio" v-model="form.gender" value="L" class="hidden" />
                    <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center group-hover:border-sigma-emerald/50 transition-all" :class="{'border-sigma-emerald bg-sigma-emerald/20': form.gender === 'L'}">
                      <div v-if="form.gender === 'L'" class="w-2 h-2 rounded-full bg-sigma-emerald"></div>
                    </div>
                    <span class="text-sm font-bold" :class="form.gender === 'L' ? 'text-sigma-text' : 'text-sigma-muted'">Laki-laki</span>
                  </label>
                  <label class="flex items-center gap-2 cursor-pointer group">
                    <input type="radio" v-model="form.gender" value="P" class="hidden" />
                    <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center group-hover:border-pink-500/50 transition-all" :class="{'border-pink-500 bg-pink-500/20': form.gender === 'P'}">
                      <div v-if="form.gender === 'P'" class="w-2 h-2 rounded-full bg-pink-500"></div>
                    </div>
                    <span class="text-sm font-bold" :class="form.gender === 'P' ? 'text-sigma-text' : 'text-sigma-muted'">Perempuan</span>
                  </label>
                </div>
              </div>

              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Lengkap</label>
                <input v-model="form.name" type="text" placeholder="Entry Full Name" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Angkatan (Tahun)</label>
                <input v-model="form.entry_year" type="text" placeholder="Ex: 2024" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tempat Lahir</label>
                <input v-model="form.birth_place" type="text" placeholder="Birth City" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tanggal Lahir</label>
                <input v-model="form.birth_date" type="date" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all dark:[color-scheme:dark] text-sm font-bold text-sigma-text" />
              </div>

              <!-- Section: Identifiers -->
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nomor Induk (NIS)</label>
                <input v-model="form.nis" type="text" placeholder="Student ID" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nomor Induk Kependudukan (NIK)</label>
                <input v-model="form.nik" type="text" placeholder="NK Nasional" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">NISN Nasional</label>
                <input v-model="form.nisn" type="text" placeholder="National ID" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
              </div>

              <!-- Section: Address -->
              <div class="md:col-span-2 space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Alamat Domisili Terkini</label>
                <textarea v-model="form.address" rows="2" placeholder="Street, City, Province..." class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all resize-none placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text"></textarea>
              </div>

              <!-- Section: Family -->
              <div class="md:col-span-2 flex items-center gap-4 mt-4 mb-2">
                <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Informasi Keluarga</span>
                <div class="h-[1px] flex-1 bg-sigma-border"></div>
              </div>

              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Orang Tua / Wali</label>
                <input v-model="form.parent_name" type="text" placeholder="Guardian Full Name" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Contact WhatsApp</label>
                <input v-model="form.parent_phone" type="text" placeholder="08xxxxxxxxxx" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
              </div>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="p-10 pt-6 border-t border-sigma-border bg-gradient-to-br from-transparent to-sigma-surface-alt flex gap-4">
            <button @click="isModalOpen = false" :disabled="isLoading" class="flex-1 py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-surface-alt hover:bg-sigma-emerald/10 transition-all text-sigma-muted hover:text-sigma-emerald disabled:opacity-50">
              Cancel
            </button>
            <button @click="handleSubmit" :disabled="isLoading" class="flex-1 [flex-grow:2] py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-emerald hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20 disabled:opacity-50 text-white">
              {{ isLoading ? 'Processing...' : 'Sync Database' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Graduation Modal -->
    <Transition name="fade">
      <div v-if="showGradModal" class="fixed inset-0 z-[110] flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-sigma-app/90 backdrop-blur-md" @click="showGradModal = false"></div>
        
        <div class="bg-sigma-surface w-full max-w-2xl rounded-[3rem] border border-sigma-border relative z-10 overflow-hidden shadow-2xl">
          <div class="p-10 pb-6 border-b border-sigma-border flex justify-between items-start">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-sigma-emerald/10 flex items-center justify-center text-sigma-emerald">
                <GraduationCap class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-2xl font-black text-white italic uppercase tracking-tight">FORM <span class="text-sigma-emerald">KELULUSAN</span></h3>
                <p class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest mt-1">Migrasi Santri Aktif ke Database Alumni</p>
              </div>
            </div>
            <button @click="showGradModal = false" class="p-3 hover:bg-sigma-emerald/10 rounded-full transition-all text-sigma-muted hover:text-sigma-emerald">
              <X class="w-6 h-6" />
            </button>
          </div>

          <div class="p-10 space-y-6">
            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Tahun Lulus</label>
                <input v-model="gradForm.graduation_year" type="number" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none text-sm font-bold" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Angkatan (Batch)</label>
                <input v-model="gradForm.batch" type="text" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none text-sm font-bold" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Nomor WhatsApp</label>
                <input v-model="gradForm.whatsapp" type="text" placeholder="08..." class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none text-sm font-bold" />
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Email Alumni</label>
                <input v-model="gradForm.email" type="email" placeholder="alumni@sigma.com" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none text-sm font-bold" />
              </div>
            </div>

            <div class="space-y-2">
              <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Status Pengabdian</label>
              <div class="grid grid-cols-3 gap-3">
                <button v-for="st in ['Tidak Mengabdi', 'Mengabdi', 'Lainnya']" :key="st" 
                        @click="gradForm.service_status = st"
                        class="py-3 px-4 rounded-xl border text-[10px] font-black uppercase tracking-widest transition-all"
                        :class="gradForm.service_status === st ? 'bg-sigma-emerald border-sigma-emerald text-white shadow-lg shadow-sigma-emerald/20' : 'bg-sigma-surface-alt border-sigma-border text-sigma-muted hover:border-sigma-emerald/50'">
                  {{ st }}
                </button>
              </div>
              <input v-if="gradForm.service_status === 'Lainnya'" v-model="gradForm.other_status" type="text" placeholder="Isi status lainnya..." class="w-full mt-3 p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-emerald/30 focus:border-sigma-emerald outline-none text-sm font-bold animate-in slide-in-from-top-2 duration-300" />
            </div>

            <div class="space-y-2">
              <label class="text-[10px] font-bold uppercase tracking-widest text-sigma-muted ml-1">Foto Terakhir Alumni</label>
              <div class="flex items-center gap-6 p-4 bg-sigma-surface-alt border border-dashed border-sigma-border rounded-3xl">
                <div class="w-20 h-20 rounded-2xl bg-sigma-app overflow-hidden border border-sigma-border flex-shrink-0">
                  <img v-if="gradPhotoPreview" :src="gradPhotoPreview" class="w-full h-full object-cover">
                  <div v-else class="w-full h-full flex items-center justify-center text-sigma-muted/30">
                    <Camera class="w-8 h-8" />
                  </div>
                </div>
                <label class="cursor-pointer px-6 py-2 bg-sigma-emerald/10 hover:bg-sigma-emerald/20 text-sigma-emerald rounded-xl text-[10px] font-black uppercase tracking-widest border border-sigma-emerald/20 transition-all">
                  Pilih Foto
                  <input type="file" @change="handleGradPhoto" accept="image/*" class="hidden">
                </label>
              </div>
            </div>
          </div>

          <div class="p-10 pt-6 border-t border-sigma-border flex gap-4 bg-sigma-surface-alt/30">
            <button @click="showGradModal = false" class="flex-1 py-5 rounded-2xl font-black uppercase tracking-widest text-[10px] text-sigma-muted hover:bg-sigma-emerald/10 hover:text-sigma-emerald border border-sigma-border transition-all">Batal</button>
            <button @click="handleGraduation" class="flex-1 [flex-grow:2] py-5 rounded-2xl font-black uppercase tracking-widest text-[10px] bg-sigma-emerald text-white shadow-xl shadow-sigma-emerald/20 hover:bg-emerald-600 transition-all">Proses Kelulusan</button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.1);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
