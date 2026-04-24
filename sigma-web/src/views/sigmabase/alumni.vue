<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Search, Filter, 
  Edit, Trash2, X, Plus, 
  FileSpreadsheet, FileText, Upload, Mail, Phone, Camera, GraduationCap, Download
} from 'lucide-vue-next'
import axios from 'axios'
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

// Config API
const API_URL = '/api/v1/admin/base/alumni'

// State
const alumniList = ref<any[]>([])
const isLoading = ref(false)
const searchQuery = ref('')
const isYearDropdownOpen = ref(false)
const selectedYear = ref('')
const showDataMenu = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const photoInput = ref<HTMLInputElement | null>(null)

// Modal State
const isModalOpen = ref(false)
const isEditing = ref(false)
const currentId = ref<number | null>(null)
const photoPreview = ref('')
const photoFile = ref<File | null>(null)

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

const uniqueYears = computed(() => {
  const years = alumniList.value.map(a => a.graduation_year).filter(y => y)
  return [...new Set(years)].sort((a, b) => String(b).localeCompare(String(a)))
})

// Methods
const fetchAlumni = async () => {
  try {
    isLoading.value = true
    const response = await axios.get(API_URL)
    alumniList.value = response.data
  } catch (error) {
    console.error('Error fetching alumni:', error)
  } finally {
    isLoading.value = false
  }
}

const filteredAlumni = computed(() => {
  let result = alumniList.value
  
  if (selectedYear.value) {
    result = result.filter(item => String(item.graduation_year) === String(selectedYear.value))
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(item => 
      (item.name && item.name.toLowerCase().includes(q)) ||
      (item.nis && item.nis.toLowerCase().includes(q)) ||
      (item.batch && item.batch.toLowerCase().includes(q))
    )
  }
  
  return result
})

const openNewModal = () => {
  isEditing.value = false
  currentId.value = null
  photoPreview.value = ''
  photoFile.value = null
  form.value = {
    name: '', gender: 'L', nis: '', nik: '', birth_place: '', birth_date: '',
    graduation_year: new Date().getFullYear().toString(), batch: '',
    email: '', address: '', whatsapp: '', service_status: 'Tidak Mengabdi', photo: ''
  }
  isModalOpen.value = true
}

const handleEdit = (alumni: any) => {
  isEditing.value = true
  currentId.value = alumni.ID || alumni.id
  form.value = { 
    ...alumni, 
    birth_date: alumni.birth_date ? alumni.birth_date.split('T')[0] : '',
    graduation_year: String(alumni.graduation_year)
  }
  photoPreview.value = alumni.photo || ''
  photoFile.value = null
  isModalOpen.value = true
}

const triggerPhotoSelect = () => {
  photoInput.value?.click()
}

const handlePhotoSelect = (event: any) => {
  const file = event.target.files[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    alert('Foto terlalu besar (maks 2MB)')
    return
  }
  photoFile.value = file
  photoPreview.value = URL.createObjectURL(file)
}

