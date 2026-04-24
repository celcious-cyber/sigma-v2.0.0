<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  Users, Search, Building2, 
  ArrowLeft, Check, X, Pencil,
  Phone, User
} from 'lucide-vue-next'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// State
const students = ref<any[]>([])
const classrooms = ref<any[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const selectedClassFilter = ref('')

// Edit State
const isEditModalOpen = ref(false)
const isSubmitting = ref(false)
const selectedStudent = ref<any>(null)
const newClassroomID = ref<number | null>(null)

const fetchData = async () => {
  isLoading.value = true
  try {
    const [stRes, clRes] = await Promise.all([
      axios.get('/api/v1/admin/base/students'),
      axios.get('/api/v1/admin/edu/classrooms')
    ])
    students.value = stRes.data
    classrooms.value = clRes.data
  } catch (err) {
    console.error('Gagal mengambil data:', err)
  } finally {
    isLoading.value = false
  }
}

const openEditModal = (student: any) => {
  selectedStudent.value = student
  newClassroomID.value = student.classroom_id || null
  isEditModalOpen.value = true
}

const handleUpdateClass = async () => {
  if (!selectedStudent.value) return
  isSubmitting.value = true
  try {
    await axios.post('/api/v1/admin/edu/students/assign', {
      student_ids: [selectedStudent.value.ID],
      classroom_id: newClassroomID.value
    })
    await fetchData()
    isEditModalOpen.value = false
  } catch (err) {
    alert('Gagal memperbarui kelas santri.')
  } finally {
    isSubmitting.value = false
  }
}

const filteredStudents = computed(() => {
  return students.value.filter(s => {
    const search = searchQuery.value.toLowerCase()
    const matchSearch = s.name.toLowerCase().includes(search) || 
                       s.nis.includes(search) ||
                       (s.nisn && s.nisn.includes(search))
    
    const matchClass = !selectedClassFilter.value || s.classroom?.name === selectedClassFilter.value
    return matchSearch && matchClass
  })
})

onMounted(fetchData)
</script>

<template>
  <div class="p-8 space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="h-10 w-1 bg-emerald-500 rounded-full shadow-[0_0_15px_rgba(16,185,129,0.5)]"></div>
        <div>
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-sigma-muted mb-1">Manajemen Akademik</h2>
          <h1 class="text-2xl font-black italic uppercase tracking-tight">Data <span class="text-emerald-500">Santri</span></h1>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button @click="router.push('/portal')" 
                class="flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-sigma-surface border border-sigma-border text-sigma-muted hover:text-sigma-text hover:border-emerald-500/30 transition-all group shadow-sm">
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span class="text-[10px] font-black uppercase tracking-widest">Portal</span>
        </button>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <div class="lg:col-span-2 relative">
        <Search class="absolute left-6 top-1/2 -translate-y-1/2 w-5 h-5 text-sigma-muted" />
        <input v-model="searchQuery" type="text" placeholder="Cari Nama, NIS, atau NISN..." 
               class="w-full pl-16 pr-8 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight shadow-sm transition-all" />
      </div>
      
      <div class="relative">
        <Building2 class="absolute left-6 top-1/2 -translate-y-1/2 w-5 h-5 text-sigma-muted" />
        <select v-model="selectedClassFilter" 
                class="w-full pl-16 pr-8 py-5 bg-sigma-surface border border-sigma-border rounded-[2rem] focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight shadow-sm appearance-none cursor-pointer">
          <option value="">Semua Kelas</option>
          <option v-for="cls in classrooms" :key="cls.ID" :value="cls.name">{{ cls.name }}</option>
        </select>
      </div>

      <div class="bg-emerald-500 text-white rounded-[2rem] p-5 flex items-center justify-between shadow-xl shadow-emerald-500/20">
        <div class="flex items-center gap-4">
          <div class="w-10 h-10 rounded-2xl bg-white/20 flex items-center justify-center backdrop-blur-md">
            <Users class="w-5 h-5" />
          </div>
          <div>
            <p class="text-[8px] font-black uppercase tracking-widest text-white/60">Total Terfilter</p>
            <p class="text-xl font-black tracking-tighter">{{ filteredStudents.length }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-sigma-surface border border-sigma-border rounded-[3rem] overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-sigma-surface-alt/50 border-b border-sigma-border">
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Santri</th>
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">NIS / NISN</th>
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Kelas</th>
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Wali / Orang Tua</th>
              <th class="p-8 text-left text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Telp. Wali</th>
              <th class="p-8 text-right text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted">Penempatan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-sigma-border">
            <tr v-for="student in filteredStudents" :key="student.ID" 
                class="group hover:bg-sigma-surface-alt/30 transition-all">
              <td class="p-8">
                <div class="flex items-center gap-5">
                  <div class="w-12 h-12 rounded-2xl bg-sigma-surface-alt border border-sigma-border flex items-center justify-center font-black text-emerald-500 text-lg group-hover:scale-110 transition-transform shadow-sm">
                    {{ student.name.charAt(0) }}
                  </div>
                  <div>
                    <h4 class="text-sm font-black uppercase tracking-tighter italic text-sigma-text">{{ student.name }}</h4>
                    <span class="px-2 py-0.5 rounded-lg bg-sigma-surface-alt text-[8px] font-black uppercase tracking-widest border border-sigma-border text-sigma-muted">{{ student.gender === 'L' ? 'Laki-laki' : 'Perempuan' }}</span>
                  </div>
                </div>
              </td>
              <td class="p-8">
                <div class="space-y-1">
                  <p class="text-xs font-bold text-sigma-text tracking-tighter">NIS: {{ student.nis }}</p>
                  <p class="text-[10px] font-black text-sigma-muted tracking-widest">NISN: {{ student.nisn || '-' }}</p>
                </div>
              </td>
              <td class="p-8">
                <div class="flex items-center gap-2">
                  <div class="w-2 h-2 rounded-full bg-emerald-500" v-if="student.classroom"></div>
                  <span class="text-xs font-black uppercase tracking-widest italic" :class="student.classroom ? 'text-sigma-text' : 'text-sigma-muted'">
                    {{ student.classroom?.name || 'Belum Ada Kelas' }}
                  </span>
                </div>
              </td>
              <td class="p-8">
                <div class="flex items-center gap-3 text-sigma-text">
                  <User class="w-3.5 h-3.5 text-sigma-muted" />
                  <span class="text-xs font-bold">{{ student.parent_name || '-' }}</span>
                </div>
              </td>
              <td class="p-8">
                <div class="flex items-center gap-3 text-sigma-muted">
                  <Phone class="w-3.5 h-3.5" />
                  <span class="text-[10px] font-black tracking-widest">{{ student.parent_phone || '-' }}</span>
                </div>
              </td>
              <td class="p-8 text-right">
                <button @click="openEditModal(student)" 
                        class="px-5 py-2.5 rounded-xl bg-sigma-surface-alt border border-sigma-border text-sigma-muted hover:text-emerald-500 hover:border-emerald-500/30 transition-all text-[9px] font-black uppercase tracking-widest flex items-center gap-2 ml-auto">
                  <Pencil class="w-3 h-3" /> Atur Kelas
                </button>
              </td>
            </tr>
            <tr v-if="filteredStudents.length === 0">
              <td colspan="6" class="p-32 text-center opacity-30 italic">
                <Users class="w-16 h-16 mx-auto mb-4" />
                <p class="font-black uppercase tracking-[0.3em] text-[10px]">Data santri tidak ditemukan</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Edit Kelas Saja -->
    <div v-if="isEditModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-6 backdrop-blur-xl bg-black/40">
      <div class="bg-sigma-surface w-full max-w-md rounded-[3rem] border border-sigma-border shadow-2xl overflow-hidden animate-in zoom-in duration-300">
        <div class="p-8 bg-emerald-600 flex justify-between items-center border-b border-white/10 text-white">
          <div>
            <h3 class="font-black uppercase tracking-widest text-xs">Atur Kelas Santri</h3>
            <p class="text-emerald-100/60 text-[10px] mt-1 font-medium">Hanya mengubah penempatan kelas</p>
          </div>
          <button @click="isEditModalOpen = false" class="p-2 rounded-xl bg-white/10 text-white hover:bg-white/20 transition-all">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-10 space-y-8">
          <!-- Student Info (Read Only) -->
          <div class="flex items-center gap-5 p-6 rounded-[2rem] bg-sigma-surface-alt border border-sigma-border">
            <div class="w-14 h-14 rounded-2xl bg-sigma-surface flex items-center justify-center font-black text-emerald-500 text-xl border border-sigma-border">
              {{ selectedStudent?.name.charAt(0) }}
            </div>
            <div>
              <h4 class="text-sm font-black uppercase tracking-tighter italic text-sigma-text">{{ selectedStudent?.name }}</h4>
              <p class="text-[10px] font-black tracking-widest text-sigma-muted">NIS: {{ selectedStudent?.nis }}</p>
            </div>
          </div>

          <!-- Edit Classroom -->
          <div class="space-y-3">
            <label class="text-[10px] font-black uppercase tracking-[0.2em] text-sigma-muted ml-1 flex items-center gap-2">
              <Building2 class="w-3.5 h-3.5" /> Pilih Kelas Baru
            </label>
            <select v-model="newClassroomID" 
                    class="w-full px-6 py-4 bg-sigma-surface-alt border border-sigma-border rounded-2xl focus:outline-none focus:border-emerald-500/50 text-sm font-bold tracking-tight appearance-none">
              <option :value="null">-- Belum Ada Kelas --</option>
              <option v-for="cls in classrooms" :key="cls.ID" :value="cls.ID">{{ cls.name }}</option>
            </select>
          </div>

          <div class="flex gap-4 pt-4">
            <button @click="isEditModalOpen = false" class="flex-1 px-8 py-4 bg-sigma-surface-alt border border-sigma-border text-sigma-text rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-sigma-border transition-all">
              Batal
            </button>
            <button @click="handleUpdateClass" :disabled="isSubmitting"
                    class="flex-1 px-8 py-4 bg-emerald-600 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:bg-emerald-700 transition-all shadow-xl shadow-emerald-500/20 disabled:opacity-50 flex items-center justify-center gap-2">
              <Check v-if="!isSubmitting" class="w-4 h-4" />
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              Simpan Perubahan
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(16, 185, 129, 0.1); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(16, 185, 129, 0.2); }
</style>
