<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Search, Plus, Trash2, 
  Building2, ArrowLeft, Check, X,
  Users
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// State
const classrooms = ref<any[]>([])
const teachers = ref<any[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  level: '10',
  type: 'KMI',
  gender: 'Laki-laki',
  wali_kelas_id: null as number | null
})

const fetchClassrooms = async () => {
  try {
    const response = await axios.get('/api/v1/admin/edu/classrooms')
    classrooms.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  } finally {
    isLoading.value = false
  }
}

const fetchTeachers = async () => {
  try {
    const response = await axios.get('/api/v1/admin/base/teachers')
    teachers.value = response.data
  } catch (err) {
    console.error('Gagal mengambil data guru:', err)
  }
}

const openModal = (classroom: any = null) => {
  if (classroom) {
    editingId.value = classroom.ID
    form.value = { 
      name: classroom.name, 
      level: classroom.level, 
      type: classroom.type, 
      gender: classroom.gender,
      wali_kelas_id: classroom.wali_kelas_id 
    }
  } else {
    editingId.value = null
    form.value = { 
      name: '', 
      level: '10', 
      type: 'KMI', 
      gender: 'Laki-laki', 
      wali_kelas_id: null 
    }
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const handleSubmit = async () => {
  isSubmitting.value = true
  try {
    const payload = {
      ...form.value,
      level: String(form.value.level)
    }

    if (editingId.value) {
      await axios.put(`/api/v1/admin/edu/classrooms/${editingId.value}`, payload)
    } else {
      await axios.post('/api/v1/admin/edu/classrooms', payload)
    }
    await fetchClassrooms()
    closeModal()
  } catch (err) {
    alert('Gagal menyimpan data kelas.')
  } finally {
    isSubmitting.value = false
  }
}

const deleteClassroom = async (id: number) => {
  if (!confirm('Yakin ingin menghapus kelas ini?')) return
  try {
    await axios.delete(`/api/v1/admin/edu/classrooms/${id}`)
    await fetchClassrooms()
  } catch (err) {
    alert('Gagal menghapus data.')
  }
}

const filteredClassrooms = computed(() => {
  return classrooms.value.filter(c => 
    c.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    c.level.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const levelOptions = computed(() => {
  if (form.value.type === 'KMI') {
    return [1, 2, 3, 4, 5, 6]
  } else {
    return [1, 2, 3, 4]
  }
})

// Student Management State
const isStudentModalOpen = ref(false)
const selectedClassroom = ref<any>(null)
const classStudents = ref<any[]>([])
const unassignedStudents = ref<any[]>([])
const studentSearch = ref('')

const openStudentModal = async (classroom: any) => {
  selectedClassroom.value = classroom
  isStudentModalOpen.value = true
  await fetchStudentData()
}

const fetchStudentData = async () => {
  if (!selectedClassroom.value) return
  try {
    const [classRes, unRes] = await Promise.all([
      axios.get(`/api/v1/admin/edu/classrooms/${selectedClassroom.value.ID}/students`),
      axios.get('/api/v1/admin/edu/students/unassigned')
    ])
    classStudents.value = classRes.data
    unassignedStudents.value = unRes.data
  } catch (err) {
    console.error('Gagal memuat data santri:', err)
  }
}

const addToClass = async (studentId: number) => {
  try {
    await axios.post('/api/v1/admin/edu/students/assign', {
      student_ids: [studentId],
      classroom_id: selectedClassroom.value.ID
    })
    await fetchStudentData()
  } catch (err) {
    alert('Gagal menambahkan santri ke kelas.')
  }
}

const removeFromClass = async (studentId: number) => {
  if (!confirm('Keluarkan santri dari kelas ini?')) return
  try {
    await axios.post('/api/v1/admin/edu/students/assign', {
      student_ids: [studentId],
      classroom_id: null
    })
    await fetchStudentData()
  } catch (err) {
    alert('Gagal mengeluarkan santri.')
  }
}

const filteredUnassigned = computed(() => {
  return unassignedStudents.value.filter(s => 
    s.name.toLowerCase().includes(studentSearch.value.toLowerCase())
  )
})

onMounted(() => {
  fetchClassrooms()
  fetchTeachers()
})
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-blue-500 rounded-full"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Ruang Kelas</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Data <span class="text-blue-500">Kelas</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-blue-500/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="flex flex-col md:flex-row gap-4 justify-between items-center bg-sigma-surface p-4 rounded-[2rem] border border-sigma-border shadow-sm">
      <div class="relative w-full md:w-96">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-sigma-muted" />
        <input v-model="searchQuery" 
               type="text" 
               placeholder="Cari Nama Kelas atau Level..." 
               class="w-full pl-12 pr-4 py-3 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight transition-all" />
      </div>
      <button @click="openModal()" 
              class="w-full md:w-auto px-8 py-3.5 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 flex items-center justify-center gap-3">
        <Plus class="w-4 h-4" /> Tambah Kelas Baru
      </button>
    </div>

    <!-- Classes Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="classroom in filteredClassrooms" :key="classroom.ID" 
           class="bg-sigma-surface border border-sigma-border p-8 rounded-[2.5rem] relative overflow-hidden group hover:border-blue-500/30 transition-all shadow-sm">
        
        <!-- Decoration -->
        <div class="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-all">
          <Building2 class="w-24 h-24 rotate-12" />
        </div>

        <div class="flex flex-col h-full space-y-6">
          <div class="flex items-start justify-between">
            <div class="flex flex-col">
              <span class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Kelas {{ classroom.level }}</span>
              <h3 class="text-2xl font-black uppercase italic tracking-tighter text-sigma-text group-hover:text-blue-500 transition-colors">
                {{ classroom.name }}
              </h3>
            </div>
            <div class="px-3 py-1 rounded-xl bg-blue-500/10 text-blue-500 text-[8px] font-black uppercase tracking-widest border border-blue-500/10">
              {{ classroom.type }}
            </div>
          </div>

          <div class="space-y-4 flex-1">
            <div class="flex items-center gap-3 p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border">
              <div class="w-10 h-10 rounded-xl bg-sigma-surface flex items-center justify-center text-sigma-muted">
                <Users class="w-5 h-5" />
              </div>
              <div class="flex flex-col">
                <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Wali Kelas</span>
                <span class="text-xs font-bold text-sigma-text">{{ classroom.wali_kelas?.name || 'Belum Ditentukan' }}</span>
              </div>
            </div>

            <div class="flex items-center gap-4 px-4">
              <div class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 rounded-full" 
                     :class="classroom.gender === 'Laki-laki' ? 'bg-blue-500' : (classroom.gender === 'Perempuan' ? 'bg-pink-500' : 'bg-emerald-500')"></div>
                <span class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">{{ classroom.gender }}</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col gap-2 pt-4 border-t border-sigma-border">
            <button @click="openStudentModal(classroom)" 
                    class="w-full py-3 rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition-all text-[10px] font-black uppercase tracking-widest shadow-lg shadow-blue-500/20 flex items-center justify-center gap-2">
              <Users class="w-3.5 h-3.5" /> Kelola Santri
            </button>
            <div class="flex gap-2">
              <button @click="openModal(classroom)" class="flex-1 py-3 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-blue-500 hover:bg-blue-500/5 transition-all text-[10px] font-black uppercase tracking-widest border border-sigma-border">Edit</button>
              <button @click="deleteClassroom(classroom.ID)" class="flex-1 py-3 rounded-xl bg-sigma-surface-alt text-sigma-muted hover:text-red-500 hover:bg-red-500/5 transition-all text-[10px] font-black uppercase tracking-widest border border-sigma-border">Hapus</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="filteredClassrooms.length === 0 && !isLoading" 
           class="md:col-span-2 lg:col-span-4 py-20 bg-sigma-surface border border-sigma-border border-dashed rounded-[3rem] flex flex-col items-center justify-center gap-4 opacity-30">
        <Building2 class="w-16 h-16" />
        <p class="font-black uppercase tracking-widest text-xs">Belum ada data kelas</p>
      </div>
    </div>

    <!-- Modal Add/Edit (Original) -->
    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-lg rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-gradient-to-r from-blue-600 to-blue-900 flex justify-between items-center border-b border-white/10">
          <div>
            <h3 class="text-white font-black uppercase tracking-widest text-xs">{{ editingId ? 'Update' : 'Tambah' }} Kelas</h3>
            <p class="text-blue-100/60 text-[10px] mt-1 font-medium">Manajemen Struktur Kelas Sigmaedu</p>
          </div>
          <button @click="closeModal" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Nama Kelas</label>
            <input v-model="form.name" type="text" placeholder="cnth. 1 Putra" 
                   class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Level / Tingkatan</label>
              <select v-model="form.level" 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
                <option v-for="lvl in levelOptions" :key="lvl" :value="lvl">Kelas {{ lvl }}</option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Jenis Kurikulum</label>
              <select v-model="form.type" 
                      class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
                <option value="KMI">KMI</option>
                <option value="Eksperiment">Eksperiment</option>
              </select>
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Wali Kelas</label>
            <select v-model="form.wali_kelas_id" 
                    class="w-full px-5 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-blue-500/50 text-sm font-bold tracking-tight appearance-none">
              <option :value="null">-- Pilih Wali Kelas --</option>
              <option v-for="t in teachers" :key="t.ID" :value="t.user_id">{{ t.name }}</option>
            </select>
          </div>

          <div class="flex items-center gap-4">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1">Gender:</label>
            <div class="flex items-center gap-4">
              <label class="flex items-center gap-2 cursor-pointer group">
                <input type="radio" v-model="form.gender" value="Laki-laki" class="hidden" />
                <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center transition-all"
                     :class="form.gender === 'Laki-laki' ? 'border-blue-500 bg-blue-500' : 'group-hover:border-blue-500/50'">
                  <div class="w-1.5 h-1.5 bg-white rounded-full" v-if="form.gender === 'Laki-laki'"></div>
                </div>
                <span class="text-xs font-bold" :class="form.gender === 'Laki-laki' ? 'text-blue-500' : 'text-sigma-muted'">Laki-laki</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer group">
                <input type="radio" v-model="form.gender" value="Perempuan" class="hidden" />
                <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center transition-all"
                     :class="form.gender === 'Perempuan' ? 'border-pink-500 bg-pink-500' : 'group-hover:border-pink-500/50'">
                  <div class="w-1.5 h-1.5 bg-white rounded-full" v-if="form.gender === 'Perempuan'"></div>
                </div>
                <span class="text-xs font-bold" :class="form.gender === 'Perempuan' ? 'text-pink-500' : 'text-sigma-muted'">Perempuan</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer group">
                <input type="radio" v-model="form.gender" value="Campuran" class="hidden" />
                <div class="w-5 h-5 rounded-full border-2 border-sigma-border flex items-center justify-center transition-all"
                     :class="form.gender === 'Campuran' ? 'border-emerald-500 bg-emerald-500' : 'group-hover:border-emerald-500/50'">
                  <div class="w-1.5 h-1.5 bg-white rounded-full" v-if="form.gender === 'Campuran'"></div>
                </div>
                <span class="text-xs font-bold" :class="form.gender === 'Campuran' ? 'text-emerald-500' : 'text-sigma-muted'">Campuran</span>
              </label>
            </div>
          </div>

          <div class="flex gap-4 pt-4">
            <button @click="closeModal" class="flex-1 px-8 py-4 bg-sigma-surface-alt border border-sigma-border text-sigma-text rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-border transition-all">
              Batal
            </button>
            <button @click="handleSubmit" :disabled="isSubmitting"
                    class="flex-1 px-8 py-4 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20 disabled:opacity-50 flex items-center justify-center gap-2">
              <Check v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ editingId ? 'Simpan Perubahan' : 'Tambah Kelas' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Kelola Santri (New) -->
    <div v-if="isStudentModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-4xl rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300 flex flex-col max-h-[90vh]">
        <div class="p-8 bg-blue-600 flex justify-between items-center border-b border-white/10">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-white/10 flex items-center justify-center text-white border border-white/10">
              <Users class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-white font-black uppercase tracking-widest text-xs">Anggota Kelas: {{ selectedClassroom?.name }}</h3>
              <p class="text-blue-100/60 text-[10px] mt-1 font-medium uppercase tracking-widest">Kelola Daftar Santri di Kelas Ini</p>
            </div>
          </div>
          <button @click="isStudentModalOpen = false" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="flex-1 overflow-hidden grid grid-cols-2">
          <!-- Current Members -->
          <div class="p-8 border-r border-sigma-border flex flex-col gap-6">
            <div class="flex items-center justify-between">
              <h4 class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Santri Terdaftar ({{ classStudents.length }})</h4>
            </div>
            <div class="flex-1 overflow-y-auto custom-scrollbar pr-4 space-y-3">
              <div v-for="s in classStudents" :key="s.ID" 
                   class="flex items-center justify-between p-4 rounded-2xl bg-sigma-surface-alt border border-sigma-border group hover:border-red-500/30 transition-all">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-sigma-surface flex items-center justify-center text-[10px] font-black border border-sigma-border">{{ s.name.charAt(0) }}</div>
                  <span class="text-xs font-bold text-sigma-text">{{ s.name }}</span>
                </div>
                <button @click="removeFromClass(s.ID)" class="p-2 rounded-lg text-sigma-muted hover:text-red-500 hover:bg-red-500/10 transition-all opacity-0 group-hover:opacity-100">
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
              <div v-if="classStudents.length === 0" class="py-10 text-center opacity-30 italic text-xs uppercase tracking-widest">Kelas Kosong</div>
            </div>
          </div>

          <!-- Unassigned Students -->
          <div class="p-8 flex flex-col gap-6 bg-sigma-surface-alt/20">
            <div class="flex flex-col gap-4">
              <h4 class="text-[10px] font-black uppercase tracking-widest text-sigma-muted">Tersedia (Belum Ada Kelas)</h4>
              <div class="relative">
                <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-3 h-3 text-sigma-muted" />
                <input v-model="studentSearch" type="text" placeholder="Cari santri..." class="w-full pl-10 pr-4 py-2.5 bg-sigma-surface border border-sigma-border rounded-xl text-[10px] font-bold focus:outline-none" />
              </div>
            </div>
            <div class="flex-1 overflow-y-auto custom-scrollbar pr-4 space-y-3">
              <div v-for="s in filteredUnassigned" :key="s.ID" 
                   class="flex items-center justify-between p-4 rounded-2xl bg-sigma-surface border border-sigma-border group hover:border-blue-500/30 transition-all">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-sigma-surface-alt flex items-center justify-center text-[10px] font-black border border-sigma-border text-sigma-muted">{{ s.name.charAt(0) }}</div>
                  <span class="text-xs font-bold text-sigma-text">{{ s.name }}</span>
                </div>
                <button @click="addToClass(s.ID)" class="p-2 rounded-lg bg-blue-500/10 text-blue-500 hover:bg-blue-500 hover:text-white transition-all shadow-sm">
                  <Plus class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="p-6 bg-sigma-surface-alt/50 border-t border-sigma-border flex justify-end">
          <button @click="isStudentModalOpen = false" class="px-8 py-3 bg-blue-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-blue-700 shadow-lg shadow-blue-500/20">Selesai</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 130, 246, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(59, 130, 246, 0.2); }
</style>