const uploadPhoto = async (alumniId: number) => {
  if (!photoFile.value) return
  const formData = new FormData()
  formData.append('photo', photoFile.value)
  try {
    await axios.post(`${API_URL}/${alumniId}/photo`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  } catch (err) {
    console.error('Gagal upload foto:', err)
  }
}

const handleSubmit = async () => {
  try {
    isLoading.value = true
    let alumniId: number
    
    if (isEditing.value && currentId.value) {
      await axios.put(`${API_URL}/${currentId.value}`, form.value)
      alumniId = currentId.value
    } else {
      const res = await axios.post(API_URL, form.value)
      alumniId = res.data.ID
    }

    if (photoFile.value) {
      await uploadPhoto(alumniId)
    }

    isModalOpen.value = false
    await fetchAlumni()
  } catch (error: any) {
    console.error('Error saving alumni:', error)
    alert('Gagal menyimpan data: ' + (error.response?.data?.error || error.message))
  } finally {
    isLoading.value = false
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus data alumni ini?')) return
  try {
    isLoading.value = true
    await axios.delete(`${API_URL}/${id}`)
    await fetchAlumni()
  } catch (error) {
    console.error('Error deleting alumni:', error)
    alert('Gagal menghapus data')
  } finally {
    isLoading.value = false
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
    startY: 20,
    theme: 'grid',
    headStyles: { fillColor: [16, 185, 129] }
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
    const rawData = XLSX.utils.sheet_to_json(worksheet)
    
    const formattedData = rawData.map((row: any) => ({
      name: row['Nama Lengkap'] || row['Nama'],
      gender: row['L/P'] || row['Jenis Kelamin'] || 'L',
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
      isLoading.value = true
      await axios.post(`${API_URL}/bulk`, formattedData)
      await fetchAlumni()
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

onMounted(fetchAlumni)
</script>

<template>
  <div class="flex flex-col">
    <div class="flex-1 flex flex-col">
      <!-- Sticky Header Container -->
      <div class="sticky top-0 z-40 bg-sigma-app/80 backdrop-blur-xl border-b border-sigma-border px-8 py-6 space-y-6 shadow-sm">
        <!-- Page Title & Stats -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-emerald-500/10 rounded-2xl flex items-center justify-center text-emerald-400">
              <GraduationCap class="w-6 h-6" />
            </div>
            <div>
              <h2 class="text-3xl font-black tracking-tight italic">DATABASE <span class="text-sigma-emerald not-italic uppercase text-2xl">ALUMNI</span></h2>
              <p class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest mt-1 hidden md:block">Database Lulusan & Pengabdian</p>
            </div>
          </div>

          <div class="text-[10px] text-sigma-muted font-bold uppercase tracking-[0.2em] px-4 py-2 bg-sigma-surface-alt rounded-full border border-sigma-border hidden lg:block">
            Verified: {{ filteredAlumni.length }} alumni
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
              placeholder="Cari nama alumni, angkatan, atau tahun lulus..."
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
                  <span>{{ selectedYear ? 'Tahun ' + selectedYear : 'Semua Tahun' }}</span>
                </div>
                <div class="w-1 h-1 rounded-full bg-sigma-muted"></div>
              </button>

              <div v-if="isYearDropdownOpen" class="absolute top-full left-0 mt-3 w-full bg-sigma-surface border border-sigma-border rounded-2xl shadow-2xl z-50 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                <div class="max-h-60 overflow-y-auto custom-scrollbar">
                  <button @click="selectedYear = ''; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-sigma-border" :class="selectedYear === '' ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-sigma-muted'">Semua Tahun</button>
                  <button v-for="year in uniqueYears" :key="year" @click="selectedYear = year as string; isYearDropdownOpen = false" class="w-full text-left px-6 py-4 text-xs font-bold uppercase tracking-widest hover:bg-sigma-emerald/10 hover:text-sigma-emerald transition-all border-b border-sigma-border last:border-0" :class="selectedYear === year ? 'text-sigma-emerald bg-sigma-emerald/5' : 'text-sigma-muted'">Tahun {{ year }}</button>
                </div>
              </div>
            </div>

            <!-- Data Operations -->
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

            <button @click="openNewModal" class="flex items-center gap-2 px-8 py-4 bg-sigma-emerald text-white rounded-2xl font-bold uppercase tracking-widest text-xs hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20">
              <Plus class="w-4 h-4" /> Tambah Alumni
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content Scrolling Area -->
      <div class="p-8 space-y-10">
        <!-- Alumni Table -->
        <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.3em] bg-sigma-surface-alt">
                  <th class="p-8 border-b border-sigma-border">Identitas Alumni</th>
                  <th class="p-8 border-b border-sigma-border text-center">Tahun Lulus</th>
                  <th class="p-8 border-b border-sigma-border text-center">Angkatan</th>
                  <th class="p-8 border-b border-sigma-border">Status Pengabdian</th>
                  <th class="p-8 border-b border-sigma-border text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border">
                <tr v-for="alumni in filteredAlumni" :key="alumni.ID || alumni.id" class="hover:bg-sigma-surface-alt transition-colors group">
                  <td class="p-8">
                    <div class="flex items-center gap-5">
                      <div class="w-14 h-16 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center overflow-hidden shadow-sm group-hover:scale-110 transition-transform">
                        <img v-if="alumni.photo" :src="alumni.photo" class="w-full h-full object-cover" />
                        <div v-else class="font-black text-sigma-emerald text-xl uppercase">
                          {{ alumni.name.charAt(0) }}
                        </div>
                      </div>
                      <div>
                        <div class="flex items-center gap-2">
                          <span class="font-bold text-sigma-text text-lg group-hover:text-sigma-emerald transition-colors leading-tight">{{ alumni.name }}</span>
                          <span class="px-2 py-0.5 rounded text-[9px] font-black uppercase tracking-widest" 
                                :class="alumni.gender === 'L' ? 'bg-blue-500/10 text-blue-500' : 'bg-pink-500/10 text-pink-500'">
                            {{ alumni.gender }}
                          </span>
                        </div>
                        <div class="text-[10px] font-mono text-sigma-muted space-y-0.5 mt-2">
                          <div class="flex items-center gap-2">
                            <span class="w-8 opacity-50 uppercase tracking-tighter">NIS:</span>
                            <span class="px-1.5 py-0.5 rounded-md bg-sigma-surface border border-sigma-border text-sigma-text/80">{{ alumni.nis || '-' }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="w-8 opacity-50 uppercase tracking-tighter">NIK:</span>
                            <span class="text-sigma-text/60">{{ alumni.nik || '-' }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="p-8 text-center">
                    <span class="inline-block px-4 py-2 rounded-xl bg-sigma-surface-alt border border-sigma-border font-black text-sm text-sigma-muted group-hover:text-sigma-text transition-all">
                      {{ alumni.graduation_year }}
                    </span>
                  </td>
                  <td class="p-8 text-center">
                    <span class="inline-block px-4 py-2 rounded-xl bg-sigma-surface-alt border border-sigma-border font-black text-xs text-sigma-muted group-hover:text-sigma-text transition-all uppercase tracking-widest">
                      {{ alumni.batch || 'GENERAL' }}
                    </span>
                  </td>
                  <td class="p-8">
                    <div class="space-y-1">
                      <span class="text-sm font-bold" 
                            :class="alumni.service_status === 'Mengabdi' ? 'text-sigma-emerald' : alumni.service_status === 'Tidak Mengabdi' ? 'text-sigma-muted' : 'text-amber-400'">
                        {{ alumni.service_status }}
                      </span>
                      <div class="flex flex-col gap-1 mt-1">
                        <div class="flex items-center gap-2 text-[10px] font-bold text-sigma-muted uppercase tracking-tight">
                          <Mail class="w-3 h-3" /> {{ alumni.email || '-' }}
                        </div>
                        <div class="flex items-center gap-2 text-[10px] font-bold text-sigma-muted uppercase tracking-tight">
                          <Phone class="w-3 h-3" /> {{ alumni.whatsapp || '-' }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="p-8 text-right">
                    <div class="flex justify-end gap-3 translate-x-4 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-500">
                      <button @click="handleEdit(alumni)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-sigma-emerald text-sigma-muted hover:text-white rounded-xl transition-all shadow-md hover:shadow-sigma-emerald/20"><Edit class="w-4 h-4" /></button>
                      <button @click="handleDelete(alumni.ID || alumni.id)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-red-500 text-sigma-muted hover:text-white rounded-xl transition-all shadow-md hover:shadow-red-500/20"><Trash2 class="w-4 h-4" /></button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-if="filteredAlumni.length === 0" class="p-20 text-center space-y-4">
              <div class="w-20 h-20 bg-sigma-surface-alt rounded-full flex items-center justify-center mx-auto">
                <Search class="w-10 h-10 text-sigma-muted/20" />
              </div>
              <p class="text-sigma-muted font-bold uppercase tracking-widest text-sm">Tidak ada alumni yang cocok dengan pencarian Anda.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Registration Modal -->
      <Transition name="fade">
        <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-sigma-app/80 backdrop-blur-md" @click="isModalOpen = false"></div>
          
          <div class="bg-sigma-surface w-full max-w-4xl rounded-[3rem] border border-sigma-border relative z-10 overflow-hidden shadow-2xl animate-in zoom-in-95 duration-500">
            <!-- Modal Header -->
            <div class="p-10 pb-6 border-b border-sigma-border flex justify-between items-start bg-gradient-to-br from-sigma-surface-alt to-transparent">
              <div>
                <div class="flex items-center gap-3 mb-3">
                  <span class="px-3 py-1 rounded-full bg-sigma-emerald text-white text-[10px] font-black uppercase tracking-[0.2em]">{{ isEditing ? 'Update' : 'New Archive' }}</span>
                </div>
                <h3 class="text-3xl font-black text-sigma-text italic uppercase tracking-tight">{{ isEditing ? 'Edit' : 'Tambah' }} <span class="text-sigma-emerald not-italic uppercase">Alumni</span></h3>
                <p class="text-sigma-muted text-sm mt-1 uppercase font-bold tracking-widest">Sigmabase Module / Alumni Central</p>
              </div>
              <button @click="isModalOpen = false" class="p-3 bg-sigma-surface-alt hover:bg-sigma-emerald/10 rounded-full transition-all text-sigma-muted hover:text-sigma-emerald">
                <X class="w-6 h-6" />
              </button>
            </div>

            <!-- Modal Content -->
            <div class="p-10 max-h-[65vh] overflow-y-auto custom-scrollbar">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <!-- Section: Identity with Photo -->
                <div class="md:col-span-2 flex items-center gap-8 mb-6 bg-sigma-surface-alt/50 p-6 rounded-[2rem] border border-sigma-border">
                  <div class="relative group">
                    <input type="file" ref="photoInput" class="hidden" accept="image/*" @change="handlePhotoSelect" />
                    <div 
                      @click="triggerPhotoSelect"
                      class="w-32 h-40 rounded-2xl border-2 border-dashed border-sigma-border flex flex-col items-center justify-center gap-2 cursor-pointer group-hover:border-sigma-emerald transition-all overflow-hidden bg-sigma-surface"
                    >
                      <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
                      <template v-else>
                        <Camera class="w-8 h-8 text-sigma-muted group-hover:text-sigma-emerald" />
                        <span class="text-[9px] font-black uppercase tracking-tighter text-sigma-muted">Upload Photo</span>
                      </template>
                    </div>
                  </div>

                  <div class="flex-1 space-y-4">
                    <div class="flex items-center gap-4 mb-2">
                      <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Biodata Alumni</span>
                      <div class="h-[1px] flex-1 bg-sigma-border"></div>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div class="space-y-2">
                        <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Lengkap</label>
                        <input v-model="form.name" type="text" placeholder="Entry Full Name" class="w-full p-4 rounded-2xl bg-sigma-surface border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-bold text-sigma-text" />
                      </div>
                      <div class="space-y-2">
                        <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Jenis Kelamin</label>
                        <div class="flex gap-6 p-4">
                          <label class="flex items-center gap-2 cursor-pointer group">
                            <input type="radio" v-model="form.gender" value="L" class="hidden" />
                            <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center" :class="{'border-sigma-emerald bg-sigma-emerald/20': form.gender === 'L'}">
                              <div v-if="form.gender === 'L'" class="w-2 h-2 rounded-full bg-sigma-emerald"></div>
                            </div>
                            <span class="text-sm font-bold">L</span>
                          </label>
                          <label class="flex items-center gap-2 cursor-pointer group">
                            <input type="radio" v-model="form.gender" value="P" class="hidden" />
                            <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center" :class="{'border-pink-500 bg-pink-500/20': form.gender === 'P'}">
                              <div v-if="form.gender === 'P'" class="w-2 h-2 rounded-full bg-pink-500"></div>
                            </div>
                            <span class="text-sm font-bold">P</span>
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nomor Induk (NIS)</label>
                  <input v-model="form.nis" type="text" placeholder="NIS" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">NIK (Nomor Induk Kependudukan)</label>
                  <input v-model="form.nik" type="text" placeholder="NIK 16 Digit" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                   <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tempat & Tanggal Lahir</label>
                   <div class="grid grid-cols-2 gap-3">
                     <input v-model="form.birth_place" type="text" placeholder="City" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                     <input v-model="form.birth_date" type="date" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all dark:[color-scheme:dark] text-sm font-bold text-sigma-text" />
                   </div>
                </div>

                <!-- Academic Section -->
                <div class="md:col-span-2 flex items-center gap-4 mt-4 mb-2">
                  <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Data Kelulusan</span>
                  <div class="h-[1px] flex-1 bg-sigma-border"></div>
                </div>

                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tahun Lulus</label>
                  <input v-model="form.graduation_year" type="number" placeholder="2024" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Angkatan (Batch)</label>
                  <input v-model="form.batch" type="text" placeholder="Ex: Generation of Light" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Status Pengabdian</label>
                  <div class="relative group">
                    <select v-model="form.service_status" class="w-full p-4 pr-12 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text appearance-none cursor-pointer hover:border-sigma-emerald/30">
                      <option value="Tidak Mengabdi" class="bg-slate-900 text-white">Tidak Mengabdi</option>
                      <option value="Mengabdi" class="bg-slate-900 text-white">Mengabdi</option>
                      <option value="Lainnya" class="bg-slate-900 text-white">Lainnya</option>
                    </select>
                    <Download class="absolute right-5 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted rotate-90 pointer-events-none group-hover:text-sigma-emerald transition-colors" />
                  </div>
                </div>

                <!-- Section: Contact -->
                <div class="md:col-span-2 flex items-center gap-4 mt-4 mb-2">
                  <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Kontak Terkini</span>
                  <div class="h-[1px] flex-1 bg-sigma-border"></div>
                </div>

                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Alamat Email</label>
                  <input v-model="form.email" type="email" placeholder="alumni@sigma.com" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Nomor WhatsApp</label>
                  <input v-model="form.whatsapp" type="text" placeholder="08xxxxxxxxxx" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="md:col-span-2 space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Alamat Domisili</label>
                  <textarea v-model="form.address" rows="2" placeholder="Street, City, Province..." class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all resize-none text-sm font-bold text-sigma-text"></textarea>
                </div>
              </div>
            </div>

            <div class="p-10 pt-6 border-t border-sigma-border bg-gradient-to-br from-transparent to-sigma-surface-alt flex gap-4">
              <button @click="isModalOpen = false" :disabled="isLoading" class="flex-1 py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-surface-alt hover:bg-sigma-emerald/10 transition-all text-sigma-muted hover:text-sigma-emerald disabled:opacity-50">
                Cancel
              </button>
              <button @click="handleSubmit" :disabled="isLoading" class="flex-1 [flex-grow:2] py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs bg-sigma-emerald hover:bg-emerald-600 transition-all shadow-xl shadow-sigma-emerald/20 disabled:opacity-50 text-white">
                {{ isLoading ? 'Syncing...' : 'Simpan Data Alumni' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </div>
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
