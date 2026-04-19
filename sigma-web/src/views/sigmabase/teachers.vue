<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Search, 
  Download, Edit, Trash2, X, Plus, 
  FileSpreadsheet, FileText, Upload, Archive, Users, Mail, Phone, Camera, ChevronRight
} from 'lucide-vue-next'
import SigmabaseSidebar from '../../components/SigmabaseSidebar.vue'
import axios from 'axios'
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

// Config API
const API_URL = '/api/v1/admin/base/teachers'

// State
const isModalOpen = ref(false)
const isEditing = ref(false)
const currentId = ref<number | null>(null)
const searchQuery = ref('')
const isLoading = ref(false)
const isActiveFilter = ref('active') // 'active', 'archived', 'all'
const showDataMenu = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const photoInput = ref<HTMLInputElement | null>(null)
const photoPreview = ref('')
const photoFile = ref<File | null>(null)

// Handle click outside to close dropdowns
onMounted(() => {
  window.addEventListener('click', (e: MouseEvent) => {
    const target = e.target as HTMLElement
    if (!target.closest('.data-menu-container')) {
      showDataMenu.value = false
    }
  })
})

// Form State
const form = ref({
  name: '',
  gender: 'L',
  nik: '',
  nip: '',
  birth_place: '',
  birth_date: '',
  position: '',
  status: 'Guru Tetap',
  education: '',
  graduation_year: '',
  email: '',
  address: '',
  whatsapp: '',
  photo: '',
  entry_year: new Date().getFullYear().toString(),
  is_active: true
})

// Database Data
const teachers = ref<any[]>([])

const fetchTeachers = async () => {
  try {
    isLoading.value = true
    const response = await axios.get(API_URL)
    teachers.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data guru:', err)
  } finally {
    isLoading.value = false
  }
}

// Client-side Search & Filtering
const filteredTeachers = computed(() => {
  let result = teachers.value
  
  // Filter by Active Status
  if (isActiveFilter.value === 'active') {
    result = result.filter(t => t.is_active === true)
  } else if (isActiveFilter.value === 'archived') {
    result = result.filter(t => t.is_active === false)
  }

  // Filter by Search Query
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(t => 
      t.name.toLowerCase().includes(q) || 
      (t.nip && t.nip.toLowerCase().includes(q)) ||
      (t.nik && t.nik.toLowerCase().includes(q))
    )
  }
  
  return result
})

const openNewModal = () => {
  photoPreview.value = ''
  photoFile.value = null
  isEditing.value = false
  currentId.value = null
  form.value = {
    name: '',
    gender: 'L',
    nik: '',
    nip: '',
    birth_place: '',
    birth_date: '',
    position: '',
    status: 'Guru Tetap',
    education: '',
    graduation_year: '',
    email: '',
    address: '',
    whatsapp: '',
    photo: '',
    entry_year: new Date().getFullYear().toString(),
    is_active: true
  }
  isModalOpen.value = true
}

const handleEdit = (teacher: any) => {
  isEditing.value = true
  currentId.value = teacher.ID
  form.value = {
    name: teacher.name,
    gender: teacher.gender || 'L',
    nik: teacher.nik,
    nip: teacher.nip || '',
    birth_place: teacher.birth_place || '',
    birth_date: teacher.birth_date ? teacher.birth_date.split('T')[0] : '',
    position: teacher.position || '',
    status: teacher.status || 'Guru Tetap',
    education: teacher.education || '',
    graduation_year: teacher.graduation_year || '',
    email: teacher.email || '',
    address: teacher.address || '',
    whatsapp: teacher.whatsapp || '',
    photo: teacher.photo || '',
    entry_year: teacher.entry_year || '',
    is_active: teacher.is_active
  }
  photoPreview.value = teacher.photo ? teacher.photo : ''
  photoFile.value = null
  isModalOpen.value = true
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    }).format(date)
  } catch (e) {
    return dateString
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('Apakah Anda yakin ingin menghapus data guru ini?')) return
  
  try {
    isLoading.value = true
    await axios.delete(`${API_URL}/${id}`)
    await fetchTeachers()
  } catch (err) {
    alert('Gagal menghapus data.')
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const toggleStatus = async (teacher: any) => {
  try {
    isLoading.value = true
    await axios.put(`${API_URL}/${teacher.ID}`, { is_active: !teacher.is_active })
    await fetchTeachers()
  } catch (err) {
    alert('Gagal mengubah status aktif.')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchTeachers()
})

// Export & Import Logic
const exportToExcel = () => {
  const data = teachers.value.map(t => ({
    'Nama Lengkap': t.name,
    'L/P': t.gender,
    'NIK': t.nik,
    'NIP': t.nip,
    'Status': t.status,
    'Jabatan': t.position,
    'Pendidikan': t.education,
    'Tahun Lulus': t.graduation_year,
    'Tempat Lahir': t.birth_place,
    'Tanggal Lahir': t.birth_date ? t.birth_date.split('T')[0] : '',
    'Email': t.email,
    'Telepon': t.whatsapp,
    'Alamat': t.address,
    'Tahun Masuk': t.entry_year,
    'Aktif': t.is_active ? 'Ya' : 'Tidak'
  }))
  
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, "Data Guru")
  XLSX.writeFile(workbook, `Data_Guru_${new Date().toLocaleDateString()}.xlsx`)
  showDataMenu.value = false
}

const exportToPDF = () => {
  const doc = new jsPDF()
  doc.text("Laporan Data Guru Al-Hikmah", 14, 15)
  
  const tableData = teachers.value.map(t => [
    t.name, t.nip || t.nik, t.status, t.position, t.whatsapp, t.is_active ? 'Aktif' : 'Arsip'
  ])
  
  autoTable(doc, {
    head: [['Nama', 'Identitas', 'Status', 'Jabatan', 'Telepon', 'Keterangan']],
    body: tableData,
    startY: 20,
    theme: 'grid',
    headStyles: { fillColor: [16, 185, 129] }
  })
  
  doc.save(`Laporan_Guru_${new Date().toLocaleDateString()}.pdf`)
  showDataMenu.value = false
}

const downloadTemplate = () => {
  const headers = [
    ['Nama Lengkap', 'NIK', 'NIP', 'Tempat Lahir', 'Tanggal Lahir', 'Jabatan', 'Status', 'Pendidikan', 'Tahun Lulus', 'Email', 'Telepon', 'Alamat', 'Tahun Masuk']
  ]
  
  const worksheet = XLSX.utils.aoa_to_sheet(headers)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, "Template")
  XLSX.writeFile(workbook, "Template_Import_Guru.xlsx")
}

const triggerImport = () => {
  fileInput.value?.click()
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

const uploadPhoto = async (teacherId: number) => {
  if (!photoFile.value) return
  
  const formData = new FormData()
  formData.append('photo', photoFile.value)
  
  try {
    await axios.post(`${API_URL}/${teacherId}/photo`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  } catch (err) {
    console.error('Gagal upload foto:', err)
  }
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
    
    const formattedData = jsonData.map((row: any) => ({
      name: row['Nama Lengkap'],
      gender: row['L/P'] || row['Jenis Kelamin'] || 'L',
      nik: String(row['NIK'] || ''),
      nip: String(row['NIP'] || ''),
      birth_place: row['Tempat Lahir'] || '',
      birth_date: row['Tanggal Lahir'] || '',
      position: row['Jabatan'] || '',
      status: row['Status'] || 'Guru Tetap',
      education: row['Pendidikan'] || '',
      graduation_year: String(row['Tahun Lulus'] || ''),
      email: row['Email'] || '',
      whatsapp: String(row['Telepon'] || ''),
      address: row['Alamat'] || '',
      entry_year: String(row['Tahun Masuk'] || ''),
      is_active: true
    }))

    try {
      isLoading.value = true
      await axios.post(`${API_URL}/bulk`, formattedData)
      await fetchTeachers()
      alert('Import berhasil!')
    } catch (err) {
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
      if (photoFile.value) await uploadPhoto(currentId.value)
    } else {
      const res = await axios.post(API_URL, form.value)
      if (photoFile.value) await uploadPhoto(res.data.ID)
    }
    
    await fetchTeachers()
    isModalOpen.value = false
  } catch (err) {
    alert('Gagal menyimpan data guru. Periksa kembali NIK/NIP.')
    console.error(err)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="flex h-screen bg-sigma-app text-sigma-text overflow-hidden font-sans transition-colors duration-300">
    
    <SigmabaseSidebar activeItem="Data Guru" />

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto custom-scrollbar flex flex-col">
      
      <!-- Sticky Header -->
      <div class="sticky top-0 z-40 bg-sigma-app/80 backdrop-blur-xl border-b border-sigma-border px-8 py-6 space-y-6 shadow-sm">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-400">
              <Users class="w-6 h-6" />
            </div>
            <div>
              <h2 class="text-3xl font-black tracking-tight italic">DATA <span class="text-sigma-emerald not-italic uppercase text-2xl">GURU</span></h2>
              <p class="text-[10px] text-sigma-muted font-bold uppercase tracking-widest mt-1 hidden md:block">SDM Al-Hikmah Professional Staff</p>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div class="text-[10px] text-sigma-muted font-bold uppercase tracking-[0.2em] px-4 py-2 bg-sigma-surface-alt rounded-full border border-sigma-border hidden lg:block">
              Active SDM: {{ teachers.filter(t => t.is_active).length }} Personnel
            </div>
          </div>
        </div>

        <!-- Unified Single-Row Action Bar -->
        <div class="flex flex-col lg:flex-row gap-4 items-center">
          <!-- Search -->
          <div class="relative flex-1 w-full lg:w-auto">
            <input type="file" ref="fileInput" class="hidden" accept=".xlsx, .xls" @change="handleImport" />
            <Search class="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted" />
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Cari Nama, NIK atau NIP Guru..." 
              class="w-full pl-14 pr-6 py-4 rounded-2xl bg-sigma-surface border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/50 text-sm font-medium"
            />
          </div>
          
          <div class="flex flex-wrap items-center gap-3 w-full lg:w-auto">
            <!-- Archiving Filter -->
            <div class="flex rounded-2xl bg-sigma-surface border border-sigma-border p-1">
              <button 
                @click="isActiveFilter = 'active'"
                class="px-4 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="isActiveFilter === 'active' ? 'bg-sigma-emerald text-white shadow-lg' : 'text-sigma-muted hover:text-sigma-text'"
              >
                Aktif
              </button>
              <button 
                @click="isActiveFilter = 'archived'"
                class="px-4 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all"
                :class="isActiveFilter === 'archived' ? 'bg-red-500 text-white shadow-lg' : 'text-sigma-muted hover:text-sigma-text'"
              >
                Arsip
              </button>
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
              <Plus class="w-4 h-4" /> Tambah Guru
            </button>
          </div>
        </div>
      </div>

      <div class="p-8 space-y-10">
        <div class="bg-sigma-surface border border-sigma-border rounded-[2.5rem] overflow-hidden shadow-sm">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="text-[10px] text-sigma-muted uppercase tracking-[0.3em] bg-sigma-surface-alt">
                  <th class="p-8 border-b border-sigma-border">Identitas Guru</th>
                  <th class="p-8 border-b border-sigma-border text-center">Lahir & Gender</th>
                  <th class="p-8 border-b border-sigma-border">Jabatan</th>
                  <th class="p-8 border-b border-sigma-border">Pendidikan</th>
                  <th class="p-8 border-b border-sigma-border">Kontak & Alamat</th>
                  <th class="p-8 border-b border-sigma-border text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-sigma-border">
                <tr v-for="teacher in filteredTeachers" :key="teacher.ID" class="hover:bg-sigma-surface-alt transition-all group">
                  <td class="p-8">
                    <div class="flex items-center gap-5">
                      <div class="w-14 h-16 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center overflow-hidden shadow-sm group-hover:scale-110 transition-transform">
                        <img v-if="teacher.photo" :src="teacher.photo" class="w-full h-full object-cover" />
                        <div v-else class="font-black text-sigma-emerald text-xl uppercase">
                          {{ teacher.name.charAt(0) }}
                        </div>
                      </div>
                      <div>
                        <div class="font-bold text-sigma-text text-lg group-hover:text-sigma-emerald transition-colors leading-tight">{{ teacher.name }}</div>
                        <div class="text-[10px] font-mono text-sigma-muted space-y-0.5 mt-2">
                          <div class="flex items-center gap-2">
                            <span class="w-8 opacity-50 uppercase tracking-tighter">NIP:</span>
                            <span class="px-1.5 py-0.5 rounded-md bg-sigma-surface border border-sigma-border text-sigma-text/80">{{ teacher.nip || '-' }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="w-8 opacity-50 uppercase tracking-tighter">NIK:</span>
                            <span class="text-sigma-text/60">{{ teacher.nik }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="p-8">
                    <div class="flex flex-col items-center gap-2">
                      <span class="px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="teacher.gender === 'P' ? 'border-pink-500/30 bg-pink-500/10 text-pink-400' : 'border-blue-500/30 bg-blue-500/10 text-blue-400'">
                        {{ teacher.gender === 'L' ? 'Laki-laki' : 'Perempuan' }}
                      </span>
                      <div class="text-center">
                        <div class="text-xs font-bold text-sigma-text/80">{{ teacher.birth_place || '-' }}</div>
                        <div class="text-[10px] text-sigma-muted font-bold mt-0.5">{{ formatDate(teacher.birth_date) }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="p-8">
                    <div class="space-y-1.5">
                      <div class="text-sm font-bold text-sigma-text uppercase tracking-tight">{{ teacher.position }}</div>
                      <div class="flex items-center gap-2">
                        <span class="px-2 py-0.5 rounded-md bg-sigma-emerald/10 text-sigma-emerald text-[9px] font-black uppercase tracking-widest border border-sigma-emerald/20">
                          {{ teacher.status }}
                        </span>
                        <span class="text-[10px] text-sigma-muted font-bold italic">Masuk: {{ teacher.entry_year }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="p-8">
                    <div class="space-y-1">
                      <div class="text-xs font-black text-sigma-emerald uppercase">{{ teacher.education }}</div>
                      <div class="flex items-center gap-2">
                        <div class="w-1.5 h-1.5 rounded-full bg-sigma-border"></div>
                        <span class="text-[10px] text-sigma-muted uppercase tracking-widest font-black">Lulus {{ teacher.graduation_year }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="p-8">
                    <div class="space-y-2">
                      <div class="flex flex-col gap-1">
                        <div class="flex items-center gap-2 text-xs font-bold text-sigma-text/80">
                          <Mail class="w-3.5 h-3.5 text-sigma-emerald" /> {{ teacher.email || '-' }}
                        </div>
                        <div class="flex items-center gap-2 text-xs font-bold text-sigma-text/60">
                          <Phone class="w-3.5 h-3.5 text-sigma-emerald" /> {{ teacher.whatsapp }}
                        </div>
                      </div>
                      <div class="max-w-[180px] text-[10px] text-sigma-muted leading-relaxed font-medium bg-sigma-surface-alt/50 p-2 rounded-xl border border-sigma-border/50">
                        {{ teacher.address || 'Alamat belum diatur' }}
                      </div>
                    </div>
                  </td>
                  <td class="p-8 text-right">
                    <div class="flex justify-end gap-3 translate-x-4 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-500">
                      <button @click="toggleStatus(teacher)" 
                        :class="teacher.is_active ? 'hover:bg-orange-500' : 'hover:bg-sigma-emerald'"
                        class="p-3 bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-white rounded-xl transition-all shadow-md"
                        :title="teacher.is_active ? 'Arsipkan Guru' : 'Aktifkan Guru'"
                      >
                        <Archive class="w-4 h-4" />
                      </button>
                      <button @click="handleEdit(teacher)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-sigma-emerald text-sigma-muted hover:text-white rounded-xl transition-all shadow-md"><Edit class="w-4 h-4" /></button>
                      <button @click="handleDelete(teacher.ID)" class="p-3 bg-sigma-surface border border-sigma-border hover:bg-red-500 text-sigma-muted hover:text-white rounded-xl transition-all shadow-md"><Trash2 class="w-4 h-4" /></button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Registration Modal -->
      <Transition name="fade">
        <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-sigma-app/80 backdrop-blur-md" @click="isModalOpen = false"></div>
          
          <div class="bg-sigma-surface w-full max-w-4xl rounded-[3rem] border border-sigma-border relative z-10 overflow-hidden shadow-2xl animate-in zoom-in-95 duration-500">
            <div class="p-10 pb-6 border-b border-sigma-border flex justify-between items-start bg-gradient-to-br from-sigma-surface-alt to-transparent">
              <div>
                <div class="flex items-center gap-3 mb-3">
                  <span class="px-3 py-1 rounded-full bg-sigma-emerald text-white text-[10px] font-black uppercase tracking-[0.2em]">{{ isEditing ? 'Update Profile' : 'New Teacher' }}</span>
                  <span v-if="!form.is_active" class="px-3 py-1 rounded-full bg-red-500 text-white text-[10px] font-black uppercase tracking-[0.2em]">Archived</span>
                </div>
                <h3 class="text-3xl font-black text-sigma-text">{{ isEditing ? 'Edit Data Pendidik' : 'Pendaftaran Guru Baru' }}</h3>
                <p class="text-sigma-muted text-sm mt-1 uppercase font-bold tracking-widest">Sigmabase Module / Personnel Central</p>
              </div>
              <button @click="isModalOpen = false" class="p-3 bg-sigma-surface-alt hover:bg-sigma-emerald/10 rounded-full transition-all text-sigma-muted hover:text-sigma-emerald">
                <X class="w-6 h-6" />
              </button>
            </div>

            <div class="p-10 max-h-[65vh] overflow-y-auto custom-scrollbar">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <!-- Section: Personal Info with Photo -->
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
                        <span class="text-[9px] font-black uppercase tracking-tighter text-sigma-muted">Upload 4x6</span>
                      </template>
                    </div>
                  </div>

                  <div class="flex-1 space-y-4">
                    <div class="flex items-center gap-4 mb-2">
                      <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Biodata Personel</span>
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

                <!-- Section: Identity -->
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">NIK (Nomor Induk Kependudukan)</label>
                  <input v-model="form.nik" type="text" placeholder="NIK 16 Digit" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">NIP (Nomor Induk Pegawai)</label>
                  <input v-model="form.nip" type="text" placeholder="NIP (Optional)" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all placeholder:text-sigma-muted/30 text-sm font-mono font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                   <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tempat & Tanggal Lahir</label>
                   <div class="grid grid-cols-2 gap-3">
                     <input v-model="form.birth_place" type="text" placeholder="City" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                     <input v-model="form.birth_date" type="date" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all dark:[color-scheme:dark] text-sm font-bold text-sigma-text" />
                   </div>
                </div>

                <!-- Section: Job & Education -->
                <div class="md:col-span-2 flex items-center gap-4 mt-4 mb-2">
                  <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Profesional & Akademik</span>
                  <div class="h-[1px] flex-1 bg-sigma-border"></div>
                </div>

                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Jabatan / Posisi</label>
                  <input v-model="form.position" type="text" placeholder="Ex: Guru Mata Pelajaran, Musyrif" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Status Kepegawaian</label>
                  <div class="relative group">
                    <select v-model="form.status" class="w-full p-4 pr-12 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text appearance-none cursor-pointer hover:border-sigma-emerald/30">
                      <option value="Guru Tetap" class="bg-slate-900 text-white">Guru Tetap</option>
                      <option value="Guru Tidak Tetap" class="bg-slate-900 text-white">Guru Tidak Tetap</option>
                      <option value="Pengabdian" class="bg-slate-900 text-white">Pengabdian</option>
                      <option value="Honor" class="bg-slate-900 text-white">Honor</option>
                    </select>
                    <ChevronRight class="absolute right-5 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted rotate-90 pointer-events-none group-hover:text-sigma-emerald transition-colors" />
                  </div>
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Pendidikan Terakhir</label>
                  <input v-model="form.education" type="text" placeholder="Ex: S1 Pendidikan Agama Islam" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Tahun Lulus & Tahun Masuk</label>
                  <div class="grid grid-cols-2 gap-3">
                    <input v-model="form.graduation_year" type="text" placeholder="Thn Lulus" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                    <input v-model="form.entry_year" type="text" placeholder="Thn Masuk" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
                  </div>
                </div>

                <!-- Section: Contact -->
                <div class="md:col-span-2 flex items-center gap-4 mt-4 mb-2">
                  <span class="text-xs font-black uppercase tracking-[0.2em] text-sigma-emerald">Kontak & Alamat</span>
                  <div class="h-[1px] flex-1 bg-sigma-border"></div>
                </div>

                <div class="space-y-2">
                  <label class="text-[10px] font-bold uppercase tracking-[0.2em] text-sigma-muted ml-1">Alamat Email</label>
                  <input v-model="form.email" type="email" placeholder="example@sigma.edu" class="w-full p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border focus:border-sigma-emerald outline-none transition-all text-sm font-bold text-sigma-text" />
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
                {{ isLoading ? 'Syncing...' : 'Simpan Data Guru' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </main>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.05); border-radius: 10px; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
